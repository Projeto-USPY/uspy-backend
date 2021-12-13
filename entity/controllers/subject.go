package controllers

// Subject is the object used for looking up subject data. It is probably one of the most used controllers
//
// It also holds a property "limit", used for filtering offering results
type Subject struct {
	Code           string `form:"code" binding:"required,alphanum"`
	CourseCode     string `form:"course" binding:"required,alphanum"`
	Specialization string `form:"specialization" binding:"required,alphanum"`

	Limit int `form:"limit"`
}
