package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Field represents a contact field
type Field struct {
	assets.Field
}

// NewField creates a new field from the given asset
func NewField(asset assets.Field) *Field {
	return &Field{Field: asset}
}

// Asset returns the underlying asset
func (f *Field) Asset() assets.Field { return f.Field }

// Value represents a value in each of the field types
type Value struct {
	Text     types.XText      `json:"text" validate:"required"`
	Datetime *types.XDateTime `json:"datetime,omitempty"`
	Number   *types.XNumber   `json:"number,omitempty"`
	State    LocationPath     `json:"state,omitempty"`
	District LocationPath     `json:"district,omitempty"`
	Ward     LocationPath     `json:"ward,omitempty"`
}

// NewValue creates an empty value
func NewValue(text types.XText, datetime *types.XDateTime, number *types.XNumber, state LocationPath, district LocationPath, ward LocationPath) *Value {
	return &Value{
		Text:     text,
		Datetime: datetime,
		Number:   number,
		State:    state,
		District: district,
		Ward:     ward,
	}
}

// Equals determines whether two values are equal
func (v *Value) Equals(o *Value) bool {
	if v == nil && o == nil {
		return true
	}
	if (v == nil && o != nil) || (v != nil && o == nil) {
		return false
	}

	dateEqual := (v.Datetime == nil && o.Datetime == nil) || (v.Datetime != nil && o.Datetime != nil && v.Datetime.Equals(*o.Datetime))
	numEqual := (v.Number == nil && o.Number == nil) || (v.Number != nil && o.Number != nil && v.Number.Equals(*o.Number))

	return v.Text.Equals(o.Text) && dateEqual && numEqual && v.State == o.State && v.District == o.District && v.Ward == o.Ward
}

// FieldValue represents a field and a set of values for that field
type FieldValue struct {
	field *Field
	*Value
}

// NewFieldValue creates a new field value
func NewFieldValue(field *Field, value *Value) *FieldValue {
	return &FieldValue{field: field, Value: value}
}

// TypedValue returns the value in its proper type or nil if there is no value in that type
func (v *FieldValue) TypedValue() types.XValue {
	// the typed value of no value is nil
	if v == nil {
		return nil
	}

	switch v.field.Type() {
	case assets.FieldTypeText:
		return v.Text
	case assets.FieldTypeDatetime:
		if v.Datetime != nil {
			return *v.Datetime
		}
	case assets.FieldTypeNumber:
		if v.Number != nil {
			return *v.Number
		}
	case assets.FieldTypeState:
		if v.State != "" {
			return v.State
		}
	case assets.FieldTypeDistrict:
		if v.District != "" {
			return v.District
		}
	case assets.FieldTypeWard:
		if v.Ward != "" {
			return v.Ward
		}
	}
	return nil
}

// Resolve resolves the given key when this field value is referenced in an expression
func (v *FieldValue) Resolve(env utils.Environment, key string) types.XValue {
	switch strings.ToLower(key) {
	case "text":
		return v.Text
	}
	return types.NewXResolveError(v, key)
}

// Describe returns a representation of this type for error messages
func (v *FieldValue) Describe() string { return "field value" }

// Reduce is called when this object needs to be reduced to a primitive
func (v *FieldValue) Reduce(env utils.Environment) types.XPrimitive {
	return types.Reduce(env, v.TypedValue())
}

// ToXJSON is called when this type is passed to @(json(...))
func (v *FieldValue) ToXJSON(env utils.Environment) types.XText {
	j, _ := types.ToXJSON(env, v.Reduce(env))
	return j
}

var _ types.XValue = (*FieldValue)(nil)
var _ types.XResolvable = (*FieldValue)(nil)

// FieldValues is the set of all field values for a contact
type FieldValues map[string]*FieldValue

// NewFieldValues creates a new field value map
func NewFieldValues(a SessionAssets, values map[string]*Value, strict bool) (FieldValues, error) {
	allFields := a.Fields().All()
	fieldValues := make(FieldValues, len(allFields))
	for _, field := range allFields {
		value := values[field.Key()]
		if value != nil {
			if value.Text.Empty() {
				return nil, fmt.Errorf("field values can't be empty")
			}
			fieldValues[field.Key()] = NewFieldValue(field, value)
		} else {
			fieldValues[field.Key()] = nil
		}
	}

	if strict {
		for key := range values {
			_, valid := fieldValues[key]
			if !valid {
				return nil, fmt.Errorf("invalid field key: %s", key)
			}
		}
	}

	return fieldValues, nil
}

// Clone returns a clone of this set of field values
func (f FieldValues) clone() FieldValues {
	clone := make(FieldValues, len(f))
	for k, v := range f {
		clone[k] = v
	}
	return clone
}

// Get gets the field value set for the given field
func (f FieldValues) Get(field *Field) *Value {
	fieldVal := f[field.Key()]
	if fieldVal != nil {
		return fieldVal.Value
	}
	return nil
}

// Clear clears the field value set for the given field
func (f FieldValues) Clear(field *Field) {
	delete(f, field.Key())
}

