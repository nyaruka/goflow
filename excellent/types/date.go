package types

import (
	"encoding/json"
	"fmt"

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
	native utils.Date
}

// NewXDate creates a new date
func NewXDate(value utils.Date) XDate {
	return XDate{native: value}
}

// Describe returns a representation of this type for error messages
func (x XDate) Describe() string { return "date" }

// ToXText converts this type to text
func (x XDate) ToXText(env utils.Environment) XText { return NewXText(x.Native().String()) }

// ToXBoolean converts this type to a bool
func (x XDate) ToXBoolean() XBoolean {
	return NewXBoolean(x != XDateZero)
}

// Native returns the native value of this type
func (x XDate) Native() utils.Date { return x.native }

// MarshalJSON is called when a struct containing this type is marshaled
func (x XDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Native().String())
}

// String returns the native string representation of this type
func (x XDate) String() string {
	return fmt.Sprintf(`XDate(%d, %d, %d)`, x.native.Year, x.native.Month, x.native.Day)
}

// Equals determines equality for this type
func (x XDate) Equals(other XDate) bool {
	return x.Native().Equal(other.Native())
}

// Compare compares this date to another
func (x XDate) Compare(other XDate) int {
	return x.Native().Compare(other.Native())
}

// XDateZero is the zero time value
var XDateZero = NewXDate(utils.ZeroDate)
var _ XValue = XDateZero

// ToXDate converts the given value to a time or returns an error if that isn't possible
func ToXDate(env utils.Environment, x XValue) (XDate, XError) {
	if !utils.IsNil(x) {
		switch typed := x.(type) {
		case XError:
			return XDateZero, typed
		case XDate:
			return typed, nil
		case XDateTime:
			return typed.In(env.Timezone()).Date(), nil
		case XText:
			parsed, err := utils.DateFromString(env, typed.Native())
			if err == nil {
				return NewXDate(parsed), nil
			}
		}
	}

	return XDateZero, NewXErrorf("unable to convert %s to a date", Describe(x))
}
