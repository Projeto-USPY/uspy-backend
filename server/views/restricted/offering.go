package restricted

import (
	"net/http"
	"sort"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

func GetOfferingComments(ctx *gin.Context, comments []*models.Comment) {
	sort.Slice(comments, func(i, j int) bool {
		if comments[i].Upvotes == comments[j].Upvotes {
			return comments[i].Downvotes < comments[j].Downvotes
		}

		return comments[i].Upvotes > comments[j].Upvotes
	})

	results := make([]*views.Comment, 0)
	for _, c := range comments {
		results = append(results, views.NewCommentFromModel(c))
	}

	ctx.JSON(http.StatusOK, results)
}

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
		sizeI, sizeJ := len(results[i].Years), len(results[j].Years)
		if results[i].Years[sizeI-1] == results[j].Years[sizeJ-1] {
			return len(results[i].Years) > len(results[j].Years)
		}

		return results[i].Years[sizeI-1] > results[j].Years[sizeJ-1]
	})

	ctx.JSON(http.StatusOK, results[:utils.Min(limit, len(results))])
}
