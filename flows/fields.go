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

// ParseValue returns a parsed field value for the given input
func (f *Field) ParseValue(env utils.Environment, value string) (interface{}, error) {
	switch f.valueType {
	case FieldValueTypeText:
		return value, nil
	case FieldValueTypeDatetime:
		return utils.DateFromString(env, value)
	case FieldValueTypeDecimal:
		return decimal.NewFromString(value)
	case FieldValueTypeState, FieldValueTypeDistrict, FieldValueTypeWard:
		locationID := utils.LocationID(value)
		locationLevel := fieldLocationLevels[f.valueType]
		locations, err := env.Locations()
		if err != nil {
			return nil, err
		}
		if locations == nil {
			return nil, fmt.Errorf("can't parse field '%s' (type %s) in environment which is not location enabled", f.key, f.valueType)
		}
		return locations.FindByID(locationID, locationLevel), nil
	}

	return nil, fmt.Errorf("field %s has invalid value type: '%s'", f.key, f.valueType)
}

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

// Default returns the value of this field value when it is the result of an expression
func (v *FieldValue) Default() interface{} {
	return v.TypedValue()
}

// String returns the string representation of this field value
func (v *FieldValue) String() string {
	return v.text
}

// FieldValues is the set of all field values for a contact
type FieldValues map[FieldKey]*FieldValue

// Clone returns a clone of this set of field values
func (f FieldValues) Clone() FieldValues {
	clone := make(FieldValues, len(f))
	for k, v := range f {
		clone[k] = v
	}
	return clone
}

// Save saves a new field value
func (f FieldValues) Save(env utils.Environment, field *Field, rawValue string) error {
	var asDatetime *time.Time
	var asDecimal *decimal.Decimal

	if parsedDecimal, err := decimal.NewFromString(rawValue); err == nil {
		asDecimal = &parsedDecimal
	}

	if parsedDatetime, err := utils.DateFromString(env, rawValue); err == nil {
		asDatetime = &parsedDatetime
	}

	// TODO parse as locations

	f[field.key] = &FieldValue{
		field:    field,
		text:     rawValue,
		datetime: asDatetime,
		decimal:  asDecimal,
	}
	return nil
}

// Resolve resolves the given key when this set of field values is referenced in an expression
func (f FieldValues) Resolve(key string) interface{} {
	return f[FieldKey(key)]
}

// Default returns the value of this set of field values when it is the result of an expression
func (f FieldValues) Default() interface{} {
	return f
}

// String returns the string representation of these Fields, which is our JSON representation
func (f FieldValues) String() string {
	fields := make([]string, 0, len(f))
	for k, v := range f {
		// TODO serilalize field value according to type
		fields = append(fields, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return strings.Join(fields, ", ")
}

var _ utils.VariableResolver = (FieldValues)(nil)

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
