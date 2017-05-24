package utils

import validator "gopkg.in/go-playground/validator.v9"

var Validator = validator.New()

// ValidateAll validates all the passed in arguments, failing fast on an error
func ValidateAll(args ...interface{}) (err error) {
	for _, arg := range args {
		err = Validator.Struct(arg)
		if err != nil {
			return err
		}
	}
	return err
}
