package controllers

// Offering is the object used for looking up offerings, e.g., subjects in the context of a given professor.
type Offering struct {
	Subject
	Hash string `form:"professor" binding:"required,len=64,alphanum"` // sha256
}
