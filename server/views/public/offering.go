package public

import (
	"net/http"
	"sort"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

func GetOfferings(ctx *gin.Context, IDs []string, offerings []*models.Offering) {
	results := make([]*views.Offering, 0, 20)

	for i := 0; i < len(offerings); i++ {
		results = append(results, views.NewPartialOfferingFromModel(IDs[i], offerings[i]))
	}

	sort.Slice(results, func(i, j int) bool {
		sizeI, sizeJ := len(results[i].Years), len(results[j].Years)
		if results[i].Years[sizeI-1] == results[j].Years[sizeJ-1] {
			return len(results[i].Years) > len(results[j].Years)
		}

		return results[i].Years[sizeI-1] > results[j].Years[sizeJ-1]
	})

	// output only the first three
	ctx.JSON(http.StatusOK, results[:utils.Min(len(results), 3)])
}
