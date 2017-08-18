package waits

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const TypeMsg string = "msg"

type MsgWait struct {
	BaseWait
}

func NewMsgWait(timeout *int) *MsgWait {
	return &MsgWait{BaseWait{Timeout: timeout}}
}

func (w *MsgWait) Type() string { return TypeMsg }

func (w *MsgWait) Begin(run flows.FlowRun, step flows.Step) {
	w.BaseWait.begin(run)

	run.ApplyEvent(step, nil, events.NewMsgWait(w.TimeoutOn))
}

// CanResume returns true for a message wait if a message has now been received on this step
func (w *MsgWait) CanResume(run flows.FlowRun, step flows.Step) bool {
	for _, event := range step.Events() {
		_, isMsg := event.(*events.MsgReceivedEvent)
		if isMsg {
			return true
		}
	}

	return false
}

var _ flows.Wait = (*MsgWait)(nil)
