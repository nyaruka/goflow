package types

import (
	"encoding/json"
	"strings"

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
func (a *xarray) Reduce(env utils.Environment) XPrimitive { return a }

// ToXText converts this type to text
func (a *xarray) ToXText(env utils.Environment) XText {
	parts := make([]string, a.Length())
	for i, v := range a.values {
		vAsText, xerr := ToXText(env, v)
		if xerr != nil {
			vAsText = xerr.ToXText(env)
		}
		parts[i] = vAsText.Native()
	}
	return NewXText("[" + strings.Join(parts, ", ") + "]")
}

// ToXBoolean converts this type to a bool
func (a *xarray) ToXBoolean(env utils.Environment) XBoolean {
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
func (a *xarray) String() string { return a.ToXText(nil).Native() }

var XArrayEmpty = NewXArray()
var _ XArray = (*xarray)(nil)
var _ json.Marshaler = (*xarray)(nil)

// ToXArray converts the given value to an array
func ToXArray(env utils.Environment, x XValue) (XArray, XError) {
	x = Reduce(env, x)

	if utils.IsNil(x) {
		return XArrayEmpty, nil
	}
	if IsXError(x) {
		return XArrayEmpty, x.(XError)
	}

	asArray, isArray := x.(XArray)
	if isArray {
		return asArray, nil
	}

	return XArrayEmpty, NewXErrorf("unable to convert %s to an array", Describe(x))
}
