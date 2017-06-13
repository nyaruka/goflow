package utils

import (
	"encoding/json"
	"strconv"

	"github.com/buger/jsonparser"
)

// NewJSONFragment creates a new json fragment for the passed in byte array
func NewJSONFragment(json []byte) JSONFragment {
	return JSONFragment{json: json}
}

// EmptyJSONFragment is a fragment which has no values
var EmptyJSONFragment = JSONFragment{nil}

// JSONFragment is a thin wrapper around a byte array that takes care of allow key lookups
// into the json in that byte array
type JSONFragment struct {
	json []byte
}

// Default returns the default value for this JSON, which is the JSON itself
func (j JSONFragment) Default() interface{} {
	return j
}

// Resolve resolves the passed in key, which is expected to be either an integer in the case
// that our JSON is an array or a key name if it is a map
func (j JSONFragment) Resolve(key string) interface{} {
	_, err := strconv.Atoi(key)

	// this is a numerical index, convert to jsonparser format
	if err == nil {
		jIdx := "[" + key + "]"
		val, valType, _, err := jsonparser.Get(j.json, jIdx)
		if err == nil {
			if err == nil {
				if valType == jsonparser.String {
					strVal, err := jsonparser.ParseString(val)
					if err == nil {
						return strVal
					}
				}
				return JSONFragment{val}
			}
		}
	}
	val, valType, _, err := jsonparser.Get(j.json, key)
	if err != nil {
		return err
	}

	if valType == jsonparser.String {
		strVal, err := jsonparser.ParseString(val)
		if err == nil {
			return strVal
		}
	}
	return JSONFragment{val}
}

var _ VariableResolver = EmptyJSONFragment

// String returns the string representation of this JSON, which is just the JSON itself
func (j JSONFragment) String() string {
	return string(j.json)
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
	j.json = data
	return nil
}

// MarshalJSON returns the JSON representation of our fragment, which is just our internal byte array
func (j *JSONFragment) MarshalJSON() ([]byte, error) {
	return j.json, nil
}
