package runs

import (
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type runContext struct {
	run flows.FlowRun

	extra *legacyExtra
}

// creates a new evaluation context for the passed in run
func newRunContext(run flows.FlowRun) types.XValue {
	return &runContext{run: run, extra: newLegacyExtra(run)}
}

// Resolve resolves the given top-level key in an expression
func (c *runContext) Resolve(env utils.Environment, key string) types.XValue {
	switch strings.ToLower(key) {
	// the different runs accessible
	case "run":
		return c.run
	case "child":
		return newRelatedRunContext(c.run.Session().GetCurrentChild(c.run))
	case "parent":
		return newRelatedRunContext(c.run.Parent())

	// shortcuts to things on the current run
	case "contact":
		return c.run.Contact()
	case "results":
		return c.run.Results()

	case "urns":
		if c.run.Contact() != nil {
			return c.run.Contact().URNs().MapContext()
		}
		return nil
	case "fields":
		if c.run.Contact() != nil {
			return c.run.Contact().Fields().Context(env)
		}
		return nil

	// other
	case "trigger":
		return c.run.Session().Trigger()
	case "input":
		return c.run.Session().Input()
	case "legacy_extra":
		c.extra.update()
		return c.extra
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
	switch strings.ToLower(key) {
	case "uuid":
		return types.NewXText(string(c.run.UUID()))

	case "contact":
		return c.run.Contact()
	case "urns":
		if c.run.Contact() != nil {
			return c.run.Contact().URNs().MapContext()
		}
		return nil
	case "fields":
		if c.run.Contact() != nil {
			return c.run.Contact().Fields().Context(env)
		}
		return nil

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
