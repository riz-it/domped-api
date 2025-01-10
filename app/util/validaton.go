package util

import (
	"github.com/go-playground/validator/v10"
	"riz.it/domped/app/domain"
)

func Validate(v *validator.Validate, i interface{}) []domain.ValidationError {
	validationErrors := []domain.ValidationError{}

	errs := v.Struct(i)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			param := err.Param()
			value := err.Value()

			validationErrors = append(validationErrors, domain.ValidationError{
				FailedField: ConvertToSpaced(err.Field()),
				Tag:         err.Tag(),
				Value:       &value,
				Param:       &param,
			})
		}
	}

	return validationErrors
}
