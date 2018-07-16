package utils

import (
	"fmt"
	"reflect"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

// Validator is our system validator, it can be shared across threads
var Validator = validator.New()

func init() {
	// use JSON tags as field names in validation error messages
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			return "-"
		}
		return name
	})
}

// ValidationErrors combines multiple validation errors as a single error
type ValidationErrors []error

// Error returns a string representation of these validation errors
func (e ValidationErrors) Error() string {
	errs := make([]string, len(e))
	for i := range e {
		errs[i] = e[i].Error()
	}
	return strings.Join(errs, ", ")
}

// Validate will run validation on the given object and return a set of field specific errors in the format:
// field <fieldname> <tag specific message>
//
// For example: "field 'flows' is required"
//
func Validate(obj interface{}) error {
	err := Validator.Struct(obj)
	if err == nil {
		return nil
	}
	validationErrs, isValidationErr := err.(validator.ValidationErrors)
	if !isValidationErr {
		return err
	}

	newErrors := make([]error, len(validationErrs))

	for v, fieldErr := range validationErrs {
		location := fieldErr.Namespace()

		// the first part of the namespace is always the struct name so we remove it
		parts := strings.Split(location, ".")[1:]

		// and ignore any parts called - as these come from composition
		newParts := make([]string, 0)
		for _, part := range parts {
			if part != "-" {
				newParts = append(newParts, part)
			}
		}

		location = strings.Join(newParts, ".")

		// generate a more user friendly description of the problem
		var problem string
		switch fieldErr.Tag() {
		case "required":
			problem = "is required"
		case "uuid":
			problem = "must be a valid UUID"
		case "uuid4":
			problem = "must be a valid UUID4"
		case "min":
			problem = fmt.Sprintf("must have a minimum of %s items", fieldErr.Param())
		case "max":
			problem = fmt.Sprintf("must have a maximum of %s items", fieldErr.Param())
		case "mutually_exclusive":
			problem = fmt.Sprintf("is mutually exclusive with '%s'", fieldErr.Param())
		case "http_method":
			problem = "is not a valid HTTP method"
		case "date_format":
			problem = "is not a valid date format"
		case "time_format":
			problem = "is not a valid time format"
		default:
			problem = fmt.Sprintf("failed tag '%s'", fieldErr.Tag())
		}

		newErrors[v] = fmt.Errorf("field '%s' %s", location, problem)
	}
	return ValidationErrors(newErrors)
}
