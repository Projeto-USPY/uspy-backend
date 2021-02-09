package entity

// entity.Reset is a struct used for binding json data when a user wants to reset their password
// see /server/controllers/account.ChangePassword more info
type Reset struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,validatePassword,nefield=OldPassword"`
}
