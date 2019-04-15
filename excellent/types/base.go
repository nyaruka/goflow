package types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/nyaruka/goflow/utils"
)

// XValue is the base interface of all excellent types
type XValue interface {
	// How type is rendered in console for debugging
	fmt.Stringer

	// How the type JSONifies
	json.Marshaler

	// Describe returns a representation for use in error messages
	Describe() string

	// Truthy determines truthiness for this type
	Truthy() bool

	// Render returns the canonical text representation
	Render() string

	// Format returns the pretty text representation
	Format(env utils.Environment) string
}

// XCountable is the interface for types which can be counted
type XCountable interface {
	Count() int
}

// Equals checks for equality between the two give values. This is only used for testing as x = y
// specifically means text(x) == text(y)
func Equals(x1 XValue, x2 XValue) bool {
	// nil == nil
	if utils.IsNil(x1) && utils.IsNil(x2) {
		return true
	} else if utils.IsNil(x1) || utils.IsNil(x2) {
		return false
	}

	// different types aren't equal
	if reflect.TypeOf(x1) != reflect.TypeOf(x2) {
		return false
	}

	// common types, do real comparisons
	switch typed := x1.(type) {
	case *XArray:
		return typed.Equals(x2.(*XArray))
	case XBoolean:
		return typed.Equals(x2.(XBoolean))
	case XDate:
		return typed.Equals(x2.(XDate))
	case XDateTime:
		return typed.Equals(x2.(XDateTime))
	case XError:
		return typed.Equals(x2.(XError))
	case XFunction:
		return typed.Equals(x2.(XFunction))
	case XNumber:
		return typed.Equals(x2.(XNumber))
	case *XObject:
		return typed.Equals(x2.(*XObject))
	case XText:
		return typed.Equals(x2.(XText))
	case XTime:
		return typed.Equals(x2.(XTime))
	default:
		panic(fmt.Sprintf("can't compare equality of instances of %T", x1))
	}
}

// IsEmpty determines if the given value is empty
func IsEmpty(x XValue) bool {
	// nil is empty
	if utils.IsNil(x) {
		return true
	}

	// empty string is empty
	text, isText := x.(XText)
	if isText && text.Length() == 0 {
		return true
	}

	// anything with count of zero is empty
	countable, isCountable := x.(XCountable)
	if isCountable && countable != nil && countable.Count() == 0 {
		return true
	}

	return false
}

// Describe returns a representation of the given value for use in error messages
func Describe(x XValue) string {
	if utils.IsNil(x) {
		return "null"
	}
	return x.Describe()
}

// Render returns the canonical text representation
func Render(x XValue) string {
	if utils.IsNil(x) {
		return ""
	}
	return x.Render()
}

// Format returns the pretty text representation
func Format(env utils.Environment, x XValue) string {
	if utils.IsNil(x) {
		return ""
	}
	return x.Format(env)
}

// String returns a representation of the given value for use in debugging
func String(x XValue) string {
	if utils.IsNil(x) {
		return "nil"
	}
	return x.String()
}
