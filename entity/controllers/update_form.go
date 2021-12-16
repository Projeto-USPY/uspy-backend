package controllers

// UpdateForm is used for updating a user's profile.
//
// It contains information that is used to obtain user records (access key and captcha)
type UpdateForm struct {
	AccessKey string `json:"access_key" binding:"required,validateAccessKey"`
	Captcha   string `json:"captcha" binding:"required,alphanum,len=4"`
}
