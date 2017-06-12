package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"math"

	"github.com/buger/jsonparser"
	"github.com/shopspring/decimal"
)

// IsSlice returns whether the passed in interface is a slice
func IsSlice(v interface{}) bool {
	val := reflect.ValueOf(v)
	return val.Kind() == reflect.Slice
}

// SliceLength returns the length of the passed in slice, it returns an error if the argument is not a slice
func SliceLength(v interface{}) (int, error) {
	if v == nil {
		return 0, fmt.Errorf("Cannot convert nil to slice")
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice {
		return val.Len(), nil
	}

	json, isJSON := v.(JSONFragment)
	if isJSON {
		count := 0
		_, err := jsonparser.ArrayEach(json.json, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			count++
		})
		if err == nil {
			return count, nil
		}
	}

	return 0, fmt.Errorf("Unable to convert %s to slice", val)
}

// MapLength returns the length of the passed in map, it returns an error if the argument is not a map
func MapLength(v interface{}) (int, error) {
	if v == nil {
		return 0, fmt.Errorf("Cannot convert nil to map")
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Map {
		return 0, fmt.Errorf("Unable to convert %s to map", val)
	}

	return val.Len(), nil
}

func IsMap(v interface{}) bool {
	val := reflect.ValueOf(v)
	return val.Kind() == reflect.Map
}

// LookupIndex tries to look up the interface at the passed in index for the passed in slice
func LookupIndex(v interface{}, idx int) (interface{}, error) {
	if v == nil {
		return nil, fmt.Errorf("Cannot convert nil to interface array")
	}

	// deal with a passed in error, we just return it out
	err, isErr := v.(error)
	if isErr {
		return nil, err
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Cannot convert non-array to interface array: %v", v)
	}

	if idx >= val.Len() || idx < -val.Len() {
		return nil, fmt.Errorf("Index %d out of range for slice of length %d", idx, val.Len())
	}

	if idx < 0 {
		idx += val.Len()
	}

	return val.Index(idx).Interface(), nil
}

// LookupKey tries to look up the interface at the passed in key for the passed in slice
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

	defer func() {
		if recover() != nil {
			value = nil
			err = fmt.Errorf("Invalid key type for map: %v", key)
		}
	}()

	mapValue := val.MapIndex(reflect.ValueOf(key))
	if !mapValue.IsValid() {
		return nil, nil
	}

	value = mapValue.Interface()
	return value, nil
}

// ToStringArray tries to turn the passed in interface (which must be an underlying slice) to a string array
func ToStringArray(env Environment, v interface{}) ([]string, error) {
	if v == nil {
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
		return NewJSONFragment(bytes), err
	}

	// null is null
	if val == nil {
		return ToFragment(json.Marshal(nil))
	}

	switch val := val.(type) {

	case error:
		return EmptyJSONFragment, val

	case string:
		return ToFragment(json.Marshal(val))

	case bool:
		return ToFragment(json.Marshal(val))

	case int, int32, int64:
		return ToFragment(json.Marshal(val))

	case float32, float64:
		return ToFragment(json.Marshal(val))

	case decimal.Decimal:
		floatVal, _ := val.Float64()
		return ToFragment(json.Marshal(floatVal))

	case time.Time:
		return ToFragment(json.Marshal(DateToISO(val)))

	case JSONFragment:
		return val, nil

	case fmt.Stringer:
		return ToFragment(json.Marshal(val.String()))

	case VariableResolver:
		// this checks that we aren't getting into an infinite loop
		valDefault := val.Default()
		valResolver, isResolver := valDefault.(VariableResolver)
		if isResolver && reflect.DeepEqual(valResolver, val) {
			return EmptyJSONFragment, fmt.Errorf("Loop found in ToJSON of '%s' with value '%+v'", reflect.TypeOf(val), val)
		}
		return ToJSON(env, valDefault)

	case []string:
		return ToFragment(json.Marshal(val))

	case []bool:
		return ToFragment(json.Marshal(val))

	case []time.Time:
		times := make([]string, len(val))
		for i := range val {
			times[i] = DateToISO(val[i])
		}
		return ToFragment(json.Marshal(times))

	case []decimal.Decimal:
		return ToFragment(json.Marshal(val))

	case []int:
		return ToFragment(json.Marshal(val))

	case map[string]string:
		return ToFragment(json.Marshal(val))

	case map[string]bool:
		return ToFragment(json.Marshal(val))

	case map[string]int:
		return ToFragment(json.Marshal(val))

	case map[string]interface{}:
		return ToFragment(json.Marshal(val))
	}

	// welp, we give up, this isn't something we can convert, return an error
	return EmptyJSONFragment, fmt.Errorf("ToString unknown type '%s' with value '%+v'", reflect.TypeOf(val), val)
}

