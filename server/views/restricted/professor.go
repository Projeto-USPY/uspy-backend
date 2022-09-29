package restricted

import (
	"net/http"
	"sort"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

// GetProfessorComments takes the comments models and returns their view response objects
//
// It also sorts comments to display most upvoted (and least downvoted) on top
func GetProfessorComments(ctx *gin.Context, comments []*models.Comment) {
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

// GetProfessorOfferings takes the offering models and returns their view response objects
func GetProfessorOfferings(ctx *gin.Context, professorID string, offerings []*models.Offering) {
	results := make([]*views.Offering, 0, 20)
	for _, off := range offerings {
		if off.Professor == "" {
			continue
		}

		results = append(results, views.NewPartialOfferingFromModel(professorID, off))
	}

	views.SortOfferings(results)
	ctx.JSON(http.StatusOK, results)
}
