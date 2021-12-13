package controllers

// SignupForm is used for signing up a new user.
//
// It contains information that is used to obtain user records (access key and captcha)
// It also contains information that will be used for later verification and authentication (email and password)
type SignupForm struct {
	AccessKey string `json:"access_key" binding:"required,validateAccessKey"`
	Password  string `json:"password" binding:"required,validatePassword"`
	Captcha   string `json:"captcha" binding:"required,alphanum,len=4"`
	Email     string `json:"email" binding:"required,email,validateEmail"`
}
