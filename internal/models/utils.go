package models

import (
	"errors"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

// Validator is implementation of validation of request values.
type Validator struct {
	Trans     ut.Translator
	Validator *validator.Validate
}

// Validate do validation for request value.
func (v *Validator) Validate(i interface{}) error {
	err := v.Validator.Struct(i)
	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	// return pretty errors
	msg := ""
	for _, v := range errs.Translate(v.Trans) {
		if msg != "" {
			msg += ", "
		}
		msg += v
	}
	return errors.New(msg)
}
