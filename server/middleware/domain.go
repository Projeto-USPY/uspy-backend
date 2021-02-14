package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func DefineDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("LOCAL") == "TRUE" {
			c.Set("front_domain", "127.0.0.1")
			c.Next()
			return
		}

		if mode, ok := os.LookupEnv("MODE"); ok {
			var frontURL string
			if mode == "prod" {
				frontURL = "https://uspy.me"
			} else {
				frontURL = "https://frontdev.uspy.me"
			}

			c.Set("front_domain", frontURL)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
