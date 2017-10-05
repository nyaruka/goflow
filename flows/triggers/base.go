package triggers

import "time"

type baseTrigger struct {
	triggeredOn time.Time
}

func (t *RunTrigger) TriggeredOn() time.Time { return t.triggeredOn }
