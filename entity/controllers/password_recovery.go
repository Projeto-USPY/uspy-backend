package controllers

type PasswordRecovery struct {
	Token    string `json:"token" binding:"required,validateRecoveryToken"`
	Password string `json:"password" binding:"required,validatePassword"`
}
