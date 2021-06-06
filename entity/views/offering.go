package views

type Offering struct {
	ProfessorName string `json:"professor"`
	ProfessorCode string `json:"code"`
	Year          string `json:"year"`

	Approval    float64 `json:"approval,omitempty"`
	Neutral     float64 `json:"neutral,omitempty"`
	Disapproval float64 `json:"disapproval,omitempty"`
}
