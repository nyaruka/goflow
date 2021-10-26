package types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/nyaruka/goflow/envs"
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
	Format(env envs.Environment) string

	Equals(o XValue) bool
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

	return x1.Equals(x2)
}

// Describe returns a representation of the given value for use in error messages
func Describe(x XValue) string {
	if utils.IsNil(x) {
		return "null"
	}
	return x.Describe()
}

// Truthy determines truthiness for the given value
func Truthy(x XValue) bool {
	if utils.IsNil(x) {
		return false
	}
	return x.Truthy()
}

// Render returns the canonical text representation
func Render(x XValue) string {
	if utils.IsNil(x) {
		return ""
	}
	return x.Render()
}

// Format returns the pretty text representation
func Format(env envs.Environment, x XValue) string {
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
