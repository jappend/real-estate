package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type customValidationTests struct {
	have any
	want bool
}

var testValidator *Validator

func TestMain(m *testing.M) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	testValidator = New(validate)
	m.Run()
}

func Test_passwordValidation(t *testing.T) {
	passwordItems := []customValidationTests{
		{"123", false},
		{"abc", false},
		{"123456789", false},
		{"ATestar123", false},
		{"ATestar123!", true},
		{"atestar123!", false},
		{"ATESTAR123!", false},
		{"ATESTAr222@", true},
		{"!!!@%%%%%%%%%%", false},
		{"1223123!21212", false},
	}

	for _, testCase := range passwordItems {
		err := testValidator.Validate.Var(testCase.have, "password")

		if testCase.want {
			assert.NoError(t, err, "Unexpected error for testCase: %v", testCase)
		} else {
			assert.Error(t, err, "Expected error for testCase: %v", testCase)
		}
	}
}
