package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

// SubjectReview is the response view object for a subject review
//
// It only contains a map with review categories and their associated values. These values are usually boolean
type SubjectReview struct {
	Review map[string]interface{} `json:"categories"`
}

// NewSubjectReviewFromModel is a constructor. It takes a SubjectReview model and returns its view response object.
func NewSubjectReviewFromModel(model *models.SubjectReview) *SubjectReview {
	return &SubjectReview{Review: model.Review}
}
