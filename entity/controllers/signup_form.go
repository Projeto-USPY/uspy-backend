package controllers

// CompleteSignupForm is used for completing a user's signup.
//
// Signup is a two-step process. First the user provides the access key which is used to
// obtain the signup token. Then the user provides the signup token and the rest of the
// information to complete the signup.
type CompleteSignupForm struct {
	Password    string `json:"password" binding:"required,validatePassword"`
	Email       string `json:"email" binding:"required,email,validateEmail"`
	SignupToken string `json:"signup_token" binding:"required,validateSignupToken"`
}

// AuthForm is used for the first step of the signup process.
//
// It contains information that is used to obtain the users' records through their access key.
type AuthForm struct {
	AccessKey string `json:"access_key" binding:"required,validateAccessKey"`
}
