package types

import (
	"time"

	"github.com/nyaruka/goflow/utils"
)

// XDateTime is a datetime value
type XDateTime struct {
	native time.Time
}

// NewXDateTime creates a new date
func NewXDateTime(value time.Time) XDateTime {
	return XDateTime{native: value}
}

// Repr returns the representation of this type
func (x XDateTime) Repr() string { return "datetime" }

// Reduce returns the primitive version of this type (i.e. itself)
func (x XDateTime) Reduce() XPrimitive { return x }

// ToXText converts this type to text
func (x XDateTime) ToXText() XText { return NewXText(utils.DateToISO(x.Native())) }

// ToXBoolean converts this type to a bool
func (x XDateTime) ToXBoolean() XBoolean { return NewXBoolean(!x.Native().IsZero()) }

// ToXJSON is called when this type is passed to @(json(...))
func (x XDateTime) ToXJSON() XText { return MustMarshalToXText(utils.DateToISO(x.Native())) }

// Native returns the native value of this type
func (x XDateTime) Native() time.Time { return x.native }

// String returns the native string representation of this type
func (x XDateTime) String() string { return x.ToXText().Native() }

// Equals determines equality for this type
func (x XDateTime) Equals(other XDateTime) bool {
	return x.Native().Equal(other.Native())
}

// Compare compares this date to another
func (x XDateTime) Compare(other XDateTime) int {
	switch {
	case x.Native().Before(other.Native()):
		return -1
	case x.Native().After(other.Native()):
		return 1
	default:
		return 0
	}
}

// MarshalJSON is called when a struct containing this type is marshaled
func (x XDateTime) MarshalJSON() ([]byte, error) {
	return x.Native().MarshalJSON()
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XDateTime) UnmarshalJSON(data []byte) error {
	nativePtr := &x.native
	return nativePtr.UnmarshalJSON(data)
}

// XDateTimeZero is the zero time value
var XDateTimeZero = NewXDateTime(time.Time{})
var _ XPrimitive = XDateTimeZero

// ToXDateTime converts the given value to a time or returns an error if that isn't possible
func ToXDateTime(env utils.Environment, x XValue) (XDateTime, XError) {
	if !utils.IsNil(x) {
		x = x.Reduce()

		switch typed := x.(type) {
		case XError:
			return XDateTimeZero, typed
		case XDateTime:
			return typed, nil
		case XText:
			parsed, err := utils.DateFromString(env, typed.Native())
			if err == nil {
				return NewXDateTime(parsed), nil
			}
		}
	}

	return XDateTimeZero, NewXErrorf("unable to convert %s to a datetime", Repr(x))
}
