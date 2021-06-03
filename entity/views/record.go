package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type Record struct {
	Grade     float64 `json:"grade"`
	Status    string  `json:"status"`
	Frequency int     `json:"frequency"`
}

func NewRecordFromModel(model *models.Record) *Record {
	return &Record{
		Grade:     model.Grade,
		Status:    model.Status,
		Frequency: model.Frequency,
	}
}
