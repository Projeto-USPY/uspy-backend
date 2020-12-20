package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/populator"
	"github.com/tpreischadt/ProjetoJupiter/server/api"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/server"
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

	apiGroup := r.Group("/api")
	subjectAPI := apiGroup.Group("/subject")
	{
		subjectAPI.GET("/all", api.GetSubjects(DB))
		subjectAPI.GET("", api.GetSubjectByCode(DB))
	}

	account := r.Group("/account")
	{
		account.POST("/login", api.Login(DB))
		account.POST("/create", api.Signup(DB))
	}

	fmt.Println(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
	err := r.Run(os.Getenv("DOMAIN") + ":" + os.Getenv("PORT"))
	panic(err)
}
