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

	createJobNodeQuery = "CREATE (a: Job $PROP)"
	createJobSCRel = "MATCH (j:Job), (sc:ServiceCategory) WHERE (j.ID = $JID) AND (sc.id in $SCIDS) MERGE (j)-[r:HasCategory]->(sc) RETURN count(*)"
	createJobSKCRel = "MATCH (j:Job), (skc:SkillCategory) WHERE (j.ID = $JID) AND (skc.id in $SKCIDS) MERGE (j)-[r:HasSkillCategory]->(skc) RETURN count(*)"
	assignJobToEmployer = "MATCH (j: Job {ID: $JID}), (e: Employer {ID: $euid}) MERGE (j)<-[r:Posted]-(e) RETURN 'OK'"

	getJobByIDQuery = "match (j:Job {ID: $ID}), (j)-[]->(sc:ServiceCategory), (j)-[]->(skc:SkillCategory), " +
		"(e:Employer)-[p:Posted]->(j)" +
		"optional match (j)-[]->(sk:Skill) optional match (e)-[p1:Posted]->(j1: Job {Status: $history}) " +
		"with j, sc, skc, e, sk, sum(j1.Wage) as money_spent " +
		"optional match (f:Freelancer)-[a:AppliedTo]->(j) " +
		"return properties(j), collect(distinct properties(sc)), collect(distinct properties(skc))," +
		"collect(distinct properties(sk)), toFloat(money_spent), toFloat(e.Rating), count(a)"

	getJobsWithFilter = "MATCH (j: Job {Status: $status}), (j)-[]->(sc:ServiceCategory), (j)-[]->(skc:SkillCategory), " +
		"(e:Employer)-[p:Posted]->(j) " +
		"WHERE(j.Title contains $title AND j.Wage >= $minWage AND e.Rating >= $erating) " +
		"optional match (j)-[]->(sk:Skill) optional match (e)-[p1:Posted]->(j1: Job {Status: $history}) " +
		"with j, sc, skc, e, sk, sum(j1.Wage) as money_spent " +
		"optional match (f:Freelancer)-[a:AppliedTo]->(j) " +
		"return properties(j), collect(distinct properties(sc)), collect(distinct properties(skc)), " +
		"collect(distinct properties(sk)), toFloat(money_spent), toFloat(e.Rating), count(a)"

	createJobApplication = "MATCH (j:Job {ID: $jid}), (f:Freelancer {ID: $fid}) MERGE (j)<-[r:AppliedTo]-(f) RETURN 'OK'"
	getJobApplicants = "MATCH (j:Job {ID: $jid})<-[r:AppliedTo]-(f:Freelancer) optional match (f)-[]->(skc:SkillCategory) optional match (f)-[]->(sk:Skill) RETURN properties(f), collect(distinct properties(skc)), collect(distinct properties(sk))"
	selectFreelancerForJob = "MATCH (j:Job {ID: $jid})<-[a:AppliedTo]-(f:Freelancer {ID: $fid}) MERGE (j)<-[r:IsEmployedTo]-(f) DELETE a RETURN 'OK'"
	getAcceptedFreelancersForJob = "MATCH (j:Job {ID: $jid})<-[r:IsEmployedTo]-(f:Freelancer) optional match (f)-[]->(skc:SkillCategory) optional match (f)-[]->(sk:Skill) RETURN properties(f), collect(distinct properties(skc)), collect(distinct properties(sk))"
	getAcceptedJobsForFreelancer = "match (j:Job)<-[r:IsEmployedTo]-(f1:Freelancer {ID: $fid}), " +
	    "(j)-[]->(sc:ServiceCategory), (j)-[]->(skc:SkillCategory)," +
		"(e:Employer)-[p:Posted]->(j)" +
		"optional match (j)-[]->(sk:Skill) optional match (e)-[p1:Posted]->(j1: Job {Status: $history}) " +
		"with j, sc, skc, e, sk, sum(j1.Wage) as money_spent " +
		"optional match (f:Freelancer)-[a:AppliedTo]->(j) " +
		"return properties(j), collect(distinct properties(sc)), collect(distinct properties(skc))," +
		"collect(distinct properties(sk)), toFloat(money_spent), toFloat(e.Rating), count(a)"
	getEmployerHistoryJobs = "MATCH (j:Job {Status: $history})<-[r:Posted]-(e:Employer {ID: $euid}), " +
	    "(j)-[]->(sc:ServiceCategory), (j)-[]->(skc:SkillCategory) " +
		"optional match (j)-[]->(sk:Skill) " +
		"with j, sc, skc, e, sk, sum(j.Wage) as money_spent " +
		"optional match (f:Freelancer)-[a:AppliedTo]->(j) " +
		"return properties(j), collect(distinct properties(sc)), collect(distinct properties(skc))," +
		"collect(distinct properties(sk)), toFloat(money_spent), toFloat(e.Rating), count(a)"
	getFreelancerHistoryJobs = "MATCH (j:Job {Status: $history})<-[r:IsEmployedTo]-(f:Freelancer {ID: $fuid}), " +
	    "(j)-[]->(sc:ServiceCategory), (j)-[]->(skc:SkillCategory), " +
		"(e:Employer)-[p:Posted]->(j)" +
		"optional match (j)-[]->(sk:Skill) optional match (e)-[p1:Posted]->(j1: Job {Status: $history}) " +
		"with j, sc, skc, e, sk, sum(j1.Wage) as money_spent " +
		"optional match (f1:Freelancer)-[a:AppliedTo]->(j) " +
		"return properties(j), collect(distinct properties(sc)), collect(distinct properties(skc))," +
		"collect(distinct properties(sk)), toFloat(money_spent), toFloat(e.Rating), count(a)"
	finishJob = "MATCH (j:Job {ID: $jid}) SET j.Status = $history RETURN 'OK'"
)

