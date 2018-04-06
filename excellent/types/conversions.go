package types

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

// ToXJSON converts the given value to a JSON string
func ToXJSON(x XValue) XString {
	if utils.IsNil(x) {
		return NewXString(`null`)
	}

	return x.ToJSON()
}

// ToXString converts the given value to a string
func ToXString(x XValue) XString {
	if utils.IsNil(x) {
		return XStringEmpty
	}

	return x.Reduce().ToString()
}

// ToXBool converts the given value to a boolean
func ToXBool(x XValue) XBool {
	if utils.IsNil(x) {
		return XBoolFalse
	}

	primitive, isPrimitive := x.(XPrimitive)
	if isPrimitive {
		return primitive.ToBool()
	}

	lengthable, isLengthable := x.(XLengthable)
	if isLengthable {
		return lengthable.Length() > 0
	}

	return x.Reduce().ToBool()
}

// ToXNumber converts the given value to a number or returns an error if that isn't possible
func ToXNumber(x XValue) (XNumber, error) {
	if utils.IsNil(x) {
		return XNumberZero, nil
	}

	x = x.Reduce()

	switch typed := x.(type) {
	case XError:
		return XNumberZero, typed
	case XNumber:
		return typed, nil
	case XString:
		parsed, err := parseDecimalFuzzy(typed.Native())
		if err == nil {
			return NewXNumber(parsed), nil
		}
	}

	return XNumberZero, fmt.Errorf("unable to convert value '%s' to a number", x)
}

// ToXTime converts the given value to a time or returns an error if that isn't possible
func ToXTime(env utils.Environment, x XValue) (XTime, error) {
	if utils.IsNil(x) {
		return XTimeZero, nil
	}

	x = x.Reduce()

	switch typed := x.(type) {
	case XError:
		return XTimeZero, typed
	case XTime:
		return typed, nil
	case XString:
		parsed, err := utils.DateFromString(env, typed.Native())
		if err == nil {
			return NewXTime(parsed), nil
		}
	}

	return XTimeZero, fmt.Errorf("unable to convert value '%v' of type '%s' to a time", x, reflect.TypeOf(x))
}

// ToInteger tries to convert the passed in value to an integer or returns an error if that isn't possible
func ToInteger(x XValue) (int, error) {
	number, err := ToXNumber(x)
	if err != nil {
		return 0, err
	}

	intPart := number.Native().IntPart()

	if intPart < math.MinInt32 && intPart > math.MaxInt32 {
		return 0, fmt.Errorf("number value %s is out of range for an integer", string(number.ToString()))
	}

	return int(intPart), nil
}

// RequireMarshalToXString calls json.Marshal in the given value and panics in the case of an error
func RequireMarshalToXString(x interface{}) XString {
	j, err := json.Marshal(x)
	if err != nil {
		panic(fmt.Sprintf("unable to marshal %v to JSON", x))
	}
	return XString(j)
}

func parseDecimalFuzzy(val string) (decimal.Decimal, error) {
	// common SMS foibles
	val = strings.ToLower(val)
	val = strings.Replace(val, "o", "0", -1)
	val = strings.Replace(val, "l", "1", -1)
	return decimal.NewFromString(val)
}
