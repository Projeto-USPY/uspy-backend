package public

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/public"
	"github.com/gin-gonic/gin"
)

// GetOfferingComments is a closure for the GET /api/restricted/subject/offerings/comments endpoint
func GetOfferingComments(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		off := ctx.MustGet("Offering").(*controllers.Offering)
		off.Subject = *sub

		public.GetOfferingComments(ctx, DB, off)
	}
}

// GetOfferingsWithStats is a closure for the GET /api/restricted/subject/offerings endpoint
//
// It differs from GetOfferings because it includes user ratings (approval, disapproval) for each offering
func GetOfferingsWithStats(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(*controllers.Subject)
		public.GetOfferingsWithStats(ctx, DB, sub)
	}
}