type NoDataFoundError struct {
	Err string
}

func (err NoDataFoundError) Error() string {
	return err.Err
}

type Database interface {
	GetAllServiceCategories() ([]*models.ServiceCategory, error)
	GetSkillCategoriesByServiceCategory(*models.ServiceCategory) ([]*models.SkillCategory, error)
	GetSkillsByCategory(category *models.SkillCategory) ([]*models.Skill, error)
	CreateJob(job *models.NewJob, scs []*models.ServiceCategory, skcs []*models.SkillCategory) error
	GetJobByID(jobID string) (*models.Job, error)
	GetJobsWithFilter(title, status string, wageMin, eRatingMin float32) ([]*models.Job, error)
	CreateJobApplication(jid, fid string) error
	GetJobApplicants(jid string) ([]*models.Freelancer, error)
	SelectFreelancerForJob(jid, fid string) error
	GetAcceptedFreelancers(jid string) ([]*models.Freelancer, error)
	GetAcceptedJobsForFreelancer(fid string) ([]*models.Job, error)
	GetEmployerHistoryJobs(euid string) ([]*models.Job, error)
	GetFreelancerHistoryJobs(fuid string) ([]*models.Job, error)
	FinishJob(jid string) error
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

func (nj *Neo4jDatabase) GetAllServiceCategories() ([]*models.ServiceCategory, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	serviceCategories, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getAllServiceCategories, map[string]interface{} {})
		if err != nil {
			return nil, err
		}

		serviceCategories := make([]*models.ServiceCategory, 0)
		for qResult.Next() {
			serviceCategoryData := &models.ServiceCategory{}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), serviceCategoryData); err != nil {
				return nil, err
			}

			serviceCategories = append(serviceCategories, serviceCategoryData)
		}

		return serviceCategories, nil
	}); if err != nil {
		return nil, err
	}

	return serviceCategories.([]*models.ServiceCategory), nil
}

func (nj *Neo4jDatabase) GetSkillCategoriesByServiceCategory(category *models.ServiceCategory) ([]*models.SkillCategory, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	skillCategories, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getAllSkillCategories, map[string]interface{} {
			"ID": category.ID,
		}); if err != nil {
			return nil, err
		}

		skillCategories := make([]*models.SkillCategory, 0)

		for qResult.Next() {
			skillCategoryData := &models.SkillCategory{}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), skillCategoryData); err != nil {
				return nil, err
			}

			skillCategories = append(skillCategories, skillCategoryData)
		}

		return skillCategories, nil
	}); if err != nil {
		return nil, err
	}

	return skillCategories.([]*models.SkillCategory), nil
}

