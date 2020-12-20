package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Todo (return default page)
// Todo2 move this to a separate go file (server.go)
func DefaultPage(c *gin.Context) {
	c.String(http.StatusNotFound, "TODO: Default Page")
}
