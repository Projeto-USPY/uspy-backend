package entity

import "github.com/go-playground/validator/v10"

var Validators = map[string]func(validator.FieldLevel) bool{
	"validatePassword":  ValidatePassword,
	"validateAccessKey": ValidateAccessKey,
}
