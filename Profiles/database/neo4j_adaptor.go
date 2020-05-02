package database

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"idp/Profiles/models"
	"log"
	"time"
)

const (
	editEmployerProfile = "MATCH (e:Employer {ID: $euid}) SET e += $props RETURN 'OK'"

	editFreelancerProfile = "MATCH (f:Freelancer {ID: $fuid}) SET f += $props RETURN 'OK'"
	addFreelancerSkillCategories = "MATCH (f: Freelancer {ID: $fuid}), (skc: SkillCategory) WHERE (skc.id in $SKCIDS) MERGE (f)-[r:HasSkillCategory]->(skc) RETURN 'OK'"
	addSkill = "MERGE (s: Skill $skill)"
	addFreelancerSkills = "MATCH (f: Freelancer {ID: $fuid}), (s: Skill) WHERE (s.id in $SKIDS) MERGE (f)-[r:HasSkill]->(sk) RETURN 'OK'"

	getEmployerProfile = "MATCH (e: Employer {ID: $euid}), (e)-[r:PostedJob]->(j:Job), (e)-[r:PostedJob]->(h: Job {status = $history}), RETURN properties(e), count(r), sum(h.Wage)"
	getFreelancerProfile = "MATCH (f: Freelancer {ID: $fuid}), (f)-[]->(skc: SkillCategory), (f)-[]->(sk: Skill) RETURN properties(f), collect(distinct properties(skc)), collect(distinct properties(sk))"
)

type Database interface {
	EditEmployerProfile(*models.EditEmployerRequest) error
	EditFreelancerProfile(*models.EditFreelancerRequest, []*models.Skill, []*models.SkillCategory) error
	GetEmployerProfile(euid string) (*models.Employer, error)
	GetFreelancerProfile(fuid string) (*models.Freelancer, error)
}

type Neo4jDatabase struct {
	Driver neo4j.Driver
}

func NewNeo4jDatabase(uri, username, password string) *Neo4jDatabase {
	time.Sleep(10 * time.Second)
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
		skills = append(skills, sk.ID)
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		qResult, err := tx.Run(editFreelancerProfile, map[string]interface{} {
			"props": freelancerProfile,
		})

		if err != nil || !qResult.Next() {
			return nil, err
		}

		qResult, err = tx.Run(addFreelancerSkillCategories, map[string]interface{} {
			"fuid": f.ID,
			"SKCIDS": skillsCategories,
		}); if err != nil || !qResult.Next() {
			tx.Rollback()
			return nil, err
		}

		for _, sk := range sks {
			qResult, err = tx.Run(addSkill, map[string]interface{} {
				"skill": map[string]interface{}{"id": sk.ID, "Skill": sk.Skill},
			}); if err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		qResult, err = tx.Run(addFreelancerSkills, map[string]interface{} {
			"fuid": f.ID,
			"SKIDS": skills,
		}); if err != nil || !qResult.Next() {
			tx.Rollback()
			return nil, err
		}

		return nil, nil
	}); if err != nil {
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
		}); if err != nil || !qResult.Next() {
			return nil, err
		}

		var employer *models.Employer
		if qResult.Next() {
			employer = &models.Employer{}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), employer); err != nil {
				return nil, err
			}

			employer.JobsPosted = qResult.Record().GetByIndex(1).(int32)
			employer.MoneySpent = qResult.Record().GetByIndex(2).(float32)
		}

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
		}); if err != nil || !qResult.Next() {
			return nil, err
		}

		var freelancer *models.Freelancer
		if qResult.Next() {
			freelancer = &models.Freelancer{}
			skillCategories := make([]*models.SkillCategory, 0)
			skills := make([]*models.Skill, 0)

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), freelancer); err != nil {
				return nil, err
			}
			if err = mapstructure.Decode(qResult.Record().GetByIndex(1), skillCategories); err != nil {
				return nil, err
			}
			if err = mapstructure.Decode(qResult.Record().GetByIndex(2), skills); err != nil {
				return nil, err
			}

			freelancer.SkillCategories = skillCategories
			freelancer.Skills = skills
		}

		return freelancer, nil
	}); if err != nil {
		return nil, err
	}

	return freelancer.(*models.Freelancer), nil
}
