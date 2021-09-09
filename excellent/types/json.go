package types

import (
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// JSONToXValue returns an X type from the given JSON
func JSONToXValue(data []byte) XValue {
	if len(data) == 0 {
		return nil
	}

	if !json.Valid(data) {
		return NewXErrorf("invalid JSON")
	}

	val, valType, _, err := jsonparser.Get(data)
	if err != nil {
		panic(fmt.Sprintf("jsonparser errored parsing valid JSON: %s", err))
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
		return jsonToObject(data)
	}

	return NewXError(errors.Errorf("unknown JSON parsing error"))
}

func jsonToObject(data []byte) *XObject {
	return NewXLazyObject(func() map[string]XValue {
		properties := make(map[string]XValue)

		jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			properties[string(key)] = jsonTypeToXValue(value, dataType)
			return nil
		})
		return properties
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

	marshaled, err := jsonx.Marshal(x)
	if err != nil {
		return XTextEmpty, NewXError(err)
	}

	return NewXText(string(marshaled)), nil
}
