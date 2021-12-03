// package middleware contains useful middleware handlers
package middleware

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
)

// JWT is used to ensure authorization with the JWT Access Cookie.
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := utils.ValidateJWT(cookie, config.Env.JWTSecret)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		ctx.Set("access_token", token)
		ctx.Set("userID", userID)
	}
}
