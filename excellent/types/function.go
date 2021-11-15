package types

import (
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
)

type XFunc func(env envs.Environment, args ...XValue) XValue

// XFunction is a callable function.
//
//   @(upper) -> upper
//   @(array(upper)[0]("abc")) -> ABC
//   @(json(upper)) -> null
//
// @type function
type XFunction struct {
	name string
	fn   XFunc
}

// NewXFunction creates a new XFunction
func NewXFunction(name string, fn XFunc) *XFunction {
	return &XFunction{name: name, fn: fn}
}

func (x *XFunction) Call(env envs.Environment, params []XValue) XValue {
	val := x.fn(env, params...)

	// if function returned an error, wrap the error with the function name
	if IsXError(val) {
		return NewXErrorf("error calling %s: %s", x.Describe(), val.(XError).Error())
	}

	return val
}

// Name returns the name or <anon> for anonymous functions
func (x *XFunction) Name() string {
	if x.name != "" {
		return x.name
	}
	return "<anon>"
}

// Describe returns a representation of this type for error messages
func (x *XFunction) Describe() string {
	return fmt.Sprintf("%s(...)", x.Name())
}

// Truthy determines truthiness for this type
func (x *XFunction) Truthy() bool { return true }

// Render returns the canonical text representation
func (x *XFunction) Render() string {
	return x.Name()
}

// Format returns the pretty text representation
func (x *XFunction) Format(env envs.Environment) string {
	return x.Render()
}

// MarshalJSON converts this type to JSON
func (x *XFunction) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(nil)
}

// String returns the native string representation of this type
func (x *XFunction) String() string {
	return fmt.Sprintf("XFunction[%s]", x.Name())
}

// Equals determines equality for this type
func (x *XFunction) Equals(o XValue) bool {
	other := o.(*XFunction)

	// functions are equal if they have the same name but anon functions are never equal
	return x.name != "" && other.name != "" && x.name == other.name
}

var _ XValue = (*XFunction)(nil)
