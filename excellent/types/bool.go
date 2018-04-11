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

// ToXString converts this type to a string
func (x XBool) ToXString() XString { return NewXString(strconv.FormatBool(x.Native())) }

// ToXBool converts this type to a bool
func (x XBool) ToXBool() XBool { return x }

// ToXJSON converts this type to JSON
func (x XBool) ToXJSON() XString { return MustMarshalToXString(x.Native()) }

// Native returns the native value of this type
func (x XBool) Native() bool { return bool(x) }

// Compare compares this bool to another
func (x XBool) Compare(other XBool) int {
	switch {
	case !x.Native() && other.Native():
		return -1
	case x.Native() && !other.Native():
		return 1
	default:
		return 0
	}
}

// XBoolFalse is the false boolean value
var XBoolFalse = NewXBool(false)

// XBoolTrue is the true boolean value
var XBoolTrue = NewXBool(true)

var _ XPrimitive = XBoolFalse
