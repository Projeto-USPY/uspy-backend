package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type Record struct {
	Subject        string `json:"subject,omitempty"`
	Course         string `json:"course,omitempty"`
	Specialization string `json:"specialization,omitempty"`

	Grade     float64 `json:"grade"`
	Status    string  `json:"status"`
	Frequency int     `json:"frequency"`

	Semester int `json:"semester,omitempty"`
	Year     int `json:"year,omitempty"`
}

func NewRecordFromModel(model *models.Record) *Record {
	return &Record{
		Grade:     model.Grade,
		Status:    model.Status,
		Frequency: model.Frequency,
	}
}
