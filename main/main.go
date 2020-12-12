package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server"
	auth "github.com/tpreischadt/ProjetoJupiter/server/auth"
)

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
			if err := c.ShouldBind(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			}

			if err := auth.Login(user); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{})
			}
		})

		account.POST("/create", func(c *gin.Context) {
			var user entity.User
			if err := c.ShouldBind(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
				return
			}

			if err := auth.CreateAccount(user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{})
			}

			c.JSON(http.StatusOK, gin.H{})
		})
	}

	r.Run()
}
