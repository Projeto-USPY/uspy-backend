package restricted

import (
	"net/http"
	"sort"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

// GetOfferingComments takes the comments models and returns their view response objects
//
// It also sorts comments to display most upvoted (and least downvoted) on top
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

// GetOfferingsWithStats takes the offering models and returns their view response objects
//
// It also takes IDs of the professors for each offering (given that this information is not stored in the model directly)
// Finally, it sorts the offerings to display most interesting ones on top
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

	views.SortOfferings(results)

	ctx.JSON(http.StatusOK, results[:utils.Min(limit, len(results))])
}
