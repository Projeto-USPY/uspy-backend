package public

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

// GetProfessors takes the professor models and returns its response view objects
func GetProfessors(ctx *gin.Context, professors []*models.Professor) {
	viewProfessors := make([]*views.Professor, 0, len(professors))
	for _, prof := range professors {
		viewProfessors = append(viewProfessors, views.NewProfessorFromModel(prof))
	}

	ctx.JSON(http.StatusOK, viewProfessors)
}
