package controllers

// Professor is the object used for looking up professor data.
type Professor struct {
	Hash      string `form:"professor" binding:"required,len=64,alphanum"`
	Institute string `form:"institute" binding:"required,alphanum"`
}
