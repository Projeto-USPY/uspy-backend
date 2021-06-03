package controllers

type Subject struct {
	Code           string `json:"code" form:"code" binding:"required,alphanum"`
	CourseCode     string `json:"course" form:"course" binding:"required,alphanum"`
	Specialization string `json:"specialization" form:"specialization" binding:"required,alphanum"`
}
