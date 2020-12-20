package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/populator"
	"github.com/tpreischadt/ProjetoJupiter/server/data/professor"
	"github.com/tpreischadt/ProjetoJupiter/server/data/subject"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server"
	"github.com/tpreischadt/ProjetoJupiter/server/auth"
	"github.com/tpreischadt/ProjetoJupiter/server/middleware"
)

var (
	DB db.Env
)

func init() {
	_ = godotenv.Load(".env")
	DB = db.InitFireStore(os.Getenv("MODE"))

	if os.Getenv("MODE") == "build" { // populate and exit
		func() {
			cnt, err := populator.PopulateICMCOfferings(DB)
			if err != nil {
				_ = DB.Client.Close()
				log.Fatalln("failed to build: ", err)
			} else {
				log.Println("total: ", cnt)
			}

			cnt, err = populator.PopulateICMCSubjects(DB)
			if err != nil {
				_ = DB.Client.Close()
				log.Fatalln("failed to build: ", err)
			} else {
				log.Println("total: ", cnt)
			}
		}()

		_ = DB.Client.Close()
		os.Exit(0)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered in main", r)
			_ = DB.Client.Close()
			os.Exit(-1)
		}
	}()

	r := gin.Default()
	r.NoRoute(server.DefaultPage)

	api := r.Group("/api")

	professorAPI := api.Group("/professor")
	{
		professorAPI.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")

			if id == "all" {
				res, _ := professor.GetProfessors(DB)
				c.JSON(http.StatusOK, res)
			} else {
				prof := professor.GetProfessorByID(id)
				c.JSON(http.StatusOK, prof)
			}
		})
	}

	subjectAPI := api.Group("/subject")
	{
		subjectAPI.GET("/all", func(c *gin.Context) {

		})

		subjectAPI.GET("", func(c *gin.Context) {
			sub := entity.Subject{}
			bindErr := c.BindQuery(&sub)
			if bindErr != nil {
				return
			}

			sub, err := subject.GetByCode(DB, sub.Code)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}

			c.JSON(http.StatusOK, sub)
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
	err := r.Run(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
	log.Print(err)
}
