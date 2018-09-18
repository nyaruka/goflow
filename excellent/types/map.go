package types

import (
	"encoding/json"

	"github.com/nyaruka/goflow/utils"
)

// XMap is a map primitive in Excellent expressions
type XMap interface {
	XPrimitive
	XResolvable
	XLengthable

	Put(string, XValue)
	Keys() []string
}

type xmap struct {
	values map[string]XValue
}

// NewXMap returns a new map with the given items
func NewXMap(values map[string]XValue) XMap {
	return &xmap{
		values: values,
	}
}

// NewEmptyXMap returns a new empty map
func NewEmptyXMap() XMap {
	return &xmap{
		values: make(map[string]XValue),
	}
}

// Describe returns a representation of this type for error messages
func (m *xmap) Describe() string { return "map" }

// Reduce returns the primitive version of this type (i.e. itself)
func (m *xmap) Reduce(env utils.Environment) XPrimitive { return m }

// ToXText converts this type to text
func (m *xmap) ToXText(env utils.Environment) XText {
	primitives := make(map[string]XValue, len(m.values))
	for k, v := range m.values {
		primitives[k] = Reduce(env, v)
	}
	return MustMarshalToXText(primitives)
}

// ToXBoolean converts this type to a bool
func (m *xmap) ToXBoolean(env utils.Environment) XBoolean {
	return NewXBoolean(len(m.values) > 0)
}

// ToXJSON is called when this type is passed to @(json(...))
func (m *xmap) ToXJSON(env utils.Environment) XText {
	marshaled := make(map[string]json.RawMessage, len(m.values))
	for k, v := range m.values {
		asJSON, err := ToXJSON(env, v)
		if err == nil {
			marshaled[k] = json.RawMessage(asJSON.Native())
		}
	}
	return MustMarshalToXText(marshaled)
}

// MarshalJSON converts this type to internal JSON
func (m *xmap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.values)
}

// Length is called when the length of this object is requested in an expression
func (m *xmap) Length() int {
	return len(m.values)
}

func (m *xmap) Resolve(env utils.Environment, key string) XValue {
	val, found := m.values[key]
	if !found {
		return NewXResolveError(m, key)
	}
	return val
}

// Put adds the given item to this map
func (m *xmap) Put(key string, value XValue) {
	m.values[key] = value
}

// Keys returns the keys of this map
func (m *xmap) Keys() []string {
	keys := make([]string, 0, len(m.values))
	for key := range m.values {
		keys = append(keys, key)
	}
	return keys
}

// String returns the native string representation of this type
func (m *xmap) String() string { return m.ToXText(nil).Native() }

var _ XMap = (*xmap)(nil)
var _ json.Marshaler = (*xmap)(nil)
