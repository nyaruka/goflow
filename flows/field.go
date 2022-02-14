package flows

import (
	"fmt"
	"sort"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
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

// Reference returns a reference to this field
func (f *Field) Reference() *assets.FieldReference {
	if f == nil {
		return nil
	}
	return assets.NewFieldReference(f.Key(), f.Name())
}

// Value represents a value in each of the field types
type Value struct {
	Text     types.XText       `json:"text" validate:"required"`
	Datetime *types.XDateTime  `json:"datetime,omitempty"`
	Number   *types.XNumber    `json:"number,omitempty"`
	State    envs.LocationPath `json:"state,omitempty"`
	District envs.LocationPath `json:"district,omitempty"`
	Ward     envs.LocationPath `json:"ward,omitempty"`
}

// NewValue creates an empty value
func NewValue(text types.XText, datetime *types.XDateTime, number *types.XNumber, state envs.LocationPath, district envs.LocationPath, ward envs.LocationPath) *Value {
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

// ToXValue returns a representation of this object for use in expressions
func (v *FieldValue) ToXValue(env envs.Environment) types.XValue {
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
			return types.NewXText(string(v.State))
		}
	case assets.FieldTypeDistrict:
		if v.District != "" {
			return types.NewXText(string(v.District))
		}
	case assets.FieldTypeWard:
		if v.Ward != "" {
			return types.NewXText(string(v.Ward))
		}
	}
	return nil
}

// QueryValue returns the value for use in contact queries
func (v *FieldValue) QueryValue() interface{} {
	// the typed value of no value is nil
	if v == nil {
		return nil
	}

	switch v.field.Type() {
	case assets.FieldTypeText:
		return v.Text.Native()
	case assets.FieldTypeDatetime:
		if v.Datetime != nil {
			return (*v.Datetime).Native()
		}
	case assets.FieldTypeNumber:
		if v.Number != nil {
			return (*v.Number).Native()
		}

	// we only search against location names and not full paths
	case assets.FieldTypeState:
		if v.State != "" {
			return v.State.Name()
		}
	case assets.FieldTypeDistrict:
		if v.District != "" {
			return v.District.Name()
		}
	case assets.FieldTypeWard:
		if v.Ward != "" {
			return v.Ward.Name()
		}
	}
	return nil
}

// FieldValues is the set of all field values for a contact
type FieldValues map[string]*FieldValue

// NewFieldValues creates a new field value map
func NewFieldValues(a SessionAssets, values map[string]*Value, missing assets.MissingCallback) FieldValues {
	allFields := a.Fields().All()
	fieldValues := make(FieldValues, len(allFields))
	for _, field := range allFields {
		value := values[field.Key()]
		fieldValues.Set(field, value)
	}

	// log any unmatched field keys as missing assets
	for key := range values {
		_, valid := fieldValues[key]
		if !valid {
			missing(assets.NewFieldReference(key, ""), nil)
		}
	}

	return fieldValues
}

// Clone returns a clone of this set of field values
func (f FieldValues) clone() FieldValues {
	clone := make(FieldValues, len(f))
	for k, v := range f {
		clone[k] = v
	}
	return clone
}

// Get gets the value set for the given field
func (f FieldValues) Get(field *Field) *Value {
	fieldVal := f[field.Key()]
	if fieldVal != nil {
		return fieldVal.Value
	}
	return nil
}

// Set sets the value for the given field (can be null to clear it)
func (f FieldValues) Set(field *Field, value *Value) {
	var fv *FieldValue
	if value != nil && !value.Text.Empty() {
		fv = NewFieldValue(field, value)
	}
	f[field.Key()] = fv
}

