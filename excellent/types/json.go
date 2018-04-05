package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/shopspring/decimal"
)

type XJSON []byte

func (x XJSON) ToJSON() XString { return NewXString(string(x)) }

func (x XJSON) Reduce() XPrimitive { return x.ToJSON() }

type XJSONObject struct {
	XJSON
}

func NewXJSONObject(data []byte) XJSONObject {
	return XJSONObject{XJSON: data}
}

func (x XJSONObject) Resolve(key string) XValue {
	val, valType, _, err := jsonparser.Get(x.XJSON, key)
	if err != nil {
		return NewXError(fmt.Errorf("can't resolve '%s'", key))
	}

	return jsonTypeToXValue(val, valType)
}

var _ XValue = XJSONObject{}
var _ XResolvable = XJSONObject{}

type XJSONArray struct {
	XJSON
}

func NewXJSONArray(data []byte) XJSONArray {
	return XJSONArray{XJSON: data}
}

func (x XJSONArray) Length() int {
	length := 0
	jsonparser.ArrayEach(x.XJSON, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		length++
	})
	return length
}

func (x XJSONArray) Index(index int) XValue {
	val, valType, _, err := jsonparser.Get(x.XJSON, fmt.Sprintf("[%d]", index))
	if err != nil {
		return NewXError(err)
	}
	return jsonTypeToXValue(val, valType)
}

var _ XValue = XJSONArray{}
var _ XIndexable = XJSONArray{}

func JSONToXValue(data []byte) XValue {
	val, valType, _, err := jsonparser.Get(data)
	if err != nil {
		return NewXError(err)
	}
	return jsonTypeToXValue(val, valType)
}

func jsonTypeToXValue(data []byte, valType jsonparser.ValueType) XValue {
	switch valType {
	case jsonparser.Null:
		return nil
	case jsonparser.String:
		strVal, err := jsonparser.ParseString(data)
		if err == nil {
			return NewXString(strVal)
		}
	case jsonparser.Number:
		decimalVal, err := decimal.NewFromString(string(data))
		if err == nil {
			return NewXNumber(decimalVal)
		}
	case jsonparser.Boolean:
		boolVal, err := jsonparser.ParseBoolean(data)
		if err == nil {
			return NewXBool(boolVal)
		}
	case jsonparser.Array:
		return NewXJSONArray(data)
	case jsonparser.Object:
		return NewXJSONObject(data)
	}

	return NewXError(fmt.Errorf("unknown JSON parsing error"))
}

// Legacy...

// EmptyJSONFragment is a fragment which has no values
var EmptyJSONFragment = JSONFragment{}

// JSONFragment is a thin wrapper around a byte array that takes care of allow key lookups
// into the json in that byte array
type JSONFragment []byte

// Resolve resolves the given key when this JSON fragment is referenced in an expression
func (j JSONFragment) Resolve(key string) interface{} {
	_, isIndex := strconv.Atoi(key)

	// this is a numerical index, convert to jsonparser format
	if isIndex == nil {
		jIdx := "[" + key + "]"
		val, valType, _, err := jsonparser.Get(j, jIdx)
		if err == nil {
			return jsonTypeToXAtom(val, valType)
		}
	}

	val, valType, _, err := jsonparser.Get(j, key)
	if err != nil {
		return fmt.Errorf("no such variable: %s", key)
	}

	return jsonTypeToXAtom(val, valType)
}

// Atomize is called when this object needs to be reduced to a primitive
func (j JSONFragment) Atomize() interface{} {
	return string(j)
}

var _ Atomizable = EmptyJSONFragment
var _ Resolvable = EmptyJSONFragment

// JSONArray is a JSON fragment containing an array
type JSONArray JSONFragment

// Atomize is called when this object needs to be reduced to a primitive
func (j JSONArray) Atomize() interface{} {
	return string(j)
}

// Index is called when this object is indexed into in an expression
func (j JSONArray) Index(index int) interface{} {
	val, valType, _, err := jsonparser.Get(j, fmt.Sprintf("[%d]", index))
	if err != nil {
		return err
	}
	return jsonTypeToXAtom(val, valType)
}

// Length is called when the length of this object is requested in an expression
func (j JSONArray) Length() int {
	length := 0
	jsonparser.ArrayEach(j, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		length++
	})
	return length
}

var _ Atomizable = JSONArray{}
var _ Indexable = JSONArray{}

func jsonTypeToXAtom(data []byte, valType jsonparser.ValueType) interface{} {
	switch valType {
	case jsonparser.Null:
		return nil
	case jsonparser.String:
		strVal, err := jsonparser.ParseString(data)
		if err == nil {
			return strVal
		}
	case jsonparser.Number:
		decimalVal, err := decimal.NewFromString(string(data))
		if err == nil {
			return decimalVal
		}
	case jsonparser.Boolean:
		boolVal, err := jsonparser.ParseBoolean(data)
		if err == nil {
			return boolVal
		}
	case jsonparser.Array:
		return JSONArray(data)
	}

	return JSONFragment(data)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// UnmarshalJSON reads a new JSONFragment from the passed in byte stream. We validate it looks
// like valid JSON then set our internal byte structure
func (j *JSONFragment) UnmarshalJSON(data []byte) error {
	// try to parse the passed in data as JSON
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return err
	}
	*j = data
	return nil
}

// MarshalJSON returns the JSON representation of our fragment, which is just our internal byte array
func (j JSONFragment) MarshalJSON() ([]byte, error) {
	return j, nil
}
