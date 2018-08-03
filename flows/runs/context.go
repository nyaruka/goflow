package runs

import (
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// RunContextTopLevels are the allowed top-level variables for expression evaluations
var RunContextTopLevels = []string{"contact", "child", "parent", "run", "trigger"}

type runContext struct {
	run flows.FlowRun
}

// creates a new evaluation context for the passed in run
func newRunContext(run flows.FlowRun) types.XValue {
	return &runContext{run: run}
}

// Resolve resolves the given top-level key in an expression
func (c *runContext) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "contact":
		return c.run.Contact()
	case "child":
		return newRelatedRunContext(c.run.Session().GetCurrentChild(c.run))
	case "parent":
		return newRelatedRunContext(c.run.Parent())
	case "run":
		return c.run
	case "trigger":
		return c.run.Session().Trigger()
	}

	return types.NewXResolveError(c, key)
}

// Describe returns a representation of this type for error messages
func (c *runContext) Describe() string { return "context" }

func (c *runContext) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(string(c.run.UUID()))
}

// ToXJSON can never actually be called on the context root
func (c *runContext) ToXJSON(env utils.Environment) types.XText {
	panic("shouldn't be possible to call ToXJSON on the context root")
}

var _ types.XValue = (*runContext)(nil)
var _ types.XResolvable = (*runContext)(nil)

// wraps parent/child runs and provides a reduced set of keys in the context
type relatedRunContext struct {
	run flows.RunSummary
}

// creates a new context for related runs
func newRelatedRunContext(run flows.RunSummary) *relatedRunContext {
	if run != nil {
		return &relatedRunContext{run: run}
	}
	return nil
}

// Resolve resolves the given key when this related run is referenced in an expression
func (c *relatedRunContext) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXText(string(c.run.UUID()))
	case "contact":
		return c.run.Contact()
	case "flow":
		return c.run.Flow()
	case "status":
		return types.NewXText(string(c.run.Status()))
	case "results":
		return c.run.Results()
	}

	return types.NewXResolveError(c, key)
}

// Describe returns a representation of this type for error messages
func (c *relatedRunContext) Describe() string { return "related run" }

func (c *relatedRunContext) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(string(c.run.UUID()))
}

// ToXJSON is called when this type is passed to @(json(...))
func (c *relatedRunContext) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, c, "uuid", "contact", "flow", "status", "results").ToXJSON(env)
}

var _ types.XValue = (*relatedRunContext)(nil)
var _ types.XResolvable = (*relatedRunContext)(nil)
