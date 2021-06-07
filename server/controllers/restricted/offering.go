package restricted

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/restricted"
	"github.com/gin-gonic/gin"
)

// GetOfferings is a closure for the GET /api/restricted/offerings endpoint
func GetOfferingsWithStats(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		restricted.GetOfferingsWithStats(ctx, DB, sub)
	}
}
