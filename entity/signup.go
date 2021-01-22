package entity

/* TODO add other input requirements, such as
- Password must be at least 8 characters long
- Password must contain at least one number
- Password must contain at least one letter
- Password must contain at least one special character
- Captcha must be alphanumeric
*/
type Signup struct {
	AccessKey string `json:"access_key" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
}
