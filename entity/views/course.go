package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

// Course is the response view object for a course
type Course struct {
	Name           string            `json:"name"`
	Code           string            `json:"code"`
	Specialization string            `json:"specialization"`
	Shift          string            `json:"shift"`
	SubjectCodes   map[string]string `json:"subjects,omitempty"`
}

// NewCourseFromModel is a constructor. It takes a course model and returns its response view object.
func NewCourseFromModel(course *models.Course) *Course {
	c := Course{
		Name:           course.Name,
		Code:           course.Code,
		Specialization: course.Specialization,
		Shift:          course.Shift,

		SubjectCodes: make(map[string]string),
	}

	for k, v := range course.SubjectCodes {
		c.SubjectCodes[k] = v
	}

	return &c
}
