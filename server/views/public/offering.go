package public

import (
	"net/http"

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

	views.SortOfferings(results)

	// output only the first three
	ctx.JSON(http.StatusOK, results[:utils.Min(len(results), 3)])
}
