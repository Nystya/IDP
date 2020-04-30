package usecases

import (
	"github.com/google/uuid"
	"idp/Jobs/database"
	"idp/Jobs/models"
)

type JobManager interface {
	GetAllServiceCategories() ([]*models.ServiceCategory, error)
	GetSkillCategoriesByServiceCategory(category *models.ServiceCategory) ([]*models.SkillCategory, error)
	GetSkillByCategory(category *models.SkillCategory) ([]*models.Skill, error)
	AddJob(job *models.Job) error
	GetJobs(qFilter *models.Filter) ([]*models.Job, error)
	GetJob(jobID string) (*models.Job, error)
}

type JobManagerConfig struct {
	DBURL string
	DBUser string
	DBPass string
}

type JobManagerImpl struct {
	conf JobManagerConfig
	db database.Database
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

func NewJobManagerImpl(conf JobManagerConfig) *JobManagerImpl {
	return &JobManagerImpl{
		conf: conf,
		db: database.NewNeo4jDatabase(conf.DBURL, conf.DBUser, conf.DBPass),
	}
}

func (jm *JobManagerImpl) AddJob(job *models.Job) error {
	id := uuid.New().String()

	job.ID = id

	newJob := &models.NewJob{
		ID:          job.ID,
		Title:       job.Title,
		Experience:  job.Experience,
		Wage:        job.Wage,
		Places:      job.Places,
		Description: job.Description,
		PostTime:    job.PostTime,
	}

	if err := jm.db.CreateJob(newJob, &job.Service, &job.Category); err != nil {
		return err
	}

	return nil
}

func (jm *JobManagerImpl) GetJobs(qFilter *models.Filter) ([]*models.Job, error) {
	jobsData, err := jm.db.GetJobsWithFilter(qFilter.Title, qFilter.WageMin, qFilter.ERating)
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
