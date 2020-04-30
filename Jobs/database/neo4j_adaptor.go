package database

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"log"
	"time"
)
import "idp/Jobs/models"

const (
	getAllServiceCategories = "MATCH (sc:ServiceCategory) RETURN properties(sc)"
	getAllSkillCategories = "MATCH (sc:ServiceCategory)<-[r:IsSubclassOf]-(sk:SkillCategory) WHERE (sc.id = $ID) RETURN properties(sk)"
	getSkillsByCategory = "MATCH (sk:Skill)<-[r:IsSkillOf]-(skc:SkillCategory) WHERE (skc.id = $ID) RETURN properties(sk)"

	createJobNodeQuery = "CREATE (a: Job {$prop})"
	createJobSCRel = "MATCH (j:Job), (sc:ServiceCategory) WHERE (j.id = $JID) AND (sc.id = $SCID) CREATE (j)-[r:HasCategory]->(sc)"
	createJobSKCRel = "MATCH (j:Job), (skc:SkillCategory) WHERE (j.id = $JID) AND (skc.id = $SKCID) CREATE (j)-[r:HasSkillCategory]->(sc)"

	getJobByIDQuery = "MATCH(j: Job) WHERE(j.ID=$ID) RETURN properties(j)"
	getJobsWithFilter = "MATCH(j: job) WHERE(j.title contains $title AND wage >= $minWage) RETURN properties(j)"
)

type Database interface {
	GetAllServiceCategories() ([]*models.ServiceCategory, error)
	GetSkillCategoriesByServiceCategory(*models.ServiceCategory) ([]*models.SkillCategory, error)
	GetSkillsByCategory(category *models.SkillCategory) ([]*models.Skill, error)
	CreateJob(job *models.NewJob, sc *models.ServiceCategory, skc *models.SkillCategory) error
	GetJobByID(jobID string) (*models.Job, error)
	GetJobsWithFilter(title string, wageMin, eRatingMin float32) ([]*models.Job, error)
}

type Neo4jDatabase struct {
	Driver neo4j.Driver
}

func NewNeo4jDatabase(uri, username, password string) *Neo4jDatabase {
	time.Sleep(6 * time.Second)
	log.Println("Connecting to Neo4j...")
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Println("Could not connect to Neo4j: ", err.Error())
		panic(err.Error())
	}
	log.Println("Connected succesfully to Neo4j!")

	return &Neo4jDatabase{
		Driver: driver,
	}
}

func (nj *Neo4jDatabase) GetAllServiceCategories() ([]*models.ServiceCategory, error) {
	serviceCategories := make([]*models.ServiceCategory, 0)

	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	qResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getAllServiceCategories, map[string]interface{} {})

		if err != nil {
			return nil, err
		}

		return qResult, nil
	})

	qServiceCategoriesResult := qResult.(neo4j.Result)

	var serviceCategoryData *models.ServiceCategory

	for qServiceCategoriesResult.Next() {
		serviceCategoryData = &models.ServiceCategory{}

		if err = mapstructure.Decode(qServiceCategoriesResult.Record().GetByIndex(0), serviceCategoryData); err != nil {
			return nil, err
		}

		serviceCategories = append(serviceCategories, serviceCategoryData)
	}

	return serviceCategories, nil
}

func (nj *Neo4jDatabase) GetSkillCategoriesByServiceCategory(category *models.ServiceCategory) ([]*models.SkillCategory, error) {
	skillCategories := make([]*models.SkillCategory, 0)

	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	qResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getAllSkillCategories, map[string]interface{} {
			"ID": category.ID,
		})

		if err != nil {
			return nil, err
		}

		return qResult, nil
	})

	qSkillCategoriesResult := qResult.(neo4j.Result)

	var skillCategoryData *models.SkillCategory

	for qSkillCategoriesResult.Next() {
		skillCategoryData = &models.SkillCategory{}

		if err = mapstructure.Decode(qSkillCategoriesResult.Record().GetByIndex(0), skillCategoryData); err != nil {
			return nil, err
		}

		skillCategories = append(skillCategories, skillCategoryData)
	}

	return skillCategories, nil
}

func (nj *Neo4jDatabase) GetSkillsByCategory(category *models.SkillCategory) ([]*models.Skill, error) {
	skills := make([]*models.Skill, 0)

	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	qResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getSkillsByCategory, map[string]interface{} {
			"ID": category.ID,
		})

		if err != nil {
			return nil, err
		}

		return qResult, nil
	})

	qSkillsResult := qResult.(neo4j.Result)

	var skillData *models.Skill

	for qSkillsResult.Next() {
		skillData = &models.Skill{}

		if err = mapstructure.Decode(qSkillsResult.Record().GetByIndex(0), skillData); err != nil {
			return nil, err
		}

		skills = append(skills, skillData)
	}

	return skills, nil
}

func (nj *Neo4jDatabase) CreateJob(job *models.NewJob, sc *models.ServiceCategory, skc *models.SkillCategory) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	newJobDetails := make(map[string]interface{})
	err = mapstructure.Decode(job, newJobDetails)
	if err != nil {
		return err
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		_, err = tx.Run(createJobNodeQuery, map[string]interface{}{
			"prop": newJobDetails,
		})

		if err != nil {
			return nil, err
		}

		_, err = tx.Run(createJobSCRel, map[string]interface{}{
			"JID": job.ID,
			"SCID": sc.ID,
		})

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = tx.Run(createJobSKCRel, map[string]interface{}{
			"JID": job.ID,
			"SKCID": skc.ID,
		})

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (nj *Neo4jDatabase) GetJobByID(jobID string) (*models.Job, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	qResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getJobByIDQuery, map[string]interface{}{
			"ID": jobID,
		})

		if err != nil {
			return nil, err
		}

		return qResult, nil
	})

	qJobResult := qResult.(neo4j.Result)

	var jobData *models.Job
	for qJobResult.Next() {
		jobData = &models.Job{}

		if err = mapstructure.Decode(qJobResult.Record().GetByIndex(0), jobData); err != nil {
			return nil, err
		}
	}

	return jobData, nil
}

func (nj *Neo4jDatabase) GetJobsWithFilter(title string, wageMin, eRatingMin float32) ([]*models.Job, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	qFilterArg := make(map[string]interface{})
	qFilterArg["title"] = title
	qFilterArg["wageMin"] = wageMin

	qResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getJobsWithFilter, qFilterArg)

		if err != nil {
			return nil, err
		}

		return qResult, nil
	})

	qJobsResults := qResult.(neo4j.Result)

	jobsData := make([]*models.Job, 0)
	for qJobsResults.Next() {
		jobData := &models.Job{}

		if err = mapstructure.Decode(qJobsResults.Record().GetByIndex(0), jobData); err != nil {
			return nil, err
		}

		jobsData = append(jobsData, jobData)
	}
	
	return jobsData, nil
}