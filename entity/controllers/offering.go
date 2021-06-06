package controllers

type Offering struct {
	Subject   Subject
	Professor string `form:"professor" binding:"required,len=32,alphanum"` // sha256
}
