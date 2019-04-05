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
		return c.run.ToXValue(env)
	case "child":
		return RelatedRunToXValue(env, c.run.Session().GetCurrentChild(c.run))
	case "parent":
		return RelatedRunToXValue(env, c.run.Parent())

	// shortcuts to things on the current run
	case "contact":
		return types.ToXValue(env, c.run.Contact())
	case "results":
		return c.run.Results().ToSimpleXDict(env)

	case "urns":
		if c.run.Contact() != nil {
			return c.run.Contact().URNs().MapContext(env)
		}
		return nil
	case "fields":
		if c.run.Contact() != nil {
			return c.run.Contact().Fields().ToXValue(env)
		}
		return nil

	// other
	case "trigger":
		return c.run.Session().Trigger().ToXValue(env)
	case "input":
		return types.ToXValue(env, c.run.Session().Input())
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

// RelatedRunToXValue returns a representation of a related run for use in expressions
func RelatedRunToXValue(env utils.Environment, run flows.RunSummary) types.XValue {
	if run == nil {
		return nil
	}

	var urns, fields types.XValue
	if run.Contact() != nil {
		urns = run.Contact().URNs().MapContext(env)
	}
	if run.Contact() != nil {
		fields = run.Contact().Fields().ToXValue(env)
	}

	return types.NewXDict(map[string]types.XValue{
		"uuid":    types.NewXText(string(run.UUID())),
		"contact": types.ToXValue(env, run.Contact()),
		"urns":    urns,
		"fields":  fields,
		"flow":    run.Flow().ToXValue(env),
		"status":  types.NewXText(string(run.Status())),
		"results": run.Results().ToXValue(env),
	})
}
