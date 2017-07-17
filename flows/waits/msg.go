package waits

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const TypeMsg string = "msg"

type MsgWait struct {
	Timeout int `json:"timeout,omitempty"`
}

func (w *MsgWait) Resolve(key string) interface{} {
	switch key {
	case "timeout":
		return w.Timeout
	}

	return nil
}

func (w *MsgWait) Default() interface{} {
	return TypeMsg
}

func (w *MsgWait) String() string {
	return TypeMsg
}

var _ utils.VariableResolver = (*MsgWait)(nil)

func (w *MsgWait) Type() string { return TypeMsg }

func (w *MsgWait) Begin(run flows.FlowRun, step flows.Step) error {
	run.ApplyEvent(step, &events.MsgWaitEvent{Timeout: w.Timeout})
	run.SetWait(w)
	return nil
}

func (w *MsgWait) GetEndEvent(run flows.FlowRun, step flows.Step) (flows.Event, error) {
	return nil, nil
}

func (w *MsgWait) End(run flows.FlowRun, step flows.Step, event flows.Event) error {
	msgEvent, isMsg := event.(*events.MsgReceivedEvent)
	if !isMsg {
		return fmt.Errorf("Must end MsgWait with MsgInEvent, got: %#v", event)
	}

	// add our msg to our step
	run.ApplyEvent(step, msgEvent)

	// and set our input @input.text
	run.SetInput(msgEvent)

	// and clear our wait
	run.SetWait(nil)

	return nil
}
