package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

// Major is the response view object for a user's major / enrolled courses
//
// It unites model information from a course and a major (major models don't have the name property)
type Major struct {
	Code           string `json:"code"`
	Specialization string `json:"specialization"`
	Name           string `json:"name"`
}

// NewMajorFromModels takes the course and major models and unites them into a view object
func NewMajorFromModels(major *models.Major, course *models.Course) *Major {
	return &Major{
		Name:           course.Name,
		Code:           major.Code,
		Specialization: major.Specialization,
	}
}
