package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// IsMap returns whether the given object is a map
func IsMap(v interface{}) bool {
	val := reflect.ValueOf(v)
	return val.Kind() == reflect.Map
}

// IsNil returns whether the given object is nil or an interface to a nil
func IsNil(v interface{}) bool {
	// if v doesn't have a type or value then v == nil
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)

	// if v is a typed nil pointer then v != nil but the value is nil
	if val.Kind() == reflect.Ptr {
		return val.IsNil()
	}

	return false
}

// Tries to use golang reflection to lookup the key in the passed in map value
func attemptMapLookup(valMap reflect.Value, key interface{}) (value interface{}, err error) {
	defer func() {
		if recover() != nil {
			value = nil
			err = fmt.Errorf("Invalid key type for map: %v", key)
		}
	}()

	mapValue := valMap.MapIndex(reflect.ValueOf(key))
	if !mapValue.IsValid() {
		return nil, nil
	}

	value = mapValue.Interface()
	return value, nil
}

// LookupKey tries to look up the interface at the passed in key for the passed in map
func LookupKey(v interface{}, key string) (value interface{}, err error) {
	if v == nil {
		return nil, fmt.Errorf("Cannot convert nil to interface map")
	}

	// deal with a passed in error, we just return it out
	err, isErr := v.(error)
	if isErr {
		return nil, err
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Map {
		return nil, fmt.Errorf("Cannot convert non-map to map: %v", v)
	}

	// try to look up the value with our key as a string
	result, err := attemptMapLookup(val, key)
	if err == nil {
		return result, err
	}

	// no luck, try to convert to an integer instead
	intKey, intErr := strconv.Atoi(key)
	if intErr != nil {
		return result, err
	}

	return attemptMapLookup(val, intKey)
}

// ToStringArray tries to turn the passed in interface (which must be an underlying slice) to a string array
func ToStringArray(env Environment, v interface{}) ([]string, error) {
	if IsNil(v) {
		return nil, fmt.Errorf("Cannot convert nil to string array")
	}

	// deal with a passed in error, we just return it out
	err, isErr := v.(error)
	if isErr {
		return nil, err
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Cannot convert non-array to string array")
	}

	strArr := make([]string, val.Len())
	for i := range strArr {
		str, err := ToString(env, val.Index(i).Interface())
		if err != nil {
			return nil, err
		}

		strArr[i] = str
	}

	return strArr, nil
}

// ToJSON tries to turn the passed in interface to a JSON fragment
func ToJSON(env Environment, val interface{}) (JSONFragment, error) {
	ToFragment := func(bytes []byte, err error) (JSONFragment, error) {
		if bytes == nil {
			return EmptyJSONFragment, err
		}
		return JSONFragment(bytes), err
	}

	// null is null
	if IsNil(val) {
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
		return ToFragment(json.Marshal(DateToISO(val)))

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
func ToString(env Environment, val interface{}) (string, error) {
	// Strings are always defined, just empty
	if IsNil(val) {
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
		return DateToISO(val), nil

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
func ToInt(env Environment, val interface{}) (int, error) {
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
func ToDecimal(env Environment, val interface{}) (decimal.Decimal, error) {
	if IsNil(val) {
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
func ToDate(env Environment, val interface{}) (time.Time, error) {
	if IsNil(val) {
		return time.Time{}, fmt.Errorf("Cannot convert nil to date")
	}

	switch val := val.(type) {
	case error:
		return time.Time{}, val

	case string:
		return DateFromString(env, val)

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
func ToBool(env Environment, test interface{}) (bool, error) {
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

// XType is an an enumeration of the possible types we can deal with
type XType int

// primitive types we convert to
const (
	XTypeNil = iota
	XTypeError
	XTypeString
	XTypeDecimal
	XTypeTime
	XTypeBool
	XTypeArray
)

// ToXAtom figures out the raw type of the passed in interface, returning that type
func ToXAtom(env Environment, val interface{}) (interface{}, XType, error) {
	if val == nil {
		return val, XTypeNil, nil
	}

	switch val := val.(type) {
	case error:
		return val, XTypeError, nil

	case string:
		return val, XTypeString, nil

	case decimal.Decimal:
		return val, XTypeDecimal, nil
	case int:
		decVal, err := ToDecimal(env, val)
		if err != nil {
			return val, XTypeNil, err
		}
		return decVal, XTypeDecimal, nil

	case time.Time:
		return val, XTypeTime, nil

	case bool:
		return val, XTypeBool, nil

	case Array:
		return val, XTypeArray, nil
	}

	return val, XTypeNil, fmt.Errorf("Unknown type '%s' with value '%+v'", reflect.TypeOf(val), val)
}

// Compare returns the difference between the two passed interfaces which must be of the same type
func Compare(env Environment, arg1 interface{}, arg2 interface{}) (int, error) {
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
