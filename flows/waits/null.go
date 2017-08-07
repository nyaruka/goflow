package waits

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const TypeNothing string = "nothing"

type NothingWait struct {
}

func (w *NothingWait) Type() string { return TypeNothing }

func (w *NothingWait) Apply(run flows.FlowRun, step flows.Step) {
	run.ApplyEvent(step, nil, &events.NothingWaitEvent{})
}

// CanResume always returns true for a nothing wait because it's not waiting for anything
func (w *NothingWait) CanResume(run flows.FlowRun, step flows.Step) bool {
	return true
}

var _ flows.Wait = (*NothingWait)(nil)
