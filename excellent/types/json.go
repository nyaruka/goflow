package types

import (
	"fmt"

	"github.com/nyaruka/goflow/utils"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func JSONToXValue(data []byte) XValue {
	if len(data) == 0 {
		return nil
	}

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
		return jsonToArray(data)
	case jsonparser.Object:
		return jsonToDict(data)
	}

	return NewXError(errors.Errorf("unknown JSON parsing error"))
}

func jsonToDict(data []byte) XDict {
	dict := NewEmptyXDict()
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		dict.Put(string(key), jsonTypeToXValue(value, dataType))
		return nil
	})
	return dict
}

func jsonToArray(data []byte) XArray {
	array := NewXArray()
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		array.Append(jsonTypeToXValue(value, dataType))
	})
	return array
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
