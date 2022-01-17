package public

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/server/models/public"
	"github.com/gin-gonic/gin"
)

// GetStats is a closure for the GET /stats endpoint
func GetStats(DB db.Env) func(*gin.Context) {
	return func(ctx *gin.Context) {
		public.GetStats(ctx, DB)
	}
}
