package models

//go:generate easytags $GOFILE

type ServiceCategory struct {
	ID       string `json:"id"`
	Category string `json:"category"`
}

type SkillCategory struct {
	ID       string `json:"id"`
	Category string `json:"category"`
}

type Skill struct {
	ID    string        `json:"id"`
	SCID  SkillCategory `json:"scid"`
	Skill string        `json:"skill"`
}
