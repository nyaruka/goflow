package types

import (
	"encoding/json"
)

// XMap is a map primitive in Excellent expressions
type XMap interface {
	XPrimitive
	XResolvable
	XLengthable
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

// ToString converts this type to a string
func (m *xmap) ToString() XString {
	strs := make(map[string]string, len(m.values))
	for k, v := range m.values {
		strs[k] = string(ToXString(v))
	}
	return RequireMarshalToXString(strs)
}

// ToBool converts this type to a bool
func (m *xmap) ToBool() XBool {
	return len(m.values) > 0
}

// ToJSON converts this type to JSON
func (m *xmap) ToJSON() XString {
	marshaled := make(map[string]json.RawMessage, len(m.values))
	for k, v := range m.values {
		marshaled[k] = json.RawMessage(v.ToJSON())
	}
	return RequireMarshalToXString(marshaled)
}

// Length is called when the length of this object is requested in an expression
func (m *xmap) Length() int {
	return len(m.values)
}

func (m *xmap) Resolve(key string) XValue {
	val, found := m.values[key]
	if !found {
		return NewXResolveError(m, key)
	}
	return val
}

var _ XMap = (*xmap)(nil)
