package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// XValue is the base interface of all Excellent types
type XValue interface {
	ToXJSON() XString
	Reduce() XPrimitive
}

// XPrimitive is the base interface of all Excellent primitive types
type XPrimitive interface {
	XValue

	ToXString() XString
	ToXBool() XBool
}

// XResolvable is the interface for types which can be keyed into, e.g. foo.bar
type XResolvable interface {
	Resolve(key string) XValue
}

// XLengthable is the interface for types which have a length
type XLengthable interface {
	Length() int
}

// XIndexable is the interface for types which can be indexed into, e.g. foo.0. Such objects
// also need to be lengthable so that the engine knows what is a valid index and what isn't.
type XIndexable interface {
	XLengthable

	Index(index int) XValue
}

type baseXPrimitive struct {
	XPrimitive
}

func (x *baseXPrimitive) String() string {
	return x.ToXString().Native()
}

// ResolveKeys is a utility function that resolves multiple keys on an XResolvable and returns the results as a map
func ResolveKeys(resolvable XResolvable, keys ...string) XMap {
	values := make(map[string]XValue, len(keys))
	for _, key := range keys {
		values[key] = resolvable.Resolve(key)
	}
	return NewXMap(values)
}

// Compare returns the difference between two given values
func Compare(x1 XValue, x2 XValue) (int, error) {
	if utils.IsNil(x1) && utils.IsNil(x2) {
		return 0, nil
	} else if utils.IsNil(x1) || utils.IsNil(x2) {
		return 0, fmt.Errorf("can't compare non-nil and nil values: %T{%s} and %T{%s}", x1, x1, x2, x2)
	}

	x1 = x1.Reduce()
	x2 = x2.Reduce()

	if reflect.TypeOf(x1) != reflect.TypeOf(x2) {
		return 0, fmt.Errorf("can't compare different types of %T and %T", x1, x2)
	}

	// common types, do real comparisons
	switch typed := x1.(type) {
	case XError:
		return strings.Compare(typed.Error(), x2.(error).Error()), nil
	case XNumber:
		return typed.Compare(x2.(XNumber)), nil
	case XBool:
		return typed.Compare(x2.(XBool)), nil
	case XDate:
		return typed.Compare(x2.(XDate)), nil
	case XString:
		return typed.Compare(x2.(XString)), nil
	}

	// TODO: find better fallback
	return strings.Compare(x1.ToXJSON().Native(), x2.ToXJSON().Native()), nil
}