// Parse parses a raw string field value into the different possible types
func (f FieldValues) Parse(env envs.Environment, fields *FieldAssets, field *Field, rawValue string) *Value {
	if rawValue == "" {
		return nil
	}

	var asText = types.NewXText(rawValue)
	var asDateTime *types.XDateTime
	var asNumber *types.XNumber

	if parsedNumber, xerr := types.ToXNumber(env, asText); xerr == nil {
		asNumber = &parsedNumber
	}

	if parsedDate, xerr := types.ToXDateTimeWithTimeFill(env, asText); xerr == nil {
		asDateTime = &parsedDate
	}

	var asLocation *envs.Location

	// to do location parsing we need an environment with location support
	locations := env.LocationResolver()

	if locations != nil {
		// for locations, if it has a '>' then it is explicit, look it up that way
		if envs.IsPossibleLocationPath(rawValue) {
			asLocation = locations.LookupLocation(envs.LocationPath(rawValue))
		} else {
			var matchingLocations []*envs.Location

			if field.Type() == assets.FieldTypeWard {
				parent := f.getFirstLocationValue(env, fields, assets.FieldTypeDistrict)
				if parent != nil {
					matchingLocations = locations.FindLocationsFuzzy(rawValue, LocationLevelWard, parent)
				}
			} else if field.Type() == assets.FieldTypeDistrict {
				parent := f.getFirstLocationValue(env, fields, assets.FieldTypeState)
				if parent != nil {
					matchingLocations = locations.FindLocationsFuzzy(rawValue, LocationLevelDistrict, parent)
				}
			} else if field.Type() == assets.FieldTypeState {
				matchingLocations = locations.FindLocationsFuzzy(rawValue, LocationLevelState, nil)
			}

			if len(matchingLocations) > 0 {
				asLocation = matchingLocations[0]
			}
		}
	}

	var asState, asDistrict, asWard envs.LocationPath
	if asLocation != nil {
		switch asLocation.Level() {
		case LocationLevelState:
			asState = envs.LocationPath(asLocation.Path())
		case LocationLevelDistrict:
			asState = envs.LocationPath(asLocation.Parent().Path())
			asDistrict = envs.LocationPath(asLocation.Path())
		case LocationLevelWard:
			asState = envs.LocationPath(asLocation.Parent().Parent().Path())
			asDistrict = envs.LocationPath(asLocation.Parent().Path())
			asWard = envs.LocationPath(asLocation.Path())
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

// Context returns the properties available in expressions
func (f FieldValues) Context(env envs.Environment) map[string]types.XValue {
	entries := make(map[string]types.XValue, len(f)+1)
	lines := make([]string, 0, len(f))

	for k, v := range f {
		val := v.ToXValue(env)
		entries[string(k)] = val

		if !utils.IsNil(val) {
			lines = append(lines, fmt.Sprintf("%s: %s", v.field.Name(), types.Render(val)))
		}
	}

	sort.Strings(lines)
	entries["__default__"] = types.NewXText(strings.Join(lines, "\n"))

	return entries
}

func (f FieldValues) getFirstLocationValue(env envs.Environment, fields *FieldAssets, valueType assets.FieldType) *envs.Location {
	// do we have a field of this type?
	field := fields.FirstOfType(valueType)
	if field == nil {
		return nil
	}
	// does this contact have a value for that field?
	value := f[field.Key()].ToXValue(env)
	if value == nil {
		return nil
	}

	return env.LocationResolver().LookupLocation(envs.LocationPath(value.(types.XText).Native()))
}

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
	for i, asset := range fields {
		field := NewField(asset)
		s.all[i] = field
		s.byKey[field.Key()] = field
	}
	return s
}

// Get returns the contact field with the given key
func (s *FieldAssets) Get(key string) *Field {
	return s.byKey[key]
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

// Query based groups don't allow conditions based on flows or other groups and therefore
// can be validated with just field assets. Thus we allow FieldAssets to operate as a
// contactql.Resolver for the purposes of loading groups.

func (s *FieldAssets) ResolveField(key string) assets.Field {
	f := s.byKey[key]
	if f != nil {
		return f
	}
	return nil
}

func (s *FieldAssets) ResolveGroup(name string) assets.Group {
	return nil
}

func (s *FieldAssets) ResolveFlow(name string) assets.Flow {
	return nil
}
