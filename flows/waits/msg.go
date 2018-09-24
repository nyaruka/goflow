package waits

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeMsg, func() flows.Wait { return &MsgWait{} })
}

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

	run.ApplyEvent(step, nil, events.NewMsgWait(w.TimeoutOn_))
}

// CanResume returns true if a message event has been received
func (w *MsgWait) CanResume(callerEvents []flows.CallerEvent) bool {
	if containsEventOfType(callerEvents, events.TypeMsgReceived) {
		return true
	}
	return w.baseTimeoutWait.CanResume(callerEvents)
}

var _ flows.Wait = (*MsgWait)(nil)
