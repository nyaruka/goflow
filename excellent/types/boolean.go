package types

import (
	"strconv"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

// XBoolean is a boolean `true` or `false`.
//
//   @(true) -> true
//   @(1 = 1) -> true
//   @(1 = 2) -> false
//   @(json(true)) -> true
//
// @type boolean
type XBoolean struct {
	native bool
}

// NewXBoolean creates a new boolean value
func NewXBoolean(value bool) XBoolean {
	return XBoolean{native: value}
}

// Describe returns a representation of this type for error messages
func (x XBoolean) Describe() string { return strconv.FormatBool(x.Native()) }

// Truthy determines truthiness for this type
func (x XBoolean) Truthy() bool { return x.Native() }

// Render returns the canonical text representation
func (x XBoolean) Render() string {
	return strconv.FormatBool(x.Native())
}

// Format returns the pretty text representation
func (x XBoolean) Format(env envs.Environment) string {
	return x.Render()
}

// MarshalJSON is called when a struct containing this type is marshaled
func (x XBoolean) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(x.Native())
}

// String returns the native string representation of this type for debugging
func (x XBoolean) String() string { return `XBoolean(` + strconv.FormatBool(x.Native()) + `)` }

// Native returns the native value of this type
func (x XBoolean) Native() bool { return x.native }

// Equals determines equality for this type
func (x XBoolean) Equals(o XValue) bool {
	other := o.(XBoolean)

	return x.Native() == other.Native()
}

// Compare compares this bool to another
func (x XBoolean) Compare(o XValue) int {
	other := o.(XBoolean)

	switch {
	case !x.Native() && other.Native():
		return -1
	case x.Native() && !other.Native():
		return 1
	default:
		return 0
	}
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XBoolean) UnmarshalJSON(data []byte) error {
	return jsonx.Unmarshal(data, &x.native)
}

// XBooleanFalse is the false boolean value
var XBooleanFalse = NewXBoolean(false)

// XBooleanTrue is the true boolean value
var XBooleanTrue = NewXBoolean(true)

var _ XValue = XBooleanFalse

// ToXBoolean converts the given value to a boolean
func ToXBoolean(x XValue) (XBoolean, XError) {
	if utils.IsNil(x) {
		return XBooleanFalse, nil
	}
	if IsXError(x) {
		return XBooleanFalse, x.(XError)
	}

	return NewXBoolean(x.Truthy()), nil
}
