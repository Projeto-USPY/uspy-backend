package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

func GetComment(ctx *gin.Context, comment *models.Comment) {
	result := views.NewCommentFromModel(comment)
	ctx.JSON(http.StatusOK, result)
}

func GetCommentRating(ctx *gin.Context, model *models.CommentRating) {
	ctx.JSON(http.StatusOK, views.NewCommentRatingFromModel(model))
}

func RateComment(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func ReportComment(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func PublishComment(ctx *gin.Context, model *models.Comment) {
	comment := views.NewCommentFromModel(model)
	ctx.JSON(http.StatusOK, comment)
}
