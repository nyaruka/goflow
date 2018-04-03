package triggers

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseTrigger struct {
	environment utils.Environment
	flow        flows.Flow
	contact     *flows.Contact
	params      types.JSONFragment
	triggeredOn time.Time
}

func (t *baseTrigger) Environment() utils.Environment { return t.environment }
func (t *baseTrigger) Flow() flows.Flow               { return t.flow }
func (t *baseTrigger) Contact() *flows.Contact        { return t.contact }
func (t *baseTrigger) Params() types.JSONFragment     { return t.params }
func (t *baseTrigger) TriggeredOn() time.Time         { return t.triggeredOn }

// Resolve resolves the given key when this trigger is referenced in an expression
func (t *baseTrigger) Resolve(key string) interface{} {
	switch key {
	case "params":
		return t.params
	}

	return fmt.Errorf("No such field '%s' on trigger", key)
}

// Atomize is called when this object needs to be reduced to a primitive
func (t *baseTrigger) Atomize() interface{} {
	return string(t.flow.UUID())
}

var _ utils.Atomizable = (*baseTrigger)(nil)
var _ utils.Resolvable = (*baseTrigger)(nil)
