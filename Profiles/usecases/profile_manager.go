package usecases

import (
    "github.com/google/uuid"
	"idp/Profiles/database"
	"idp/Profiles/models"
)

type ProfileManager interface {
	CreateEmployerProfile(employer *models.EditEmployerRequest) error
	CreateFreelancerProfile(freelancer *models.EditFreelancerRequest) error
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

func (p ProfileManagerImpl) CreateEmployerProfile(employer *models.EditEmployerRequest) error {
	return p.db.CreateEmployerProfile(employer)
}

func (p ProfileManagerImpl) CreateFreelancerProfile(freelancer *models.EditFreelancerRequest) error {
	return p.db.CreateFreelancerProfile(freelancer)
}

func (p ProfileManagerImpl) EditEmployerProfile(employer *models.EditEmployerRequest) error {
	return p.db.EditEmployerProfile(employer)
}

func (p ProfileManagerImpl) EditFreelancerProfile(freelancer *models.EditFreelancerRequest) error {
    for _, skill := range freelancer.Skills {
        id := uuid.New().String()
        skill.ID = id
    }

	return p.db.EditFreelancerProfile(freelancer, freelancer.Skills, freelancer.SkillCategories)
}

func (p ProfileManagerImpl) GetEmployerProfile(euid string) (*models.Employer, error) {
	return p.db.GetEmployerProfile(euid)
}

func (p ProfileManagerImpl) GetFreelancerProfile(fuid string) (*models.Freelancer, error) {
	return p.db.GetFreelancerProfile(fuid)
}