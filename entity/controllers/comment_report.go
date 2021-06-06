package controllers

type CommentReport struct {
	Offering Offering
	Comment  string `form:"comment" binding:"required,len=32,alphanum"`
	Body     string `json:"body" binding:"required,max=300"`
}
