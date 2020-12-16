package main

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server"
	"github.com/tpreischadt/ProjetoJupiter/server/auth"
	"github.com/tpreischadt/ProjetoJupiter/server/data"
	"github.com/tpreischadt/ProjetoJupiter/server/middleware"
)

var (
	client *firestore.Client
	ctx    context.Context
)

func init() {
	client, ctx = db.InitFireStore(".env")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered in main", r)
			_ = client.Close()
			os.Exit(-1)
		}
	}()

	r := gin.Default()
	r.NoRoute(server.DefaultPage)

	api := r.Group("/api")

	err := data.LoadData()
	if err != nil {
		panic(err)
	}

	professorAPI := api.Group("/professor")
	{
		professorAPI.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")

			if id == "all" {
				res := data.GetProfessors()
				c.JSON(http.StatusOK, res)
			} else {
				//prof := data.GetProfessorByID(id)
				//c.JSON(http.StatusOK, prof)
			}
		})
	}

	subjectAPI := api.Group("/subject")
	{
		subjectAPI.GET("/:code", func(c *gin.Context) {
			code := c.Param("code")

			if code == "all" {
				res := data.GetSubjects()
				c.JSON(http.StatusOK, res)
			} else {
				subject := data.GetSubjectByCode(code)
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
				c.Status(http.StatusOK)
			}
		})

		account.POST("/create", func(c *gin.Context) {
			var user entity.User
			if err := c.ShouldBind(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if err := auth.CreateAccount(user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{})
			}

			c.JSON(http.StatusOK, gin.H{})
		})
	}

	r.GET("/profile", middleware.JWTMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	fmt.Println(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
	err = r.Run(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
	log.Print(err)
}
