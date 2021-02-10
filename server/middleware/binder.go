// package middleware contains useful middleware handlers
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"net/http"
)

// Subject is a middleware for binding subject data
func Subject() gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := entity.Subject{}
		if err := c.BindQuery(&subject); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Set("Subject", subject)
	}
}