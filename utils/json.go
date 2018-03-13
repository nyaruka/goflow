package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/buger/jsonparser"
)

// UnmarshalAndValidate is a convenience function to unmarshal an object and validate it
func UnmarshalAndValidate(data []byte, obj interface{}, objName string) error {
	err := json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	err = ValidateAs(obj, objName)
	if err != nil {
		return err
	}

	return nil
}

// UnmarshalArray unmarshals an array of objects from the given JSON
func UnmarshalArray(data json.RawMessage) ([]json.RawMessage, error) {
	var items []json.RawMessage
	err := json.Unmarshal(data, &items)
	return items, err
}

// EmptyJSONFragment is a fragment which has no values
var EmptyJSONFragment = JSONFragment{}

// JSONFragment is a thin wrapper around a byte array that takes care of allow key lookups
// into the json in that byte array
type JSONFragment []byte

// Resolve resolves the given key when this JSON fragment is referenced in an expression
func (j JSONFragment) Resolve(key string) interface{} {
	_, err := strconv.Atoi(key)

	// this is a numerical index, convert to jsonparser format
	if err == nil {
		jIdx := "[" + key + "]"
		val, valType, _, err := jsonparser.Get(j, jIdx)
		if err == nil {
			if err == nil {
				if valType == jsonparser.String {
					strVal, err := jsonparser.ParseString(val)
					if err == nil {
						return strVal
					}
				}
				return JSONFragment(val)
			}
		}
	}
	val, valType, _, err := jsonparser.Get(j, key)
	if err != nil {
		return fmt.Errorf("no such variable: %s", key)
	}

	if valType == jsonparser.String {
		strVal, err := jsonparser.ParseString(val)
		if err == nil {
			return strVal
		}
	}
	return JSONFragment(val)
}

// Default returns the value of this JSON fragment when it is the result of an expression
func (j JSONFragment) Default() interface{} {
	return j
}

var _ VariableResolver = EmptyJSONFragment

// String returns the string representation of this JSON, which is just the JSON itself
func (j JSONFragment) String() string {
	return string(j)
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
