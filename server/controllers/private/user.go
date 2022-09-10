package private

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/private"
	"github.com/gin-gonic/gin"
)

// GetSubjectGrade is a closure for the GET /private/subject/grade endpoint
func GetSubjectGrade(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// get user and subject info
		userID := ctx.MustGet("userID").(string)
		sub := ctx.MustGet("Subject").(*controllers.Subject)

		private.GetSubjectGrade(ctx, DB, userID, sub)
	}
}

// GetSubjectReview is a closure for the GET /private/subject/review endpoint
func GetSubjectReview(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// get user and subject info
		userID := ctx.MustGet("userID").(string)
		sub := ctx.MustGet("Subject").(*controllers.Subject)

		private.GetSubjectReview(ctx, DB, userID, sub)
	}
}

// UpdateSubjectReview is a closure for the POST /private/subject/review endpoint
func UpdateSubjectReview(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)

		// get subject and review data
		sr := controllers.SubjectReview{Subject: *sub}
		if err := ctx.ShouldBindJSON(&sr); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// get user data
		userID := ctx.MustGet("userID").(string)

		private.UpdateSubjectReview(ctx, DB, userID, &sr)
	}
}
