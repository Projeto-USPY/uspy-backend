package controllers

type Subject struct {
	Code           string `form:"code" binding:"required,alphanum"`
	CourseCode     string `form:"course" binding:"required,alphanum"`
	Specialization string `form:"specialization" binding:"required,alphanum"`

	Limit int `form:"limit"`
}
