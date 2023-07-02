//Package server contains basic setup functions to start up the web server
package server

import (
	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/Projeto-USPY/uspy-backend/entity/validation"
	"github.com/Projeto-USPY/uspy-backend/server/controllers/account"
	"github.com/Projeto-USPY/uspy-backend/server/controllers/private"
	"github.com/Projeto-USPY/uspy-backend/server/controllers/public"
	"github.com/Projeto-USPY/uspy-backend/server/controllers/restricted"
	"github.com/Projeto-USPY/uspy-backend/server/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setupAccount(DB db.Database, accountGroup *gin.RouterGroup) {
	accountGroup.DELETE("", middleware.JWT(), account.Delete(DB))
	accountGroup.GET("/captcha", account.SignupCaptcha())
	accountGroup.GET("/logout", middleware.JWT(), account.Logout())
	accountGroup.POST("/login", account.Login(DB))
	accountGroup.POST("/login-with-google", account.LoginWithGoogle(DB))
	accountGroup.POST("/create", account.Signup(DB))
	accountGroup.PUT("/password_change", middleware.JWT(), account.ChangePassword(DB))
	accountGroup.PUT("/password_reset", account.ResetPassword(DB))
	accountGroup.GET("/verify", account.VerifyAccount(DB))
	accountGroup.PUT("/update", middleware.JWT(), account.Update(DB))

	emailGroup := accountGroup.Group("/email")
	{
		emailGroup.POST("/verification", account.VerifyEmail(DB))
		emailGroup.POST("/password_reset", account.RequestPasswordReset(DB))
	}

	profileGroup := accountGroup.Group("/profile", middleware.JWT())
	{
		profileGroup.GET("", account.Profile(DB))
		profileGroup.GET("/majors", account.GetMajors(DB))
		profileGroup.GET("/curriculum", entity.CurriculumQueryBinder, account.SearchCurriculum(DB))

		transcriptGroup := profileGroup.Group("/transcript")
		{
			transcriptGroup.GET("", entity.TranscriptQueryBinder, account.SearchTranscript(DB))
			transcriptGroup.GET("/years", account.GetTranscriptYears(DB))
		}
	}
}

func setupPublic(DB db.Database, apiGroup *gin.RouterGroup) {
	apiGroup.GET("/stats", public.GetStats(DB))
	apiGroup.GET("/subject/all", entity.CourseBinder, public.GetSubjects(DB))
	apiGroup.GET("/institutes", public.GetInstitutes(DB))
	apiGroup.GET("/courses", entity.InstituteBinder, public.GetCourses(DB))
	apiGroup.GET("/professors", entity.InstituteBinder, public.GetProfessors(DB))
	subjectAPI := apiGroup.Group("/subject", entity.SubjectBinder)
	{
		subjectAPI.GET("", public.GetSubjectByCode(DB))
		subjectAPI.GET("/relations", public.GetRelations(DB))
		subjectAPI.GET("/offerings", public.GetOfferings(DB))
	}
}

func setupRestricted(DB db.Database, restrictedGroup *gin.RouterGroup) {
	subjectAPI := restrictedGroup.Group("/subject", entity.SubjectBinder)
	{
		subjectAPI.GET("/grades", restricted.GetGrades(DB))
		subjectAPI.GET("/offerings", restricted.GetOfferingsWithStats(DB))

		offeringsAPI := subjectAPI.Group("/offerings", entity.OfferingBinder)
		{
			offeringsAPI.GET("/comments", restricted.GetOfferingComments(DB))
		}
	}
}

func setupPrivate(DB db.Database, privateGroup *gin.RouterGroup) {
	subjectAPI := privateGroup.Group("/subject", entity.SubjectBinder)
	{
		subjectAPI.GET("/grade", private.GetSubjectGrade(DB))
		subjectAPI.GET("/review", private.GetSubjectReview(DB))
		subjectAPI.POST("/review", private.UpdateSubjectReview(DB))

		offeringsAPI := subjectAPI.Group("/offerings", entity.OfferingBinder)
		{
			offeringsAPI.GET("/comments", private.GetComment(DB))
			offeringsAPI.PUT("/comments", private.PublishComment(DB))
			offeringsAPI.DELETE("/comments", private.DeleteComment(DB))

			commentsAPI := offeringsAPI.Group("/comments", entity.CommentRatingBinder)
			{
				commentsAPI.GET("/rating", private.GetCommentRating(DB))
				commentsAPI.PUT("/rating", private.RateComment(DB))
				commentsAPI.PUT("/report", private.ReportComment(DB))
			}
		}
	}
}

// SetupRouter takes the database and sets up all routes with their associated closure callbacks
//
// It returns an error if any validation fails to be registered
func SetupRouter(DB db.Database) (*gin.Engine, error) {
	log.Info("setting up router...")

	r := gin.New() // Create web-server object

	err := validation.SetupValidators()
	if err != nil {
		return nil, err
	}

	r.Use(gin.Recovery(), middleware.Logger(), middleware.DefineDomain(), middleware.DumpErrors())
	gin.ForceConsoleColor()

	if config.Env.IsLocal() {
		r.Use(middleware.AllowAnyOrigin())
	} else {
		if limiter := middleware.RateLimiter(config.Env.RateLimit); limiter != nil {
			r.Use(limiter)
		}
		r.Use(middleware.AllowUSPYOrigin())
	}

	// Login, Logout, Sign-in and other account related operations
	setupAccount(DB, r.Group("/account"))

	// Public endpoints: available for all users, including guests
	setupPublic(DB, r.Group("/api"))

	// Restricted endpoints: available only for registered users
	setupRestricted(DB, r.Group("/api/restricted", middleware.JWT()))

	// Private endpoints: every endpoint related to operations that the user utilizes their own data
	setupPrivate(DB, r.Group("/private", middleware.JWT()))

	return r, nil
}
