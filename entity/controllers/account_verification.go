package controllers

// AccountVerification is the object used to verify a recently-registered account
type AccountVerification struct {
	Token string `form:"token" binding:"required,validateVerificationToken"`
}
