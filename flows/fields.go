package flows

import (
	"encoding/json"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
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

var fieldLocationLevels = map[FieldValueType]LocationLevel{
	FieldValueTypeState:    LocationLevel(1),
	FieldValueTypeDistrict: LocationLevel(2),
	FieldValueTypeWard:     LocationLevel(3),
}

// Field represents a contact field
type Field struct {
	key       FieldKey
	name      string
	valueType FieldValueType
}

// NewField returns a new field object with the passed in uuid, key and value type
func NewField(key FieldKey, name string, valueType FieldValueType) *Field {
	return &Field{key: key, name: name, valueType: valueType}
}

// Key returns the key of the field
func (f *Field) Key() FieldKey { return f.key }

// FieldValue represents a contact's value for a specific field
type FieldValue struct {
	field    *Field
	text     types.XString
	datetime *types.XDate
	decimal  *types.XNumber
	state    *Location
	district *Location
	ward     *Location
}

func (v *FieldValue) IsEmpty() bool {
	return !(!v.text.Empty() || v.datetime != nil || v.decimal != nil || v.state != nil || v.district != nil || v.ward != nil)
}

func (v *FieldValue) TypedValue() types.XValue {
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
func (v *FieldValue) Resolve(key string) types.XValue {
	switch key {
	case "text":
		return v.text
	}
	return types.NewXResolveError(v, key)
}

// Reduce is called when this object needs to be reduced to a primitive
func (v *FieldValue) Reduce() types.XPrimitive {
	return v.TypedValue().Reduce()
}

// ToXJSON is called when this type is passed to @(to_json(...))
func (v *FieldValue) ToXJSON() types.XString { return v.Reduce().ToXJSON() }

var _ types.XValue = (*FieldValue)(nil)
var _ types.XResolvable = (*FieldValue)(nil)

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

func (f FieldValues) setValue(env utils.Environment, field *Field, rawValue types.XString) {
	var asDate *types.XDate
	var asNumber *types.XNumber

	if parsedNumber, xerr := types.ToXNumber(rawValue); xerr == nil {
		asNumber = &parsedNumber
	}

	if parsedDate, xerr := types.ToXDate(env, rawValue); xerr == nil {
		asDate = &parsedDate
	}

	// TODO parse as locations

	f[field.key] = &FieldValue{
		field:    field,
		text:     rawValue,
		datetime: asDate,
		decimal:  asNumber,
	}
}

// Length is called to get the length of this object
func (f FieldValues) Length() int {
	return len(f)
}

// Resolve resolves the given key when this set of field values is referenced in an expression
func (f FieldValues) Resolve(key string) types.XValue {
	val, exists := f[FieldKey(key)]
	if !exists {
		return types.NewXResolveError(f, key)
	}
	return val
}

// Reduce is called when this object needs to be reduced to a primitive
func (f FieldValues) Reduce() types.XPrimitive {
	values := types.NewXEmptyMap()
	for k, v := range f {
		values.Put(string(k), v)
	}
	return values
}

// ToXJSON is called when this type is passed to @(to_json(...))
func (f FieldValues) ToXJSON() types.XString {
	return f.Reduce().ToXJSON()
}

var _ types.XValue = (FieldValues)(nil)
var _ types.XLengthable = (FieldValues)(nil)
var _ types.XResolvable = (FieldValues)(nil)

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
	Name      string         `json:"name"`
	ValueType FieldValueType `json:"value_type,omitempty"`
}

// ReadField reads a contact field from the given JSON
func ReadField(data json.RawMessage) (*Field, error) {
	var fe fieldEnvelope
	if err := utils.UnmarshalAndValidate(data, &fe, "field"); err != nil {
		return nil, err
	}

	return NewField(fe.Key, fe.Name, fe.ValueType), nil
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
