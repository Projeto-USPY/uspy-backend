package controllers

type Login struct {
	ID       string `json:"login" binding:"required,numeric"`
	Password string `json:"pwd" binding:"required,validatePassword"`
	Remember bool   `json:"remember"`
}
