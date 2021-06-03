// for backend-db communication, see /server/models/restricted
package restricted

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/restricted"
	"github.com/gin-gonic/gin"
)

// GetGrades is a closure for the GET /api/restricted/subject/grades endpoint
func GetGrades(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sub := ctx.MustGet("Subject").(controllers.Subject)
		restricted.GetGrades(ctx, DB, &sub)
	}
}
