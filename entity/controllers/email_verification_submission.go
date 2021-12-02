package controllers

type EmailVerificationSubmission struct {
	Email string `json:"email" binding:"required,email,validateEmail"`
}
