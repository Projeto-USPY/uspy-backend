package account

import (
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/server/models/account"
	"github.com/gin-gonic/gin"
)

// Profile is a closure for the GET /account/profile endpoint
func Profile(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)

		account.Profile(ctx, DB, userID)
	}
}

// GetMajors is a closure for the GET /account/profile/majors endpoint
func GetMajors(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)

		account.GetMajors(ctx, DB, userID)
	}
}
