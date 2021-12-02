package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyEmail(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func RequestPasswordReset(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
