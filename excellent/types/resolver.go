package types

import (
	"fmt"
	"reflect"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

// Resolvable is the interface for objects in the context which can be keyed into, e.g. foo.bar
type Resolvable interface {
	Resolve(key string) interface{}
}

// Lengthable is the interface for objects in the context which have a length
type Lengthable interface {
	Length() int
}

// Indexable is the interface for objects in the context which can be indexed into, e.g. foo.0. Such objects
// also need to be lengthable so that the engine knows what is a valid index and what isn't.
type Indexable interface {
	Lengthable

	Index(index int) interface{}
}

// Atomizable is the interface for objects in the context which can reduce themselves to an XAtom primitive
type Atomizable interface {
	Atomize() interface{}
}

type mapResolver struct {
	values map[string]interface{}
}

// NewMapResolver returns a simple resolver that resolves variables according to the values
// passed in
func NewMapResolver(values map[string]interface{}) Resolvable {
	return &mapResolver{
		values: values,
	}
}

// Resolve resolves the given key when this map is referenced in an expression
func (r *mapResolver) Resolve(key string) interface{} {
	val, found := r.values[key]
	if !found {
		return fmt.Errorf("no key '%s' in map", key)
	}
	return val
}

// Atomize is called when this object needs to be reduced to a primitive
func (r *mapResolver) Atomize() interface{} { return fmt.Sprintf("%s", r.values) }

var _ Atomizable = (*mapResolver)(nil)
var _ Resolvable = (*mapResolver)(nil)

// ToXAtom figures out the raw type of the passed in interface, returning that type
func ToXAtom(env utils.Environment, val interface{}) (interface{}, XType, error) {
	if val == nil {
		return val, XTypeNil, nil
	}

	switch val := val.(type) {
	case error:
		return val, XTypeError, nil

	case string:
		return val, XTypeString, nil

	case decimal.Decimal:
		return val, XTypeNumber, nil
	case int:
		return decimal.New(int64(val), 0), XTypeNumber, nil

	case time.Time:
		return val, XTypeTime, nil

	case bool:
		return val, XTypeBool, nil

	case Array:
		return val, XTypeArray, nil
	}

	return val, XTypeNil, fmt.Errorf("Unknown type '%s' with value '%+v'", reflect.TypeOf(val), val)
}
