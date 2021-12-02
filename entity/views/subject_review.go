package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type SubjectReview struct {
	Review map[string]interface{} `json:"categories"`
}

func NewSubjectReviewFromModel(model *models.SubjectReview) *SubjectReview {
	return &SubjectReview{Review: model.Review}
}
