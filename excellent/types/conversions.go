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
func ToXJSON(x XValue) (XText, XError) {
	if utils.IsNil(x) {
		return NewXText(`null`), nil
	}
	if IsXError(x) {
		return XTextEmpty, x.(XError)
	}

	return x.ToXJSON(), nil
}

// ToXText converts the given value to a string
func ToXText(x XValue) (XText, XError) {
	if utils.IsNil(x) {
		return XTextEmpty, nil
	}
	if IsXError(x) {
		return XTextEmpty, x.(XError)
	}

	return x.Reduce().ToXText(), nil
}

// ToXBool converts the given value to a boolean
func ToXBool(x XValue) (XBoolean, XError) {
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

// ToXNumber converts the given value to a number or returns an error if that isn't possible
func ToXNumber(x XValue) (XNumber, XError) {
	if utils.IsNil(x) {
		return XNumberZero, nil
	}

	x = x.Reduce()

	switch typed := x.(type) {
	case XError:
		return XNumberZero, typed
	case XNumber:
		return typed, nil
	case XText:
		parsed, err := parseDecimalFuzzy(typed.Native())
		if err == nil {
			return NewXNumber(parsed), nil
		}
	}

	return XNumberZero, NewXErrorf("unable to convert value '%s' to a number", x)
}

// ToXDate converts the given value to a time or returns an error if that isn't possible
func ToXDate(env utils.Environment, x XValue) (XDateTime, XError) {
	if utils.IsNil(x) {
		return XDateTimeZero, nil
	}

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

	return XDateTimeZero, NewXErrorf("unable to convert value '%s' of type '%s' to a date", x, reflect.TypeOf(x))
}

// ToInteger tries to convert the passed in value to an integer or returns an error if that isn't possible
func ToInteger(x XValue) (int, XError) {
	number, err := ToXNumber(x)
	if err != nil {
		return 0, err
	}

	intPart := number.Native().IntPart()

	if intPart < math.MinInt32 && intPart > math.MaxInt32 {
		return 0, NewXErrorf("number value %s is out of range for an integer", number.ToXText().Native())
	}

	return int(intPart), nil
}

// MustMarshalToXText calls json.Marshal in the given value and panics in the case of an error
func MustMarshalToXText(x interface{}) XText {
	j, err := json.Marshal(x)
	if err != nil {
		panic(fmt.Sprintf("unable to marshal %s to JSON", x))
	}
	return NewXText(string(j))
}

func parseDecimalFuzzy(val string) (decimal.Decimal, error) {
	// common SMS foibles
	val = strings.ToLower(val)
	val = strings.Replace(val, "o", "0", -1)
	val = strings.Replace(val, "l", "1", -1)
	return decimal.NewFromString(val)
}
