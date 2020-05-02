package models

//go:generate easytags $GOFILE

const (
	JobStatusOpen    = "OPEN"
	JobStatusHistory = "HISTORY"
)

type NewJob struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Experience  string  `json:"experience"`
	Wage        float32 `json:"wage"`
	Places      int32   `json:"places"`
	Description string  `json:"description"`
	PostTime    string  `json:"post_time"`
}

type Job struct {
	ID                string             `json:"id"`
	EUID              string             `json:"euid"`
	Status            string             `json:"status"`
	Title             string             `json:"title"`
	ServiceCategories []*ServiceCategory `json:"service_categories"`
	SkillCategories   []*SkillCategory   `json:"skill_categories"`
	Experience        string             `json:"experience"`
	Wage              float32            `json:"wage"`
	Places            int32              `json:"places"`
	Description       string             `json:"description"`
	Skills            []*Skill           `json:"skills"`
	PostTime          string             `json:"post_time"`
	ERating           float32            `json:"e_rating"`
	NrCandidates      int                `json:"nr_candidates"`
	MoneySpent        float64            `json:"money_spent"`
}

type Filter struct {
	ID      string  `json:"id"`
	Status  string  `json:"status"`
	Title   string  `json:"title"`
	WageMin float32 `json:"wage_min"`
	ERating float32 `json:"e_rating"`
}
