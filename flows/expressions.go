package flows

import (
	"reflect"
	"strconv"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Contextable is an object that can accessed in expressions as a object with properties
type Contextable interface {
	Context(env envs.Environment) map[string]types.XValue
}

// Context generates a lazy object for use in expressions
func Context(env envs.Environment, contextable Contextable) types.XValue {
	// we allow passing nil pointers which will become non-nil Contextables
	if contextable == nil || reflect.ValueOf(contextable).IsNil() {
		return nil
	}

	return types.NewXLazyObject(func() map[string]types.XValue {
		return contextable.Context(env)
	})

}

// ContextFunc generates a lazy object for use in expressions
func ContextFunc(env envs.Environment, fn func(envs.Environment) map[string]types.XValue) *types.XObject {
	return types.NewXLazyObject(func() map[string]types.XValue {
		return fn(env)
	})
}

// RunContextTopLevels are the allowed top-level variables for expression evaluations
var RunContextTopLevels = []string{
	"child",
	"contact",
	"fields",
	"globals",
	"input",
	"legacy_extra",
	"node",
	"parent",
	"locals",
	"results",
	"resume",
	"run",
	"ticket",
	"trigger",
	"urns",
	"webhook",
}

// ContactQueryEscaping is the escaping function used for expressions in contact queries
func ContactQueryEscaping(s string) string {
	return strconv.Quote(s)
}
