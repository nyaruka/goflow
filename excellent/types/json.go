package types

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/shopspring/decimal"
)

// XJSON is the base type for XJSONObject and XJSONArray
type XJSON []byte

func (x XJSON) ToJSON() XString { return NewXString(string(x)) }

func (x XJSON) Reduce() XPrimitive { return x.ToJSON() }

type XJSONObject struct {
	XJSON
}

func NewXJSONObject(data []byte) XJSONObject {
	return XJSONObject{XJSON: data}
}

func (x XJSONObject) Length() int {
	length := 0
	jsonparser.ObjectEach(x.XJSON, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		length++
		return nil
	})
	return length
}

func (x XJSONObject) Resolve(key string) XValue {
	val, valType, _, err := jsonparser.Get(x.XJSON, key)
	if err != nil {
		return NewXError(fmt.Errorf("can't resolve '%s'", key))
	}

	return jsonTypeToXValue(val, valType)
}

var _ XValue = XJSONObject{}
var _ XLengthable = XJSONObject{}
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
