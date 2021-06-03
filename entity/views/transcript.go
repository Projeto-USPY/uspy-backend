package views

import "github.com/Projeto-USPY/uspy-backend/iddigital"

type Transcript struct {
	Grades []Record `json:"grades"`
	Name   string   `json:"name"`
	Nusp   string   `json:"nusp"`
}

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
