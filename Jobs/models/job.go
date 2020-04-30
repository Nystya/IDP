package models

import "time"

//go:generate easytags $GOFILE

type ServiceCategory struct {
	ID      string `json:"id"`
	Service string `json:"service"`
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

type NewJob struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Experience  string    `json:"experience"`
	Wage        float32   `json:"wage"`
	Places      int32     `json:"places"`
	Description string    `json:"description"`
	PostTime    time.Time `json:"post_time"`
}

type Job struct {
	ID           string          `json:"id"`
	EUID         string          `json:"euid"`
	Title        string          `json:"title"`
	Service      ServiceCategory `json:"service"`
	Category     SkillCategory   `json:"category"`
	Experience   string          `json:"experience"`
	Wage         float32         `json:"wage"`
	Places       int32           `json:"places"`
	Description  string          `json:"description"`
	Skills       []Skill         `json:"skills"`
	PostTime     time.Time       `json:"post_time"`
	ERating      float32         `json:"e_rating"`
	NrCandidates int             `json:"nr_candidates"`
	MoneySpent   float64         `json:"money_spent"`
}

type Filter struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	WageMin float32 `json:"wage_min"`
	ERating float32 `json:"e_rating"`
}
