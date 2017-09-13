package flows

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/nyaruka/goflow/utils"
)

// FieldReference is a reference to field used in a flow
type FieldReference struct {
	UUID FieldUUID `json:"uuid" validate:"uuid4"`
	Key  string    `json:"key"`
}

// NewFieldReference creates a new field reference with the given UUID and key
func NewFieldReference(uuid FieldUUID, key string) *FieldReference {
	return &FieldReference{UUID: uuid, Key: key}
}

// FieldValueType is the data type of values for each field
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

func (f *Field) ParseValue(env utils.Environment, value string) (interface{}, error) {
	switch f.valueType {
	case FieldValueTypeText:
		return value, nil
	case FieldValueTypeDecimal:
		return decimal.NewFromString(value)
	case FieldValueTypeDatetime:
		return utils.DateFromString(env, value)
	}

	// TODO location field values

	return nil, fmt.Errorf("field %s has invalid value type: '%s'", f.key, f.valueType)
}

// FieldValue represents a contact's value for a specific field
type FieldValue struct {
	field     *Field
	value     interface{}
	createdOn time.Time
}

func NewFieldValue(field *Field, value interface{}, createdOn time.Time) *FieldValue {
	return &FieldValue{field: field, value: value, createdOn: createdOn}
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
	// TODO serialize field value according to type
	return fmt.Sprintf("%s", v.value)
}

// String returns the string representation of this field value
func (v *FieldValue) JSON() string {
	switch v.field.valueType {
	case FieldValueTypeText:
		return v.value.(string)
	case FieldValueTypeDecimal:
		return v.value.(decimal.Decimal).String()
	case FieldValueTypeDatetime:
		return utils.DateToISO(v.value.(time.Time))
	}

	// TODO location field values

	return fmt.Sprintf("%v", v.value)
}

type FieldValues map[string]*FieldValue

func (f FieldValues) Save(env utils.Environment, field *Field, rawValue string) error {
	value, err := field.ParseValue(env, rawValue)
	if err != nil {
		return err
	}

	f[field.key] = NewFieldValue(field, value, time.Now().UTC())
	return nil
}

func (f FieldValues) Resolve(key string) interface{} {
	fmt.Printf("FieldValues.Resolve(%s) -> %v\n", key, f[key])
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
