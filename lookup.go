package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	envTag               = "env"
	envSep               = ","
	validatorArgsPattern = `(\w+)(?:=([-\w\s]+))?` // (e.g., min=0, expectedValues=development production)
	requiredFlag         = "required"
	expectedValuesFlag   = "expectedValues"
)

var (
	// Regex to capture dynamic validation arguments
	validatorArgsRegex = regexp.MustCompile(validatorArgsPattern)
)

// ValidatorRegistry holds registered validators for easy lookup.
var validatorRegistry = map[string]ValidatorFactory{
	requiredFlag:       newRequiredValidator,
	expectedValuesFlag: newExpectedValueValidator,
}

// AddValidator allows users to add custom validators to the registry
func AddValidator(name string, fn ValidatorFactory) {
	validatorRegistry[name] = fn
}

// Lookup loads environment variables into the provided struct.
func Lookup(m any) error {
	return lookupValue(reflect.ValueOf(m).Elem())
}

// lookupValue recursively processes struct fields for environment variables.
func lookupValue(v reflect.Value) error {
	var errs error

	for i := 0; i < v.Type().NumField(); i++ {
		field := v.Type().Field(i)
		valueField := v.Field(i)

		// Skip unexported fields or unsettable fields
		if !v.Field(i).CanSet() {
			continue
		}

		// Handle struct fields recursively
		if v.Field(i).Kind() == reflect.Struct {
			if err := lookupValue(valueField); err != nil {
				errs = errors.Join(errs, err)
			}

			continue
		}

		// Read the env tag configuration
		envName, validators := parseEnvTag(field.Tag.Get(envTag), envSep)
		if envName == "" {
			continue // No env tag, skip
		}

		// Read env from name
		envValue := os.Getenv(envName)

		// Validate required flags
		for _, validator := range validators {
			if err := validator.Validate(envName, envValue); err != nil {
				errs = errors.Join(errs, err)
			}
		}

		// Parse and set the environment variable value
		if err := setFieldValue(valueField, envName, envValue); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

// setFieldValue sets the value of a struct field based on the environment variable.
func setFieldValue(valueField reflect.Value, envName string, envValue string) error {
	// Skip if envValue empty
	if envValue == "" {
		return nil
	}

	switch valueField.Kind() {
	case reflect.String:
		valueField.SetString(strings.TrimSpace(envValue))
	case reflect.Int:
		n, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing int for %s: %w", envName, err)
		}

		valueField.SetInt(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(envValue)
		if err != nil {
			return fmt.Errorf("error parsing bool for %s: %w", envName, err)
		}

		valueField.SetBool(b)
	default:
		return fmt.Errorf("unsupported type %s for %s", valueField.Kind(), envName)
	}

	return nil
}

// parseEnvTag parses the tag and returns the environment variable name and a list of validators.
func parseEnvTag(tag string, sep string) (string, []Validator) {
	if tag == "-" || tag == "" {
		return "", nil
	}

	sepIndex := strings.Index(tag, sep)
	if sepIndex == -1 {
		return tag, nil
	}

	envName := tag[:sepIndex]
	options := tag[sepIndex+len(sep):]

	validators := make([]Validator, 0, strings.Count(options, sep)+1)

	// Iterate over each validation option (after the env name)
	for _, option := range strings.Split(options, sep) {
		matches := validatorArgsRegex.FindStringSubmatch(option)
		if len(matches) > 0 {
			validatorName := matches[1]
			validatorArgs := matches[2]

			// Lookup the validator in the registry and create it
			if createValidator, ok := validatorRegistry[validatorName]; ok {
				validators = append(validators, createValidator(validatorArgs))
			}

			continue
		}

		validators = append(validators, validatorRegistry[option](""))
	}

	return envName, validators
}
