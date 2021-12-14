package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// VerifyEmail is a dummy view method
func VerifyEmail(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// RequestPasswordReset is a dummy view method
func RequestPasswordReset(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
