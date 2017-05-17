package events

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
)

const FLOW_EXIT string = "flow_exit"

type FlowExitEvent struct {
	Flow     flows.FlowUUID    `json:"flow"       validate:"nonzero"`
	Status   flows.RunStatus   `json:"status"     validate:"nonzero"`
	Contact  flows.ContactUUID `json:"contact"    validate:"nonzero"`
	ExitedOn *time.Time        `json:"exited_on"  validate:"nonzero"`
	BaseEvent
}

func NewFlowExitEvent(run flows.FlowRunReference) *FlowExitEvent {
	event := FlowExitEvent{Flow: run.FlowUUID(), Status: run.Status(), Contact: run.ContactUUID(), ExitedOn: run.ExitedOn()}
	return &event
}

func (e *FlowExitEvent) Type() string { return FLOW_EXIT }

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

	return fmt.Errorf("No such field '%s' on Flow Exit event", key)
}

func (e *FlowExitEvent) Default() interface{} {
	return e.Flow
}
