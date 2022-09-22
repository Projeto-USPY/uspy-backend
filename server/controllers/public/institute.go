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

// GetProfessors is a closure for the GET /professors endpoint
func GetProfessors(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		inst := ctx.MustGet("Institute").(*controllers.Institute)
		public.GetProfessors(ctx, DB, inst)
	}
}
