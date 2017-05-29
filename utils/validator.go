package utils

import (
	"fmt"
	"regexp"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

// Validator is our system validator, it can be shared across threads
var Validator = validator.New()

// ValidateAll validates all the passed in arguments, failing fast on an error
func ValidateAll(args ...interface{}) (err error) {
	for _, arg := range args {
		err = Validator.Struct(arg)
		if err != nil {
			break
		}
	}

	// if we got a validation error, rewrite our fields to be snake-case (underscores)
	// as our client is always JSON
	vErrs, isValidation := err.(validator.ValidationErrors)
	if isValidation {
		snakeErr := ValidationError(make([]error, len(vErrs)))

		for i := range vErrs {
			fieldname := strings.TrimRight(underscore(vErrs[i].Field()), "_")
			snakeErr[i] = fmt.Errorf("field '%s' %s", fieldname, vErrs[i].Tag())
		}

		err = snakeErr
	}
	return err
}

// ValidationError is our error type for validation errors
type ValidationError []error

// Error returns a string representation of these validation errors
func (e ValidationError) Error() string {
	errs := make([]string, len(e))
	for i := range e {
		errs[i] = e[i].Error()
	}
	return strings.Join(errs, ", ")
}

// Utility function to convert CamelCase to snake_case for our field names in validation errors

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func underscore(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}
