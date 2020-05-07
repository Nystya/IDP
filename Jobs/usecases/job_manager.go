package usecases

import (
	"github.com/google/uuid"
	"idp/Jobs/database"
	"idp/Jobs/models"
	"idp/Jobs/messaging"
	"time"
)

type JobManager interface {
	GetAllServiceCategories() ([]*models.ServiceCategory, error)
	GetSkillCategoriesByServiceCategory(category *models.ServiceCategory) ([]*models.SkillCategory, error)
	GetSkillByCategory(category *models.SkillCategory) ([]*models.Skill, error)
	AddJob(job *models.Job) (*models.Job, error)
	GetJobs(qFilter *models.Filter) ([]*models.Job, error)
	GetJob(jobID string) (*models.Job, error)
	ApplyForJob(jobID, freelancerID string) error
	GetApplicants(jobID string) ([]*models.Freelancer, error)
	SelectFreelancerForJob(jobID, freelancerID string) error
	GetAcceptedFreelancers(jobId string) ([]*models.Freelancer, error)
	GetAcceptedJobsForFreelancer(fid string) ([]*models.Job, error)
	GetHistoryJobs(uid, actorType string) ([]*models.Job, error)
	FinishJob(jid string) error
}

const (
	MQTTHost     = "tcp://mqtt:1883"
	MQTTClient   = "adaptor_client"
	MQTTUsername = ""
	MQTTPassword = ""
	MQTTTopic    = "#"
	MQTTQOS      = 0
)

type JobManagerConfig struct {
	DBURL string
	DBUser string
	DBPass string
}

type JobManagerImpl struct {
	conf JobManagerConfig
	db database.Database
	mq *messaging.MQTTConnection
}

func NewJobManagerImpl(conf JobManagerConfig) *JobManagerImpl {
	return &JobManagerImpl{
		conf: conf,
		db: database.NewNeo4jDatabase(conf.DBURL, conf.DBUser, conf.DBPass),
		mq: messaging.NewMQTTConnection(MQTTHost, MQTTClient, MQTTUsername, MQTTPassword, MQTTTopic, MQTTQOS),
	}
}

func (jm *JobManagerImpl) GetAllServiceCategories() ([]*models.ServiceCategory, error) {
	return jm.db.GetAllServiceCategories()
}

func (jm *JobManagerImpl) GetSkillCategoriesByServiceCategory(category *models.ServiceCategory) ([]*models.SkillCategory, error) {
	return jm.db.GetSkillCategoriesByServiceCategory(category)
}

func (jm *JobManagerImpl) GetSkillByCategory(category *models.SkillCategory) ([]*models.Skill, error) {
	return jm.db.GetSkillsByCategory(category)
}

func (jm *JobManagerImpl) AddJob(job *models.Job) (*models.Job, error) {
	id := uuid.New().String()

	job.ID = id
	job.PostTime = time.Now().String()

	newJob := &models.NewJob{
		ID:          job.ID,
		EUID: 		 job.EUID,
		Title:       job.Title,
		Experience:  job.Experience,
		Wage:        job.Wage,
		Places:      job.Places,
		Description: job.Description,
		PostTime:    job.PostTime,
	}

	if err := jm.db.CreateJob(newJob, job.ServiceCategories, job.SkillCategories); err != nil {
		return nil, err
	}

	msg := map[string]interface{}{
		"Wage": newJob.Wage,
		"Places": newJob.Places,
	}

	go jm.mq.Client.Publish("jobs/" + newJob.EUID, MQTTQOS, false, msg)

	return job, nil
}

func (jm *JobManagerImpl) GetJobs(qFilter *models.Filter) ([]*models.Job, error) {
	jobsData, err := jm.db.GetJobsWithFilter(qFilter.Title, qFilter.Status, qFilter.WageMin, qFilter.ERating)
	if err != nil {
		return nil, err
	}

	return jobsData, nil
}

func (jm *JobManagerImpl) GetJob(jobID string) (*models.Job, error) {
	jobData, err := jm.db.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}

	return jobData, nil
}

func (jm *JobManagerImpl) ApplyForJob(jobID, freelancerID string) error {
	return jm.db.CreateJobApplication(jobID, freelancerID)
}

func (jm *JobManagerImpl) GetApplicants(jobID string) ([]*models.Freelancer, error) {
	return jm.db.GetJobApplicants(jobID)
}

func (jm *JobManagerImpl) SelectFreelancerForJob(jobID, freelancerID string) error {
	return jm.db.SelectFreelancerForJob(jobID, freelancerID)
}

func (jm *JobManagerImpl) GetAcceptedFreelancers(jobId string) ([]*models.Freelancer, error) {
	return jm.db.GetAcceptedFreelancers(jobId)
}

func (jm *JobManagerImpl) GetAcceptedJobsForFreelancer(fid string) ([]*models.Job, error) {
	return jm.db.GetAcceptedJobsForFreelancer(fid)
}


func (jm *JobManagerImpl) GetHistoryJobs(uid, actorType string) ([]*models.Job, error) {
	if actorType == models.ActorTypeEmployer {
		return jm.db.GetEmployerHistoryJobs(uid)
	} else if actorType == models.ActorTypeFreelancer {
		return jm.db.GetFreelancerHistoryJobs(uid)
	}

	return nil, nil
}

func (jm *JobManagerImpl) FinishJob(jid string) error {
	return jm.db.FinishJob(jid)
}