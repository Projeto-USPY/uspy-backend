package controllers

// Login is the object used for user authentication
type Login struct {
	ID       string `json:"login" binding:"required,numeric"`
	Password string `json:"pwd" binding:"required,validatePassword"`
	Remember bool   `json:"remember"`
}

type LoginWithGoogle struct {
	Token    string `json:"token" binding:"required"`
	Remember bool   `json:"remember"`
}
