package types

import (
	"encoding/json"
)

// XArray is an array primitive in Excellent expressions
type XArray interface {
	XPrimitive
	XIndexable

	Append(XValue)
}

type xarray struct {
	values []XValue
}

// NewXArray returns a new array with the given items
func NewXArray(values ...XValue) XArray {
	return &xarray{values: values}
}

// ToString converts this type to a string
func (a *xarray) ToString() XString {
	strs := make([]string, len(a.values))
	for i := range a.values {
		strs[i] = string(ToXString(a.values[i]))
	}
	return RequireMarshalToXString(strs)
}

// ToBool converts this type to a bool
func (a *xarray) ToBool() XBool {
	return len(a.values) > 0
}

// ToJSON converts this type to JSON
func (a *xarray) ToJSON() XString {
	marshaled := make([]json.RawMessage, len(a.values))
	for i := range a.values {
		marshaled[i] = json.RawMessage(a.values[i].ToJSON())
	}
	return RequireMarshalToXString(marshaled)
}

// Index is called when this object is indexed into in an expression
func (a *xarray) Index(index int) XValue {
	return a.values[index]
}

// Length is called when the length of this object is requested in an expression
func (a *xarray) Length() int {
	return len(a.values)
}

// Append adds the given item to this array
func (a *xarray) Append(value XValue) {
	a.values = append(a.values, value)
}

var _ XArray = (*xarray)(nil)

// Legacy...

// Array is a simple data structure which is indexable in expressions
type Array interface {
	Indexable

	Append(interface{})
}

type array struct {
	values []interface{}
}

// NewArray returns a new array with the given items
func NewArray(values ...interface{}) Array {
	return &array{values: values}
}

// Index is called when this object is indexed into in an expression
func (a *array) Index(index int) interface{} {
	return a.values[index]
}

// Length is called when the length of this object is requested in an expression
func (a *array) Length() int {
	return len(a.values)
}

// Append adds the given item to this array
func (a *array) Append(value interface{}) {
	a.values = append(a.values, value)
}

func (a *array) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.values)
}

var _ Indexable = (*array)(nil)
