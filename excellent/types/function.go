package types

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
)

// XFunction is a callable function.
//
//   @(upper) -> function
//   @(array(upper)[0]("abc")) -> ABC
//   @(json(upper)) -> null
//
// @type function
type XFunction func(env envs.Environment, args ...XValue) XValue

// Describe returns a representation of this type for error messages
func (x XFunction) Describe() string { return "function" }

// Truthy determines truthiness for this type
func (x XFunction) Truthy() bool { return true }

// Render returns the canonical text representation
func (x XFunction) Render() string {
	return "function"
}

// Format returns the pretty text representation
func (x XFunction) Format(env envs.Environment) string {
	return x.Render()
}

// MarshalJSON converts this type to JSON
func (x XFunction) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(nil)
}

// String returns the native string representation of this type
func (x XFunction) String() string {
	return `XFunction`
}

// Equals determines equality for this type
func (x XFunction) Equals(other XFunction) bool {
	return true // TODO
}

var _ XValue = XFunction(nil)
