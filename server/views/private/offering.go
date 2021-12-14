package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/gin-gonic/gin"
)

// GetComment takes the model comment and presents its response view object.
func GetComment(ctx *gin.Context, comment *models.Comment) {
	result := views.NewCommentFromModel(comment)
	ctx.JSON(http.StatusOK, result)
}

// GetCommentRating takes the model comment rating and presents its response view object.
func GetCommentRating(ctx *gin.Context, model *models.CommentRating) {
	ctx.JSON(http.StatusOK, views.NewCommentRatingFromModel(model))
}

// RateComment is a dummy view method
func RateComment(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// ReportComment is a dummy view method
func ReportComment(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// PublishComment takes the model comment and presents its response view object
func PublishComment(ctx *gin.Context, model *models.Comment) {
	comment := views.NewCommentFromModel(model)
	ctx.JSON(http.StatusOK, comment)
}

// DeleteComment is a dummy view method
func DeleteComment(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
