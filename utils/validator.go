package utils

import (
	"reflect"
	"strings"

	"fmt"

	"errors"

	validator "gopkg.in/go-playground/validator.v9"
)

// Validator is our system validator, it can be shared across threads
var Validator = validator.New()

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

type ValidationErrors []error

func NewValidationErrors(messages ...string) ValidationErrors {
	errs := make([]error, len(messages))
	for m, msg := range messages {
		errs[m] = errors.New(msg)
	}
	return ValidationErrors(errs)
}

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
		fieldName := getFieldName(obj, fieldErr.Field())
		var location string
		var problem string

		location = fmt.Sprintf("'%s'", fieldName)
		if objName != "" {
			location = fmt.Sprintf("'%s' on '%s'", fieldName, objName)
		} else {
			location = fmt.Sprintf("'%s'", fieldName)
		}

		switch fieldErr.Tag() {
		case "required":
			problem = "is required"
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

		newErrors[v] = fmt.Errorf("field %s %s", location, problem)
	}
	return ValidationErrors(newErrors)
}

// utilty to get the name used when marshaling a field to JSON. Returns an empty string if field has no json tag
func getFieldName(obj interface{}, fieldName string) string {
	objType := reflect.TypeOf(obj)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	field, _ := objType.FieldByName(fieldName)
	jsonTag, found := field.Tag.Lookup("json")
	if !found {
		return fieldName
	}

	tagParts := strings.Split(jsonTag, ",")
	return tagParts[0]
}
