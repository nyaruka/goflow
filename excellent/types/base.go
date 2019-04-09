package types

import (
	"fmt"
	"reflect"

	"github.com/nyaruka/goflow/utils"
)

// XValue is the base interface of all Excellent types
type XValue interface {
	fmt.Stringer

	Describe() string
	ToXText(env utils.Environment) XText
	ToXBoolean(env utils.Environment) XBoolean
	ToXJSON(env utils.Environment) XText
}

// XLengthable is the interface for types which have a length
type XLengthable interface {
	Length() int
}

// Equals checks for equality between the two give values
func Equals(env utils.Environment, x1 XValue, x2 XValue) bool {
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

	// for complex objects, use equality of their JSON representation
	return x1.ToXJSON(env).Native() == x2.ToXJSON(env).Native()
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

// Describe returns a representation of the given value for use in error messages
func Describe(x XValue) string {
	if utils.IsNil(x) {
		return "null"
	}
	return x.Describe()
}

// XRepresentable is the interface for any object which can be represented in an expression
type XRepresentable interface {
	ToXValue(env utils.Environment) XValue
}

// ToXValue is a utility to convert the given XRepresentable to an XValue
func ToXValue(env utils.Environment, obj XRepresentable) XValue {
	if utils.IsNil(obj) {
		return nil
	}
	return obj.ToXValue(env)
}
