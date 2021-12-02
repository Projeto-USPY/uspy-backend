package controllers

type AccountVerification struct {
	Token string `form:"token" binding:"required,validateVerificationToken"`
}
