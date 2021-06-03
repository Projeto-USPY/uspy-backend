package account

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/account"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Profile is a closure for the GET /account/profile endpoint
func Profile(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		account.Profile(ctx, DB, userID)
	}
}

// ResetPassword is a closure for PUT /account/password_reset
// It differs from ChangePassword because the user does not have to be logged in.
func ResetPassword(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// validate user data
		var signupForm controllers.SignupForm
		if err := ctx.ShouldBindJSON(&signupForm); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.ResetPassword(ctx, DB, &signupForm)
	}
}

// ChangePassword is a closure for PUT /account/password_change
// It differs from ResetPassword because the user must be logged in.
func ChangePassword(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// get user info
		token := ctx.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		var reset controllers.PasswordChange
		// bind old and new password
		if err := ctx.ShouldBindJSON(&reset); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.ChangePassword(ctx, DB, userID, &reset)
	}
}

// Logout is a closure for the GET /account/logout endpoint
func Logout() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		account.Logout(ctx)
	}
}

// Login is a closure for the POST /account/login endpoint
func Login(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var login controllers.Login

		// validate login data
		if err := ctx.ShouldBindJSON(&login); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.Login(ctx, DB, &login)
	}
}

// Signup is a closure for the POST /account/create endpoint
func Signup(DB db.Env) func(g *gin.Context) {
	return func(ctx *gin.Context) {
		// validate user data
		var signupForm controllers.SignupForm
		if err := ctx.ShouldBindJSON(&signupForm); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.Signup(ctx, DB, &signupForm)
	}
}

// SignupCaptcha is a closure for the GET /account/captcha endpoint
func SignupCaptcha() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		account.SignupCaptcha(ctx)
	}
}

// Signup is a closure for the DELETE /account endpoint
func Delete(DB db.Env) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		account.Delete(ctx, DB, userID)
	}
}
