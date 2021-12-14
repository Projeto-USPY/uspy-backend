package public

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

// GetOfferings takes the offering models and returns its response view objects
//
// It also takes the IDs of the professors for each offering, since that information is not stored in the DTO.
// Besides that, it sorts the offerings to rank most interesting ones on top
// Since this is a public endpoint, results size is limited to 3
func GetOfferings(ctx *gin.Context, IDs []string, offerings []*models.Offering) {
	results := make([]*views.Offering, 0, 20)

	for i := 0; i < len(offerings); i++ {
		results = append(results, views.NewPartialOfferingFromModel(IDs[i], offerings[i]))
	}

	views.SortOfferings(results)

	// output only the first three
	ctx.JSON(http.StatusOK, results[:utils.Min(len(results), 3)])
}