func (nj *Neo4jDatabase) GetSkillsByCategory(category *models.SkillCategory) ([]*models.Skill, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	skills, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getSkillsByCategory, map[string]interface{} {
			"ID": category.ID,
		}); if err != nil {
			return nil, err
		}

		skills := make([]*models.Skill, 0)
		for qResult.Next() {
			skillData := &models.Skill{}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), skillData); err != nil {
				return nil, err
			}

			skills = append(skills, skillData)
		}

		return skills, nil
	}); if err != nil {
		return nil, err
	}

	return skills.([]*models.Skill), nil
}

func (nj *Neo4jDatabase) CreateJob(job *models.NewJob, scs []*models.ServiceCategory, skcs []*models.SkillCategory) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Println("Could not get a session to database")
		return err
	}
	defer session.Close()

	newJobDetails := make(map[string]interface{})
	err = mapstructure.Decode(job, &newJobDetails)
	if err != nil {
		log.Println("Could not decode job")
		return err
	}

	newJobDetails["Status"] = models.JobStatusOpen

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		_, err = tx.Run(createJobNodeQuery, map[string]interface{}{
			"PROP": newJobDetails,
		}); if err != nil {
			log.Println("Could not create job node.")
			return nil, err
		}

		scids := make([]string, 0)
		for _, sc := range scs {
			scids = append(scids, sc.ID)
		}

		skcids := make([]string, 0)
		for _, skc := range skcs {
			skcids = append(skcids, skc.ID)
		}

		qResult, err := tx.Run(createJobSCRel, map[string]interface{}{
			"JID": job.ID,
			"SCIDS": scids,
		}); if err != nil || !qResult.Next(){
			log.Println("Could not create Job-SC relationship")
			tx.Rollback()
			return nil, err
		}

		qResult, err = tx.Run(createJobSKCRel, map[string]interface{}{
			"JID": job.ID,
			"SKCIDS": skcids,
		}); if err != nil || !qResult.Next(){
			log.Println("Could not get create Job-SKC relationship")
			tx.Rollback()
			return nil, err
		}

		qResult, err = tx.Run(assignJobToEmployer, map[string]interface{}{
			"JID": job.ID,
			"euid": job.EUID,
		}); if err != nil || !qResult.Next(){
			log.Println("Could not assign job to employer")
			tx.Rollback()
			return nil, err
		}

		return nil, nil
	}); if err != nil {
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

	jobData, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getJobByIDQuery, map[string]interface{}{
			"ID": jobID,
			"history": models.JobStatusHistory,
		}); if err != nil {
			return nil, err
		}

		var jobData *models.Job
		scDatas := make([]*models.ServiceCategory, 0)
		skcDatas := make([]*models.SkillCategory, 0)
		skDatas := make([]*models.Skill, 0)

		for qResult.Next() {
			jobData = &models.Job{}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), jobData); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(1), &scDatas); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(2), &skcDatas); err != nil {
				return nil, err
			}
			if err = mapstructure.Decode(qResult.Record().GetByIndex(3), &skDatas); err != nil {
				return nil, err
			}

			jobData.MoneySpent = qResult.Record().GetByIndex(4).(float64)
			jobData.ERating = float32(qResult.Record().GetByIndex(5).(float64))
			jobData.NrCandidates = int(qResult.Record().GetByIndex(6).(int64))

			jobData.ServiceCategories = scDatas
			jobData.SkillCategories = skcDatas
		}

		return jobData, nil
	}); if err != nil {
		return nil, err
	}; if jobData == nil {
		return nil, NoDataFoundError{Err: "Job not found."}
	}

	return jobData.(*models.Job), nil
}

