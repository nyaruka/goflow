package flows

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"
)

// FieldReference is a reference to field used in a flow
type FieldReference struct {
	UUID FieldUUID `json:"uuid" validate:"uuid4"`
	Key  string    `json:"key"`
}

func NewFieldReference(uuid FieldUUID, key string) *FieldReference {
	return &FieldReference{UUID: uuid, Key: key}
}

type FieldValueType string

const (
	FieldValueTypeText     FieldValueType = "text"
	FieldValueTypeDecimal  FieldValueType = "decimal"
	FieldValueTypeDatetime FieldValueType = "datetime"
	FieldValueTypeWard     FieldValueType = "ward"
	FieldValueTypeDistrict FieldValueType = "district"
	FieldValueTypeState    FieldValueType = "state"
)

// Field represents a contact field
type Field struct {
	uuid      FieldUUID
	key       string
	valueType FieldValueType
}

// NewField returns a new field object with the passed in uuid, key and value type
func NewField(uuid FieldUUID, key string, valueType FieldValueType) *Field {
	return &Field{uuid: uuid, key: key, valueType: valueType}
}

// UUID returns the UUID of the field
func (f *Field) UUID() FieldUUID { return f.uuid }

// Key returns the key of the field
func (f *Field) Key() string { return f.key }

// FieldValue represents a contact's value for a specific field
type FieldValue struct {
	field     *Field
	value     interface{}
	createdOn time.Time
}

func (v *FieldValue) Resolve(key string) interface{} {
	switch key {
	case "value":
		return v.value
	case "created_on":
		return v.createdOn
	}
	return fmt.Errorf("no field '%s' on field value", key)
}

// Default returns the default value for FieldValue, which is the value
func (v *FieldValue) Default() interface{} {
	return v.value
}

// String returns the string representation of this field value
func (v *FieldValue) String() string {
	// TODO serilalize field value according to type
	return fmt.Sprintf("%s", v.value)
}

type FieldValues map[string]*FieldValue

func (f FieldValues) Save(field *Field, value string) {
	// TODO deserialize non-string values
	f[field.key] = &FieldValue{field: field, value: value, createdOn: time.Now().UTC()}
}

func (f FieldValues) Resolve(key string) interface{} {
	return f[key]
}

// Default returns the default value for FieldValues, which is ourselves
func (f FieldValues) Default() interface{} {
	return f
}

// String returns the string representation of these Fields, which is our JSON representation
func (f FieldValues) String() string {
	fields := make([]string, 0, len(f))
	for k, v := range f {
		// TODO serilalize field value according to type
		fields = append(fields, fmt.Sprintf("%s: %s", k, v.value))
	}
	return strings.Join(fields, ", ")
}

var _ utils.VariableResolver = (FieldValues)(nil)

// FieldSet defines the unordered set of all fields for a session
type FieldSet struct {
	fields       []*Field
	fieldsByUUID map[FieldUUID]*Field
}

func NewFieldSet(fields []*Field) *FieldSet {
	s := &FieldSet{fields: fields, fieldsByUUID: make(map[FieldUUID]*Field, len(fields))}
	for _, field := range s.fields {
		s.fieldsByUUID[field.uuid] = field
	}
	return s
}

func (s *FieldSet) FindByUUID(uuid FieldUUID) *Field {
	return s.fieldsByUUID[uuid]
}

// FindByKey looks for a field with the given key
func (s *FieldSet) FindByKey(key string) *Field {
	for _, field := range s.fields {
		if strings.ToLower(field.key) == key {
			return field
		}
	}
	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldEnvelope struct {
	UUID      FieldUUID      `json:"uuid" validate:"required,uuid4"`
	Key       string         `json:"key"`
	ValueType FieldValueType `json:"value_type,omitempty"`
}

func ReadField(data json.RawMessage) (*Field, error) {
	var fe fieldEnvelope
	if err := utils.UnmarshalAndValidate(data, &fe, "field"); err != nil {
		return nil, err
	}

	return NewField(fe.UUID, fe.Key, fe.ValueType), nil
}

func ReadFieldSet(data json.RawMessage) (*FieldSet, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	fields := make([]*Field, len(items))
	for d := range items {
		if fields[d], err = ReadField(items[d]); err != nil {
			return nil, err
		}
	}

	return NewFieldSet(fields), nil
}
