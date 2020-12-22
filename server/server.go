package server

import (
	"github.com/joho/godotenv"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/populator"
	"github.com/tpreischadt/ProjetoJupiter/server/api"
	"github.com/tpreischadt/ProjetoJupiter/server/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Todo (return default page)
// Todo2 move this to a separate go file (server.go)
func DefaultPage(c *gin.Context) {
	c.String(http.StatusNotFound, "TODO: Default Page")
}

func SetupDB(envPath string) db.Env {
	_ = godotenv.Load(envPath)
	DB := db.InitFireStore(os.Getenv("MODE"))

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

	return DB
}

func SetupRouter(DB db.Env) *gin.Engine {
	r := gin.Default() // Create web-server object
	r.Use(gin.Recovery())
	r.NoRoute(DefaultPage) // Create a fallback, in case no route matches

	if os.Getenv("MODE") == "dev" {
		r.Use(middleware.AllowAnyOriginMiddleware())
	}

	// Login, Logout, Sign-in and other account related operations
	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login", api.Login(DB))
		accountGroup.POST("/create", api.Signup(DB))
	}

	apiGroup := r.Group("/api")
	subjectAPI := apiGroup.Group("/subject")
	{
		// Available for guests
		subjectAPI.GET("", api.GetSubjectByCode(DB))
		subjectAPI.GET("/relations", api.GetSubjectGraph(DB))
		subjectAPI.GET("/all", api.GetSubjects(DB))

		// Restricted means all registered users can see.
		restrictedGroup := apiGroup.Group("/restricted")
		restrictedGroup.Use(middleware.JWTMiddleware())
		{
			restrictedGroup.GET("/grades", api.GetSubjectGrades(DB))
		}
	}

	// Private means the user can only interact with data related to them.
	privateGroup := r.Group("/private")
	privateGroup.Use(middleware.JWTMiddleware())
	{
		reviewGroup := privateGroup.Group("/review")
		reviewGroup.GET("/subject")
		reviewGroup.POST("/subject")
	}

	return r
}
