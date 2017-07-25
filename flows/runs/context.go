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

func (c *runContext) Validate() error {
	// TODO: do some validation here
	return nil
}

func (c *runContext) Resolve(key string) interface{} {
	switch key {

	case "contact":
		return c.run.Contact()

	case "child":
		return c.run.Child()

	case "input":
		return c.run.Input()

	case "parent":
		return c.run.Parent()

	case "path":
		return c.run.Path()

	case "run":
		return c.run

	case "webhook":
		return c.run.Webhook()

	}

	return fmt.Errorf("No field '%s' on context", key)
}

func (c *runContext) Default() interface{} {
	return c
}

func (c *runContext) String() string {
	return c.run.UUID().String()
}

var _ utils.VariableResolver = (*runContext)(nil)
