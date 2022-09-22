package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

// Professor is the response view object for a professor
type Professor struct {
	Name string `json:"name"`
	Hash string `json:"code"`
}

// NewProfessorFromModel is a constructor. It takes a professor model and returns its view object
func NewProfessorFromModel(prof *models.Professor) *Professor {
	return &Professor{
		Name: prof.Name,
		Hash: prof.Hash(),
	}
}
