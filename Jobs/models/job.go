package models

import "time"

type Service struct {
	ID string
	Service string
}

type SkillCategory struct {
	ID string
	Category string
}

type Skill struct {
	ID string
	SCID SkillCategory
	Skill string
}

type Job struct {
	ID string
	EUID string
	Title string
	Service Service
	Category SkillCategory
	Experience string
	Wage float32
	Places int32
	Description string
	Skills []Skill
	PostTime time.Time
	ERating float32
	NrCandidates int
	MoneySpent float64
}


