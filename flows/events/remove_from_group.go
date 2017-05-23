package events

import "github.com/nyaruka/goflow/flows"

const REMOVE_FROM_GROUP string = "remove_from_group"

type RemoveFromGroupEvent struct {
	Group flows.GroupUUID `json:"group"  validate:"required"`
	Name  string          `json:"name"   validate:"required"`
	BaseEvent
}

func (e *RemoveFromGroupEvent) Type() string { return REMOVE_FROM_GROUP }
