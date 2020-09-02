package types

import (
	"encoding/json"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

// XArray is an array of items.
//
//   @(array(1, "x", true)) -> [1, x, true]
//   @(array(1, "x", true)[1]) -> x
//   @(count(array(1, "x", true))) -> 3
//   @(json(array(1, "x", true))) -> [1,"x",true]
//
// @type array
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

// Get is called when this object is indexed
func (x *XArray) Get(index int) XValue {
	return x.values()[index]
}

// Count is called when the length of this object is requested in an expression
func (x *XArray) Count() int {
	return len(x.values())
}

// Describe returns a representation of this type for error messages
func (x *XArray) Describe() string { return "array" }

// Truthy determines truthiness for this type
func (x *XArray) Truthy() bool {
	return x.Count() > 0
}

// Render returns the canonical text representation
func (x *XArray) Render() string {
	parts := make([]string, x.Count())
	for i, v := range x.values() {
		parts[i] = Render(v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// Format returns the pretty text representation
func (x *XArray) Format(env envs.Environment) string {
	parts := make([]string, x.Count())
	multiline := false

	for i, v := range x.values() {
		parts[i] = Format(env, v)
		if strings.ContainsRune(parts[i], '\n') {
			multiline = true
		}
	}

	if multiline {
		for i, p := range parts {
			p = utils.Indent(p, "  ")
			parts[i] = "-" + p[1:]
		}

		return strings.Join(parts, "\n")
	}

	return strings.Join(parts, ", ")
}

// MarshalJSON converts this type to internal JSON
func (x *XArray) MarshalJSON() ([]byte, error) {
	marshaled := make([]json.RawMessage, x.Count())
	for i, v := range x.values() {
		asJSON, err := ToXJSON(v)
		if err == nil {
			marshaled[i] = json.RawMessage(asJSON.Native())
		}
	}
	return jsonx.Marshal(marshaled)
}

// String returns the native string representation of this type
func (x *XArray) String() string {
	parts := make([]string, x.Count())
	for i, v := range x.values() {
		parts[i] = String(v)
	}
	return `XArray[` + strings.Join(parts, ", ") + `]`
}

// Equals determines equality for this type
func (x *XArray) Equals(other *XArray) bool {
	if x.Count() != other.Count() {
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

// ToXArray converts the given value to an array
func ToXArray(env envs.Environment, x XValue) (*XArray, XError) {
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
