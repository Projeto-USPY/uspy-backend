package controllers

type SignupForm struct {
	AccessKey string `json:"access_key" binding:"required,validateAccessKey"`
	Password  string `json:"password" binding:"required,validatePassword"`
	Captcha   string `json:"captcha" binding:"required,alphanum,len=4"`
}
