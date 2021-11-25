package views

import "github.com/Projeto-USPY/uspy-backend/entity/models"

type Offering struct {
	ProfessorName string   `json:"professor"`
	ProfessorCode string   `json:"code"`
	Years         []string `json:"years"`

	Approval    float64 `json:"approval"`
	Neutral     float64 `json:"neutral"`
	Disapproval float64 `json:"disapproval"`
}

func NewOfferingFromModel(ID string, model *models.Offering, approval, disapproval, neutral int) *Offering {
	total := (approval + disapproval + neutral)

	approvalRate := 0.0
	disapprovalRate := 0.0
	neutralRate := 0.0

	if total != 0 {
		approvalRate = float64(approval) / float64(total)
		disapprovalRate = float64(disapproval) / float64(total)
		neutralRate = float64(neutral) / float64(total)
	}

	return &Offering{
		ProfessorName: model.Professor,
		ProfessorCode: ID,
		Years:         model.Years,
		Approval:      approvalRate,
		Disapproval:   disapprovalRate,
		Neutral:       neutralRate,
	}
}

func NewPartialOfferingFromModel(ID string, model *models.Offering) *Offering {
	return &Offering{
		ProfessorName: model.Professor,
		ProfessorCode: ID,
		Years:         model.Years,
	}
}
