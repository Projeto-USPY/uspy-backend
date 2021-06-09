package private

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PublishComment(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
