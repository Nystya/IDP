package usecases

import (
	"idp/Profiles/database"
	"idp/Profiles/models"
)

type ProfileManager interface {
	EditEmployerProfile(employer *models.EditEmployerRequest) error
	EditFreelancerProfile(freelancer *models.EditFreelancerRequest) error
	GetEmployerProfile(euid string) (*models.Employer, error)
	GetFreelancerProfile(fuid string) (*models.Freelancer, error)
}

type ProfileManagerConfig struct {
	DBURL string
	DBUser string
	DBPass string
}

type ProfileManagerImpl struct {
	conf ProfileManagerConfig
	db database.Database
}

func NewProfileManagerImpl(conf ProfileManagerConfig) *ProfileManagerImpl {
	return &ProfileManagerImpl{
		conf: conf,
		db: database.NewNeo4jDatabase(conf.DBURL, conf.DBUser, conf.DBPass),
	}
}

func (p ProfileManagerImpl) EditEmployerProfile(employer *models.EditEmployerRequest) error {
	return p.db.EditEmployerProfile(employer)
}

func (p ProfileManagerImpl) EditFreelancerProfile(freelancer *models.EditFreelancerRequest) error {
	return p.db.EditFreelancerProfile(freelancer, freelancer.Skills, freelancer.SkillCategories)
}

func (p ProfileManagerImpl) GetEmployerProfile(euid string) (*models.Employer, error) {
	return p.db.GetEmployerProfile(euid)
}

func (p ProfileManagerImpl) GetFreelancerProfile(fuid string) (*models.Freelancer, error) {
	return p.db.GetFreelancerProfile(fuid)
}