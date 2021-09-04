package account

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/account"
	"github.com/gin-gonic/gin"
)

// VerifyEmail is a closure for the POST /account/email/verification endpoint
func VerifyEmail(DB db.Env) func(g *gin.Context) {
	return func(ctx *gin.Context) {
		var form controllers.EmailVerificationSubmission
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.VerifyEmail(ctx, DB, &form)
	}
}
