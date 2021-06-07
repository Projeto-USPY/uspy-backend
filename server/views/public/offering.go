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

	sort.Slice(offerings, func(i, j int) bool {
		return offerings[i].Year > offerings[j].Year
	})

	for i := 0; i < utils.Min(3, len(offerings)); i++ { // only return 3 most recent to guests
		results = append(results, views.NewPartialOfferingFromModel(IDs[i], offerings[i]))
	}

	ctx.JSON(http.StatusOK, results)
}
