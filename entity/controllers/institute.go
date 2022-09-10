package controllers

// Institute is a query parameter to specify which institute to get course data from
type Institute struct {
	Code string `form:"institute" binding:"required,alphanum"`
}