func (nj *Neo4jDatabase) GetJobsWithFilter(title, status string, wageMin, eRatingMin float32) ([]*models.Job, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	qFilterArg := make(map[string]interface{})
	qFilterArg["title"] = title
	qFilterArg["minWage"] = wageMin
	qFilterArg["erating"] = eRatingMin
	qFilterArg["history"] = models.JobStatusHistory
	qFilterArg["status"] = models.JobStatusOpen

	jobsData, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getJobsWithFilter, qFilterArg)
		if err != nil {
			return nil, err
		}

		jobsData := make([]*models.Job, 0)

		for qResult.Next() {
			jobData := &models.Job{}
			scDatas := make([]*models.ServiceCategory, 0)
			skcDatas := make([]*models.SkillCategory, 0)
			skDatas := make([]*models.Skill, 0)

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), jobData); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(1), &scDatas); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(2), &skcDatas); err != nil {
				return nil, err
			}
			if err = mapstructure.Decode(qResult.Record().GetByIndex(3), &skDatas); err != nil {
				return nil, err
			}
			jobData.MoneySpent = qResult.Record().GetByIndex(4).(float64)
			jobData.ERating = float32(qResult.Record().GetByIndex(5).(float64))
			jobData.NrCandidates = int(qResult.Record().GetByIndex(6).(int64))

			jobData.ServiceCategories = scDatas
			jobData.SkillCategories = skcDatas

			jobsData = append(jobsData, jobData)
		}

		return jobsData, nil
	}); if err != nil {
		return nil, err
	}
	
	return jobsData.([]*models.Job), nil
}

func (nj *Neo4jDatabase) CreateJobApplication(jid, fid string) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(createJobApplication, map[string]interface{}{
			"jid": jid,
			"fid": fid,
		}); if err != nil {
			tx.Rollback()
			return nil, err
		}; if !qResult.Next() {
			tx.Rollback()
			return nil, NoDataFoundError{Err: "Job or freelancer not found"}
		}

		return qResult, nil
	}); if err != nil {
		return err
	}

	return nil
}

func (nj *Neo4jDatabase) GetJobApplicants(jid string) ([]*models.Freelancer, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	freelancersData, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getJobApplicants, map[string]interface{}{
			"jid": jid,
		}); if err != nil {
			return nil, err
		}

		freelancersData := make([]*models.Freelancer, 0)
		for qResult.Next() {
			freelancerData := &models.Freelancer{}
			skcDatas := make([]*models.SkillCategory, 0)
			skDatas := make([]*models.Skill, 0)


			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), freelancerData); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(1), &skcDatas); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(2), &skDatas); err != nil {
				return nil, err
			}

			freelancerData.SkillCategories = skcDatas
			freelancerData.Skills = skDatas

			freelancersData = append(freelancersData, freelancerData)
		}

		return freelancersData, nil
	}); if err != nil {
		return nil, err
	}

	return freelancersData.([]*models.Freelancer), nil
}

func (nj *Neo4jDatabase) SelectFreelancerForJob(jid, fid string) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(selectFreelancerForJob, map[string]interface{}{
			"jid": jid,
			"fid": fid,
		}); if err != nil {
			tx.Rollback()
			return nil, err
		}; if !qResult.Next() {
			tx.Rollback()
			return nil, NoDataFoundError{Err: "Job or freelancer not found"}
		}

		return qResult, nil
	}); if err != nil {
		return err
	}

	return nil
}

func (nj *Neo4jDatabase) GetAcceptedFreelancers(jid string) ([]*models.Freelancer, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	freelancersData, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getAcceptedFreelancersForJob, map[string]interface{}{
			"jid": jid,
		}); if err != nil {
			return nil, err
		}

		freelancersData := make([]*models.Freelancer, 0)
		for qResult.Next() {
			freelancerData := &models.Freelancer{}
			skcDatas := make([]*models.SkillCategory, 0)
			skDatas := make([]*models.Skill, 0)


			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), freelancerData); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(1), &skcDatas); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(2), &skDatas); err != nil {
				return nil, err
			}

			freelancerData.SkillCategories = skcDatas
			freelancerData.Skills = skDatas

			freelancersData = append(freelancersData, freelancerData)
		}

		return freelancersData, nil
	}); if err != nil {
		return nil, err
	}

	return freelancersData.([]*models.Freelancer), nil
}

