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
	baseXPrimitive

	values []XValue
}

// NewXArray returns a new array with the given items
func NewXArray(values ...XValue) XArray {
	if values == nil {
		values = []XValue{}
	}
	return &xarray{values: values}
}

// Reduce returns the primitive version of this type (i.e. itself)
func (a *xarray) Reduce() XPrimitive { return a }

// ToXText converts this type to text
func (a *xarray) ToXText() XText {
	strs := make([]XText, len(a.values))
	for i := range a.values {
		strs[i] = a.values[i].Reduce().ToXText()
	}
	return MustMarshalToXText(strs)
}

// ToXBoolean converts this type to a bool
func (a *xarray) ToXBoolean() XBoolean {
	return NewXBoolean(len(a.values) > 0)
}

// ToXJSON is called when this type is passed to @(json(...))
func (a *xarray) ToXJSON() XText {
	marshaled := make([]json.RawMessage, len(a.values))
	for i, v := range a.values {
		asJSON, err := ToXJSON(v)
		if err == nil {
			marshaled[i] = json.RawMessage(asJSON.Native())
		}
	}
	return MustMarshalToXText(marshaled)
}

// MarshalJSON converts this type to internal JSON
func (a *xarray) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.values)
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
var _ json.Marshaler = (*xarray)(nil)
