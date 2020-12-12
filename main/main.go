package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/server"
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

	api.POST("/login", func(c *gin.Context) {

	})

	r.Run()
}
