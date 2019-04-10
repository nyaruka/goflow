package types

import (
	"fmt"

	"github.com/nyaruka/goflow/utils"
)

// XTime is a time of day value
type XTime struct {
	native utils.TimeOfDay
}

// NewXTime creates a new time
func NewXTime(value utils.TimeOfDay) XTime {
	return XTime{native: value}
}

// Describe returns a representation of this type for error messages
func (x XTime) Describe(env utils.Environment) string { return "time" }

// ToXText converts this type to text
func (x XTime) ToXText(env utils.Environment) XText { return NewXText(x.Native().String()) }

// ToXBoolean converts this type to a bool
func (x XTime) ToXBoolean(env utils.Environment) XBoolean {
	return NewXBoolean(x != XTimeZero)
}

// ToXJSON is called when this type is passed to @(json(...))
func (x XTime) ToXJSON(env utils.Environment) XText {
	return MustMarshalToXText(x.Native().String())
}

// Native returns the native value of this type
func (x XTime) Native() utils.TimeOfDay { return x.native }

// String returns the native string representation of this type
func (x XTime) String() string {
	return fmt.Sprintf(`XTime(%d, %d, %d, %d)`, x.native.Hour, x.native.Minute, x.native.Second, x.native.Nanos)
}

// Equals determines equality for this type
func (x XTime) Equals(other XTime) bool {
	return x.Native().Equal(other.Native())
}

// Compare compares this date to another
func (x XTime) Compare(other XTime) int {
	return x.Native().Compare(other.Native())
}

// XTimeZero is the zero time value
var XTimeZero = NewXTime(utils.ZeroTimeOfDay)
var _ XValue = XTimeZero

// ToXTime converts the given value to a time or returns an error if that isn't possible
func ToXTime(env utils.Environment, x XValue) (XTime, XError) {
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
				return NewXTime(utils.NewTimeOfDay(int(asInt), 0, 0, 0)), nil
			} else if asInt == 24 {
				return XTimeZero, nil
			}
		case XText:
			parsed, err := utils.TimeFromString(typed.Native())
			if err == nil {
				return NewXTime(parsed), nil
			}
		}
	}

	return XTimeZero, NewXErrorf("unable to convert %s to a time", Describe(env, x))
}
