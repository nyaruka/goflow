package utils

import (
	"fmt"
	"regexp"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

// Validator is our system validator, it can be shared across threads
var Validator = validator.New()

// ValidateAllUnlessErr validates all the passed in arguments, failing fast on any error including the passed in one
func ValidateAllUnlessErr(err error, args ...interface{}) error {
	if err != nil {
		return err
	}
	return ValidateAll(args...)
}

// ValidateAll validates all the passed in arguments, failing fast on an error
func ValidateAll(args ...interface{}) (err error) {
	for _, arg := range args {
		if arg == nil {
			continue
		}

		err = Validator.Struct(arg)
		if err != nil {
			errFormat := "%sfield '%s' %s"

			// see if we are a typed envelope, if so can provide better errors
			typeDesc := ""
			typed, isTyped := arg.(Typed)
			fmt.Printf("%#v is typed: %s\n", arg, isTyped)
			if isTyped {
				errFormat = "%s: field '%s' %s"
				typeDesc = typed.Type()
			}

			// if we got a validation error, rewrite our fields to be snake-case (underscores)
			// as our client is always JSON
			vErrs, isValidation := err.(validator.ValidationErrors)
			if isValidation {
				snakeErr := ValidationError(make([]error, len(vErrs)))

				for i := range vErrs {
					fieldname := strings.TrimRight(underscore(vErrs[i].Field()), "_")
					snakeErr[i] = fmt.Errorf(errFormat, typeDesc, fieldname, vErrs[i].Tag())
				}

				err = snakeErr
			}
			return err
		}
	}
	return nil
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
