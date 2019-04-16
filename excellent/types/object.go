package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// XObject is an object with named properties.
//
//   @(object("foo", 1, "bar", "x")) -> {bar: x, foo: 1}
//   @(object("foo", 1, "bar", "x").bar) -> x
//   @(object("foo", 1, "bar", "x")["bar"]) -> x
//   @(count(object("foo", 1, "bar", "x"))) -> 2
//   @(json(object("foo", 1, "bar", "x"))) -> {"bar":"x","foo":1}
//
// @type object
type XObject struct {
	XValue
	XCountable

	data   map[string]XValue
	source func() map[string]XValue
}

// NewXObject returns a new object with the given properties
func NewXObject(properties map[string]XValue) *XObject {
	return &XObject{
		data: properties,
	}
}

// NewXLazyObject returns a new lazy object with the source function
func NewXLazyObject(source func() map[string]XValue) *XObject {
	return &XObject{
		source: source,
	}
}

// Describe returns a representation of this type for error messages
func (x *XObject) Describe() string { return "object" }

// Truthy determines truthiness for this type
func (x *XObject) Truthy() bool {
	return x.Count() > 0
}

// Render returns the canonical text representation
func (x *XObject) Render() string {
	pairs := make([]string, 0, x.Count())
	for _, k := range x.keys(true) {
		rendered := Render(x.values()[k])
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, rendered))
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

// Format returns the pretty text representation
func (x *XObject) Format(env utils.Environment) string {
	pairs := make([]string, 0, x.Count())
	for _, k := range x.keys(true) {
		formatted := Format(env, x.values()[k])
		if strings.ContainsRune(formatted, '\n') {
			formatted = utils.Indent(formatted, "  ")
			formatted = fmt.Sprintf("%s:\n%s", k, formatted)
		} else {
			formatted = fmt.Sprintf("%s: %s", k, formatted)
		}
		pairs = append(pairs, formatted)
	}
	return strings.Join(pairs, "\n")
}

// MarshalJSON converts this type to internal JSON
func (x *XObject) MarshalJSON() ([]byte, error) {
	marshaled := make(map[string]json.RawMessage, x.Count())
	for k, v := range x.values() {
		asJSON, err := ToXJSON(v)
		if err == nil {
			marshaled[k] = json.RawMessage(asJSON.Native())
		}
	}
	return json.Marshal(marshaled)
}

// String returns the native string representation of this type for debugging
func (x *XObject) String() string {
	pairs := make([]string, 0, x.Count())
	for _, k := range x.keys(true) {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, String(x.values()[k])))
	}
	return "XObject{" + strings.Join(pairs, ", ") + "}"
}

// Count is called when the length of this object is requested in an expression
func (x *XObject) Count() int {
	return len(x.values())
}

// Get retrieves the named property
func (x *XObject) Get(key string) (XValue, bool) {
	key = strings.ToLower(key)
	for k, v := range x.values() {
		if strings.ToLower(k) == key {
			return v, true
		}
	}

	return nil, false
}

// Keys returns the properties of this object
func (x *XObject) Keys() []string {
	return x.keys(false)
}

// Equals determines equality for this type
func (x *XObject) Equals(other *XObject) bool {
	keys1 := x.keys(true)
	keys2 := other.keys(true)

	if len(keys1) != len(keys2) {
		return false
	}

	for k, key := range keys1 {
		if key != keys2[k] {
			return false
		}

		if !Equals(x.values()[key], other.values()[key]) {
			return false
		}
	}

	return true
}

func (x *XObject) keys(sorted bool) []string {
	keys := make([]string, 0, x.Count())
	for key := range x.values() {
		keys = append(keys, key)
	}
	if sorted {
		sort.Strings(keys)
	}
	return keys
}

func (x *XObject) values() map[string]XValue {
	if x.data == nil {
		x.data = x.source()
	}
	return x.data
}

// XObjectEmpty is the empty empty
var XObjectEmpty = NewXObject(map[string]XValue{})

var _ json.Marshaler = (*XObject)(nil)

// ToXObject converts the given value to an object
func ToXObject(env utils.Environment, x XValue) (*XObject, XError) {
	if utils.IsNil(x) {
		return XObjectEmpty, nil
	}
	if IsXError(x) {
		return XObjectEmpty, x.(XError)
	}

	object, isObject := x.(*XObject)
	if isObject && object != nil {
		return object, nil
	}

	return XObjectEmpty, NewXErrorf("unable to convert %s to an object", Describe(x))
}
