package public

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/public"
	"github.com/gin-gonic/gin"
)

// GetOfferings is a closure for the GET /api/offerings endpoint
func GetOfferings(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		public.GetOfferings(ctx, DB, sub)
	}
}
