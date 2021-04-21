package middleware

import (
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/gin-gonic/gin"
)

// DefineDomain is a middleware for setting the cookie domain values
func DefineDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Env.IsLocal() {
			c.Set("front_domain", "127.0.0.1")
			c.Next()
			return
		}

		c.Set("front_domain", "uspy.me")
	}
}
