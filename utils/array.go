package utils

import (
	"encoding/json"
)

type Array interface {
	VariableIndexer

	Append(interface{})
}

type array struct {
	values []interface{}
}

func NewArray(values ...interface{}) Array {
	return &array{values: values}
}

func (a *array) Index(index int) interface{} {
	return a.values[index]
}

func (a *array) Length() int {
	return len(a.values)
}

func (a *array) Append(value interface{}) {
	a.values = append(a.values, value)
}

func (a *array) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.values)
}

var _ VariableIndexer = (*array)(nil)
