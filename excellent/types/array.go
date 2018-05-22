package types

import (
	"encoding/json"

	"github.com/nyaruka/goflow/utils"
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
	if values == nil {
		values = []XValue{}
	}
	return &xarray{values: values}
}

// Describe returns a representation of this type for error messages
func (a *xarray) Describe() string { return "array" }

// Reduce returns the primitive version of this type (i.e. itself)
func (a *xarray) Reduce() XPrimitive { return a }

// ToXText converts this type to text
func (a *xarray) ToXText() XText {
	texts := make([]XText, len(a.values))
	for i, v := range a.values {
		asText, err := ToXText(v)
		if err == nil {
			texts[i] = asText
		}
	}
	return MustMarshalToXText(texts)
}

// ToXBoolean converts this type to a bool
func (a *xarray) ToXBoolean() XBoolean {
	return NewXBoolean(len(a.values) > 0)
}

// ToXJSON is called when this type is passed to @(json(...))
func (a *xarray) ToXJSON(env utils.Environment) XText {
	marshaled := make([]json.RawMessage, len(a.values))
	for i, v := range a.values {
		asJSON, err := ToXJSON(env, v)
		if err == nil {
			marshaled[i] = json.RawMessage(asJSON.Native())
		}
	}
	return MustMarshalToXText(marshaled)
}

// MarshalJSON converts this type to internal JSON
func (a *xarray) MarshalJSON() ([]byte, error) {
	return utils.JSONMarshal(a.values)
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

// String returns the native string representation of this type
func (a *xarray) String() string { return a.ToXText().Native() }

var _ XArray = (*xarray)(nil)
var _ json.Marshaler = (*xarray)(nil)
