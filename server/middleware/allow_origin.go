package middleware

import (
	"github.com/gin-gonic/gin"
)

// AllowAnyOrigin enables CORS for testing purposes
func AllowAnyOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}
