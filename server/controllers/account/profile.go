package account

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/account"
	"github.com/gin-gonic/gin"
)

// Profile is a closure for the GET /account/profile endpoint
func Profile(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)

		account.Profile(ctx, DB, userID)
	}
}

// GetMajors is a closure for the GET /account/profile/majors endpoint
func GetMajors(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)

		account.GetMajors(ctx, DB, userID)
	}
}

// SearchCurriculum is a closure for the GET /account/profile/curriculum endpoint
//
// It takes an optional query parameter called "optional", which enforces that queried subjects are not obligatory
func SearchCurriculum(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)

		curriculumQuery := ctx.MustGet("CurriculumQuery").(*controllers.CurriculumQuery)
		account.SearchCurriculum(ctx, DB, userID, curriculumQuery)
	}
}

// GetTranscriptYears is a closure for the GET /account/profile/transcript/years endpoint
//
// It retrieves the last years the user has been in USP
func GetTranscriptYears(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)

		account.GetTranscriptYears(ctx, DB, userID)
	}
}

// SearchTranscript is a closure for the GET /account/profile/transcript endpoint
//
// It looks up into the user's transcript and returns the subjects they took in a given year and semester, along with their grades data
func SearchTranscript(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)
		transcriptQuery := ctx.MustGet("TranscriptQuery").(*controllers.TranscriptQuery)

		account.SearchTranscript(ctx, DB, userID, transcriptQuery)
	}
}
