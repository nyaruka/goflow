package waits

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const MSG string = "msg"

type MsgWait struct {
	Timeout int `json:"timeout"`
}

func (w *MsgWait) Resolve(key string) interface{} {
	switch key {
	case "timeout":
		return w.Timeout
	}

	return nil
}

func (w *MsgWait) Default() interface{} {
	return MSG
}

func (w *MsgWait) Type() string { return MSG }

func (w *MsgWait) Begin(run flows.FlowRun, step flows.Step) error {
	run.AddEvent(step, &events.MsgWaitEvent{Timeout: w.Timeout})
	run.SetWait(w)
	return nil
}

func (w *MsgWait) ShouldEnd(run flows.FlowRun, step flows.Step) (flows.Event, error) {
	return nil, nil
}

func (w *MsgWait) End(run flows.FlowRun, step flows.Step, event flows.Event) error {
	msgEvent, isMsg := event.(*events.MsgInEvent)
	if !isMsg {
		return fmt.Errorf("Must end MsgWait with MsgInEvent, got: %#v", event)
	}

	// add our msg to our step
	run.AddEvent(step, msgEvent)

	// and set our input @input.text
	run.SetInput(msgEvent)

	// and clear our wait
	run.SetWait(nil)

	return nil
}
