package types

import (
	"encoding/json"
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

func jsonToDict(data []byte) *XDict {
	return NewXLazyDict(func() map[string]XValue {
		entries := make(map[string]XValue, 0)

		jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			entries[string(key)] = jsonTypeToXValue(value, dataType)
			return nil
		})
		return entries
	})
}

func jsonToArray(data []byte) *XArray {
	return NewXLazyArray(func() []XValue {

		items := make([]XValue, 0)

		jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			items = append(items, jsonTypeToXValue(value, dataType))
		})
		return items
	})
}

// ToXJSON converts the given value to a JSON string
func ToXJSON(x XValue) (XText, XError) {
	if utils.IsNil(x) {
		return NewXText(`null`), nil
	}
	if IsXError(x) {
		return XTextEmpty, x.(XError)
	}

	marshaled, err := json.Marshal(x)
	if err != nil {
		return XTextEmpty, NewXError(err)
	}

	return NewXText(string(marshaled)), nil
}

// MustMarshalToXText calls json.Marshal in the given value and panics in the case of an error
func MustMarshalToXText(x interface{}) XText {
	j, err := utils.JSONMarshal(x)
	if err != nil {
		panic(fmt.Sprintf("unable to marshal %s to JSON", x))
	}
	return NewXText(string(j))
}
