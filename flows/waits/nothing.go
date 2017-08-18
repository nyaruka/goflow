package waits

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const TypeNothing string = "nothing"

type NothingWait struct {
	BaseWait
}

func (w *NothingWait) Type() string { return TypeNothing }

func (w *NothingWait) Begin(run flows.FlowRun, step flows.Step) {
	w.BaseWait.begin(run)

	run.ApplyEvent(step, nil, events.NewNothingWait())
}

// CanResume always returns true for a nothing wait because it's not waiting for anything
func (w *NothingWait) CanResume(run flows.FlowRun, step flows.Step) bool {
	return true
}

var _ flows.Wait = (*NothingWait)(nil)
