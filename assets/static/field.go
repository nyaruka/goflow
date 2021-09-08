package static

import (
	"github.com/nyaruka/goflow/assets"
)

// Field is a JSON serializable implementation of a field asset
type Field struct {
	UUID_ assets.FieldUUID `json:"uuid"`
	Key_  string           `json:"key" validate:"required"`
	Name_ string           `json:"name"`
	Type_ assets.FieldType `json:"type" validate:"required"`
}

// NewField creates a new field from the passed in key, name and type
func NewField(uuid assets.FieldUUID, key string, name string, valueType assets.FieldType) assets.Field {
	return &Field{UUID_: uuid, Key_: key, Name_: name, Type_: valueType}
}

// UUID returns the UUID of this field
func (f *Field) UUID() assets.FieldUUID { return f.UUID_ }

// Key returns the unique key of the field
func (f *Field) Key() string { return f.Key_ }

// Name returns the name of the field
func (f *Field) Name() string { return f.Name_ }

// Type returns the value type of the field
func (f *Field) Type() assets.FieldType { return f.Type_ }
