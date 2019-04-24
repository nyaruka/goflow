package flows

import (
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Contextable is an object that can accessed in expressions as a object with properties
type Contextable interface {
	Context(env utils.Environment) map[string]types.XValue
}

// Context generates a lazy object for use in expressions
func Context(env utils.Environment, contextable Contextable) *types.XObject {
	if !utils.IsNil(contextable) {
		return types.NewXLazyObject(func() map[string]types.XValue {
			return contextable.Context(env)
		})
	}
	return nil
}

// ContextFunc generates a lazy object for use in expressions
func ContextFunc(env utils.Environment, fn func(utils.Environment) map[string]types.XValue) *types.XObject {
	return types.NewXLazyObject(func() map[string]types.XValue {
		return fn(env)
	})
}

// RunContextTopLevels are the allowed top-level variables for expression evaluations
var RunContextTopLevels = []string{
	"child",
	"contact",
	"fields",
	"input",
	"legacy_extra",
	"parent",
	"results",
	"run",
	"trigger",
	"urns",
}
