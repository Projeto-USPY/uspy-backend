// package server contains basic setup functions to start up the web server
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server/controllers/account"
	"github.com/tpreischadt/ProjetoJupiter/server/controllers/private"
	"github.com/tpreischadt/ProjetoJupiter/server/controllers/public"
	"github.com/tpreischadt/ProjetoJupiter/server/controllers/restricted"
	"github.com/tpreischadt/ProjetoJupiter/server/middleware"
	"net/http"
	"os"
)

// Todo (return default page)
// Todo2 move this to a separate go file (server.go)
func DefaultPage(c *gin.Context) {
	c.String(http.StatusNotFound, "TODO: Default Page")
}

func SetupRouter(DB db.Env) (*gin.Engine, error) {
	r := gin.Default() // Create web-server object
	r.Use(gin.Recovery())
	r.NoRoute(DefaultPage) // Create a fallback, in case no route matches

	err := entity.SetupValidators()
	if err != nil {
		return nil, err
	}

	if os.Getenv("MODE") == "dev" {
		r.Use(middleware.AllowAnyOrigin())
	}

	// Login, Logout, Sign-in and other account related operations
	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login", account.Login(DB))
		accountGroup.GET("/captcha", account.SignupCaptcha())
		accountGroup.POST("/create", account.Signup(DB))
		accountGroup.GET("/logout", middleware.JWT(), account.Logout())
		accountGroup.PUT("/password", middleware.JWT(), account.ChangePassword(DB))
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/subject/all", public.GetSubjects(DB))
		subjectAPI := apiGroup.Group("/subject", middleware.Subject())
		{
			// Available for guests
			subjectAPI.GET("", public.GetSubjectByCode(DB))
			subjectAPI.GET("/relations", public.GetSubjectGraph(DB))

			// Restricted means all registered users can see.
			restrictedGroup := apiGroup.Group("/restricted", middleware.JWT())
			{
				subRestricted := restrictedGroup.Group("/subject", middleware.Subject())
				{
					subRestricted.GET("/grades", restricted.GetSubjectGrades(DB))
				}
			}
		}
	}

	// Private means the user can only interact with data related to them.
	privateGroup := r.Group("/private", middleware.JWT())
	{
		subPrivate := privateGroup.Group("/subject", middleware.Subject())
		{
			subPrivate.GET("/review", private.GetSubjectReview(DB))
			subPrivate.POST("/review", private.UpdateSubjectReview(DB))
		}
	}

	return r, nil
}
