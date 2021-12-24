package controllers

// TranscriptQuery is the object that holds data needed to search in a user's transcript
type TranscriptQuery struct {
	// Subject query data
	Course         string `form:"course" binding:"required"`
	Specialization string `form:"specialization" binding:"required"`
	Optional       bool   `form:"optional"`
	Semester       int    `form:"semester" binding:"required"`
}
