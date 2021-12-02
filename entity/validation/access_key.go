package validation

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
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
