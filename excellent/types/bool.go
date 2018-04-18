package types

import (
	"encoding/json"
	"strconv"
)

// XBool is a boolean true or false
type XBool struct {
	baseXPrimitive

	native bool
}

// NewXBool creates a new XBool
func NewXBool(value bool) XBool {
	return XBool{native: value}
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XBool) Reduce() XPrimitive { return x }

// ToXText converts this type to text
func (x XBool) ToXText() XText { return NewXText(strconv.FormatBool(x.Native())) }

// ToXBool converts this type to a bool
func (x XBool) ToXBool() XBool { return x }

// ToXJSON is called when this type is passed to @(json(...))
func (x XBool) ToXJSON() XText { return MustMarshalToXText(x.Native()) }

// Native returns the native value of this type
func (x XBool) Native() bool { return x.native }

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

// MarshalJSON is called when a struct containing this type is marshaled
func (x XBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Native())
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XBool) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &x.native)
}

// XBoolFalse is the false boolean value
var XBoolFalse = NewXBool(false)

// XBoolTrue is the true boolean value
var XBoolTrue = NewXBool(true)

var _ XPrimitive = XBoolFalse
