package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type SubjectSibling struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Optional bool   `json:"optional"`
}

func NewSubjectSibling(model *models.Subject) *SubjectSibling {
	return &SubjectSibling{
		Code:     model.Code,
		Name:     model.Name,
		Optional: model.Optional,
	}
}
