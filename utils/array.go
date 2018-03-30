package utils

type Array struct {
	values []interface{}
}

func NewArray(values []interface{}) *Array {
	return &Array{values: values}
}

func (a *Array) Index(index int) interface{} {
	return a.values[index]
}

func (a *Array) Length() int {
	return len(a.values)
}

var _ VariableIndexer = (*Array)(nil)
