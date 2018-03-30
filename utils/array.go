package utils

import (
	"encoding/json"
)

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
