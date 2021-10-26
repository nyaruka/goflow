package types

import (
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

// XDateTime is a datetime value.
//
//   @(datetime("1979-07-18T10:30:45.123456Z")) -> 1979-07-18T10:30:45.123456Z
//   @(format_datetime(datetime("1979-07-18T10:30:45.123456Z"))) -> 18-07-1979 05:30
//   @(json(datetime("1979-07-18T10:30:45.123456Z"))) -> "1979-07-18T10:30:45.123456Z"
//
// @type datetime
type XDateTime struct {
	native time.Time
}

// NewXDateTime creates a new date
func NewXDateTime(value time.Time) XDateTime {
	return XDateTime{native: value}
}

// Describe returns a representation of this type for error messages
func (x XDateTime) Describe() string { return "datetime" }

// Truthy determines truthiness for this type
func (x XDateTime) Truthy() bool {
	return !x.Native().IsZero()
}

// Render returns the canonical text representation
func (x XDateTime) Render() string {
	return dates.FormatISO(x.Native())
}

// Format returns the pretty text representation
func (x XDateTime) Format(env envs.Environment) string {
	formatted, _ := x.FormatCustom(env, string(env.DateFormat())+" "+string(env.TimeFormat()), env.Timezone())
	return formatted
}

// FormatCustom provides customised formatting
func (x XDateTime) FormatCustom(env envs.Environment, layout string, tz *time.Location) (string, error) {
	// convert to our timezone if we have one (otherwise we remain in the date's default)
	dt := x.Native()
	if tz != nil {
		dt = dt.In(tz)
	}

	return dates.Format(dt, layout, env.DefaultLocale().ToBCP47(), dates.DateTimeLayouts)
}

// String returns the native string representation of this type
func (x XDateTime) String() string {
	return fmt.Sprintf(`XDateTime(`+x.native.Format("2006, 1, 2, 15, 4, 5, %d, MST")+`)`, x.native.Nanosecond())
}

// Native returns the native value of this type
func (x XDateTime) Native() time.Time { return x.native }

// Date returns the date part of this datetime
func (x XDateTime) Date() XDate {
	return NewXDate(dates.ExtractDate(x.Native()))
}

// Time returns the time part of this datetime
func (x XDateTime) Time() XTime {
	return NewXTime(dates.ExtractTimeOfDay(x.Native()))
}

// In returns a copy of this datetime in a different timezone
func (x XDateTime) In(tz *time.Location) XDateTime {
	return NewXDateTime(x.Native().In(tz))
}

// ReplaceTime returns the a new date time with the time part replaced by the given time
func (x XDateTime) ReplaceTime(tm XTime) XDateTime {
	d := x.Native()
	t := tm.Native()
	return NewXDateTime(time.Date(d.Year(), d.Month(), d.Day(), t.Hour, t.Minute, t.Second, t.Nanos, d.Location()))
}

// Equals determines equality for this type
func (x XDateTime) Equals(o XValue) bool {
	other := o.(XDateTime)

	return x.Native().Equal(other.Native())
}

// Compare compares this date to another
func (x XDateTime) Compare(o XValue) int {
	other := o.(XDateTime)

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
	return jsonx.Marshal(dates.FormatISO(x.Native()))
}

// UnmarshalJSON is called when a struct containing this type is unmarshaled
func (x *XDateTime) UnmarshalJSON(data []byte) error {
	nativePtr := &x.native
	return nativePtr.UnmarshalJSON(data)
}

// XDateTimeZero is the zero time value
var XDateTimeZero = NewXDateTime(envs.ZeroDateTime)
var _ XValue = XDateTimeZero

// ToXDateTime converts the given value to a time or returns an error if that isn't possible
func ToXDateTime(env envs.Environment, x XValue) (XDateTime, XError) {
	return toXDateTime(env, x, false)
}

// ToXDateTimeWithTimeFill converts the given value to a time or returns an error if that isn't possible
func ToXDateTimeWithTimeFill(env envs.Environment, x XValue) (XDateTime, XError) {
	return toXDateTime(env, x, true)
}

// converts the given value to a time or returns an error if that isn't possible
func toXDateTime(env envs.Environment, x XValue, fillTime bool) (XDateTime, XError) {
	if !utils.IsNil(x) {
		switch typed := x.(type) {
		case XError:
			return XDateTimeZero, typed
		case XDate:
			return NewXDateTime(typed.Native().Combine(dates.ZeroTimeOfDay, env.Timezone())), nil
		case XDateTime:
			return typed, nil
		case XText:
			parsed, err := envs.DateTimeFromString(env, typed.Native(), fillTime)
			if err == nil {
				return NewXDateTime(parsed), nil
			}
		case *XObject:
			if typed.hasDefault() {
				return toXDateTime(env, typed.Default(), fillTime)
			}
		}
	}

	return XDateTimeZero, NewXErrorf("unable to convert %s to a datetime", Describe(x))
}