// Set sets the field value set for the given field
func (f FieldValues) Set(env utils.Environment, field *Field, rawValue string, fields *FieldAssets) *Value {
	runEnv := env.(RunEnvironment)
	var value *Value

	// if raw value is empty string, set an empty value, other parse into different types
	if rawValue == "" {
		f.Clear(field)
		return nil
	}

	value = f.parseValue(runEnv, fields, field, rawValue)
	fieldValue := NewFieldValue(field, value)
	f[field.Key()] = fieldValue
	return fieldValue.Value
}

func (f FieldValues) parseValue(env RunEnvironment, fields *FieldAssets, field *Field, rawValue string) *Value {
	var asText = types.NewXText(rawValue)
	var asDateTime *types.XDateTime
	var asNumber *types.XNumber

	if parsedNumber, xerr := types.ToXNumber(env, asText); xerr == nil {
		asNumber = &parsedNumber
	}

	if parsedDate, xerr := types.ToXDateTimeWithTimeFill(env, asText); xerr == nil {
		asDateTime = &parsedDate
	}

	var asLocation *utils.Location

	// for locations, if it has a '>' then it is explicit, look it up that way
	if IsPossibleLocationPath(rawValue) {
		asLocation, _ = env.LookupLocation(LocationPath(rawValue))
	} else {
		var matchingLocations []*utils.Location

		if field.Type() == assets.FieldTypeWard {
			parent := f.getFirstLocationValue(env, fields, assets.FieldTypeDistrict)
			if parent != nil {
				matchingLocations, _ = env.FindLocationsFuzzy(rawValue, LocationLevelWard, parent)
			}
		} else if field.Type() == assets.FieldTypeDistrict {
			parent := f.getFirstLocationValue(env, fields, assets.FieldTypeState)
			if parent != nil {
				matchingLocations, _ = env.FindLocationsFuzzy(rawValue, LocationLevelDistrict, parent)
			}
		} else if field.Type() == assets.FieldTypeState {
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

	return &Value{
		Text:     asText,
		Datetime: asDateTime,
		Number:   asNumber,
		State:    asState,
		District: asDistrict,
		Ward:     asWard,
	}
}

func (f FieldValues) getFirstLocationValue(env RunEnvironment, fields *FieldAssets, valueType assets.FieldType) *utils.Location {
	// do we have a field of this type?
	field := fields.FirstOfType(valueType)
	if field == nil {
		return nil
	}
	// does this contact have a value for that field?
	value := f[field.Key()].TypedValue()
	if value == nil {
		return nil
	}

	location, err := env.LookupLocation(value.(LocationPath))
	if err != nil {
		return nil
	}
	return location
}

// Length is called to get the length of this object which in this case is the number of set values
func (f FieldValues) Length() int {
	count := 0
	for _, v := range f {
		if v != nil {
			count++
		}
	}
	return count
}

// Resolve resolves the given key when this set of field values is referenced in an expression
func (f FieldValues) Resolve(env utils.Environment, key string) types.XValue {
	val, exists := f[strings.ToLower(key)]
	if !exists {
		return types.NewXErrorf("no such contact field '%s'", key)
	}
	return val
}

// Describe returns a representation of this type for error messages
func (f FieldValues) Describe() string { return "field values" }

// Reduce is called when this object needs to be reduced to a primitive
func (f FieldValues) Reduce(env utils.Environment) types.XPrimitive {
	values := types.NewEmptyXMap()
	for k, v := range f {
		values.Put(string(k), v)
	}
	return values
}

// ToXJSON is called when this type is passed to @(json(...))
func (f FieldValues) ToXJSON(env utils.Environment) types.XText {
	return f.Reduce(env).ToXJSON(env)
}

var _ types.XValue = (FieldValues)(nil)
var _ types.XLengthable = (FieldValues)(nil)
var _ types.XResolvable = (FieldValues)(nil)

// FieldAssets provides access to all field assets
type FieldAssets struct {
	all   []*Field
	byKey map[string]*Field
}

// NewFieldAssets creates a new set of field assets
func NewFieldAssets(fields []assets.Field) *FieldAssets {
	s := &FieldAssets{
		all:   make([]*Field, len(fields)),
		byKey: make(map[string]*Field, len(fields)),
	}
	for f, asset := range fields {
		field := NewField(asset)
		s.all[f] = field
		s.byKey[field.Key()] = field
	}
	return s
}

// Get returns the contact field with the given key
func (s *FieldAssets) Get(key string) (*Field, error) {
	field, found := s.byKey[key]
	if !found {
		return nil, fmt.Errorf("no such field with key '%s'", key)
	}
	return field, nil
}

// All returns all the fields in this set
func (s *FieldAssets) All() []*Field {
	return s.all
}

// FirstOfType returns the first field in this set with the given value type
func (s *FieldAssets) FirstOfType(valueType assets.FieldType) *Field {
	for _, field := range s.all {
		if field.Type() == valueType {
			return field
		}
	}
	return nil
}
