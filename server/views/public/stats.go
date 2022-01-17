package public

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

// GetStats takes the calculated statistics and sets the output object as the response
func GetStats(
	ctx *gin.Context,
	statsChan <-chan *views.StatsEntry,
) {
	if ctx.IsAborted() {
		return
	}

	stats := make(map[string]int)
	for i := 0; i < 5; i++ {
		value := <-statsChan
		stats[value.Name] = value.Count
	}

	ctx.JSON(http.StatusOK, stats)
}
