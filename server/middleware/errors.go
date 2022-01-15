package middleware

import (
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// DumpErrors is a middleware for dumping all errors that are set in the context
func DumpErrors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		for _, err := range ctx.Errors {
			log.Error(err)
		}
	}
}
