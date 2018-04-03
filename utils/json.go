package utils

import (
	"encoding/json"
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
