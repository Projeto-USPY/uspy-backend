// package server contains basic setup functions to start up the web server
package server

import (
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/validation"
	"github.com/Projeto-USPY/uspy-backend/server/controllers/account"
	"github.com/Projeto-USPY/uspy-backend/server/controllers/private"
	"github.com/Projeto-USPY/uspy-backend/server/controllers/public"
	"github.com/Projeto-USPY/uspy-backend/server/controllers/restricted"
	"github.com/Projeto-USPY/uspy-backend/server/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(DB db.Env) (*gin.Engine, error) {
	r := gin.Default() // Create web-server object

	err := validation.SetupValidators()
	if err != nil {
		return nil, err
	}

	r.Use(gin.Recovery(), middleware.DefineDomain())

	if config.Env.IsLocal() {
		r.Use(middleware.AllowAnyOrigin())
	} else {
		if limiter := middleware.RateLimiter(config.Env.RateLimit); limiter != nil {
			r.Use(limiter)
		}
		r.Use(middleware.AllowUSPYOrigin())
	}

	// Login, Logout, Sign-in and other account related operations
	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login", account.Login(DB))

		accountGroup.GET("/captcha", account.SignupCaptcha())
		accountGroup.POST("/create", account.Signup(DB))

		accountGroup.GET("/logout", middleware.JWT(), account.Logout())

		accountGroup.PUT("/password_change", middleware.JWT(), account.ChangePassword(DB))
		accountGroup.PUT("/password_reset", account.ResetPassword(DB))

		accountGroup.GET("/profile", middleware.JWT(), account.Profile(DB))

		accountGroup.DELETE("", middleware.JWT(), account.Delete(DB))
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/subject/all", public.GetSubjects(DB))
		subjectAPI := apiGroup.Group("/subject", middleware.Subject())
		{
			// Available for guests
			subjectAPI.GET("", public.GetSubjectByCode(DB))
			subjectAPI.GET("/relations", public.GetRelations(DB))

			// Restricted means all registered users can see.
			restrictedGroup := apiGroup.Group("/restricted", middleware.JWT())
			{
				subRestricted := restrictedGroup.Group("/subject", middleware.Subject())
				{
					subRestricted.GET("/grades", restricted.GetGrades(DB))
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

			subPrivate.GET("/grade", private.GetSubjectGrade(DB))
		}
	}

	return r, nil
}
