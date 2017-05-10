package events

import "github.com/nyaruka/goflow/flows"

const FLOW_WAIT string = "flow_wait"

type FlowWaitEvent struct {
	Flow flows.FlowUUID `json:"flow"`
	BaseEvent
}

func (e *FlowWaitEvent) Type() string { return FLOW_WAIT }
