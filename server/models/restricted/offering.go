package restricted

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/gin-gonic/gin"
)

// GetOfferings is a closure for the GET /api/restricted/offerings endpoint
func GetOfferingsWithStats(ctx *gin.Context, DB db.Env, sub *controllers.Subject) {

}
