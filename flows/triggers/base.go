package triggers

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseTrigger struct {
	environment utils.Environment
	flow        flows.Flow
	contact     *flows.Contact
	params      utils.JSONFragment
	triggeredOn time.Time
}

func (t *baseTrigger) Environment() utils.Environment { return t.environment }
func (t *baseTrigger) Flow() flows.Flow               { return t.flow }
func (t *baseTrigger) Contact() *flows.Contact        { return t.contact }
func (t *baseTrigger) Params() utils.JSONFragment     { return t.params }
func (t *baseTrigger) TriggeredOn() time.Time         { return t.triggeredOn }

// Resolve resolves the given key when this trigger is referenced in an expression
func (t *baseTrigger) Resolve(key string) interface{} {
	switch key {
	case "params":
		return t.params
	}

	return fmt.Errorf("No such field '%s' on trigger", key)
}

func (t *baseTrigger) String() string {
	return string(t.flow.UUID())
}
