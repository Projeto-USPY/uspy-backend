package middleware

import (
	"net/http"

	"github.com/tpreischadt/ProjetoJupiter/server/auth"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware is used to ensure authorization with the JWT Access Cookie.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = auth.ValidateJWT(cookie)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
