package runs

import (
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
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
func (c *runContext) Resolve(key string) types.XValue {
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

func (c *runContext) Reduce() types.XPrimitive {
	return types.NewXString(c.run.UUID().String())
}

func (c *runContext) ToJSON() types.XString { return types.NewXString("TODO") }

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
func (c *relatedRunContext) Resolve(key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXString(string(c.run.UUID()))
	case "contact":
		return c.run.Contact()
	case "flow":
		return c.run.Flow()
	case "status":
		return types.NewXString(string(c.run.Status()))
	case "results":
		return c.run.Results()
	}

	return types.NewXResolveError(c, key)
}

func (c *relatedRunContext) Reduce() types.XPrimitive {
	return types.NewXString(c.run.UUID().String())
}

func (c *relatedRunContext) ToJSON() types.XString { return types.NewXString("TODO") }

var _ types.XValue = (*relatedRunContext)(nil)
var _ types.XResolvable = (*relatedRunContext)(nil)
