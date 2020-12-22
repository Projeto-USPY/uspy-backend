package middleware

import (
	"github.com/gin-gonic/gin"
)

// AllowAnyOriginMiddleware enables CORS for testing purposes
func AllowAnyOriginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}
