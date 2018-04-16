package types

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

// XString is a string of characters
type XString struct {
	baseXPrimitive

	native string
}

// NewXString creates a new XString
func NewXString(value string) XString {
	return XString{native: value}
}

// Reduce returns the primitive version of this type (i.e. itself)
func (x XString) Reduce() XPrimitive { return x }

// String converts this type to native string
func (x XString) String() string {
	return x.Native()
}

// ToXString converts this type to a string
func (x XString) ToXString() XString { return x }

// ToXBool converts this type to a bool
func (x XString) ToXBool() XBool {
	return NewXBool(!x.Empty() && strings.ToLower(x.Native()) != "false")
}

// ToXJSON is called when this type is passed to @(to_json(...))
func (x XString) ToXJSON() XString { return MustMarshalToXString(x.Native()) }

// Native returns the native value of this type
func (x XString) Native() string { return x.native }

// Compare compares this string to another
func (x XString) Compare(other XString) int {
	return strings.Compare(x.Native(), other.Native())
}

// Length returns the length of this string
func (x XString) Length() int { return utf8.RuneCountInString(x.Native()) }

// Empty returns whether this is an empty string
func (x XString) Empty() bool { return x.Native() == "" }

// MarshalJSON is called when a struct containing this type is marshaled
func (x XString) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Native())
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XString) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &x.native)
}

// XStringEmpty is the empty string value
var XStringEmpty = NewXString("")
var _ XPrimitive = XStringEmpty
var _ XLengthable = XStringEmpty
