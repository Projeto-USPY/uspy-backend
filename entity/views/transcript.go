package views

import "github.com/Projeto-USPY/uspy-backend/iddigital"

// Transcript is the view response object for a users' history of grades, e.g. their progress
//
// It contains basic user info that is used solely as identifiers along with their grades data that is used to generate the subjects grades distribution
type Transcript struct {
	Grades []Record `json:"grades"`
	Name   string   `json:"name"`
	Nusp   string   `json:"nusp"`
}

// NewTranscript is a constructor. It takes a model object and returns its response view object.
func NewTranscript(model *iddigital.Transcript) *Transcript {
	t := Transcript{
		Name:   model.Name,
		Nusp:   model.Nusp,
		Grades: make([]Record, 0),
	}

	for _, g := range model.Grades {
		t.Grades = append(t.Grades, Record{
			Subject:        g.Subject,
			Course:         g.Course,
			Specialization: g.Specialization,
			Grade:          g.Grade,
			Status:         g.Status,
			Frequency:      g.Frequency,
			Semester:       g.Semester,
			Year:           g.Year,
		})
	}

	return &t
}
