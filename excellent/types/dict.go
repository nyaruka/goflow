package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// XDict is a map primitive in Excellent expressions
type XDict struct {
	XValue
	XLengthable

	values map[string]XValue
}

// NewXDict returns a new map with the given items
func NewXDict(values map[string]XValue) *XDict {
	return &XDict{
		values: values,
	}
}

// Describe returns a representation of this type for error messages
func (x *XDict) Describe(env utils.Environment) string { return "dict" }

// ToXText converts this type to text
func (x *XDict) ToXText(env utils.Environment) XText {
	// get our keys sorted A-Z
	sortedKeys := x.Keys()
	sort.Strings(sortedKeys)

	pairs := make([]string, 0, x.Length())
	for _, k := range sortedKeys {
		vAsText, xerr := ToXText(env, x.values[k])
		if xerr != nil {
			vAsText = xerr.ToXText(env)
		}

		pairs = append(pairs, fmt.Sprintf("%s: %s", k, vAsText))
	}
	return NewXText("{" + strings.Join(pairs, ", ") + "}")
}

// ToXBoolean converts this type to a bool
func (x *XDict) ToXBoolean(env utils.Environment) XBoolean {
	return NewXBoolean(len(x.values) > 0)
}

// ToXJSON is called when this type is passed to @(json(...))
func (x *XDict) ToXJSON(env utils.Environment) XText {
	marshaled := make(map[string]json.RawMessage, len(x.values))
	for k, v := range x.values {
		asJSON, err := ToXJSON(env, v)
		if err == nil {
			marshaled[k] = json.RawMessage(asJSON.Native())
		}
	}
	return MustMarshalToXText(marshaled)
}

// MarshalJSON converts this type to internal JSON
func (x *XDict) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.values)
}

// Length is called when the length of this object is requested in an expression
func (x *XDict) Length() int {
	return len(x.values)
}

// Get retrieves the named item from this dict
func (x *XDict) Get(key string) (XValue, bool) {
	key = strings.ToLower(key)
	for k, v := range x.values {
		if strings.ToLower(k) == key {
			return v, true
		}
	}

	return nil, false
}

// Keys returns the keys of this dict
func (x *XDict) Keys() []string {
	keys := make([]string, 0, len(x.values))
	for key := range x.values {
		keys = append(keys, key)
	}
	return keys
}

// String returns the native string representation of this type
func (x *XDict) String() string { return x.ToXText(nil).Native() }

// XDictEmpty is the empty dict
var XDictEmpty = NewXDict(map[string]XValue{})

var _ json.Marshaler = (*XDict)(nil)

// ToXDict converts the given value to a dict
func ToXDict(env utils.Environment, x XValue) (*XDict, XError) {
	if utils.IsNil(x) {
		return XDictEmpty, nil
	}
	if IsXError(x) {
		return XDictEmpty, x.(XError)
	}

	asDict, isDict := x.(*XDict)
	if isDict {
		return asDict, nil
	}

	return XDictEmpty, NewXErrorf("unable to convert %s to a dict", Describe(env, x))
}
