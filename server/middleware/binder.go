// package middleware contains useful middleware handlers
package middleware

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/gin-gonic/gin"
)

// Subject is a middleware for binding subject data
func Subject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subject := controllers.Subject{}
		if err := ctx.BindQuery(&subject); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ctx.Set("Subject", subject)
	}
}
