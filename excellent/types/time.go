package types

import (
	"fmt"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

// XTime is a time of day.
//
//   @(time_from_parts(16, 30, 45)) -> 16:30:45.000000
//   @(format_time(time_from_parts(16, 30, 45))) -> 16:30
//   @(json(time_from_parts(16, 30, 45))) -> "16:30:45.000000"
//
// @type time
type XTime struct {
	native dates.TimeOfDay
}

// NewXTime creates a new time
func NewXTime(value dates.TimeOfDay) XTime {
	return XTime{native: value}
}

// Describe returns a representation of this type for error messages
func (x XTime) Describe() string { return "time" }

// Truthy determines truthiness for this type
func (x XTime) Truthy() bool {
	return x != XTimeZero
}

// Render returns the canonical text representation
func (x XTime) Render() string { return x.Native().String() }

// Format returns the pretty text representation
func (x XTime) Format(env envs.Environment) string {
	formatted, _ := x.FormatCustom(env, string(env.TimeFormat()))
	return formatted
}

// FormatCustom provides customised formatting
func (x XTime) FormatCustom(env envs.Environment, layout string) (string, error) {
	return x.Native().Format(layout, env.DefaultLocale().ToBCP47())
}

// MarshalJSON is called when a struct containing this type is marshaled
func (x XTime) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(x.Native().String())
}

// String returns the native string representation of this type
func (x XTime) String() string {
	return fmt.Sprintf(`XTime(%d, %d, %d, %d)`, x.native.Hour, x.native.Minute, x.native.Second, x.native.Nanos)
}

// Native returns the native value of this type
func (x XTime) Native() dates.TimeOfDay { return x.native }

// Equals determines equality for this type
func (x XTime) Equals(o XValue) bool {
	other := o.(XTime)

	return x.Native().Equal(other.Native())
}

// Compare compares this date to another
func (x XTime) Compare(o XValue) int {
	other := o.(XTime)

	c := x.Native().Compare(other.Native())
	if c > 0 {
		return 1
	} else if c < 0 {
		return -1
	}
	return 0
}

// XTimeZero is the zero time value
var XTimeZero = NewXTime(dates.ZeroTimeOfDay)
var _ XValue = XTimeZero

// ToXTime converts the given value to a time or returns an error if that isn't possible
func ToXTime(env envs.Environment, x XValue) (XTime, XError) {
	if !utils.IsNil(x) {
		switch typed := x.(type) {
		case XError:
			return XTimeZero, typed
		case XTime:
			return typed, nil
		case XDateTime:
			return typed.Time(), nil
		case XNumber:
			asInt := typed.Native().IntPart()
			if asInt >= 0 && asInt <= 23 {
				return NewXTime(dates.NewTimeOfDay(int(asInt), 0, 0, 0)), nil
			} else if asInt == 24 {
				return XTimeZero, nil
			}
		case XText:
			parsed, err := envs.TimeFromString(typed.Native())
			if err == nil {
				return NewXTime(parsed), nil
			}
		case *XObject:
			if typed.hasDefault() {
				return ToXTime(env, typed.Default())
			}
		}
	}

	return XTimeZero, NewXErrorf("unable to convert %s to a time", Describe(x))
}
