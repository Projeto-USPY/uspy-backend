package restricted

import (
	"net/http"
	"sort"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

// GetOfferings is a closure for the GET /api/restricted/offerings endpoint
func GetOfferingsWithStats(
	ctx *gin.Context,
	IDs []string,
	offerings []*models.Offering,
	stats []*models.OfferingStats,
	limit int,
) {
	results := make([]*views.Offering, 0, 20)

	for i := 0; i < len(offerings); i++ {
		results = append(results,
			views.NewOfferingFromModel(
				IDs[i],
				offerings[i],
				stats[i].Approval,
				stats[i].Disapproval,
				stats[i].Neutral,
			),
		)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Year > results[j].Year
	})

	ctx.JSON(http.StatusOK, results[:utils.Min(limit, len(results))])
}
