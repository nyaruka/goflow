package flows

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

// FieldKey is the unique key for this field
type FieldKey string

// FieldValueType is the data type of values for each field
type FieldValueType string

// field value types
const (
	FieldValueTypeText     FieldValueType = "text"
	FieldValueTypeDecimal  FieldValueType = "decimal"
	FieldValueTypeDatetime FieldValueType = "datetime"
	FieldValueTypeWard     FieldValueType = "ward"
	FieldValueTypeDistrict FieldValueType = "district"
	FieldValueTypeState    FieldValueType = "state"
)

var fieldLocationLevels = map[FieldValueType]utils.LocationLevel{
	FieldValueTypeState:    utils.LocationLevel(1),
	FieldValueTypeDistrict: utils.LocationLevel(2),
	FieldValueTypeWard:     utils.LocationLevel(3),
}

// Field represents a contact field
type Field struct {
	key       FieldKey
	label     string
	valueType FieldValueType
}

// NewField returns a new field object with the passed in uuid, key and value type
func NewField(key FieldKey, label string, valueType FieldValueType) *Field {
	return &Field{key: key, label: label, valueType: valueType}
}

// Key returns the key of the field
func (f *Field) Key() FieldKey { return f.key }

// FieldValue represents a contact's value for a specific field
type FieldValue struct {
	field    *Field
	text     string
	datetime *time.Time
	decimal  *decimal.Decimal
	state    *utils.Location
	district *utils.Location
	ward     *utils.Location
}

func (v *FieldValue) IsEmpty() bool {
	return !(v.text != "" || v.datetime != nil || v.decimal != nil || v.state != nil || v.district != nil || v.ward != nil)
}

func (v *FieldValue) TypedValue() interface{} {
	switch v.field.valueType {
	case FieldValueTypeText:
		return v.text
	case FieldValueTypeDatetime:
		if v.datetime != nil {
			return *v.datetime
		}
	case FieldValueTypeDecimal:
		if v.decimal != nil {
			return *v.decimal
		}
	case FieldValueTypeState:
		return v.state
	case FieldValueTypeDistrict:
		return v.district
	case FieldValueTypeWard:
		return v.ward
	}
	return nil
}

// Resolve resolves the given key when this field value is referenced in an expression
func (v *FieldValue) Resolve(key string) interface{} {
	switch key {
	case "text":
		return v.text
	}
	return fmt.Errorf("no field '%s' on field value", key)
}

// Atomize is called when this object needs to be reduced to a primitive
func (v *FieldValue) Atomize() interface{} {
	return v.TypedValue()
}

// FieldValues is the set of all field values for a contact
type FieldValues map[FieldKey]*FieldValue

// Clone returns a clone of this set of field values
func (f FieldValues) clone() FieldValues {
	clone := make(FieldValues, len(f))
	for k, v := range f {
		clone[k] = v
	}
	return clone
}

func (f FieldValues) setValue(env utils.Environment, field *Field, rawValue string) {
	var asDatetime *time.Time
	var asDecimal *decimal.Decimal

	if parsedDecimal, err := utils.ToDecimal(env, rawValue); err == nil {
		asDecimal = &parsedDecimal
	}

	if parsedDatetime, err := utils.ToDate(env, rawValue); err == nil {
		asDatetime = &parsedDatetime
	}

	// TODO parse as locations

	f[field.key] = &FieldValue{
		field:    field,
		text:     rawValue,
		datetime: asDatetime,
		decimal:  asDecimal,
	}
}

// Resolve resolves the given key when this set of field values is referenced in an expression
func (f FieldValues) Resolve(key string) interface{} {
	val, exists := f[FieldKey(key)]
	if !exists {
		return fmt.Errorf("no such contact field '%s'", key)
	}
	return val
}

// Atomize is called when this object needs to be reduced to a primitive
func (f FieldValues) Atomize() interface{} {
	fields := make([]string, 0, len(f))
	for k, v := range f {
		// TODO serilalize field value according to type
		fields = append(fields, fmt.Sprintf("%s: %s", k, v.TypedValue()))
	}
	return strings.Join(fields, ", ")
}

var _ utils.Atomizable = (FieldValues)(nil)
var _ utils.Resolvable = (FieldValues)(nil)

// FieldSet defines the unordered set of all fields for a session
type FieldSet struct {
	fields      []*Field
	fieldsByKey map[FieldKey]*Field
}

// NewFieldSet creates a new set of fields
func NewFieldSet(fields []*Field) *FieldSet {
	s := &FieldSet{fields: fields, fieldsByKey: make(map[FieldKey]*Field, len(fields))}
	for _, field := range s.fields {
		s.fieldsByKey[field.key] = field
	}
	return s
}

// FindByKey finds the contact field with the given key
func (s *FieldSet) FindByKey(key FieldKey) *Field {
	return s.fieldsByKey[key]
}

func (s *FieldSet) All() []*Field {
	return s.fields
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldEnvelope struct {
	Key       FieldKey       `json:"key"`
	Label     string         `json:"label"`
	ValueType FieldValueType `json:"value_type,omitempty"`
}

// ReadField reads a contact field from the given JSON
func ReadField(data json.RawMessage) (*Field, error) {
	var fe fieldEnvelope
	if err := utils.UnmarshalAndValidate(data, &fe, "field"); err != nil {
		return nil, err
	}

	return NewField(fe.Key, fe.Label, fe.ValueType), nil
}

// ReadFieldSet reads a set of contact fields from the given JSON
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
