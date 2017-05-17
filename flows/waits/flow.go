package waits

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

const FLOW string = "flow"

type FlowWait struct {
	Flow flows.FlowUUID `json:"flow"`
}

func (w *FlowWait) Resolve(key string) interface{} {
	switch key {
	case "flow":
		return w.Flow
	}

	return nil
}

func (w *FlowWait) Default() interface{} {
	return FLOW
}

func (w *FlowWait) Type() string { return FLOW }

func (w *FlowWait) Begin(run flows.FlowRun, step flows.Step) error {
	run.AddEvent(step, &events.FlowWaitEvent{Flow: w.Flow})
	run.SetWait(w)
	return nil
}

func (w *FlowWait) ShouldEnd(run flows.FlowRun, step flows.Step) (flows.Event, error) {
	child := run.Child()
	if child == nil {
		return nil, fmt.Errorf("FlowWait should always have a child run set")
	}

	// child isn't complete yet, shouldn't end
	if child.Status() == flows.RunActive {
		return nil, nil
	}

	return events.NewFlowExitEvent(child), nil
}

func (w *FlowWait) End(run flows.FlowRun, step flows.Step, event flows.Event) error {
	flowEvent, isFlow := event.(*events.FlowExitEvent)
	if !isFlow {
		return fmt.Errorf("Must end FlowWait with FlowExitEvent, got: %#v", event)
	}

	// make sure the flows match
	if flowEvent.Flow != w.Flow {
		return fmt.Errorf("Must end FlowWait with FlowExitEvent for the same flow, expected '%s', got '%s'", w.Flow, flowEvent.Flow)
	}

	// log this event
	run.AddEvent(step, flowEvent)

	// and clear our wait
	run.SetWait(nil)
	return nil
}
