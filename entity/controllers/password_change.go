package controllers

type PasswordChange struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,validatePassword,nefield=OldPassword"`
}
