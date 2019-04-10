package types

import (
	"encoding/json"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// XArray is an array primitive in Excellent expressions
type XArray struct {
	XValue

	data   []XValue
	source func() []XValue
}

// NewXArray returns a new array with the given items
func NewXArray(data ...XValue) *XArray {
	if data == nil {
		data = []XValue{}
	}
	return &XArray{data: data}
}

// NewXLazyArray returns a new lazy array with the given source function
func NewXLazyArray(source func() []XValue) *XArray {
	return &XArray{
		source: source,
	}
}

// Describe returns a representation of this type for error messages
func (x *XArray) Describe() string { return "array" }

// ToXText converts this type to text
func (x *XArray) ToXText(env utils.Environment) XText {
	parts := make([]string, x.Length())
	for i, v := range x.values() {
		vAsText, xerr := ToXText(env, v)
		if xerr != nil {
			vAsText = xerr.ToXText(env)
		}
		parts[i] = vAsText.Native()
	}
	return NewXText("[" + strings.Join(parts, ", ") + "]")
}

// ToXBoolean converts this type to a bool
func (x *XArray) ToXBoolean() XBoolean {
	return NewXBoolean(len(x.values()) > 0)
}

// ToXJSON is called when this type is passed to @(json(...))
func (x *XArray) ToXJSON() XText {
	marshaled := make([]json.RawMessage, len(x.values()))
	for i, v := range x.values() {
		asJSON, err := ToXJSON(v)
		if err == nil {
			marshaled[i] = json.RawMessage(asJSON.Native())
		}
	}
	return MustMarshalToXText(marshaled)
}

// MarshalJSON converts this type to internal JSON
func (x *XArray) MarshalJSON() ([]byte, error) {
	return utils.JSONMarshal(x.values())
}

// Get is called when this object is indexed
func (x *XArray) Get(index int) XValue {
	return x.values()[index]
}

// Length is called when the length of this object is requested in an expression
func (x *XArray) Length() int {
	return len(x.values())
}

// String returns the native string representation of this type
func (x *XArray) String() string {
	parts := make([]string, x.Length())
	for i, v := range x.values() {
		parts[i] = String(v)
	}
	return `XArray[` + strings.Join(parts, ", ") + `]`
}

// Equals determines equality for this type
func (x *XArray) Equals(other *XArray) bool {
	if x.Length() != other.Length() {
		return false
	}

	for i, v := range x.values() {
		if !Equals(v, other.values()[i]) {
			return false
		}
	}
	return true
}

func (x *XArray) values() []XValue {
	if x.data == nil {
		x.data = x.source()
	}
	return x.data
}

// XArrayEmpty is the empty array
var XArrayEmpty = NewXArray()

var _ json.Marshaler = (*XArray)(nil)

// ToXArray converts the given value to an array
func ToXArray(env utils.Environment, x XValue) (*XArray, XError) {
	if utils.IsNil(x) {
		return XArrayEmpty, nil
	}
	if IsXError(x) {
		return XArrayEmpty, x.(XError)
	}

	asArray, isArray := x.(*XArray)
	if isArray {
		return asArray, nil
	}

	return XArrayEmpty, NewXErrorf("unable to convert %s to an array", Describe(x))
}
