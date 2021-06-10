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

func PublishComment(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
