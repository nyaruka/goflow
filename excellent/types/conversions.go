package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

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
		parsed, err := parseDecimalFuzzy(string(typed))
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
		parsed, err := utils.DateFromString(env, string(typed))
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

// Legacy...

// ToJSON tries to turn the passed in interface to a JSON fragment
func ToJSON(env utils.Environment, val interface{}) (JSONFragment, error) {
	ToFragment := func(bytes []byte, err error) (JSONFragment, error) {
		if bytes == nil {
			return EmptyJSONFragment, err
		}
		return JSONFragment(bytes), err
	}

	// null is null
	if utils.IsNil(val) {
		return ToFragment(json.Marshal(nil))
	}

	switch val := val.(type) {

	case error:
		return EmptyJSONFragment, val

	case string:
		return ToFragment(json.Marshal(val))

	case bool:
		return ToFragment(json.Marshal(val))

	case int:
		return ToFragment(json.Marshal(val))

	case decimal.Decimal:
		floatVal, _ := val.Float64()
		return ToFragment(json.Marshal(floatVal))

	case time.Time:
		return ToFragment(json.Marshal(utils.DateToISO(val)))

	case JSONFragment:
		return val, nil

	case Array:
		return ToFragment(json.Marshal(val))

	case Atomizable:
		return ToJSON(env, val.Atomize())
	}

	// welp, we give up, this isn't something we can convert, return an error
	return EmptyJSONFragment, fmt.Errorf("ToJSON unknown type '%s' with value '%+v'", reflect.TypeOf(val), val)
}

// ToString tries to turn the passed in interface to a string
func ToString(env utils.Environment, val interface{}) (string, error) {
	// Strings are always defined, just empty
	if utils.IsNil(val) {
		return "", nil
	}

	switch val := val.(type) {
	case error:
		return "", val

	case string:
		return val, nil

	case bool:
		return strconv.FormatBool(val), nil

	case decimal.Decimal:
		return val.String(), nil
	case int:
		return strconv.FormatInt(int64(val), 10), nil

	case time.Time:
		return utils.DateToISO(val), nil

	case Atomizable:
		return ToString(env, val.Atomize())

	case Array:
		var output bytes.Buffer
		for i := 0; i < val.Length(); i++ {
			if i > 0 {
				output.WriteString(", ")
			}
			itemAsStr, err := ToString(env, val.Index(i))
			if err != nil {
				return "", err
			}
			output.WriteString(itemAsStr)
		}
		return output.String(), nil
	}

	// welp, we give up, this isn't something we can convert, return an error
	return "", fmt.Errorf("unable to convert value '%+v' of type '%s' to a string", val, reflect.TypeOf(val))
}

// ToInt tries to convert the passed in interface{} to an integer value, returning an error if that isn't possible
func ToInt(env utils.Environment, val interface{}) (int, error) {
	dec, err := ToDecimal(env, val)
	if err != nil {
		return 0, err
	}

	if dec.IntPart() < math.MinInt32 && dec.IntPart() > math.MaxInt32 {
		return 0, fmt.Errorf("Decimal value %d is out of range for an integer", dec.IntPart())
	}

	return int(dec.IntPart()), nil
}

// ToDecimal tries to convert the passed in interface{} to a Decimal value, returning an error if that isn't possible
func ToDecimal(env utils.Environment, val interface{}) (decimal.Decimal, error) {
	if utils.IsNil(val) {
		return decimal.Zero, nil
	}

	switch val := val.(type) {
	case error:
		return decimal.Zero, val

	case string:
		// common SMS foibles
		subbed := strings.ToLower(val)
		subbed = strings.Replace(subbed, "o", "0", -1)
		subbed = strings.Replace(subbed, "l", "1", -1)
		parsed, err := decimal.NewFromString(subbed)
		if err != nil {
			return decimal.Zero, fmt.Errorf("Cannot convert '%s' to a decimal", val)
		}
		return parsed, nil

	case decimal.Decimal:
		return val, nil
	case int:
		return decimal.NewFromString(strconv.FormatInt(int64(val), 10))

	case Atomizable:
		return ToDecimal(env, val.Atomize())
	}

	return decimal.Zero, fmt.Errorf("unable to convert value '%+v' of type '%s' to a decimal", val, reflect.TypeOf(val))
}

