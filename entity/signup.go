/* Package db contains useful functions related to the Firestore Database */
package entity

import (
	"github.com/go-playground/validator/v10"
	"strings"
	"unicode"
)

func validateAccessKey(f1 validator.FieldLevel) bool {
	auth := f1.Field().String()
	fields := strings.Split(auth, "-")

	if len(fields) != 4 {
		return false
	}

	for _, f := range fields {
		for _, r := range f {
			if !unicode.IsDigit(r) && !unicode.IsUpper(r) {
				return false
			}
		}
	}

	return true
}

func validatePassword(f1 validator.FieldLevel) bool {
	pass := f1.Field().String()
	if len(pass) < 8 || len(pass) > 20 {
		return false
	}
	letter, num, symbol := 0, 0, 0
	for _, c := range pass {
		if unicode.IsLetter(c) {
			letter++
		} else if unicode.IsDigit(c) {
			num++
		} else if unicode.IsGraphic(c) {
			symbol++
		}
	}

	return letter > 0 && num > 0 && symbol > 0
}

// entity.Signup is a struct used for signup input validation
// see /server/controllers/account.Signup for more details
type Signup struct {
	AccessKey string `json:"access_key" binding:"required,validateAccessKey"`
	Password  string `json:"password" binding:"required,validatePassword"`
	Captcha   string `json:"captcha" binding:"required,alphanum,len=4"`
}
