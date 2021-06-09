package controllers

type CommentReport struct {
	Offering Offering
	Comment  string `form:"comment" binding:"required,uuid"`
	Body     string `json:"body" binding:"required,min=10,max=300"`
}
