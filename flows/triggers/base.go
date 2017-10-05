package triggers

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

type baseTrigger struct {
	flow        flows.Flow
	triggeredOn time.Time
}

func (t *baseTrigger) Flow() flows.Flow       { return t.flow }
func (t *baseTrigger) TriggeredOn() time.Time { return t.triggeredOn }
