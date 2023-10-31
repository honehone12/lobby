package validator

import (
	lib "github.com/go-playground/validator/v10"
)

type Validator struct {
	inner *lib.Validate
}

func NewValidator() *Validator {
	return &Validator{
		inner: lib.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.inner.Struct(i)
}