// ToString tries to turn the passed in interface to a string
func ToString(env Environment, val interface{}) (string, error) {
	// Strings are always defined, just empty
	if val == nil {
		return "", nil
	}

	switch val := val.(type) {

	case error:
		return "", val

	case string:
		return val, nil

	case bool:
		return strconv.FormatBool(val), nil

	case int:
		return strconv.FormatInt(int64(val), 10), nil
	case int32:
		return strconv.FormatInt(int64(val), 10), nil
	case int64:
		return strconv.FormatInt(val, 10), nil

	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil

	case decimal.Decimal:
		return val.String(), nil

	case time.Time:
		return DateToISO(val), nil

	case fmt.Stringer:
		return val.String(), nil

	case VariableResolver:
		// this checks that we aren't getting into an infinite loop
		valDefault := val.Default()
		valResolver, isResolver := valDefault.(VariableResolver)
		if isResolver && reflect.DeepEqual(valResolver, val) {
			return "", fmt.Errorf("Loop found in ToString of '%s' with value '%+v'", reflect.TypeOf(val), val)
		}
		return ToString(env, valDefault)

	case []string:
		return strings.Join(val, ", "), nil

	case []bool:
		var output bytes.Buffer
		for i := range val {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString(strconv.FormatBool(val[i]))
		}
		return output.String(), nil

	case []time.Time:
		var output bytes.Buffer
		for i := range val {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString(DateToISO(val[i]))
		}
		return output.String(), nil

	case []decimal.Decimal:
		var output bytes.Buffer
		for i := range val {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString(val[i].String())
		}
		return output.String(), nil

	case []int:
		var output bytes.Buffer
		for i := range val {
			if i > 0 {
				output.WriteString(", ")
			}
			output.WriteString(strconv.FormatInt(int64(val[i]), 10))
		}
		return output.String(), nil

	case map[string]string:
		bytes, err := json.Marshal(val)
		return string(bytes), err
	}

	// welp, we give up, this isn't something we can convert, return an error
	return "", fmt.Errorf("ToString unknown type '%s' with value '%+v'", reflect.TypeOf(val), val)
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
	if val == nil {
		return decimal.Zero, nil
	}

	switch val := val.(type) {

	case error:
		return decimal.Zero, val

	case decimal.Decimal:
		return val, nil

	case int:
		return decimal.NewFromString(strconv.FormatInt(int64(val), 10))
	case int32:
		return decimal.NewFromString(strconv.FormatInt(int64(val), 10))
	case int64:
		return decimal.NewFromString(strconv.FormatInt(val, 10))

	case float32:
		return decimal.NewFromFloat(float64(val)), nil
	case float64:
		return decimal.NewFromFloat(val), nil

	case string:
		// common SMS foibles
		val = strings.ToLower(val)
		val = strings.Replace(val, "o", "0", -1)
		val = strings.Replace(val, "l", "1", -1)
		return decimal.NewFromString(val)
	}

	asString, err := ToString(env, val)
	if err != nil {
		return decimal.Zero, err
	}

	return ToDecimal(env, asString)
}

// ToDate tries to convert the passed in interface to a time.Time returning an error if that isn't possible
func ToDate(env Environment, val interface{}) (time.Time, error) {
	if val == nil {
		return time.Time{}, fmt.Errorf("Cannot convert nil to date")
	}

	switch val := val.(type) {

	case error:
		return time.Time{}, val

	case decimal.Decimal:
		return time.Time{}, fmt.Errorf("Cannot convert decimal to date")

	case int, int32, int64:
		return time.Time{}, fmt.Errorf("Cannot convert integer to date")

	case float32, float64:
		return time.Time{}, fmt.Errorf("Cannot convert float to date")

	case time.Time:
		return val, nil

	case string:
		return DateFromString(env, val)
	}

	asString, err := ToString(env, val)
	if err != nil {
		return time.Time{}, err
	}

	return ToDate(env, asString)
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

	case bool:
		return test, nil

	case int:
		return test != 0, nil
	case int32:
		return test != 0, nil
	case int64:
		return test != 0, nil

	case float32:
		return test != float32(0), nil
	case float64:
		return test != float64(0), nil

	case decimal.Decimal:
		return !test.Equals(decimal.Zero), nil

	case time.Time:
		return !test.IsZero(), nil

	case string:
		return test != "" && strings.ToLower(test) != "false", nil

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
	}

	asString, err := ToString(env, test)
	if err != nil {
		return false, err
	}

	return ToBool(env, asString)
}

// XType is an an enumeration of the possible types we can deal with
type XType int

const ( // primitive types we convert to
	NIL = iota
	STRING
	DECIMAL
	TIME
	BOOLEAN
	ERROR
	STRING_SLICE
	DECIMAL_SLICE
	TIME_SLICE
	BOOL_SLICE
	MAP
)

// ToXAtom figures out the raw type of the passed in interface, returning that type
func ToXAtom(env Environment, val interface{}) (interface{}, XType, error) {
	if val == nil {
		return val, NIL, nil
	}

	switch val := val.(type) {
	case error:
		return val, ERROR, nil

	case decimal.Decimal:
		return val, DECIMAL, nil

	case int, int32, int64, float32, float64:
		decVal, err := ToDecimal(env, val)
		if err != nil {
			return val, NIL, err
		}
		return decVal, DECIMAL, nil

	case time.Time:
		return val, TIME, nil

	case string:
		return val, STRING, nil

	case bool:
		return val, BOOLEAN, nil

	case []string:
		return val, STRING_SLICE, nil

	case []time.Time:
		return val, TIME_SLICE, nil

	case []decimal.Decimal:
		return val, DECIMAL_SLICE, nil

	case []bool:
		return val, BOOL_SLICE, nil
	}

	return val, NIL, fmt.Errorf("Unknown type '%s' with value '%+v'", reflect.TypeOf(val), val)
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
	case arg1Type == arg2Type && arg1Type == ERROR:
		return strings.Compare(arg1.(error).Error(), arg2.(error).Error()), nil

	case arg1Type == arg2Type && arg1Type == DECIMAL:
		return arg1.(decimal.Decimal).Cmp(arg2.(decimal.Decimal)), nil

	case arg1Type == arg2Type && arg1Type == BOOLEAN:
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

	case arg1Type == arg2Type && arg1Type == TIME:
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

	case arg1Type == arg2Type && arg1Type == STRING:
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
