package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func validateEmail(f1 validator.FieldLevel) bool {
	email := f1.Field().String()
	return strings.HasSuffix(email, "@usp.br")
}
