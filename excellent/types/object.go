package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

const serializeDefaultAs = "__default__"

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

	def            XValue
	props          map[string]XValue
	source         func() map[string]XValue
	marshalDefault bool
}

// NewXObject returns a new object with the given properties
func NewXObject(properties map[string]XValue) *XObject {
	return NewXLazyObject(func() map[string]XValue { return properties })
}

// NewXLazyObject returns a new lazy object with the source function and default
func NewXLazyObject(source func() map[string]XValue) *XObject {
	return &XObject{
		source: source,
	}
}

// Describe returns a representation of this type for error messages
func (x *XObject) Describe() string { return "object" }

// Truthy determines truthiness for this type
func (x *XObject) Truthy() bool {
	if x.hasDefault() {
		return Truthy(x.Default())
	}

	return x.Count() > 0
}

// Render returns the canonical text representation
func (x *XObject) Render() string {
	if x.hasDefault() {
		return Render(x.Default())
	}

	pairs := make([]string, 0, x.Count())
	for _, k := range x.Properties() {
		rendered := Render(x.properties()[k])
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, rendered))
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

// Format returns the pretty text representation
func (x *XObject) Format(env envs.Environment) string {
	if x.hasDefault() {
		return Format(env, x.Default())
	}

	pairs := make([]string, 0, x.Count())
	for _, k := range x.Properties() {
		formatted := Format(env, x.properties()[k])
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
	for p, v := range x.properties() {
		asJSON, err := ToXJSON(v)
		if err == nil {
			marshaled[p] = json.RawMessage(asJSON.Native())
		}
	}

	if x.hasDefault() && x.marshalDefault {
		asJSON, err := ToXJSON(x.def)
		if err == nil {
			marshaled[serializeDefaultAs] = json.RawMessage(asJSON.Native())
		}
	}

	return jsonx.Marshal(marshaled)
}

// ReadXObject reads an instance of this type from JSON
func ReadXObject(data []byte) (*XObject, error) {
	v := JSONToXValue(data)
	switch typed := v.(type) {
	case *XObject:
		return typed, nil
	case XError:
		return nil, typed
	default:
		return nil, errors.New("JSON doesn't contain an object")
	}
}

// String returns the native string representation of this type for debugging
func (x *XObject) String() string {
	pairs := make([]string, 0, x.Count())

	if x.hasDefault() {
		pairs = append(pairs, fmt.Sprintf("%s: %s", serializeDefaultAs, String(x.Default())))
	}

	for _, k := range x.Properties() {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, String(x.properties()[k])))
	}
	return "XObject{" + strings.Join(pairs, ", ") + "}"
}

// Count is called when the length of this object is requested in an expression
func (x *XObject) Count() int {
	return len(x.properties())
}

// Get retrieves the named property
func (x *XObject) Get(key string) (XValue, bool) {
	key = strings.ToLower(key)
	for p, v := range x.properties() {
		if strings.ToLower(p) == key {
			return v, true
		}
	}

	return nil, false
}

// Properties returns the sorted property names of this object
func (x *XObject) Properties() []string {
	names := make([]string, 0, x.Count())
	for name := range x.properties() {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Equals determines equality for this type
func (x *XObject) Equals(other *XObject) bool {
	if x.hasDefault() || other.hasDefault() {
		if !Equals(x.Default(), other.Default()) {
			return false
		}
	}

	props1 := x.Properties()
	props2 := other.Properties()

	if len(props1) != len(props2) {
		return false
	}

	for p, name := range props1 {
		if name != props2[p] {
			return false
		}

		if !Equals(x.properties()[name], other.properties()[name]) {
			return false
		}
	}

	return true
}

func (x *XObject) properties() map[string]XValue {
	x.ensureInitialized()
	return x.props
}

// Default returns the default value for this
func (x *XObject) Default() XValue {
	x.ensureInitialized()
	return x.def
}

func (x *XObject) SetMarshalDefault(marshal bool) {
	x.marshalDefault = marshal
}

// Default returns the default value for this
func (x *XObject) hasDefault() bool {
	return x.Default() != x
}

func (x *XObject) ensureInitialized() {
	if x.props == nil {
		props := x.source()

		x.def = x
		x.props = make(map[string]XValue, len(props))
		for p, v := range props {
			if p == serializeDefaultAs {
				x.def = v
			} else {
				x.props[p] = v
			}
		}
	}
}

// XObjectEmpty is the empty empty
var XObjectEmpty = NewXObject(map[string]XValue{})

var _ json.Marshaler = (*XObject)(nil)

// ToXObject converts the given value to an object
func ToXObject(env envs.Environment, x XValue) (*XObject, XError) {
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
