package account

import (
	"fmt"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/account"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

// VerifyEmail sends a new verification email
func VerifyEmail(ctx *gin.Context, DB db.Env, emailForm *controllers.EmailVerificationSubmission) {
	// check if email exists and if it's not already verified
	emailHash := utils.SHA256(emailForm.Email)
	docs := DB.Client.Collection("users").Where("email", "==", emailHash).Limit(1).Documents(ctx)
	var user models.User

	snaps, err := docs.GetAll()

	if err != nil { // an error happened
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not find email to resend verification:  %s", err.Error()))
		return
	} else if len(snaps) == 0 { // user not found
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("email %s not found in database", emailForm.Email))
		return
	} else if err := snaps[0].DataTo(&user); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error binding user to user object: %s", err.Error()))
		return
	}

	if user.Verified { // already verified
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	// send email
	if err := sendEmailVerification(emailForm.Email, snaps[0].Ref.ID); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to send email verification to user %s; %s", emailForm.Email, err.Error()))
		return
	}

	account.VerifyEmail(ctx)
}
