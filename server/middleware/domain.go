package middleware

import (
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/gin-gonic/gin"
)

// DefineDomain is a middleware for setting the cookie domain values
func DefineDomain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if config.Env.IsLocal() {
			ctx.Set("front_domain", config.Env.FrontOriginHost)
			ctx.Next()
			return
		}

		ctx.Set("front_domain", "uspy.me")
	}
}
