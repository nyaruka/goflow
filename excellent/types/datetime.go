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

// Describe returns a representation of this type for error messages
func (x XDateTime) Describe() string { return "datetime" }

// Reduce returns the primitive version of this type (i.e. itself)
func (x XDateTime) Reduce(env utils.Environment) XPrimitive { return x }

// ToXText converts this type to text
func (x XDateTime) ToXText(env utils.Environment) XText { return NewXText(utils.DateTimeToISO(x.Native())) }

// ToXBoolean converts this type to a bool
func (x XDateTime) ToXBoolean(env utils.Environment) XBoolean {
	return NewXBoolean(!x.Native().IsZero())
}

// ToXJSON is called when this type is passed to @(json(...))
func (x XDateTime) ToXJSON(env utils.Environment) XText {
	return MustMarshalToXText(utils.DateTimeToISO(x.Native()))
}

// Native returns the native value of this type
func (x XDateTime) Native() time.Time { return x.native }

// String returns the native string representation of this type
func (x XDateTime) String() string { return x.ToXText(nil).Native() }

// Date returns the date part of this datetime
func (x XDateTime) Date() XDate {
	return NewXDate(utils.ExtractDate(x.Native()))
}

// Time returns the time part of this datetime
func (x XDateTime) Time() XTime {
	return NewXTime(utils.ExtractTimeOfDay(x.Native()))
}

// ReplaceTime returns the a new date time with the time part replaced by the given time
func (x XDateTime) ReplaceTime(tm XTime) XDateTime {
	d := x.Native()
	t := tm.Native()
	return NewXDateTime(time.Date(d.Year(), d.Month(), d.Day(), t.Hour, t.Minute, t.Second, t.Nanos, d.Location()))
}

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
var XDateTimeZero = NewXDateTime(utils.ZeroDateTime)
var _ XPrimitive = XDateTimeZero

// ToXDateTime converts the given value to a time or returns an error if that isn't possible
func ToXDateTime(env utils.Environment, x XValue) (XDateTime, XError) {
	return toXDateTime(env, x, false)
}

// ToXDateTimeWithTimeFill converts the given value to a time or returns an error if that isn't possible
func ToXDateTimeWithTimeFill(env utils.Environment, x XValue) (XDateTime, XError) {
	return toXDateTime(env, x, true)
}

// converts the given value to a time or returns an error if that isn't possible
func toXDateTime(env utils.Environment, x XValue, fillTime bool) (XDateTime, XError) {
	if !utils.IsNil(x) {
		x = x.Reduce(env)

		switch typed := x.(type) {
		case XError:
			return XDateTimeZero, typed
		case XDateTime:
			return typed, nil
		case XText:
			parsed, err := utils.DateTimeFromString(env, typed.Native(), fillTime)
			if err == nil {
				return NewXDateTime(parsed), nil
			}
		}
	}

	return XDateTimeZero, NewXErrorf("unable to convert %s to a datetime", Describe(x))
}