// ToDate tries to convert the passed in interface to a time.Time returning an error if that isn't possible
func ToDate(env utils.Environment, val interface{}) (time.Time, error) {
	if utils.IsNil(val) {
		return time.Time{}, fmt.Errorf("Cannot convert nil to date")
	}

	switch val := val.(type) {
	case error:
		return time.Time{}, val

	case string:
		return utils.DateFromString(env, val)

	case decimal.Decimal:
		return time.Time{}, fmt.Errorf("Cannot convert decimal to date")
	case int:
		return time.Time{}, fmt.Errorf("Cannot convert integer to date")

	case time.Time:
		return val, nil

	case Atomizable:
		return ToDate(env, val.Atomize())
	}

	return time.Time{}, fmt.Errorf("unable to convert value '%+v' of type '%s' to a date", val, reflect.TypeOf(val))
}

// ToBool tests whether the passed in item should be considered True
// false, 0, "" and nil are false, everything else is true
func ToBool(env utils.Environment, test interface{}) (bool, error) {
	if test == nil {
		return false, nil
	}

	switch test := test.(type) {
	case error:
		return false, test

	case string:
		return test != "" && strings.ToLower(test) != "false", nil

	case decimal.Decimal:
		return !test.Equals(decimal.Zero), nil
	case int:
		return test != 0, nil

	case time.Time:
		return !test.IsZero(), nil

	case bool:
		return test, nil

	case JSONFragment:
		asString, err := ToString(env, test)
		if err != nil {
			return false, err
		}
		// is this a number?
		num, err := ToDecimal(env, asString)
		if err == nil {
			return !num.Equals(decimal.Zero), nil
		}

		noWhite := strings.Join(strings.Fields(asString), "")

		// empty array?
		if noWhite == "[]" {
			return false, nil
		}

		// empty dict
		if noWhite == "{}" {
			return false, nil
		}

		// finally just string version
		return asString != "" && strings.ToLower(asString) != "false", nil

	case Atomizable:
		return ToBool(env, test.Atomize())
	}

	return false, fmt.Errorf("unable to convert value '%+v' of type '%s' to a bool", test, reflect.TypeOf(test))
}

// Compare returns the difference between the two passed interfaces which must be of the same type
func Compare(env utils.Environment, arg1 interface{}, arg2 interface{}) (int, error) {
	if arg1 == nil && arg2 == nil {
		return 0, nil
	}

	arg1, arg1Type, _ := ToXAtom(env, arg1)
	arg2, arg2Type, _ := ToXAtom(env, arg2)

	// common types, do real comparisons
	switch {
	case arg1Type == arg2Type && arg1Type == XTypeError:
		return strings.Compare(arg1.(error).Error(), arg2.(error).Error()), nil

	case arg1Type == arg2Type && arg1Type == XTypeDecimal:
		return arg1.(decimal.Decimal).Cmp(arg2.(decimal.Decimal)), nil

	case arg1Type == arg2Type && arg1Type == XTypeBool:
		bool1 := arg1.(bool)
		bool2 := arg2.(bool)

		switch {
		case !bool1 && bool2:
			return -1, nil
		case bool1 == bool2:
			return 0, nil
		case bool1 && !bool2:
			return 1, nil
		}

	case arg1Type == arg2Type && arg1Type == XTypeTime:
		time1 := arg1.(time.Time)
		time2 := arg2.(time.Time)

		switch {
		case time1.Before(time2):
			return -1, nil
		case time1.Equal(time2):
			return 0, nil
		case time1.After(time2):
			return 1, nil
		}

	case arg1Type == arg2Type && arg1Type == XTypeString:
		return strings.Compare(arg1.(string), arg2.(string)), nil
	}

	if arg1Type != arg2Type {
		return 0, fmt.Errorf("Cannot compare different types of %#v and %#v", arg1, arg2)
	}

	arg1Str := fmt.Sprintf("%v", arg1)
	arg2Str := fmt.Sprintf("%v", arg2)

	cmp := strings.Compare(arg1Str, arg2Str)
	return cmp, nil
}
