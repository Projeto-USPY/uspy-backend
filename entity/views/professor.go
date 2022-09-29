package views

import (
	"fmt"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

// Professor is the response view object for a professor
type Professor struct {
	Name string `json:"name"`
	Hash string `json:"code"`

	Stats map[string]int `json:"stats,omitempty"`
}

// NewProfessorFromModel is a constructor. It takes a professor model and returns its view object
func NewProfessorFromModel(model *models.Professor) *Professor {
	view := &Professor{
		Name:  model.Name,
		Hash:  model.Hash(),
		Stats: make(map[string]int),
	}

	for _, reviewMap := range model.Reviews {
		for key, value := range reviewMap.Review {
			view.Stats[fmt.Sprintf("total_%s", key)]++ // increment total

			if value {
				view.Stats[key]++ // increment positive
			} else {
				view.Stats[key] = utils.Max(0, view.Stats[key])
			}
		}
	}

	return view
}
