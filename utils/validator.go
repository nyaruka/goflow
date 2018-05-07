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
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate will run validation on the given object and return a set of field specific errors in the format:
// field <fieldname> <tag specific message>
//
// For example: "field 'flows' is required"
//
func Validate(obj interface{}) error {
	return validate(obj, "")
}

// ValidateAs will run validation on the given object and return a set of field specific errors in the format:
// field <fieldname> [on <objName>] <tag specific message>
//
// For example: "field 'flows' on 'assets' is required"
//
func ValidateAs(obj interface{}, objName string) error {
	return validate(obj, objName)
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

func validate(obj interface{}, objName string) error {
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

		// the first part of the namespace is always the struct name so either replace that with
		// the provided path or remove it
		parts := strings.SplitN(location, ".", 2)
		if objName != "" {
			parts[0] = objName
			location = strings.Join(parts, ".")
		} else {
			location = strings.Join(parts[1:], ".")
		}

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
		default:
			problem = fmt.Sprintf("failed tag '%s'", fieldErr.Tag())
		}

		newErrors[v] = fmt.Errorf("field '%s' %s", location, problem)
	}
	return ValidationErrors(newErrors)
}
