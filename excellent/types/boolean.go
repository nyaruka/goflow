package types

import (
	"encoding/json"
	"strconv"

	"github.com/nyaruka/goflow/utils"
)

// XBoolean is a boolean true or false
type XBoolean struct {
	native bool
}

// NewXBoolean creates a new boolean value
func NewXBoolean(value bool) XBoolean {
	return XBoolean{native: value}
}

// Describe returns a representation of this type for error messages
func (x XBoolean) Describe(env utils.Environment) string { return strconv.FormatBool(x.Native()) }

// ToXText converts this type to text
func (x XBoolean) ToXText(env utils.Environment) XText {
	return NewXText(strconv.FormatBool(x.Native()))
}

// ToXBoolean converts this type to a bool
func (x XBoolean) ToXBoolean(env utils.Environment) XBoolean { return x }

// ToXJSON is called when this type is passed to @(json(...))
func (x XBoolean) ToXJSON(env utils.Environment) XText { return MustMarshalToXText(x.Native()) }

// Native returns the native value of this type
func (x XBoolean) Native() bool { return x.native }

// String returns the native string representation of this type for debugging
func (x XBoolean) String() string { return `XBoolean(` + strconv.FormatBool(x.Native()) + `)` }

// Equals determines equality for this type
func (x XBoolean) Equals(other XBoolean) bool {
	return x.Native() == other.Native()
}

// Compare compares this bool to another
func (x XBoolean) Compare(other XBoolean) int {
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
func (x XBoolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Native())
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XBoolean) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &x.native)
}

// XBooleanFalse is the false boolean value
var XBooleanFalse = NewXBoolean(false)

// XBooleanTrue is the true boolean value
var XBooleanTrue = NewXBoolean(true)

var _ XValue = XBooleanFalse

// ToXBoolean converts the given value to a boolean
func ToXBoolean(env utils.Environment, x XValue) (XBoolean, XError) {
	if utils.IsNil(x) {
		return XBooleanFalse, nil
	}
	if IsXError(x) {
		return XBooleanFalse, x.(XError)
	}

	return x.ToXBoolean(env), nil
}
