package controllers

// Course is a query parameter to specify which Course to get course data from
type Course struct {
	Institute string `form:"institute" binding:"required,alphanum"`

	Code           string `form:"course" binding:"required,alphanum"`
	Specialization string `form:"specialization" binding:"required,alphanum"`
}
