package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

// Institute is the response view object for a institute
type Institute struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// NewInstituteFromModel is a constructor. It takes an institute model and returns its response view object.
func NewInstituteFromModel(model *models.Institute) *Institute {
	return &Institute{
		Name: model.Name,
		Code: model.Code,
	}
}
