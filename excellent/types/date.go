package types

import (
	"github.com/nyaruka/goflow/utils"
)

// XDate is a date value
type XDate struct {
	native utils.Date
}

// NewXDate creates a new date
func NewXDate(value utils.Date) XDate {
	return XDate{native: value}
}

// Describe returns a representation of this type for error messages
func (x XDate) Describe(env utils.Environment) string { return "date" }

// ToXText converts this type to text
func (x XDate) ToXText(env utils.Environment) XText { return NewXText(x.Native().String()) }

// ToXBoolean converts this type to a bool
func (x XDate) ToXBoolean(env utils.Environment) XBoolean {
	return NewXBoolean(x != XDateZero)
}

// ToXJSON is called when this type is passed to @(json(...))
func (x XDate) ToXJSON(env utils.Environment) XText {
	return MustMarshalToXText(x.Native().String())
}

// Native returns the native value of this type
func (x XDate) Native() utils.Date { return x.native }

// String returns the native string representation of this type
func (x XDate) String() string { return x.ToXText(nil).Native() }

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

	return XDateZero, NewXErrorf("unable to convert %s to a date", Describe(env, x))
}
