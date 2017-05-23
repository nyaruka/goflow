package events

import "github.com/nyaruka/goflow/flows"

const ADD_TO_GROUP string = "add_to_group"

type AddToGroupEvent struct {
	Group flows.GroupUUID `json:"group"  validate:"required"`
	Name  string          `json:"name"   validate:"required"`
	BaseEvent
}

func (e *AddToGroupEvent) Type() string { return ADD_TO_GROUP }
