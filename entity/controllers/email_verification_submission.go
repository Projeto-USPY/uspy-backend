package controllers

// EmailVerificationSubmission is the object that holds the email used for account verification
type EmailVerificationSubmission struct {
	Email string `json:"email" binding:"required,email,validateEmail"`
}
