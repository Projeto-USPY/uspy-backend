package account

import (
	"net/http"
	"time"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/iddigital"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

// Signup sets the records
func Signup(ctx *gin.Context, userID string, records iddigital.Transcript) {
	ctx.JSON(http.StatusOK, views.NewTranscript(&records))
}

// VerifyAccount is a dummy view method
func VerifyAccount(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// SignupCaptcha sets the captcha and cookie data in the response
func SignupCaptcha(ctx *gin.Context, resp *http.Response) {
	defer resp.Body.Close()

	cookies := resp.Cookies()
	for _, ck := range cookies {
		domain := ctx.MustGet("front_domain").(string)
		secureCookie := !config.Env.IsLocal()
		ctx.SetCookie(ck.Name, ck.Value, ck.MaxAge, "/", domain, secureCookie, true)
	}

	ctx.DataFromReader(
		http.StatusOK,
		resp.ContentLength,
		resp.Header.Get("Content-Type"),
		resp.Body,
		map[string]string{},
	)
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
