package controllers

// TranscriptQuery is the object used for querying a user's transcript
type TranscriptQuery struct {
	Year     int `form:"year" binding:"required"`
	Semester int `form:"semester" binding:"required,oneof=1 2"`
}
