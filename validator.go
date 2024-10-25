package env

import (
	"fmt"
	"strings"
)

// Validator defines the interface for validating environment variable values.
type Validator interface {
	Validate(envName string, value string) error
}

// ValidatorFactory represents a function to build validator
type ValidatorFactory func(args string) Validator

// requiredValidator checks if a value is present.
type requiredValidator struct{}

func newRequiredValidator(_ string) Validator {
	return &requiredValidator{}
}

func (v requiredValidator) Validate(envName string, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", envName)
	}
	return nil
}

type expectedValueValidator struct {
	expectedValues []string
}

func newExpectedValueValidator(args string) Validator {
	return &expectedValueValidator{
		expectedValues: strings.Split(args, " "),
	}
}

func (v expectedValueValidator) Validate(envName string, value string) error {
	for _, expectedValue := range v.expectedValues {
		if value == expectedValue {
			return nil
		}
	}

	return fmt.Errorf("%s is unexpected value: %s", envName, value)
}
