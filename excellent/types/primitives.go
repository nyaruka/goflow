package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

// XValue is the base interface of all Excellent types
type XValue interface {
	ToJSON() XString
	Reduce() XPrimitive
}

// XPrimitive is the base interface of all Excellent primitive types
type XPrimitive interface {
	XValue

	ToString() XString
	ToBool() XBool
}

// XResolvable is the interface for types which can be keyed into, e.g. foo.bar
type XResolvable interface {
	Resolve(key string) XValue
}

// XLengthable is the interface for types which have a length
type XLengthable interface {
	Length() int
}

// XIndexable is the interface for types which can be indexed into, e.g. foo.0. Such objects
// also need to be lengthable so that the engine knows what is a valid index and what isn't.
type XIndexable interface {
	XLengthable

	Index(index int) XValue
}

// XString is a string of characters
type XString string

// NewXString creates a new XString
func NewXString(value string) XString {
	return XString(value)
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XString) Reduce() XPrimitive { return x }

// ToString converts this type to a string
func (x XString) ToString() XString { return x }

// ToBool converts this type to a bool
func (x XString) ToBool() XBool { return string(x) != "" && strings.ToLower(string(x)) != "false" }

// ToJSON converts this type to JSON
func (x XString) ToJSON() XString { return RequireMarshalToXString(x.Native()) }

// Native returns the native value of this type
func (x XString) Native() string { return string(x) }

// Length returns the length of this string
func (x XString) Length() int { return utf8.RuneCountInString(x.Native()) }

// XStringEmpty is the empty string value
var XStringEmpty = NewXString("")
var _ XPrimitive = XStringEmpty
var _ XLengthable = XStringEmpty

// XNumber is any whole or fractional number
type XNumber decimal.Decimal

// NewXNumber creates a new XNumber
func NewXNumber(value decimal.Decimal) XNumber {
	return XNumber(value)
}

// NewXNumberFromInt creates a new XNumber from the given int
func NewXNumberFromInt(value int) XNumber {
	return XNumber(decimal.New(int64(value), 0))
}

// NewXNumberFromInt64 creates a new XNumber from the given int
func NewXNumberFromInt64(value int64) XNumber {
	return XNumber(decimal.New(value, 0))
}

// RequireXNumberFromString creates a new XNumber from the given string
func RequireXNumberFromString(value string) XNumber {
	return XNumber(decimal.RequireFromString(value))
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XNumber) Reduce() XPrimitive { return x }

// ToString converts this type to a string
func (x XNumber) ToString() XString { return XString(x.Native().String()) }

// ToBool converts this type to a bool
func (x XNumber) ToBool() XBool { return XBool(!x.Native().Equals(decimal.Zero)) }

// ToJSON converts this type to JSON
func (x XNumber) ToJSON() XString { return RequireMarshalToXString(x.Native()) }

// Native returns the native value of this type
func (x XNumber) Native() decimal.Decimal { return decimal.Decimal(x) }

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XNumber) UnmarshalJSON(data []byte) error {
	nativePtr := (*decimal.Decimal)(x)
	return nativePtr.UnmarshalJSON(data)
}

// XNumberZero is the zero number value
var XNumberZero = XNumber(decimal.Zero)
var _ XPrimitive = XNumberZero

// XBool is a boolean true or false
type XBool bool

// NewXBool creates a new XBool
func NewXBool(value bool) XBool {
	return XBool(value)
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XBool) Reduce() XPrimitive { return x }

// ToString converts this type to a string
func (x XBool) ToString() XString { return XString(strconv.FormatBool(x.Native())) }

// ToBool converts this type to a bool
func (x XBool) ToBool() XBool { return x }

// ToJSON converts this type to JSON
func (x XBool) ToJSON() XString { return RequireMarshalToXString(x.Native()) }

// Native returns the native value of this type
func (x XBool) Native() bool { return bool(x) }

// XBoolFalse is the false boolean value
var XBoolFalse = NewXBool(false)

// XBoolTrue is the true boolean value
var XBoolTrue = NewXBool(true)

var _ XPrimitive = XBoolFalse

// XTime is a point in time
type XTime time.Time

// NewXTime creates a new XTime
func NewXTime(value time.Time) XTime {
	return XTime(value)
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XTime) Reduce() XPrimitive { return x }

// ToString converts this type to a string
func (x XTime) ToString() XString { return XString(utils.DateToISO(x.Native())) }

// ToBool converts this type to a bool
func (x XTime) ToBool() XBool { return XBool(!x.Native().IsZero()) }

// ToJSON converts this type to JSON
func (x XTime) ToJSON() XString { return RequireMarshalToXString(utils.DateToISO(x.Native())) }

// Native returns the native value of this type
func (x XTime) Native() time.Time { return time.Time(x) }

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XTime) UnmarshalJSON(data []byte) error {
	nativePtr := (*time.Time)(x)
	return nativePtr.UnmarshalJSON(data)
}

// XTimeZero is the zero time value
var XTimeZero = NewXTime(time.Time{})
var _ XPrimitive = XTimeZero

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
	return NewXError(fmt.Errorf("unable to resolve '%s' on %s", key, reflect.TypeOf(resolvable)))
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

func (x xerror) Error() string { return x.err.Error() }

// NilXError is the nil error value
var NilXError = NewXError(nil)
var _ XError = NilXError

// IsXError returns whether the given value is an error value
func IsXError(x XValue) bool {
	_, isError := x.(XError)
	return isError
}

// Compare returns the difference between two given values
func Compare(x1 XValue, x2 XValue) (int, error) {
	if utils.IsNil(x1) && utils.IsNil(x2) {
		return 0, nil
	} else if utils.IsNil(x1) || utils.IsNil(x2) {
		return 0, fmt.Errorf("can't compare non-nil and nil values: %v and %v", x1, x2)
	}

	x1 = x1.Reduce()
	x2 = x2.Reduce()

	if reflect.TypeOf(x1) != reflect.TypeOf(x2) {
		return 0, fmt.Errorf("can't compare different types of %#v and %#v", x1, x2)
	}

	// common types, do real comparisons
	switch typed := x1.(type) {
	case XError:
		return strings.Compare(typed.Error(), x2.(error).Error()), nil
	case XNumber:
		return typed.Native().Cmp(x2.(XNumber).Native()), nil
	case XBool:
		bool1 := typed.Native()
		bool2 := x2.(XBool).Native()

		switch {
		case !bool1 && bool2:
			return -1, nil
		case bool1 == bool2:
			return 0, nil
		case bool1 && !bool2:
			return 1, nil
		}
	case XTime:
		time1 := typed.Native()
		time2 := x2.(XTime).Native()

		switch {
		case time1.Before(time2):
			return -1, nil
		case time1.Equal(time2):
			return 0, nil
		case time1.After(time2):
			return 1, nil
		}
	case XString:
		return strings.Compare(typed.Native(), x2.(XString).Native()), nil
	}

	// TODO: find better fallback
	return strings.Compare(fmt.Sprintf("%v", x1), fmt.Sprintf("%v", x2)), nil
}
