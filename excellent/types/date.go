package types

import (
	"fmt"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

// XDate is a Gregorian calendar date value.
//
//   @(date_from_parts(2019, 4, 11)) -> 2019-04-11
//   @(format_date(date_from_parts(2019, 4, 11))) -> 11-04-2019
//   @(json(date_from_parts(2019, 4, 11))) -> "2019-04-11"
//
// @type date
type XDate struct {
	native dates.Date
}

// NewXDate creates a new date
func NewXDate(value dates.Date) XDate {
	return XDate{native: value}
}

// Describe returns a representation of this type for error messages
func (x XDate) Describe() string { return "date" }

// Truthy determines truthiness for this type
func (x XDate) Truthy() bool {
	return x != XDateZero
}

// Render returns the canonical text representation
func (x XDate) Render() string { return x.Native().String() }

// Format returns the pretty text representation
func (x XDate) Format(env envs.Environment) string {
	formatted, _ := x.FormatCustom(env, string(env.DateFormat()))
	return formatted
}

// FormatCustom provides customised formatting
func (x XDate) FormatCustom(env envs.Environment, layout string) (string, error) {
	return x.Native().Format(layout, env.DefaultLocale().ToBCP47())
}

// MarshalJSON is called when a struct containing this type is marshaled
func (x XDate) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(x.Native().String())
}

// String returns the native string representation of this type
func (x XDate) String() string {
	return fmt.Sprintf(`XDate(%d, %d, %d)`, x.native.Year, x.native.Month, x.native.Day)
}

// Native returns the native value of this type
func (x XDate) Native() dates.Date { return x.native }

// Equals determines equality for this type
func (x XDate) Equals(other XDate) bool {
	return x.Native().Equal(other.Native())
}

// Compare compares this date to another
func (x XDate) Compare(other XDate) int {
	return x.Native().Compare(other.Native())
}

// XDateZero is the zero time value
var XDateZero = NewXDate(dates.ZeroDate)
var _ XValue = XDateZero

// ToXDate converts the given value to a time or returns an error if that isn't possible
func ToXDate(env envs.Environment, x XValue) (XDate, XError) {
	if !utils.IsNil(x) {
		switch typed := x.(type) {
		case XError:
			return XDateZero, typed
		case XDate:
			return typed, nil
		case XDateTime:
			return typed.In(env.Timezone()).Date(), nil
		case XText:
			parsed, err := envs.DateFromString(env, typed.Native())
			if err == nil {
				return NewXDate(parsed), nil
			}
		case *XObject:
			if typed.hasDefault() {
				return ToXDate(env, typed.Default())
			}
		}
	}

	return XDateZero, NewXErrorf("unable to convert %s to a date", Describe(x))
}
