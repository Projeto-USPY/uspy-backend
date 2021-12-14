package restricted

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/restricted"
	"github.com/gin-gonic/gin"
)

// GetOfferingComments is a closure for the GET /api/restricted/subject/offerings/comments endpoint
func GetOfferingComments(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		off := ctx.MustGet("Offering").(*controllers.Offering)
		off.Subject = *sub

		restricted.GetOfferingComments(ctx, DB, off)
	}
}

// GetOfferingsWithStats is a closure for the GET /api/restricted/subject/offerings endpoint
//
// It differs from GetOfferings because it includes user ratings (approval, disapproval) for each offering
func GetOfferingsWithStats(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		restricted.GetOfferingsWithStats(ctx, DB, sub)
	}
}
