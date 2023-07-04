package controllers

// InstituteCourse is a controller to specify which course to get course data from
type InstituteCourse struct {
	Institute string `form:"institute" binding:"required,alphanum"`

	Course
}

// Course is a controller to specify which course to get course data from
//
// It differs from Course in that it does not require an institute
type Course struct {
	Code           string `form:"course" binding:"required,alphanum"`
	Specialization string `form:"specialization" binding:"required,alphanum"`
}
