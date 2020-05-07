package database

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"idp/Profiles/models"
	"log"
// 	"time"
)

const (
	createEmployerProfile = "CREATE (e:Employer $props) RETURN e"
	createFreelancerProfile = "CREATE (f:Freelancer $props) RETURN f"

	editEmployerProfile = "MATCH (e:Employer {ID: $euid}) SET e += $props RETURN 'OK'"

	editFreelancerProfile = "MATCH (f:Freelancer {ID: $fuid}) SET f += $props RETURN 'OK'"
	addFreelancerSkillCategories = "MATCH (f: Freelancer {ID: $fuid}), (skc: SkillCategory) WHERE (skc.id in $SKCIDS) MERGE (f)-[r:HasSkillCategory]->(skc) RETURN 'OK'"
	addSkill = "MERGE (s:Skill {ID: $sid, Skill: $skill}) RETURN s"
	testQuery = "MATCH (s:Skill {ID: $sid}) return properties(s)"
	addSkillToCategory = "MATCH (s:Skill {ID: $sid}), (skc:SkillCategory {id: $skcid}) merge (s)-[r:HasCategory]->(skc) return 'OK'"
	addFreelancerSkills = "MATCH (f: Freelancer {ID: $fuid}), (s: Skill) WHERE (s.Skill in $SKIDS) MERGE (f)-[r:HasSkill]->(s) RETURN 'OK'"

	getEmployerProfile = "MATCH (e: Employer {ID: $euid}) optional match (e)-[r:Posted]->(j:Job) optional match (e)-[a:Posted]->(h: Job {Status: $history}) with e, j, r, sum(h.Wage) as money_spent RETURN properties(e), count(r), money_spent"
	getFreelancerProfile = "MATCH (f: Freelancer {ID: $fuid}) optional match (f)-[]->(skc: SkillCategory) optional match (f)-[]->(sk: Skill) RETURN properties(f), collect(distinct properties(skc)), collect(distinct properties(sk))"
)

type NoDataFoundError struct {
	Err string
}

func (err NoDataFoundError) Error() string {
	return err.Err
}

type Database interface {
	CreateEmployerProfile(request *models.EditEmployerRequest) error
	CreateFreelancerProfile(request *models.EditFreelancerRequest) error
	EditEmployerProfile(*models.EditEmployerRequest) error
	EditFreelancerProfile(*models.EditFreelancerRequest, []*models.Skill, []*models.SkillCategory) error
	GetEmployerProfile(euid string) (*models.Employer, error)
	GetFreelancerProfile(fuid string) (*models.Freelancer, error)
}

type Neo4jDatabase struct {
	Driver neo4j.Driver
}

func NewNeo4jDatabase(uri, username, password string) *Neo4jDatabase {
// 	time.Sleep(10 * time.Second)
	log.Println("Connecting to Neo4j...")
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		log.Println("Could not connect to Neo4j: ", err.Error())
		panic(err.Error())
	}
	log.Println("Connected succesfully to Neo4j!")

	return &Neo4jDatabase{
		Driver: driver,
	}
}

func (nj *Neo4jDatabase) CreateEmployerProfile(employer *models.EditEmployerRequest) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	employerProfile := make(map[string]interface{})
	err = mapstructure.Decode(employer, &employerProfile)
	if err != nil {
		return err
	}

	employerProfile["Rating"] = 0

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(createEmployerProfile, map[string]interface{}{
			"props": employerProfile,
		}); if err != nil {
			return nil, err
		}

		return nil, nil
	}); if err != nil {
		return err
	}

	return nil
}

func (nj *Neo4jDatabase) CreateFreelancerProfile(freelancer *models.EditFreelancerRequest) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	freelancerProfile := make(map[string]interface{})
	err = mapstructure.Decode(freelancer, &freelancerProfile)
	if err != nil {
		return err
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(createFreelancerProfile, map[string]interface{}{
			"props": freelancerProfile,
		}); if err != nil {
			return nil, err
		}

		return nil, nil
	}); if err != nil {
		return err
	}

	return nil
}

func (nj *Neo4jDatabase) EditEmployerProfile(employer *models.EditEmployerRequest) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}

	employerProfile := make(map[string]interface{})
	err = mapstructure.Decode(employer, &employerProfile)
	if err != nil {
		log.Println("Could not decode employer profile")
		return err
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		qResult, err := tx.Run(editEmployerProfile, map[string]interface{} {
			"props": employerProfile,
		})

		if err != nil || !qResult.Next() {
			return nil, err
		}

		return nil, nil
	}); if err != nil {
		return err
	}

	return nil
}

