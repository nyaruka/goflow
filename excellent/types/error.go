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
	err error
}

// NewXError creates a new XError
func NewXError(err error) XError {
	return xerror{err: err}
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

// ToString converts this type to a string
func (x xerror) ToString() XString { return XString(x.Native().Error()) }

// ToBool converts this type to a bool
func (x xerror) ToBool() XBool { return XBool(false) }

// ToJSON converts this type to JSON
func (x xerror) ToJSON() XString { return RequireMarshalToXString(x.Native().Error()) }

// Native returns the native value of this type
func (x xerror) Native() error { return x.err }

func (x xerror) Error() string  { return x.err.Error() }
func (x xerror) String() string { return x.Native().Error() }

// NilXError is the nil error value
var NilXError = NewXError(nil)
var _ XError = NilXError

// IsXError returns whether the given value is an error value
func IsXError(x XValue) bool {
	_, isError := x.(XError)
	return isError
}
