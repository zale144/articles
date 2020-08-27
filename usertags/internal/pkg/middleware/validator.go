package middleware

import "github.com/go-playground/validator/v10"

func NewCustomValidator(val *validator.Validate) *CustomValidator {
	return &CustomValidator{
		validator: val,
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
