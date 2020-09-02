package types

import (
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
)

// XError is an error
type XError interface {
	error
	XValue
	Equals(XError) bool
}

type xerror struct {
	native error
}

// NewXError creates a new XError
func NewXError(err error) XError {
	return xerror{native: err}
}

// NewXErrorf creates a new XError
func NewXErrorf(format string, a ...interface{}) XError {
	return NewXError(fmt.Errorf(format, a...))
}

// Describe returns a representation of this type for error messages
func (x xerror) Describe() string { return "error" }

// Truthy determines truthiness for this type
func (x xerror) Truthy() bool { return false }

// Render returns the canonical text representation
func (x xerror) Render() string { return x.Native().Error() }

// Format returns the pretty text representation
func (x xerror) Format(env envs.Environment) string { return "" }

// MarshalJSON converts this type to JSON
func (x xerror) MarshalJSON() ([]byte, error) { return jsonx.Marshal(nil) }

// String returns the native string representation of this type for debugging
func (x xerror) String() string { return `XError("` + x.Native().Error() + `")` }

// Native returns the native value of this type
func (x xerror) Native() error { return x.native }

func (x xerror) Error() string { return x.Native().Error() }

// Equals determines equality for this type
func (x xerror) Equals(other XError) bool {
	return x.String() == other.String()
}

// NilXError is the nil error value
var NilXError = NewXError(nil)
var _ XError = NilXError

// IsXError returns whether the given value is an error value
func IsXError(x XValue) bool {
	_, isError := x.(XError)
	return isError
}
