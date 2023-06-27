package account

import (
	"net/http"
	"time"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

// PreSignup is a dummy view method
func PreSignup(ctx *gin.Context, signupToken string) {
	ctx.JSON(http.StatusOK, signupToken)
}

// CompleteSignup is a dummy view method
func CompleteSignup(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// VerifyAccount is a dummy view method
func VerifyAccount(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// Login sets the profile data once it is successful
func Login(ctx *gin.Context, id, name string, lastUpdate time.Time) {
	ctx.JSON(http.StatusOK, views.NewProfile(id, name, lastUpdate))
}

// Logout removes the access token once it is successful
func Logout(ctx *gin.Context) {
	utils.RemoveAccessToken(ctx, !config.Env.IsLocal())
	ctx.Status(http.StatusOK)
}

// ResetPassword is a dummy view method
func ResetPassword(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// ChangePassword is a dummy view method
func ChangePassword(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// Delete removes the access token once it is succesful
func Delete(ctx *gin.Context) {
	utils.RemoveAccessToken(ctx, !config.Env.IsLocal())
	ctx.Status(http.StatusOK)
}

// Update is a dummy view method
func Update(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
