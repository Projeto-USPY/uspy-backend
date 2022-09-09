package public

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/public"
	"github.com/gin-gonic/gin"
)

// GetInstitutes is a closure for the GET /institutes endpoint
func GetInstitutes(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		public.GetInstitutes(ctx, DB)
	}
}

// GetCourses is a closure for the GET /courses endpoint
func GetCourses(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		inst := ctx.MustGet("Institute").(*controllers.Institute)
		public.GetCourses(ctx, DB, inst)
	}
}

// GetSubjects is a closure for the GET /api/subject/all endpoint
func GetSubjects(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		course := ctx.MustGet("Course").(*controllers.Course)
		public.GetAllSubjects(ctx, DB, course)
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
