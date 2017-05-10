package events

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

const FLOW_ENTER string = "flow_enter"

type FlowEnterEvent struct {
	Flow    flows.FlowUUID    `json:"flow"        validate:"nonzero"`
	Contact flows.ContactUUID `json:"contact"     validate:"nonzero"`
	BaseEvent
}

func NewFlowEnterEvent(flow flows.FlowUUID, contact flows.ContactUUID) *FlowEnterEvent {
	event := FlowEnterEvent{Flow: flow, Contact: contact}
	return &event
}

func (e *FlowEnterEvent) Type() string { return FLOW_ENTER }

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

func (e *FlowEnterEvent) Default() interface{} {
	return e.Flow
}
