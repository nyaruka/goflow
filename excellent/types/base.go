package types

import (
	"fmt"
	"reflect"

	"github.com/nyaruka/goflow/utils"
)

// XValue is the base interface of all Excellent types
type XValue interface {
	Describe() string
	ToXJSON() XText
	Reduce() XPrimitive
}

// XPrimitive is the base interface of all Excellent primitive types
type XPrimitive interface {
	XValue
	fmt.Stringer

	ToXText() XText
	ToXBoolean() XBoolean
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

// ResolveKeys is a utility function that resolves multiple keys on an XResolvable and returns the results as a map
func ResolveKeys(resolvable XResolvable, keys ...string) XMap {
	values := make(map[string]XValue, len(keys))
	for _, key := range keys {
		values[key] = resolvable.Resolve(key)
	}
	return NewXMap(values)
}

// Equals checks for equality between the two give values
func Equals(x1 XValue, x2 XValue) bool {
	// nil == nil
	if utils.IsNil(x1) && utils.IsNil(x2) {
		return true
	} else if utils.IsNil(x1) || utils.IsNil(x2) {
		return false
	}

	x1 = x1.Reduce()
	x2 = x2.Reduce()

	// different types aren't equal
	if reflect.TypeOf(x1) != reflect.TypeOf(x2) {
		return false
	}

	// common types, do real comparisons
	switch typed := x1.(type) {
	case XText:
		return typed.Equals(x2.(XText))
	case XNumber:
		return typed.Equals(x2.(XNumber))
	case XBoolean:
		return typed.Equals(x2.(XBoolean))
	case XDateTime:
		return typed.Equals(x2.(XDateTime))
	case XError:
		return typed.Equals(x2.(XError))
	}

	// for arrays and maps, use equality of their JSON representation
	return x1.ToXJSON().Native() == x2.ToXJSON().Native()
}

// IsEmpty determines if the given value is empty
func IsEmpty(x XValue) bool {
	// nil is empty
	if utils.IsNil(x) {
		return true
	}

	// anything with length of zero is empty
	asLengthable, isLengthable := x.(XLengthable)
	if isLengthable && asLengthable.Length() == 0 {
		return true
	}

	return false
}

// Repr returns the representation of the given value
func Repr(x XValue) string {
	if utils.IsNil(x) {
		return "null"
	}
	return x.Describe()
}
