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
	case "flow_uuid":
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

	// we don't have a child, so stop execution
	if child == nil {
		return nil, nil
	}

	// child isn't complete yet, shouldn't end
	if child.Status() == flows.RunActive {
		return nil, nil
	}

	// see if we already have an exit event on our step for this flow
	evts := step.Events()
	for i := len(evts) - 1; i >= 0; i-- {
		evt := evts[i]
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
		return fmt.Errorf("must end flow wait with flow_exited event, got: %#v", event.Type())
	}

	// make sure the flows match
	if flowEvent.FlowUUID != w.FlowUUID {
		return fmt.Errorf("must end flow wait with flow_exited for the same flow, expected '%s', got '%s'", w.FlowUUID, flowEvent.FlowUUID)
	}

	// log this event
	run.AddEvent(step, flowEvent)

	// and clear our wait
	run.SetWait(nil)
	return nil
}
