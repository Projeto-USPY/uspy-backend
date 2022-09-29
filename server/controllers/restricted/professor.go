package restricted

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/restricted"
	"github.com/gin-gonic/gin"
)

// GetProfessorComments is a closure for the GET /api/restricted/professor/comments endpoint
func GetProfessorComments(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		prof := ctx.MustGet("Professor").(*controllers.Professor)
		restricted.GetProfessorComments(ctx, DB, prof)
	}
}

// GetProfessorOfferings is a closure for the GET /api/restricted/professor/offerings endpoint
func GetProfessorOfferings(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		prof := ctx.MustGet("Professor").(*controllers.Professor)
		restricted.GetProfessorOfferings(ctx, DB, prof)
	}
}
