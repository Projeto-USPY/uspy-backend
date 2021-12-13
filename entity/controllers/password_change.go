package controllers

// PasswordChange is the object used for changing a user's password.
//
// It differs from PasswordRecovery because the user must be signed in.
type PasswordChange struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,validatePassword,nefield=OldPassword"`
}
