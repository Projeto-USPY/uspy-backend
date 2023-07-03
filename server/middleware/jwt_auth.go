package middleware

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

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
			log.Error("missing cookie in JWT middleware")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := utils.ValidateJWT(cookie, config.Env.JWTSecret)
		if err != nil {
			log.Error(fmt.Sprintf("Error validating JWT: %v", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if cookieType, ok := claims["type"].(string); !ok || cookieType != "access" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID := claims["user"].(string)

		ctx.Set("access_token", token)
		ctx.Set("userID", userID)
	}
}
