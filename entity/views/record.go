package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

// Record is the response view object for a user subject record
//
// It contains three types of data:
// - optional basic subject data
// - record data (grade, status, frequency)
// - optinal time data
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

// NewRecordFromModel is a constructor. It takes a model and returns its response view object.
func NewRecordFromModel(model *models.Record) *Record {
	return &Record{
		Grade:     model.Grade,
		Status:    model.Status,
		Frequency: model.Frequency,
	}
}
