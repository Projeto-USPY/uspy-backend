package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server/auth"
	"net/http"
	"os"
)

func Login(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user entity.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if err := auth.Login(user); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{})
		}

		if jwt, err := auth.GenerateJWT(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			domain := os.Getenv("DOMAIN")

			// expiration date = 1 month
			c.SetCookie("access_token", jwt, 30*24*3600, "/", domain, false, true)
			c.Status(http.StatusOK)
		}
	}
}

// TODO
func Signup(DB db.Env) func(g *gin.Context) {
	return func(c *gin.Context) {

	}
}
