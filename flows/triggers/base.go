package triggers

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseTrigger struct {
	flow        flows.Flow
	contact     *flows.Contact
	params      utils.JSONFragment
	triggeredOn time.Time
}

func (t *baseTrigger) Flow() flows.Flow           { return t.flow }
func (t *baseTrigger) Contact() *flows.Contact    { return t.contact }
func (t *baseTrigger) Params() utils.JSONFragment { return t.params }
func (t *baseTrigger) TriggeredOn() time.Time     { return t.triggeredOn }

func (t *baseTrigger) Default() interface{} {
	return t
}

// Resolve resolves the passed in key to a value, returning an error if the key is unknown
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
