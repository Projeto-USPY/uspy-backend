package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// DumpErrors is a middleware for dumping all errors that are set in the context
func DumpErrors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, err := range ctx.Errors {
			log.Printf("Got error: %s\n", err.Error())
		}
	}
}
