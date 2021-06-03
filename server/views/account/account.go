package account

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/iddigital"
	"github.com/Projeto-USPY/uspy-backend/server/middleware"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

func setAccessToken(ctx *gin.Context, token string) {
	domain := ctx.MustGet("front_domain").(string)
	secureCookie := !config.Env.IsLocal()
	ctx.SetCookie("access_token", token, 0, "/", domain, secureCookie, true)
}

func removeAccessToken(ctx *gin.Context) {
	domain := ctx.MustGet("front_domain").(string)
	secureCookie := !config.Env.IsLocal()

	// delete access_token cookie
	ctx.SetCookie("access_token", "", -1, "/", domain, secureCookie, true)
}

// Profile sets the profile data once it is successful
func Profile(ctx *gin.Context, user models.User) {
	if key, ok := os.LookupEnv("AES_KEY"); ok {
		if name, err := utils.AESDecrypt(user.NameHash, key); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error decrypting nameHash: %s", err.Error()))
		} else {
			ctx.JSON(http.StatusOK, views.Profile{User: user.ID, Name: name})
		}
	} else {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("AES_KEY 128/196/256-bit key env variable was not provided"))
	}
}

// Signup sets the records and sets the access token
func Signup(ctx *gin.Context, user *models.User, records iddigital.Records) {
	jwtToken, err := middleware.GenerateJWT(user)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error generating jwt for new user: %s", err.Error()))
		return
	}

	setAccessToken(ctx, jwtToken)
	ctx.JSON(http.StatusOK, records)
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
func Login(ctx *gin.Context, id, name string) {
	ctx.JSON(http.StatusOK, views.Profile{User: id, Name: name})
}

// Logout removes the access token once it is successful
func Logout(ctx *gin.Context) {
	removeAccessToken(ctx)
	ctx.Status(http.StatusOK)
}

func ResetPassword(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func ChangePassword(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// Delete removes the access token once it is succesful
func Delete(ctx *gin.Context) {
	removeAccessToken(ctx)
	ctx.Status(http.StatusOK)
}
