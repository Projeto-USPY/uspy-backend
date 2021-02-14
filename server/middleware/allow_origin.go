// package middleware contains useful middleware handlers
package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// AllowAnyOrigin enables CORS for testing purposes
func AllowAnyOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "http://127.0.0.1")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With,observe")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}

}

// AllowUSPYOrigin enables CORS for the Frontend, according to dev/prod environment variables
func AllowUSPYOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With,observe")

		if mode, ok := os.LookupEnv("MODE"); ok {
			var frontURL string
			if mode == "dev" {
				frontURL = "https://frontdev.uspy.me"
			} else {
				frontURL = "https://uspy.me"
			}

			c.Header("Access-Control-Allow-Origin", frontURL)
			c.SetSameSite(http.SameSiteNoneMode)

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		} else {
			c.AbortWithStatus(500)
			return
		}
	}
}
