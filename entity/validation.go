/* Package db contains useful functions related to the Firestore Database */
package entity

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validators = map[string]func(validator.FieldLevel) bool{
	"validatePassword":      validatePassword,
	"validateAccessKey":     validateAccessKey,
	"validateSubjectReview": validateSubjectReview,
}

// SetupValidators registers the default validation functions designed for each entity
func SetupValidators() error {
	for key, value := range validators {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			err := v.RegisterValidation(key, value)
			if err != nil {
				return err
			}
		} else {
			return errors.New("failed to setup validators")
		}
	}

	return nil
}
