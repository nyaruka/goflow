package types

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// Field is a JSON serializable implementation of a field asset
type Field struct {
	Key_  string           `json:"key" validate:"required"`
	Name_ string           `json:"name"`
	Type_ assets.FieldType `json:"value_type" validate:"required"`
}

// NewField creates a new field from the passed in key, name and type
func NewField(key string, name string, valueType assets.FieldType) assets.Field {
	return &Field{Key_: key, Name_: name, Type_: valueType}
}

// Key returns the unique key of the field
func (f *Field) Key() string { return f.Key_ }

// Name returns the name of the field
func (f *Field) Name() string { return f.Name_ }

// Type returns the value type of the field
func (f *Field) Type() assets.FieldType { return f.Type_ }

// ReadField reads a field from the given JSON
func ReadField(data json.RawMessage) (assets.Field, error) {
	f := &Field{}
	if err := utils.UnmarshalAndValidate(data, f); err != nil {
		return nil, fmt.Errorf("unable to read field: %s", err)
	}
	return f, nil
}

// ReadFields reads fields from the given JSON
func ReadFields(data json.RawMessage) ([]assets.Field, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	fields := make([]assets.Field, len(items))
	for d := range items {
		if fields[d], err = ReadField(items[d]); err != nil {
			return nil, err
		}
	}

	return fields, nil
}
