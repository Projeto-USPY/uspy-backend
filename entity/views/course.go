package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type Course struct {
	Name           string            `json:"name"`
	Code           string            `json:"code"`
	Specialization string            `json:"specialization"`
	SubjectCodes   map[string]string `json:"subjects"`
}

func NewCourseFromModel(course *models.Course) *Course {
	c := Course{
		Name:           course.Name,
		Code:           course.Code,
		Specialization: course.Specialization,
		SubjectCodes:   make(map[string]string),
	}

	for k, v := range course.SubjectCodes {
		c.SubjectCodes[k] = v
	}

	return &c
}
