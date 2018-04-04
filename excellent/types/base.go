package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

func init() {
	decimal.MarshalJSONWithoutQuotes = true
}

// XType represents a type in Excellent
type XType int

// the supported types
const (
	XTypeString XType = iota
	XTypeNumber
	XTypeBool
	XTypeTime
	XTypeArray
	XTypeObject
	XTypeError
	XTypeNil
)

// XValue is the base interface of all Excellent types
type XValue interface {
	Type() XType
	ToJSON() XString
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

// XReducible is the interface for types which can reduce themselves to a primitive
type XReducible interface {
	Reduce() XPrimitive
}

// XString is a string of characters
type XString string

// NewXString creates a new XString
func NewXString(value string) XString {
	return XString(value)
}

// Type is the XType of this type
func (v XString) Type() XType { return XTypeString }

// ToString converts this type to a string
func (v XString) ToString() XString { return v }

// ToBool converts this type to a bool
func (v XString) ToBool() XBool { return string(v) != "" && strings.ToLower(string(v)) != "false" }

// ToJSON converts this type to JSON
func (v XString) ToJSON() XString { return RequireMarshalToXString(v.Native()) }

// Native returns the native value of this type
func (v XString) Native() string { return string(v) }

func (v XString) Length() int { return len(v) }

var NilXString = NewXString("")
var _ XPrimitive = NilXString
var _ XLengthable = NilXString

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

// RequireXNumberFromString creates a new XNumber from the given string
func RequireXNumberFromString(value string) XNumber {
	return XNumber(decimal.RequireFromString(value))
}

// Type is the XType of this type
func (v XNumber) Type() XType { return XTypeNumber }

// ToString converts this type to a string
func (v XNumber) ToString() XString { return XString(v.Native().String()) }

// ToBool converts this type to a bool
func (v XNumber) ToBool() XBool { return XBool(!v.Native().Equals(decimal.Zero)) }

// ToJSON converts this type to JSON
func (v XNumber) ToJSON() XString { return RequireMarshalToXString(v.Native()) }

// Native returns the native value of this type
func (v XNumber) Native() decimal.Decimal { return decimal.Decimal(v) }

var NilXNumber = XNumber(decimal.Zero)
var _ XPrimitive = NilXNumber

// XBool is a boolean true or false
type XBool bool

// NewXBool creates a new XBool
func NewXBool(value bool) XBool {
	return XBool(value)
}

// Type is the XType of this type
func (v XBool) Type() XType { return XTypeBool }

// ToString converts this type to a string
func (v XBool) ToString() XString { return XString(strconv.FormatBool(v.Native())) }

// ToBool converts this type to a bool
func (v XBool) ToBool() XBool { return v }

// ToJSON converts this type to JSON
func (v XBool) ToJSON() XString { return RequireMarshalToXString(v.Native()) }

// Native returns the native value of this type
func (v XBool) Native() bool { return bool(v) }

var NilXBool = NewXBool(false)
var _ XPrimitive = NilXBool

// XTime is a point in time
type XTime time.Time

// NewXTime creates a new XTime
func NewXTime(value time.Time) XTime {
	return XTime(value)
}

// Type is the XType of this type
func (v XTime) Type() XType { return XTypeTime }

// ToString converts this type to a string
func (v XTime) ToString() XString { return XString(utils.DateToISO(v.Native())) }

// ToBool converts this type to a bool
func (v XTime) ToBool() XBool { return XBool(!v.Native().IsZero()) }

// ToJSON converts this type to JSON
func (v XTime) ToJSON() XString { return RequireMarshalToXString(utils.DateToISO(v.Native())) }

// Native returns the native value of this type
func (v XTime) Native() time.Time { return time.Time(v) }

var NilXTime = NewXTime(time.Time{})
var _ XPrimitive = NilXTime

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

// Type is the XType of this type
func (v xerror) Type() XType { return XTypeError }

// ToString converts this type to a string
func (v xerror) ToString() XString { return XString(v.Native().Error()) }

// ToBool converts this type to a bool
func (v xerror) ToBool() XBool { return XBool(false) }

// ToJSON converts this type to JSON
func (v xerror) ToJSON() XString { return RequireMarshalToXString(v.Native().Error()) }

// Native returns the native value of this type
func (v xerror) Native() error { return v.err }

func (v xerror) Error() string { return v.err.Error() }

var NilXError = NewXError(nil)
var _ XError = NilXError

// XObject is the interface for any complex object in Excellent
type XObject interface {
	XValue
	XReducible
}

// BaseXObject is base of any XObject
type BaseXObject struct{}

func (v *BaseXObject) Reduce() XValue { panic("BaseXObject should implement XReducible") }

// Type is the XType of this type
func (v *BaseXObject) Type() XType { return XTypeObject }

// RequireMarshalToXString calls json.Marshal in the given value and panics in the case of an error
func RequireMarshalToXString(v interface{}) XString {
	j, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("unable to marshal %v to JSON", v))
	}
	return XString(j)
}

// ToXString converts the given value to a string
func ToXString(value XValue) XString {
	switch v := value.(type) {
	case XPrimitive:
		return v.ToString()
	case XReducible:
		return v.Reduce().ToString()
	}
	panic(fmt.Sprintf("can't convert type %v to a string", value))
}

// ToXBool converts the given value to a bool
func ToXBool(value XValue) XBool {
	switch v := value.(type) {
	case XPrimitive:
		return v.ToBool()
	case XLengthable:
		return v.Length() > 0
	case XReducible:
		return v.Reduce().ToBool()
	}
	panic(fmt.Sprintf("can't convert type %v to a bool", value))
}
