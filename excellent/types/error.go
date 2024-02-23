package types

import (
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
)

// XError is an error
type XError struct {
	baseValue

	native error
}

// NewXError creates a new XError
func NewXError(err error) *XError {
	return &XError{native: err}
}

// NewXErrorf creates a new XError
func NewXErrorf(format string, a ...any) *XError {
	return NewXError(fmt.Errorf(format, a...))
}

// Describe returns a representation of this type for error messages
func (x *XError) Describe() string { return "error" }

// Truthy determines truthiness for this type
func (x *XError) Truthy() bool { return false }

// Render returns the canonical text representation
func (x *XError) Render() string { return x.Native().Error() }

// Format returns the pretty text representation
func (x *XError) Format(env envs.Environment) string { return "" }

// MarshalJSON converts this type to JSON
func (x *XError) MarshalJSON() ([]byte, error) { return jsonx.Marshal(nil) }

// String returns the native string representation of this type for debugging
func (x *XError) String() string { return `XError("` + x.Native().Error() + `")` }

// Native returns the native value of this type
func (x *XError) Native() error {
	if x == nil {
		return nil
	}
	return x.native
}

func (x *XError) Error() string { return x.Native().Error() }

// Equals determines equality for this type
func (x *XError) Equals(o XValue) bool {
	other := o.(*XError)

	return x.String() == other.String()
}

// IsXError returns whether the given value is an error value
func IsXError(x XValue) bool {
	_, isError := x.(*XError)
	return isError
}
