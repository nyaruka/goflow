package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// FieldUUID is the UUID of a field
type FieldUUID uuids.UUID

// FieldType is the data type of values for each field
type FieldType string

// field value types
const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeDatetime FieldType = "datetime"
	FieldTypeWard     FieldType = "ward"
	FieldTypeDistrict FieldType = "district"
	FieldTypeState    FieldType = "state"
)

// Field is a custom contact property.
//
//   {
//     "uuid": "d66a7823-eada-40e5-9a3a-57239d4690bf",
//     "key": "gender",
//     "name": "Gender",
//     "type": "text"
//   }
//
// @asset field
type Field interface {
	UUID() FieldUUID
	Key() string
	Name() string
	Type() FieldType
}

// FieldReference is a reference to a field
type FieldReference struct {
	Key  string `json:"key" validate:"required"`
	Name string `json:"name"`
}

// NewFieldReference creates a new field reference with the given key and name
func NewFieldReference(key string, name string) *FieldReference {
	return &FieldReference{Key: key, Name: name}
}

// Type returns the name of the asset type
func (r *FieldReference) Type() string {
	return "field"
}

// Identity returns the unique identity of the asset
func (r *FieldReference) Identity() string {
	return string(r.Key)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *FieldReference) Variable() bool {
	return false
}

func (r *FieldReference) String() string {
	return fmt.Sprintf("%s[key=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ Reference = (*FieldReference)(nil)
