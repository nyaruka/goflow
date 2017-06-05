package waits

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const TypeFlow string = "flow"

type FlowWait struct {
	FlowUUID flows.FlowUUID `json:"flow_uuid"   validate:"required,uuid4"`
}

func (w *FlowWait) Resolve(key string) interface{} {
	switch key {
	case "flow":
		return w.FlowUUID
	}

	return nil
}

func (w *FlowWait) Default() interface{} {
	return TypeFlow
}

func (w *FlowWait) Type() string { return TypeFlow }

func (w *FlowWait) Begin(run flows.FlowRun, step flows.Step) error {
	run.AddEvent(step, &events.FlowWaitEvent{FlowUUID: w.FlowUUID})
	run.SetWait(w)
	return nil
}

func (w *FlowWait) GetEndEvent(run flows.FlowRun, step flows.Step) (flows.Event, error) {
	child := run.Child()
	if child == nil {
		return nil, fmt.Errorf("FlowWait should always have a child run set")
	}

	// child isn't complete yet, shouldn't end
	if child.Status() == flows.RunActive {
		return nil, nil
	}

	// see if we already have an exit event on our step for this flow
	for _, evt := range step.Events() {
		exit, isExit := evt.(*events.FlowExitedEvent)
		if isExit && exit.FlowUUID == w.FlowUUID {
			return exit, nil
		}
	}

	// our flow didn't exit, return nil
	return nil, nil
}

func (w *FlowWait) End(run flows.FlowRun, step flows.Step, event flows.Event) error {
	flowEvent, isFlow := event.(*events.FlowExitedEvent)
	if !isFlow {
		return fmt.Errorf("Must end FlowWait with FlowExitEvent, got: %#v", event)
	}

	// make sure the flows match
	if flowEvent.FlowUUID != w.FlowUUID {
		return fmt.Errorf("Must end FlowWait with FlowExitEvent for the same flow, expected '%s', got '%s'", w.FlowUUID, flowEvent.FlowUUID)
	}

	// log this event
	run.AddEvent(step, flowEvent)

	// and clear our wait
	run.SetWait(nil)
	return nil
}
