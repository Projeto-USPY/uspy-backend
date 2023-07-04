package public

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/public"
	"github.com/gin-gonic/gin"
)

// ListSubjects is a closure for the GET /api/subject/list endpoint
func ListSubjects(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		course := ctx.MustGet("Course").(*controllers.Course)
		public.ListSubjectsByCourse(ctx, DB, course)
	}
}

// GetSubjects is a closure for the GET /api/subject/search endpoint
func GetSubjects(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		course := ctx.MustGet("InstituteCourse").(*controllers.InstituteCourse)
		public.SearchSubjects(ctx, DB, course)
	}
}

// GetSubjectByCode is a closure for the GET /api/subject endpoint
func GetSubjectByCode(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		public.Get(ctx, DB, sub)
	}
}

// GetRelations is a closure for the GET /api/subject/relations endpoint
func GetRelations(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		public.GetRelations(ctx, DB, sub)
	}
}

// GetSiblingSubjects is a closure for the GET /api/subject/siblings endpoint
func GetSiblingSubjects(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		public.GetSiblingSubjects(ctx, DB, sub)
	}
}
