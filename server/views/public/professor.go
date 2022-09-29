package public

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

// GetProfessor takes the professor model and returns its response view object
func GetProfessor(ctx *gin.Context, model *models.Professor) {
	ctx.JSON(http.StatusOK, views.NewProfessorFromModel(model))
}

// GetProfessorOfferings takes the offering models and returns their view response objects
func GetProfessorOfferings(ctx *gin.Context, professorID string, offerings []*models.Offering) {
	results := make([]*views.Offering, 0, 20)
	for _, off := range offerings {
		if off.Professor == "" {
			continue
		}

		results = append(results, views.NewPartialOfferingFromModel(professorID, off))
	}

	views.SortOfferings(results)

	// output only the first three
	ctx.JSON(http.StatusOK, results[:utils.Min(len(results), 3)])
}
