package controllers

// Major is the object used for looking up a user's registered majors (for example BCC)
type Major struct {
	Course         string `json:"course" binding:"required"`
	Specialization string `json:"specialization" binding:"required"`
}
