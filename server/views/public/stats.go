package public

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

// GetStats takes the calculated statistics and sets the output object as the response
func GetStats(
	ctx *gin.Context,
	stats *models.Stats,
) {
	if ctx.IsAborted() {
		return
	}

	ctx.JSON(http.StatusOK, views.NewStatsFromModel(stats))
}
