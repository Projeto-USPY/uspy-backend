package account

import (
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/server/models/account"
	"github.com/gin-gonic/gin"
)

// ResetPassword is a closure for PUT /account/password_reset
// It differs from ChangePassword because the user does not have to be logged in.
func ResetPassword(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// validate user data
		var recovery controllers.PasswordRecovery
		if err := ctx.ShouldBindJSON(&recovery); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.ResetPassword(ctx, DB, &recovery)
	}
}

// ChangePassword is a closure for PUT /account/password_change
// It differs from ResetPassword because the user must be logged in.
func ChangePassword(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// get user info
		userID := ctx.MustGet("userID").(string)

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
func Login(DB db.Database) func(ctx *gin.Context) {
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

// PreSignup is a closure for the PUT /account/auth endpoint
func PreSignup(DB db.Database) func(g *gin.Context) {
	return func(ctx *gin.Context) {
		// validate user data
		var authForm controllers.AuthForm
		if err := ctx.ShouldBindJSON(&authForm); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.PreSignup(ctx, DB, &authForm)
	}
}

// CompleteSignup is a closure for the POST /account/create endpoint
func CompleteSignup(DB db.Database) func(g *gin.Context) {
	return func(ctx *gin.Context) {
		// validate user data
		var signupForm controllers.CompleteSignupForm
		if err := ctx.ShouldBindJSON(&signupForm); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.CompleteSignup(ctx, DB, &signupForm)
	}
}

// VerifyAccount is a closure for the GET /account/verify endpoint
func VerifyAccount(DB db.Database) func(g *gin.Context) {
	return func(ctx *gin.Context) {
		// validate verification token
		var verification controllers.AccountVerification
		if err := ctx.ShouldBindQuery(&verification); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.VerifyAccount(ctx, DB, &verification)
	}
}

// Delete is a closure for the DELETE /account endpoint
func Delete(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)

		account.Delete(ctx, DB, userID)
	}
}

// Update is a closure for the PUT /account/update endpoint
func Update(DB db.Database) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("userID").(string)

		// validate update data
		var updateForm controllers.UpdateForm
		if err := ctx.ShouldBindJSON(&updateForm); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		account.Update(ctx, DB, userID, &updateForm)
	}
}
