package utils

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"

	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

// our system validator, it can be shared across threads
var valx = validator.New()

// ErrorMessageFunc is the type for a function that can convert a field error to user friendly message
type ErrorMessageFunc func(validator.FieldError) string

var messageFuncs = map[string]ErrorMessageFunc{
	"required": func(e validator.FieldError) string { return "is required" },
	"uuid":     func(e validator.FieldError) string { return "must be a valid UUID" },
	"uuid4":    func(e validator.FieldError) string { return "must be a valid UUID4" },
	"url":      func(e validator.FieldError) string { return "is not a valid URL" },
	"min":      func(e validator.FieldError) string { return fmt.Sprintf("must have a minimum of %s items", e.Param()) },
	"max":      func(e validator.FieldError) string { return fmt.Sprintf("must have a maximum of %s items", e.Param()) },
	"mutually_exclusive": func(e validator.FieldError) string {
		return fmt.Sprintf("is mutually exclusive with '%s'", e.Param())
	},
}

func init() {
	// use JSON tags as field names in validation error messages
	valx.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			return "-"
		}
		return name
	})

	RegisterValidatorAlias("http_method", "eq=GET|eq=HEAD|eq=POST|eq=PUT|eq=PATCH|eq=DELETE", func(validator.FieldError) string {
		return "is not a valid HTTP method"
	})
}

// RegisterValidatorTag registers a tag
func RegisterValidatorTag(tag string, fn validator.Func, message ErrorMessageFunc) {
	valx.RegisterValidation(tag, fn)

	messageFuncs[tag] = message
}

// RegisterValidatorAlias registers a tag alias
func RegisterValidatorAlias(alias, tags string, message ErrorMessageFunc) {
	valx.RegisterAlias(alias, tags)

	messageFuncs[alias] = message
}

// RegisterStructValidator registers a struct level validator
func RegisterStructValidator(fn validator.StructLevelFunc, types ...interface{}) {
	valx.RegisterStructValidation(fn, types...)
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
	var err error

	// gets the value stored in the interface var, and if it's a pointer, dereferences it
	v := reflect.Indirect(reflect.ValueOf(obj))

	if v.Type().Kind() == reflect.Slice {
		err = valx.Var(obj, `required,dive`)
	} else {
		err = valx.Struct(obj)
	}

	if err == nil {
		return nil
	}

	validationErrs, isValidationErr := err.(validator.ValidationErrors)
	if !isValidationErr {
		return err
	}

	newErrors := make([]error, len(validationErrs))

	for i, fieldErr := range validationErrs {
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
		messageFunc := messageFuncs[fieldErr.Tag()]
		if messageFunc != nil {
			problem = messageFunc(fieldErr)
		} else {
			problem = fmt.Sprintf("failed tag '%s'", fieldErr.Tag())
		}

		newErrors[i] = errors.Errorf("field '%s' %s", location, problem)
	}
	return ValidationErrors(newErrors)
}

// UnmarshalAndValidate is a convenience function to unmarshal an object and validate it
func UnmarshalAndValidate(data []byte, obj interface{}) error {
	err := jsonx.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return Validate(obj)
}

// UnmarshalAndValidateWithLimit unmarshals a struct with a limit on how many bytes can be read from the given reader
func UnmarshalAndValidateWithLimit(reader io.ReadCloser, s interface{}, limit int64) error {
	if err := jsonx.UnmarshalWithLimit(reader, s, limit); err != nil {
		return err
	}
	return Validate(s)
}
