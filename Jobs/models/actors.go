package models

//go:generate easytags $GOFILE

const (
	ActorTypeFreelancer = "freelancer"
	ActorTypeEmployer = "employer"
)

type Freelancer struct {
	ID              string           `json:"id"`
	Phone           string           `json:"phone"`
	LastName        string           `json:"last_name"`
	FirstName       string           `json:"first_name"`
	Rating          float32          `json:"rating"`
	Balance         float32          `json:"balance"`
	Description     string           `json:"description"`
	Photo           string           `json:"photo"`
	Skills          []*Skill         `json:"skills"`
	SkillCategories []*SkillCategory `json:"skill_categories"`
}
