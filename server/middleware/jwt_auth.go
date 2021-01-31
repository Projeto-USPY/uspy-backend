package middleware

import (
	"net/http"

	"github.com/tpreischadt/ProjetoJupiter/server/models"

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

		err = models.ValidateJWT(cookie)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
