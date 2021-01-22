package models

import (
	"gopkg.in/go-playground/validator.v9"
)

// Validator is implementation of validation of rquest values.
type Validator struct {
	Validator *validator.Validate
}

// Validate do validation for request value.
func (v *Validator) Validate(i interface{}) error {
	err := v.Validator.Struct(i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	return errs
}
