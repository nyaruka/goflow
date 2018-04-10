package types

import (
	"strings"
	"unicode/utf8"
)

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
func (x XString) ToJSON() XString { return MustMarshalToXString(x.Native()) }

// Native returns the native value of this type
func (x XString) Native() string { return string(x) }

// Compare compares this string to another
func (x XString) Compare(other XString) int {
	return strings.Compare(x.Native(), other.Native())
}

// Length returns the length of this string
func (x XString) Length() int { return utf8.RuneCountInString(x.Native()) }

// XStringEmpty is the empty string value
var XStringEmpty = NewXString("")
var _ XPrimitive = XStringEmpty
var _ XLengthable = XStringEmpty
