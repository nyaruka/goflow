package waits

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const TypeMsg string = "msg"

type MsgWait struct {
	TimeoutWait
}

func NewMsgWait(timeout *int) *MsgWait {
	return &MsgWait{TimeoutWait{Timeout: timeout}}
}

func (w *MsgWait) Type() string { return TypeMsg }

func (w *MsgWait) Begin(run flows.FlowRun, step flows.Step) {
	w.TimeoutWait.Begin(run)

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

func (w *MsgWait) ResumeByTimeOut(run flows.FlowRun) {
	w.BaseWait.ResumeByTimeOut(run)

	run.SetInput(nil)
}

var _ flows.Wait = (*MsgWait)(nil)
