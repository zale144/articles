package middleware

import "github.com/go-playground/validator/v10"

func NewCustomValidator(val *validator.Validate) *CustomValidator {
	return &CustomValidator{
		validator: val,
	}
}

//CustomValidator struct
type CustomValidator struct {
	validator *validator.Validate
}

//Validate validates the structs
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
