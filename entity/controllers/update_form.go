package controllers

// UpdateForm is used for updating a user's profile.
//
// It contains information that is used to obtain user records (captcha)
type UpdateForm struct {
	AccessKey string `json:"access_key" binding:"required,validateAccessKey"`
}
