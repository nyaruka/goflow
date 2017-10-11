package runs

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type runContext struct {
	run flows.FlowRun
}

// creates a new evaluation context for the passed in run
func newRunContext(run flows.FlowRun) utils.VariableResolver {
	return &runContext{run: run}
}

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
	}

	return fmt.Errorf("no field '%s' on context", key)
}

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

// Resolve provides a more limited set of results for parent and child runs
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
