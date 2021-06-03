package views

import (
	"github.com/Projeto-USPY/uspy-backend/entity/models"
)

type Subject struct {
	Code           string          `json:"code"`
	CourseCode     string          `json:"course"`
	Specialization string          `json:"specialization"`
	Name           string          `json:"name"`
	Description    string          `json:"desc"`
	Semester       int             `json:"semester"`
	ClassCredits   int             `json:"class"`
	AssignCredits  int             `json:"assign"`
	TotalHours     string          `json:"hours"`
	Optional       bool            `json:"optional"`
	Stats          map[string]int  `json:"stats"`
	Requirements   [][]Requirement `json:"requirements"`
}

// Transforms from map[string][]models.Requirement to [][]views.Requirement
func transformRequirements(sub *models.Subject) [][]Requirement {
	requirements := [][]Requirement{}
	for _, val := range sub.Requirements {
		group := []Requirement{}
		for _, r := range val {
			group = append(group, Requirement{
				Subject: r.Subject,
				Name:    r.Name,
				Strong:  r.Strong,
			})
		}
		requirements = append(requirements, group)
	}
	return requirements
}

func NewSubjectFromModel(model *models.Subject) *Subject {
	return &Subject{
		Code:           model.Code,
		CourseCode:     model.CourseCode,
		Specialization: model.Specialization,
		Name:           model.Name,
		Description:    model.Description,
		Semester:       model.Semester,
		ClassCredits:   model.ClassCredits,
		AssignCredits:  model.AssignCredits,
		TotalHours:     model.TotalHours,
		Optional:       model.Optional,
		Stats:          model.Stats,
		Requirements:   transformRequirements(model),
	}
}
