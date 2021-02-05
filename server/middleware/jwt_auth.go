package middleware

import (
	"net/http"

	"github.com/tpreischadt/ProjetoJupiter/server/models"

	"github.com/gin-gonic/gin"
)

// JWT is used to ensure authorization with the JWT Access Cookie.
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := models.ValidateJWT(cookie)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("access_token", token)
	}
}
