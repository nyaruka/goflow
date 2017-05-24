package events

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
)

// TypeFlowExit is the type of our flow exit
const TypeFlowExit string = "flow_exit"

// FlowExitEvent events are created when a contact exits a flow. It contains not only the
// contact and flow which was exited, but also the time it was exited and the exit status.
//
// ```
//   {
//    "step": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "flow_exit",
//    "flow": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//    "contact": "95eb96df-461b-4668-b168-727f8ceb13dd",
//    "exited_on": "2006-01-02T15:04:05Z",
//    "status": "C"
//   }
// ```
//
// @event flow_exit
type FlowExitEvent struct {
	BaseEvent
	Flow     flows.FlowUUID    `json:"flow"       validate:"required"`
	Contact  flows.ContactUUID `json:"contact"    validate:"required"`
	Status   flows.RunStatus   `json:"status"     validate:"required"`
	ExitedOn *time.Time        `json:"exited_on"  validate:"required"`
}

// NewFlowExitEvent returns a new flow exit event
func NewFlowExitEvent(run flows.FlowRunReference) *FlowExitEvent {
	event := FlowExitEvent{Flow: run.FlowUUID(), Status: run.Status(), Contact: run.ContactUUID(), ExitedOn: run.ExitedOn()}
	return &event
}

// Type returns the type of our event
func (e *FlowExitEvent) Type() string { return TypeFlowExit }

// Resolve resolves the passed in key
func (e *FlowExitEvent) Resolve(key string) interface{} {
	switch key {

	case "contact":
		return e.Contact

	case "exited_on":
		return e.ExitedOn

	case "flow":
		return e.Flow

	case "status":
		return e.Status

	}

	return fmt.Errorf("no such field '%s' on Flow Exit event", key)
}

// Default returns the default value for this event
func (e *FlowExitEvent) Default() interface{} {
	return e
}

// String returns the default string value
func (e *FlowExitEvent) String() interface{} {
	return string(e.Flow)
}
