package utils

import "gopkg.in/validator.v2"

// Validates all the passed in arguments, failing fast on an error
func ValidateAll(args ...interface{}) (err error) {
	for _, arg := range args {
		err = validator.Validate(arg)
		if err != nil {
			return err
		}
	}
	return err
}
