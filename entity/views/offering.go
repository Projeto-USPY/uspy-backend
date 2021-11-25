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

	approvalRate := 0
	disapprovalRate := 0
	neutralRate := 0

	if total != 0 {
		approvalRate = approval / total
		disapprovalRate = disapproval / total
		neutralRate = neutral / total
	}

	return &Offering{
		ProfessorName: model.Professor,
		ProfessorCode: ID,
		Years:         model.Years,
		Approval:      float64(approvalRate),
		Disapproval:   float64(disapprovalRate),
		Neutral:       float64(neutralRate),
	}
}

func NewPartialOfferingFromModel(ID string, model *models.Offering) *Offering {
	return &Offering{
		ProfessorName: model.Professor,
		ProfessorCode: ID,
		Years:         model.Years,
	}
}
