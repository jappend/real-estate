package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasSymbol             = regexp.MustCompile(`[^a-zA-z0-9]`)
		hasUpperCaseCharacter = regexp.MustCompile(`[A-Z]`)
		hasLowerCaseCharacter = regexp.MustCompile(`[a-z]`)
		hasNumber             = regexp.MustCompile(`[0-9]`)
	)

	return hasSymbol.MatchString(password) && hasUpperCaseCharacter.MatchString(password) && hasLowerCaseCharacter.MatchString(password) && hasNumber.MatchString(password)
}
