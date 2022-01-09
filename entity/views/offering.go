package views

import (
	"sort"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
)

// Offering is the response view object for an offering
//
// It contains time information and some stats
type Offering struct {
	ProfessorName string   `json:"professor"`
	ProfessorCode string   `json:"code"`
	Years         []string `json:"years"`

	Approval    float64 `json:"approval"`
	Neutral     float64 `json:"neutral"`
	Disapproval float64 `json:"disapproval"`
}

// SortOfferings takes a list of offerings and sorts them
//
// It sorts based on the approval and neutral ratings values of each offering.
// If this value is equal in both objects, it uses the disapproval and amount of years as a tiebreaker
func SortOfferings(results []*Offering) {
	sort.SliceStable(results,
		func(i, j int) bool {
			if results[i].Approval == results[j].Approval {
				if results[i].Neutral == results[j].Neutral {
					if results[i].Disapproval == results[j].Disapproval {
						// if ratings are the same, show latest or most offerings
						sizeI, sizeJ := len(results[i].Years), len(results[j].Years)
						if results[i].Years[sizeI-1] == results[j].Years[sizeJ-1] {
							return len(results[i].Years) > len(results[j].Years)
						}

						return results[i].Years[sizeI-1] > results[j].Years[sizeJ-1]
					}

					return results[i].Disapproval < results[j].Disapproval
				}

				return results[i].Neutral > results[j].Neutral
			}

			return results[i].Approval > results[j].Approval
		},
	)
}

// NewOfferingFromModel is a constructor. It takes a model and returns its response view object.
//
// It also requires the ID of the professor in the current context and their approval stats
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

// NewPartialOfferingFromModel is a constructor. It takes an offering model and returns its view response object
//
// It is partial because it leaves stats data empty for the given offering
// This is used in public offering endpoints
func NewPartialOfferingFromModel(ID string, model *models.Offering) *Offering {
	return &Offering{
		ProfessorName: model.Professor,
		ProfessorCode: ID,
		Years:         model.Years,
	}
}
