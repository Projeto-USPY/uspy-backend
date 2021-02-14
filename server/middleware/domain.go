package middleware

import (
	"github.com/gin-gonic/gin"
	"os"
)

func DefineDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("LOCAL") == "TRUE" {
			c.Set("front_domain", "127.0.0.1")
			c.Next()
			return
		}

		c.Set("front_domain", "uspy.me")
	}
}
