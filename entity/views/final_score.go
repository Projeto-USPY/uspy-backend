package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type FinalScore struct {
	Grade     float64 `json:"grade"`
	Status    string  `json:"status"`
	Frequency int     `json:"frequency"`
}

func NewFinalScoreFromModel(model *models.FinalScore) *FinalScore {
	return &FinalScore{
		Grade:     model.Grade,
		Status:    model.Status,
		Frequency: model.Frequency,
	}
}
