package types

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
)

// XText is a string of characters.
//
//   @("abc") -> abc
//   @(text_length("abc")) -> 3
//   @(upper("abc")) -> ABC
//   @(json("abc")) -> "abc"
//
// @type text
type XText struct {
	native string
}

// NewXText creates a new text value
func NewXText(value string) XText {
	return XText{native: value}
}

// Describe returns a representation of this type for error messages
func (x XText) Describe() string { return fmt.Sprintf(`"%s"`, x.Native()) }

// Truthy determines truthiness for this type
func (x XText) Truthy() bool {
	return !x.Empty() && strings.ToLower(x.Native()) != "false"
}

// Render returns the canonical text representation
func (x XText) Render() string { return x.Native() }

// Format returns the pretty text representation
func (x XText) Format(env envs.Environment) string {
	return x.Render()
}

// String returns the native string representation of this type for debugging
func (x XText) String() string { return `XText("` + x.Native() + `")` }

// Native returns the native value of this type
func (x XText) Native() string { return x.native }

// Equals determines equality for this type
func (x XText) Equals(other XText) bool {
	return x.Native() == other.Native()
}

// Compare compares this string to another
func (x XText) Compare(other XText) int {
	return strings.Compare(x.Native(), other.Native())
}

// Slice returns a substring of this text
func (x XText) Slice(start, end int) XText {
	length := x.Length()
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start >= length || end < start {
		return XTextEmpty
	}

	runes := []rune(x.native)[start:end]
	return NewXText(string(runes))
}

// Length returns the length of this string
func (x XText) Length() int { return utf8.RuneCountInString(x.Native()) }

// Empty returns whether this is an empty string
func (x XText) Empty() bool { return x.Native() == "" }

// MarshalJSON is called when a struct containing this type is marshaled
func (x XText) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(x.Native())
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XText) UnmarshalJSON(data []byte) error {
	return jsonx.Unmarshal(data, &x.native)
}

// XTextEmpty is the empty text value
var XTextEmpty = NewXText("")
var _ XValue = XTextEmpty

// ToXText converts the given value to a string
func ToXText(env envs.Environment, x XValue) (XText, XError) {
	if utils.IsNil(x) {
		return XTextEmpty, nil
	}
	if IsXError(x) {
		return XTextEmpty, x.(XError)
	}

	return NewXText(x.Render()), nil
}
