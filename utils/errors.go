package utils

import (
	"errors"
	"fmt"
)

// RichError is an parameterized error
type RichError struct {
	msg    string
	Domain string
	Code   string
	Extra  map[string]string
}

// NewRichError creates a new rich error
func NewRichError(domain, code, err string, args ...any) *RichError {
	return &RichError{Domain: domain, Code: code, msg: fmt.Sprintf(err, args...)}
}

func (e *RichError) WithExtra(k, v string) *RichError {
	if e.Extra == nil {
		e.Extra = make(map[string]string)
	}
	e.Extra[k] = v
	return e
}

// Error returns the error message
func (e *RichError) Error() string {
	return e.msg
}

// IsRichError is a utility to find a RichError in the error chain
func IsRichError(err error) (bool, *RichError) {
	var rerr *RichError
	return errors.As(err, &rerr), rerr
}
