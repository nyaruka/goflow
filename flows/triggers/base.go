package triggers

import (
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseTrigger struct {
	environment utils.Environment
	flow        flows.Flow
	contact     *flows.Contact
	params      types.XValue
	triggeredOn time.Time
}

func (t *baseTrigger) Environment() utils.Environment { return t.environment }
func (t *baseTrigger) Flow() flows.Flow               { return t.flow }
func (t *baseTrigger) Contact() *flows.Contact        { return t.contact }
func (t *baseTrigger) Params() types.XValue           { return t.params }
func (t *baseTrigger) TriggeredOn() time.Time         { return t.triggeredOn }

// Resolve resolves the given key when this trigger is referenced in an expression
func (t *baseTrigger) Resolve(key string) types.XValue {
	switch key {
	case "params":
		return t.params
	}

	return types.NewXResolveError(t, key)
}

// Reduce is called when this object needs to be reduced to a primitive
func (t *baseTrigger) Reduce() types.XPrimitive {
	return types.NewXString(string(t.flow.UUID()))
}

func (t *baseTrigger) ToJSON() types.XString { return types.NewXString("TODO") }

var _ types.XValue = (*baseTrigger)(nil)
var _ types.XResolvable = (*baseTrigger)(nil)
