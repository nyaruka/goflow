package waits

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const TypeTime string = "time"

type TimeWait struct {
	Timeout int `json:"timeout" validation:"required,gte=0"`
}

func (w *TimeWait) Type() string { return TypeMsg }

func (w *TimeWait) Apply(run flows.FlowRun, step flows.Step) {
	run.ApplyEvent(step, nil, &events.TimeWaitEvent{Timeout: w.Timeout})
}

// CanResume returns true for a message wait if a message has now been received on this step
func (w *TimeWait) CanResume(run flows.FlowRun, step flows.Step) bool {
	return true
}

var _ flows.Wait = (*TimeWait)(nil)
