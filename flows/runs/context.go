package runs

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// RunContextTopLevels are the allowed top-level variables for expression evaluations
var RunContextTopLevels = []string{"contact", "child", "parent", "run", "trigger"}

type runContext struct {
	run flows.FlowRun
}

// creates a new evaluation context for the passed in run
func newRunContext(run flows.FlowRun) utils.VariableResolver {
	return &runContext{run: run}
}

// Resolve resolves the given top-level key in an expression
func (c *runContext) Resolve(key string) interface{} {
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

	return fmt.Errorf("no field '%s' on context", key)
}

// Default returns the value of this context when it is the result of an expression
func (c *runContext) Default() interface{} {
	return c
}

func (c *runContext) String() string {
	return c.run.UUID().String()
}

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
func (c *relatedRunContext) Resolve(key string) interface{} {
	switch key {
	case "uuid":
		return c.run.UUID()
	case "contact":
		return c.run.Contact()
	case "flow":
		return c.run.Flow()
	case "status":
		return c.run.Status()
	case "results":
		return c.run.Results()
	}

	return fmt.Errorf("no field '%s' on related run", key)
}

func (c *relatedRunContext) Default() interface{} {
	return c
}

func (c *relatedRunContext) String() string {
	return c.run.UUID().String()
}

var _ utils.VariableResolver = (*runContext)(nil)
var _ utils.VariableResolver = (*relatedRunContext)(nil)
