package types

import (
	"strconv"
)

// XBool is a boolean true or false
type XBool bool

// NewXBool creates a new XBool
func NewXBool(value bool) XBool {
	return XBool(value)
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XBool) Reduce() XPrimitive { return x }

// ToString converts this type to a string
func (x XBool) ToString() XString { return NewXString(strconv.FormatBool(x.Native())) }

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
