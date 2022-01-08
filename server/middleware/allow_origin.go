//Package middleware contains useful middleware handlers
package middleware

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/gin-gonic/gin"
)

// AllowAnyOrigin enables CORS for all origins testing purposes
func AllowAnyOrigin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Origin", "http://127.0.0.1")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With,observe")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}

}

// AllowUSPYOrigin enables CORS for the Frontend, according to dev/prod environment variables
func AllowUSPYOrigin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With,observe")

		var frontURL string
		if config.Env.IsDev() {
			frontURL = "https://frontdev.uspy.me"
		} else {
			frontURL = "https://uspy.me"
		}

		ctx.Header("Access-Control-Allow-Origin", frontURL)
		ctx.SetSameSite(http.SameSiteNoneMode)

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
	}
}
