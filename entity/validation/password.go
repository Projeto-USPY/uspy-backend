package validation

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

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
