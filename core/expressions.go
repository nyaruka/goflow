package core

import (
	"reflect"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Contextable is an object that can accessed in expressions as a object with properties
type Contextable interface {
	Context(env envs.Environment) map[string]types.XValue
}

// Context generates a lazy object for use in expressions
func Context(env envs.Environment, contextable Contextable) types.XValue {
	// this is the one place where callers may pass nil pointers boxed into non-nil Contextables, and we scrub
	// them to true nils here - XValues themselves are never non-nil interfaces to nil pointers
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
