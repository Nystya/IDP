package usecases

import (
	"github.com/google/uuid"
	"idp/Jobs/database"
	"idp/Jobs/models"
)

type JobManager interface {
	AddJob(job models.Job) error
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

func NewJobManagerImpl(conf JobManagerConfig) *JobManagerImpl {
	return &JobManagerImpl{
		conf: conf,
		db: database.NewNeo4jDatabase(conf.DBURL, conf.DBUser, conf.DBPass),
	}
}

func (jm *JobManagerImpl) AddJob(job models.Job) error {
	id := uuid.New().String()

	job.ID = id

	if err := jm.db.CreateJob(job); err != nil {
		return err
	}

	return nil
}