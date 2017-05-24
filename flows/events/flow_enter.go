package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

// TypeFlowEnter is the type of our flow enter event
const TypeFlowEnter string = "flow_enter"

// FlowEnterEvent events are created when a contact first enters a flow.
//
// ```
//   {
//    "step": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "flow_enter",
//    "flow": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
//    "contact": "95eb96df-461b-4668-b168-727f8ceb13dd"
//   }
// ```
//
// @event flow_enter
type FlowEnterEvent struct {
	BaseEvent
	Flow    flows.FlowUUID    `json:"flow"        validate:"required,uuid4"`
	Contact flows.ContactUUID `json:"contact"     validate:"required,uuid4"`
}

// NewFlowEnterEvent returns a new flow enter event for the passed in flow and contact
func NewFlowEnterEvent(flow flows.FlowUUID, contact flows.ContactUUID) *FlowEnterEvent {
	event := FlowEnterEvent{Flow: flow, Contact: contact}
	return &event
}

// Type returns the type of this event
func (e *FlowEnterEvent) Type() string { return TypeFlowEnter }

// Resolve resolves the passed in key for this event
func (e *FlowEnterEvent) Resolve(key string) interface{} {
	switch key {

	case "contact":
		return e.Contact

	case "created_on":
		return e.CreatedOn

	case "flow":
		return e.Flow
	}

	return fmt.Errorf("No such field '%s' on Flow Enter event", key)
}

// Default returns the default value for this event
func (e *FlowEnterEvent) Default() interface{} {
	return e
}

// String returns the default string value for this event
func (e *FlowEnterEvent) String() string {
	return string(e.Flow)
}