func (nj *Neo4jDatabase) GetAcceptedJobsForFreelancer(fid string) ([]*models.Job, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	jobDatas, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getAcceptedJobsForFreelancer, map[string]interface{}{
			"fid": fid,
			"history": models.JobStatusHistory,
		}); if err != nil {
			return nil, err
		}

		jobDatas := make([]*models.Job, 0)

		for qResult.Next() {
			jobData := &models.Job{}
			scDatas := make([]*models.ServiceCategory, 0)
			skcDatas := make([]*models.SkillCategory, 0)
			skDatas := make([]*models.Skill, 0)

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), jobData); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(1), &scDatas); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(2), &skcDatas); err != nil {
				return nil, err
			}
			if err = mapstructure.Decode(qResult.Record().GetByIndex(3), &skDatas); err != nil {
				return nil, err
			}
			jobData.MoneySpent = qResult.Record().GetByIndex(4).(float64)
			jobData.ERating = float32(qResult.Record().GetByIndex(5).(float64))
			jobData.NrCandidates = int(qResult.Record().GetByIndex(6).(int64))

			jobData.ServiceCategories = scDatas
			jobData.SkillCategories = skcDatas

			jobDatas = append(jobDatas, jobData)
		}

		log.Println(jobDatas)

		return jobDatas, nil
	}); if err != nil {
		return nil, err
	}

	return jobDatas.([]*models.Job), nil
}

func (nj *Neo4jDatabase) GetEmployerHistoryJobs(euid string) ([]*models.Job, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	jobsData, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getEmployerHistoryJobs, map[string]interface{}{
			"euid": euid,
			"history": models.JobStatusHistory,
		}); if err != nil {
			return nil, err
		}

		jobsData := make([]*models.Job, 0)
		for qResult.Next() {
			jobData := &models.Job{}
			scDatas := make([]*models.ServiceCategory, 0)
			skcDatas := make([]*models.SkillCategory, 0)
			skDatas := make([]*models.Skill, 0)

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), jobData); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(1), &scDatas); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(2), &skcDatas); err != nil {
				return nil, err
			}
			if err = mapstructure.Decode(qResult.Record().GetByIndex(3), &skDatas); err != nil {
				return nil, err
			}
			jobData.MoneySpent = qResult.Record().GetByIndex(4).(float64)
			jobData.ERating = float32(qResult.Record().GetByIndex(5).(float64))
			jobData.NrCandidates = int(qResult.Record().GetByIndex(6).(int64))

			jobData.ServiceCategories = scDatas
			jobData.SkillCategories = skcDatas

			jobsData = append(jobsData, jobData)
		}

		return jobsData, nil
	}); if err != nil {
		return nil, err
	}

	return jobsData.([]*models.Job), nil
}

func (nj *Neo4jDatabase) GetFreelancerHistoryJobs(fuid string) ([]*models.Job, error) {
	session, err := nj.Driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	jobsData, err := session.ReadTransaction(func(tx neo4j.Transaction) (i interface{}, err error) {
		qResult, err := tx.Run(getFreelancerHistoryJobs, map[string]interface{}{
			"fuid": fuid,
			"history": models.JobStatusHistory,
		}); if err != nil {
			return nil, err
		}

		jobsData := make([]*models.Job, 0)

		for qResult.Next() {
			jobData := &models.Job{}
			scDatas := make([]*models.ServiceCategory, 0)
			skcDatas := make([]*models.SkillCategory, 0)
			skDatas := make([]*models.Skill, 0)

			if err = mapstructure.Decode(qResult.Record().GetByIndex(0), jobData); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(1), &scDatas); err != nil {
				return nil, err
			}

			if err = mapstructure.Decode(qResult.Record().GetByIndex(2), &skcDatas); err != nil {
				return nil, err
			}
			if err = mapstructure.Decode(qResult.Record().GetByIndex(3), &skDatas); err != nil {
				return nil, err
			}
			jobData.MoneySpent = qResult.Record().GetByIndex(4).(float64)
			jobData.ERating = float32(qResult.Record().GetByIndex(5).(float64))
			jobData.NrCandidates = int(qResult.Record().GetByIndex(6).(int64))

			jobData.ServiceCategories = scDatas
			jobData.SkillCategories = skcDatas

			jobsData = append(jobsData, jobData)
		}

		return jobsData, nil
	}); if err != nil {
		return nil, err
	}

	return jobsData.([]*models.Job), nil
}

func (nj *Neo4jDatabase) FinishJob(jid string) error {
	session, err := nj.Driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	log.Println(jid)

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		qResult, err := tx.Run(finishJob, map[string]interface{}{
			"jid": jid,
			"history": models.JobStatusHistory,
		}); if err != nil {
			return nil, err
		}; if !qResult.Next() {
			return nil, NoDataFoundError{Err: "Job not found"}
		}

		return qResult, nil
	}); if err != nil {
		return err
	}

	return nil
}