func (nj *Neo4jDatabase) EditFreelancerProfile(f *models.EditFreelancerRequest, sks []*models.Skill, skc []*models.SkillCategory) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}

    f = &models.EditFreelancerRequest{
        ID: f.ID,
        LastName: f.LastName,
        FirstName: f.FirstName,
        Description: f.Description,
    }

	freelancerProfile := make(map[string]interface{})
	err = mapstructure.Decode(f, &freelancerProfile)
	if err != nil {
		log.Println("Could not decode freelancer profile")
		return err
	}

	skillsCategories := make([]string, 0)
	for _, s := range skc {
		skillsCategories = append(skillsCategories, s.ID)
	}

	skills := make([]string, 0)
	for _, sk := range sks {
		skills = append(skills, sk.Skill)
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		qResult, err := tx.Run(editFreelancerProfile, map[string]interface{} {
		    "fuid": f.ID,
			"props": freelancerProfile,
		}); if err != nil {
			tx.Close()
			return nil, err
		}; if !qResult.Next() {
			tx.Close()
	        return nil, NoDataFoundError{Err: "Freelancer not found"}
		}

		qResult, err = tx.Run(addFreelancerSkillCategories, map[string]interface{} {
			"fuid": f.ID,
			"SKCIDS": skillsCategories,
		}); if err != nil {
			tx.Rollback()
			tx.Close()
			return nil, err
		}; if !qResult.Next() {
		    tx.Rollback()
			tx.Close()
		    return nil, NoDataFoundError{Err: "Freelancer or skill category not found"}
		}

		return nil, nil
	}); if err != nil {
	    log.Println("207" + err.Error())
	    return err
	}


    for _, sk := range sks {
       _, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
           qResult, err := tx.Run(addSkill, map[string]interface{} {
               "sid": sk.ID,
               "skill": sk.Skill,
           }); if err != nil {
               return nil, err
           }; if !qResult.Next() {
			   return nil, NoDataFoundError{Err: "Could not add skill"}
           }

           qResult, err = tx.Run(addSkillToCategory, map[string]interface{} {
               "sid": sk.ID,
               "skcid": sk.SCID.ID,
           }); if err != nil {
               return nil, err
           }; if !qResult.Next() {
               return nil, NoDataFoundError{Err: "Could not find skill or skill category"}
           }

           return nil, nil
       }); if err != nil {
       	continue
       }
    }

    _, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		qResult, err := tx.Run(addFreelancerSkills, map[string]interface{} {
			"fuid": f.ID,
			"SKIDS": skills,
		}); if err != nil {
		    log.Println(err.Error())
			tx.Rollback()
			return nil, err
		}; if !qResult.Next() {
		    log.Println(err.Error())
            tx.Rollback()
            return nil, NoDataFoundError{Err: "Freelancer or skill not found"}
		}

		return nil, nil
	}); if err != nil {
	    log.Println("254" + err.Error())
		return err
	}

	return nil
}

func (nj *Neo4jDatabase) GetEmployerProfile(euid string) (*models.Employer, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}

	employer, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		qResult, err := tx.Run(getEmployerProfile, map[string]interface{} {
			"euid": euid,
			"history": models.JobStatusHistory,
		}); if err != nil {
			return nil, err
		}; if !qResult.Next() {
	        return nil, NoDataFoundError{Err: "Employer not found"}
		}

		var employer *models.Employer

        employer = &models.Employer{}

        if err = mapstructure.Decode(qResult.Record().GetByIndex(0), employer); err != nil {
            return nil, err
        }

        employer.JobsPosted = (int32)(qResult.Record().GetByIndex(1).(int64))
        employer.MoneySpent = (float32)(qResult.Record().GetByIndex(2).(float64))

		return employer, nil
	}); if err != nil {
		return nil, err
	}

	return employer.(*models.Employer), nil
}

func (nj *Neo4jDatabase) GetFreelancerProfile(fuid string) (*models.Freelancer, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}

	freelancer, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		qResult, err := tx.Run(getFreelancerProfile, map[string]interface{} {
			"fuid": fuid,
		}); if err != nil {
			return nil, err
		}; if !qResult.Next() {
		    return nil, NoDataFoundError{Err: "Freelancer not found"}
		}

		var freelancer *models.Freelancer

        freelancer = &models.Freelancer{}
        skillCategories := make([]*models.SkillCategory, 0)
        skills := make([]*models.Skill, 0)

        if err = mapstructure.Decode(qResult.Record().GetByIndex(0), freelancer); err != nil {
            return nil, err
        }
        if err = mapstructure.Decode(qResult.Record().GetByIndex(1), &skillCategories); err != nil {
            return nil, err
        }
        if err = mapstructure.Decode(qResult.Record().GetByIndex(2), &skills); err != nil {
            return nil, err
        }

        freelancer.SkillCategories = skillCategories
        freelancer.Skills = skills

		return freelancer, nil
	}); if err != nil {
		return nil, err
	}

	return freelancer.(*models.Freelancer), nil
}
