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

// Reduce returns the primitive version of this type (i.e. itself)
func (x XBoolean) Reduce() XPrimitive { return x }

// ToXText converts this type to text
func (x XBoolean) ToXText() XText { return NewXText(strconv.FormatBool(x.Native())) }

// ToXBoolean converts this type to a bool
func (x XBoolean) ToXBoolean() XBoolean { return x }

// ToXJSON is called when this type is passed to @(json(...))
func (x XBoolean) ToXJSON() XText { return MustMarshalToXText(x.Native()) }

// Native returns the native value of this type
func (x XBoolean) Native() bool { return x.native }

// String returns the native string representation of this type
func (x XBoolean) String() string { return x.ToXText().Native() }

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

var _ XPrimitive = XBooleanFalse

// ToXBoolean converts the given value to a boolean
func ToXBoolean(x XValue) (XBoolean, XError) {
	if utils.IsNil(x) {
		return XBooleanFalse, nil
	}
	if IsXError(x) {
		return XBooleanFalse, x.(XError)
	}

	primitive, isPrimitive := x.(XPrimitive)
	if isPrimitive {
		return primitive.ToXBoolean(), nil
	}

	lengthable, isLengthable := x.(XLengthable)
	if isLengthable {
		return NewXBoolean(lengthable.Length() > 0), nil
	}

	return x.Reduce().ToXBoolean(), nil
}
