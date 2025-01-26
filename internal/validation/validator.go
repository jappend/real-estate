package validation

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validate *validator.Validate
}

type ErrorResponse struct {
	Error        bool
	FailedField  string
	Tag          string
	ErrorMessage string
}

func New(validate *validator.Validate) *Validator {
	validate.RegisterValidation("password", passwordValidation)

	return &Validator{
		Validate: validate,
	}
}

func (v *Validator) ValidateData(data interface{}) map[string]string {
	if errors := structValidation(v.Validate, data); len(errors) > 0 && errors[0].Error {
		log.Println("Errors while validating data in the ValidateData function... ", errors)
		errMap := make(map[string]string)

		for _, err := range errors {
			errMap[err.FailedField] = err.ErrorMessage
		}

		return errMap
	}

	return nil
}

// Utilities
func structValidation(validate *validator.Validate, data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse

	errors := validate.Struct(data)
	if errors != nil {
		for _, err := range errors.(validator.ValidationErrors) {
			var errResp ErrorResponse

			errResp.FailedField = strings.ToLower(err.Field())
			errResp.Tag = err.Tag()
			errResp.Error = true

			switch err.Tag() {
			case "required":
				errResp.ErrorMessage = fmt.Sprintf("The '%s' field is required.", errResp.FailedField)
			case "email":
				errResp.ErrorMessage = "The 'email' field needs to be a valid email."
			case "password":
				errResp.ErrorMessage = "The 'password' field needs to have at least 8 characters in length, at least one symbol, one lowercased letter, one uppercased letter and one number."
			}

			validationErrors = append(validationErrors, errResp)
		}
	}

	return validationErrors
}
