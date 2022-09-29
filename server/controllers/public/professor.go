package public

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/public"
	"github.com/gin-gonic/gin"
)

// GetProfessorByCode is a closure for the GET /api/professor endpoint
func GetProfessorByCode(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		prof := ctx.MustGet("Professor").(*controllers.Professor)
		public.GetProfessor(ctx, DB, prof)
	}
}

// GetProfessorOfferings is a closure for the GET /api/professor/offerings endpoint
func GetProfessorOfferings(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		prof := ctx.MustGet("Professor").(*controllers.Professor)
		public.GetProfessorOfferings(ctx, DB, prof)
	}
}
