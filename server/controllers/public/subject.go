package public

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/public"
	"github.com/gin-gonic/gin"
)

// GetSubject is a closure for the GET /api/subject/all endpoint
func GetSubjects(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		public.GetAllSubjects(ctx, DB)
	}
}

// GetSubjectByCode is a closure for the GET /api/subject endpoint
func GetSubjectByCode(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		public.Get(ctx, DB, sub)
	}
}

// GetRelations is a closure for the GET /api/subject/relations endpoint
func GetRelations(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		public.GetRelations(ctx, DB, sub)
	}
}
