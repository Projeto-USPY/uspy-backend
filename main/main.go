package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server"
	"github.com/tpreischadt/ProjetoJupiter/server/auth"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	r := gin.Default()
	r.NoRoute(server.DefaultPage)

	api := r.Group("/api")

	err := server.LoadData()
	if err != nil {
		panic(err)
	}

	professorAPI := api.Group("/professor")
	{
		professorAPI.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")

			if id == "all" {
				res := server.GetProfessors()
				c.JSON(http.StatusOK, res)
			} else {
				prof := server.GetProfessorByID(id)
				c.JSON(http.StatusOK, prof)
			}
		})
	}

	subjectAPI := api.Group("/subject")
	{
		subjectAPI.GET("/:code", func(c *gin.Context) {
			code := c.Param("code")

			if code == "all" {
				res := server.GetSubjects()
				c.JSON(http.StatusOK, res)
			} else {
				subject := server.GetSubjectByCode(code)
				c.JSON(http.StatusOK, subject)
			}
		})
	}

	account := r.Group("/account")
	{
		account.POST("/login", func(c *gin.Context) {
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
				c.JSON(http.StatusOK, gin.H{})
			}
		})

		account.POST("/create", func(c *gin.Context) {
			var user entity.User
			if err := c.ShouldBind(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if err := auth.CreateAccount(user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{})
			}

			c.JSON(http.StatusOK, gin.H{})
		})
	}

	fmt.Println(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
	r.Run(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
}
