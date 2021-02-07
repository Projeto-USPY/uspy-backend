package middleware

import (
	"github.com/gin-gonic/gin"
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
