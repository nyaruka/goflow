package events

import "github.com/nyaruka/goflow/flows"

const REMOVE_FROM_GROUP string = "remove_from_group"

type RemoveFromGroupEvent struct {
	Group flows.GroupUUID `json:"group"  validate:"nonzero"`
	Name  string          `json:"name"   validate:"nonzero"`
	BaseEvent
}

func (e *RemoveFromGroupEvent) Type() string { return REMOVE_FROM_GROUP }
