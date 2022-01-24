package validator

import (
	"github.com/go-playground/validator"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type BodyRequestValidator struct {
	Validator *validator.Validate
}

func (v *BodyRequestValidator) Validate(schema interface{}) error {
	return v.Validator.Struct(schema)
}