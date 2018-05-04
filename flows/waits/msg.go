package waits

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const TypeMsg string = "msg"

// MsgWait is a wait which waits for an incoming message (i.e. a msg_received event)
type MsgWait struct {
	baseTimeoutWait
}

// NewMsgWait creates a new message wait
func NewMsgWait(timeout *int) *MsgWait {
	return &MsgWait{baseTimeoutWait{Timeout_: timeout}}
}

// Type returns the type of this wait
func (w *MsgWait) Type() string { return TypeMsg }

// Begin beings waiting at this wait
func (w *MsgWait) Begin(run flows.FlowRun, step flows.Step) {
	w.baseTimeoutWait.Begin(run)

	run.ApplyEvent(step, nil, events.NewMsgWait(w.TimeoutOn))
}

// CanResume returns true for a message wait if a message has now been received
func (w *MsgWait) CanResume(callerEvents []flows.Event) bool {
	for _, event := range callerEvents {
		_, isMsg := event.(*events.MsgReceivedEvent)
		if isMsg {
			return true
		}
	}

	return false
}

func (w *MsgWait) ResumeByTimeOut(run flows.FlowRun) {
	w.baseWait.ResumeByTimeOut(run)

	run.SetInput(nil)
}

var _ flows.Wait = (*MsgWait)(nil)
