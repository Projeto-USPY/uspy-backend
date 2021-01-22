package entity

import (
	"github.com/go-playground/validator/v10"
	"strings"
	"unicode"
)

func ValidateAccessKey(f1 validator.FieldLevel) bool {
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

func ValidatePassword(f1 validator.FieldLevel) bool {
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

type Signup struct {
	AccessKey string `json:"access_key" binding:"required,validateAccessKey"`
	Password  string `json:"password" binding:"required,validatePassword"`
	Captcha   string `json:"captcha" binding:"required,alphanum,len=4"`
	Terms     bool   `json:"terms" binding:"required,eq=true"`
}
