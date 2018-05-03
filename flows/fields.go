package flows

import (
	"encoding/json"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// FieldValueType is the data type of values for each field
type FieldValueType string

// field value types
const (
	FieldValueTypeText     FieldValueType = "text"
	FieldValueTypeNumber   FieldValueType = "number"
	FieldValueTypeDatetime FieldValueType = "datetime"
	FieldValueTypeWard     FieldValueType = "ward"
	FieldValueTypeDistrict FieldValueType = "district"
	FieldValueTypeState    FieldValueType = "state"
)

// Field represents a contact field
type Field struct {
	key       string
	name      string
	valueType FieldValueType
}

// NewField returns a new field object with the passed in uuid, key and value type
func NewField(key string, name string, valueType FieldValueType) *Field {
	return &Field{key: key, name: name, valueType: valueType}
}

// Key returns the key of the field
func (f *Field) Key() string { return f.key }

// FieldValue represents a contact's value for a specific field
type FieldValue struct {
	field    *Field
	text     types.XText
	datetime *types.XDateTime
	number   *types.XNumber
	state    LocationPath
	district LocationPath
	ward     LocationPath
}

// NewEmptyFieldValue creates a new empty value for the given field
func NewEmptyFieldValue(field *Field) *FieldValue {
	return &FieldValue{field: field}
}

// IsEmpty returns whether this field value is set for any type
func (v *FieldValue) IsEmpty() bool {
	return v.text.Empty() && v.datetime == nil && v.number == nil && v.state == "" && v.district == "" && v.ward == ""
}

// TypedValue returns the value in its proper type
func (v *FieldValue) TypedValue() types.XValue {
	switch v.field.valueType {
	case FieldValueTypeText:
		return v.text
	case FieldValueTypeDatetime:
		if v.datetime != nil {
			return *v.datetime
		}
	case FieldValueTypeNumber:
		if v.number != nil {
			return *v.number
		}
	case FieldValueTypeState:
		return types.NewXText(string(v.state))
	case FieldValueTypeDistrict:
		return types.NewXText(string(v.district))
	case FieldValueTypeWard:
		return types.NewXText(string(v.ward))
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

// ToXJSON is called when this type is passed to @(json(...))
func (v *FieldValue) ToXJSON() types.XText { return v.Reduce().ToXJSON() }

var _ types.XValue = (*FieldValue)(nil)
var _ types.XResolvable = (*FieldValue)(nil)

// EmptyFieldValue is used when a contact doesn't have a value set for a field
var EmptyFieldValue = &FieldValue{}

// FieldValues is the set of all field values for a contact
type FieldValues map[string]*FieldValue

// Clone returns a clone of this set of field values
func (f FieldValues) clone() FieldValues {
	clone := make(FieldValues, len(f))
	for k, v := range f {
		clone[k] = v
	}
	return clone
}

func (f FieldValues) getValue(key string) *FieldValue {
	return f[key]
}

func (f FieldValues) setValue(env RunEnvironment, fieldSet *FieldSet, key string, rawValue string) {
	field := fieldSet.FindByKey(key)

	if rawValue == "" {
		f[key] = NewEmptyFieldValue(field)
		return
	}

	var asText = types.NewXText(rawValue)
	var asDateTime *types.XDateTime
	var asNumber *types.XNumber

	if parsedNumber, xerr := types.ToXNumber(asText); xerr == nil {
		asNumber = &parsedNumber
	}

	if parsedDate, xerr := types.ToXDateTime(env, asText); xerr == nil {
		asDateTime = &parsedDate
	}

	var asLocation *utils.Location

	// for locations, if it has a '>' then it is explicit, look it up that way
	if IsPossibleLocationPath(rawValue) {
		asLocation, _ = env.LookupLocation(LocationPath(rawValue))
	} else {
		var matchingLocations []*utils.Location

		if field.valueType == FieldValueTypeWard {
			parent := f.getFirstLocationValue(env, fieldSet, FieldValueTypeDistrict)
			if parent != nil {
				matchingLocations, _ = env.FindLocationsFuzzy(rawValue, LocationLevelWard, parent)
			}
		} else if field.valueType == FieldValueTypeDistrict {
			parent := f.getFirstLocationValue(env, fieldSet, FieldValueTypeState)
			if parent != nil {
				matchingLocations, _ = env.FindLocationsFuzzy(rawValue, LocationLevelDistrict, parent)
			}
		} else if field.valueType == FieldValueTypeState {
			matchingLocations, _ = env.FindLocationsFuzzy(rawValue, LocationLevelState, nil)
		}

		if len(matchingLocations) > 0 {
			asLocation = matchingLocations[0]
		}
	}

	var asState, asDistrict, asWard LocationPath
	if asLocation != nil {
		switch asLocation.Level() {
		case LocationLevelState:
			asState = LocationPath(asLocation.Path())
		case LocationLevelDistrict:
			asState = LocationPath(asLocation.Parent().Path())
			asDistrict = LocationPath(asLocation.Path())
		case LocationLevelWard:
			asState = LocationPath(asLocation.Parent().Parent().Path())
			asDistrict = LocationPath(asLocation.Parent().Path())
			asWard = LocationPath(asLocation.Path())
		}
	}

	f[key] = &FieldValue{
		field:    field,
		text:     asText,
		datetime: asDateTime,
		number:   asNumber,
		state:    asState,
		district: asDistrict,
		ward:     asWard,
	}
}

func (f FieldValues) getFirstLocationValue(env RunEnvironment, fieldSet *FieldSet, valueType FieldValueType) *utils.Location {
	field := fieldSet.FirstOfType(FieldValueTypeDistrict)
	if field == nil {
		return nil
	}
	value := f[field.Key()].TypedValue()
	location, err := env.LookupLocation(value.(LocationPath))
	if err != nil {
		return nil
	}
	return location
}

// Length is called to get the length of this object
func (f FieldValues) Length() int {
	return len(f)
}

// Resolve resolves the given key when this set of field values is referenced in an expression
func (f FieldValues) Resolve(key string) types.XValue {
	val, exists := f[key]
	if !exists {
		return types.NewXResolveError(f, key)
	}
	return val
}

// Reduce is called when this object needs to be reduced to a primitive
func (f FieldValues) Reduce() types.XPrimitive {
	values := types.NewEmptyXMap()
	for k, v := range f {
		values.Put(string(k), v)
	}
	return values
}

// ToXJSON is called when this type is passed to @(json(...))
func (f FieldValues) ToXJSON() types.XText {
	return f.Reduce().ToXJSON()
}

var _ types.XValue = (FieldValues)(nil)
var _ types.XLengthable = (FieldValues)(nil)
var _ types.XResolvable = (FieldValues)(nil)

// FieldSet defines the unordered set of all fields for a session
type FieldSet struct {
	fields      []*Field
	fieldsByKey map[string]*Field
}

// NewFieldSet creates a new set of fields
func NewFieldSet(fields []*Field) *FieldSet {
	s := &FieldSet{
		fields:      fields,
		fieldsByKey: make(map[string]*Field, len(fields)),
	}
	for _, field := range s.fields {
		s.fieldsByKey[field.key] = field
	}
	return s
}

// FindByKey finds the contact field with the given key
func (s *FieldSet) FindByKey(key string) *Field {
	return s.fieldsByKey[key]
}

// FirstOfType returns the first field in this set with the given value type
func (s *FieldSet) FirstOfType(valueType FieldValueType) *Field {
	for _, field := range s.fields {
		if field.valueType == valueType {
			return field
		}
	}
	return nil
}

// All returns all the fields in this set
func (s *FieldSet) All() []*Field {
	return s.fields
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldEnvelope struct {
	Key       string         `json:"key"`
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
