package events

import "github.com/nyaruka/goflow/flows"

const SAVE_TO_CONTACT string = "save_to_contact"

type SaveToContactEvent struct {
	Field flows.FieldUUID `json:"field"  validate:"required"`
	Name  string          `json:"name"   validate:"required"`
	Value string          `json:"value"  validate:"required"`
	BaseEvent
}

func (e *SaveToContactEvent) Type() string { return SAVE_TO_CONTACT }
