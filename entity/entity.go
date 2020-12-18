package entity

// User represents an user
type User struct {
	ID       int
	Login    string `json:"login" binding:"required"`
	Password string `json:"pwd" binding:"required"` // used only because of REST requests, do not store in db
}
