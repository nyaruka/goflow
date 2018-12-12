package types

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/utils"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// XJSON is the base type for XJSONObject and XJSONArray
type XJSON []byte

func (x XJSON) ToXJSON(env utils.Environment) XText { return NewXText(string(x)) }

func (x XJSON) Reduce(env utils.Environment) XPrimitive { return x.ToXJSON(env) }

// String converts this type to native string
func (x XJSON) String() string {
	return string(x)
}

func (x XJSON) MarshalJSON() ([]byte, error) {
	return []byte(x), nil
}

type XJSONObject struct {
	XJSON
}

func NewXJSONObject(data []byte) XJSONObject {
	return XJSONObject{XJSON: data}
}

// Describe returns a representation of this type for error messages
func (x XJSONObject) Describe() string { return "json object" }

func (x XJSONObject) Length() int {
	length := 0
	jsonparser.ObjectEach(x.XJSON, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		length++
		return nil
	})
	return length
}

func (x XJSONObject) Resolve(env utils.Environment, key string) XValue {
	val, valType, _, err := jsonparser.Get(x.XJSON, key)
	if err != nil {
		return NewXResolveError(x, key)
	}

	return jsonTypeToXValue(val, valType)
}

var _ XValue = XJSONObject{}
var _ XLengthable = XJSONObject{}
var _ XResolvable = XJSONObject{}
var _ json.Marshaler = XJSONObject{}

type XJSONArray struct {
	XJSON
}

func NewXJSONArray(data []byte) XJSONArray {
	return XJSONArray{XJSON: data}
}

// Describe returns a representation of this type for error messages
func (x XJSONArray) Describe() string { return "json array" }

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
var _ json.Marshaler = XJSONArray{}

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
			return NewXText(strVal)
		}
	case jsonparser.Number:
		decimalVal, err := decimal.NewFromString(string(data))
		if err == nil {
			return NewXNumber(decimalVal)
		}
	case jsonparser.Boolean:
		boolVal, err := jsonparser.ParseBoolean(data)
		if err == nil {
			return NewXBoolean(boolVal)
		}
	case jsonparser.Array:
		return NewXJSONArray(data)
	case jsonparser.Object:
		return NewXJSONObject(data)
	}

	return NewXError(errors.Errorf("unknown JSON parsing error"))
}

// ToXJSON converts the given value to a JSON string
func ToXJSON(env utils.Environment, x XValue) (XText, XError) {
	if utils.IsNil(x) {
		return NewXText(`null`), nil
	}
	if IsXError(x) {
		return XTextEmpty, x.(XError)
	}

	return x.ToXJSON(env), nil
}

// MustMarshalToXText calls json.Marshal in the given value and panics in the case of an error
func MustMarshalToXText(x interface{}) XText {
	j, err := utils.JSONMarshal(x)
	if err != nil {
		panic(fmt.Sprintf("unable to marshal %s to JSON", x))
	}
	return NewXText(string(j))
}
