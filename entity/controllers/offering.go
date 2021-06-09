package controllers

type Offering struct {
	Subject Subject
	Hash    string `form:"professor" binding:"required,len=64,alphanum"` // sha256
}
