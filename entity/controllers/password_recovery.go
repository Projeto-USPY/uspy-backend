package controllers

// PasswordRecovery is the object used for recovering a user's password.
//
// It differs from PasswordRecovery because the user is not signed in.
// This is used when the user forgot/lost their password.
type PasswordRecovery struct {
	Token    string `json:"token" binding:"required,validateRecoveryToken"`
	Password string `json:"password" binding:"required,validatePassword"`
}
