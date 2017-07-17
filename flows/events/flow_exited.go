package events

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
)

// TypeFlowExited is the type of our flow exit
const TypeFlowExited string = "flow_exited"

// FlowExitedEvent events are created when a contact exits a flow. It contains not only the
// contact and flow which was exited, but also the time it was exited and the exit status.
//
// ```
//   {
//    "type": "flow_exited",
//    "step_uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "flow_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//    "contact_uuid": "95eb96df-461b-4668-b168-727f8ceb13dd",
//    "exited_on": "2006-01-02T15:04:05Z",
//    "status": "C"
//   }
// ```
//
// @event flow_exited
type FlowExitedEvent struct {
	BaseEvent
	FlowUUID    flows.FlowUUID    `json:"flow_uuid"       validate:"required"`
	ContactUUID flows.ContactUUID `json:"contact_uuid"    validate:"required"`
	Status      flows.RunStatus   `json:"status"          validate:"required"`
	ExitedOn    *time.Time        `json:"exited_on"       validate:"required"`
}

// NewFlowExitedEvent returns a new flow exit event
func NewFlowExitedEvent(run flows.FlowRunReference) *FlowExitedEvent {
	return &FlowExitedEvent{
		BaseEvent:   NewBaseEvent(),
		FlowUUID:    run.FlowUUID(),
		Status:      run.Status(),
		ContactUUID: run.ContactUUID(),
		ExitedOn:    run.ExitedOn(),
	}
}

// Type returns the type of our event
func (e *FlowExitedEvent) Type() string { return TypeFlowExited }

// Resolve resolves the passed in key
func (e *FlowExitedEvent) Resolve(key string) interface{} {
	switch key {

	case "contact_uuid":
		return e.ContactUUID

	case "exited_on":
		return e.ExitedOn

	case "flow_uuid":
		return e.FlowUUID

	case "status":
		return e.Status

	}

	return fmt.Errorf("no such field '%s' on Flow Exit event", key)
}

// Default returns the default value for this event
func (e *FlowExitedEvent) Default() interface{} {
	return e
}

// String returns the default string value
func (e *FlowExitedEvent) String() string {
	return string(e.FlowUUID)
}

// Apply applies this event to the given run
func (e *FlowExitedEvent) Apply(run flows.FlowRun, step flows.Step) {}

var _ flows.Input = (*FlowExitedEvent)(nil)
