package types

import (
	"fmt"
)

// XError is an error
type XError interface {
	XPrimitive
	error
}

type xerror struct {
	baseXPrimitive

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

// NewXResolveError creates a new XError when a key can't be resolved on an XResolvable
func NewXResolveError(resolvable XResolvable, key string) XError {
	return NewXError(fmt.Errorf("unable to resolve '%s'", key))
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x xerror) Reduce() XPrimitive { return x }

// ToXString converts this type to a string
func (x xerror) ToXString() XString { return NewXString(x.Native().Error()) }

// ToXBool converts this type to a bool
func (x xerror) ToXBool() XBool { return XBoolFalse }

// ToXJSON is called when this type is passed to @(json(...))
func (x xerror) ToXJSON() XString { return MustMarshalToXString(x.Native().Error()) }

// MarshalJSON converts this type to internal JSON
func (x xerror) MarshalJSON() ([]byte, error) {
	return nil, nil
}

// Native returns the native value of this type
func (x xerror) Native() error { return x.native }

func (x xerror) Error() string { return x.Native().Error() }

// NilXError is the nil error value
var NilXError = NewXError(nil)
var _ XError = NilXError

// IsXError returns whether the given value is an error value
func IsXError(x XValue) bool {
	_, isError := x.(XError)
	return isError
}
