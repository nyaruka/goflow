package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

// TypeContactURNAdded is the type of our add URN event
const TypeContactURNAdded string = "contact_urn_added"

// ContactURNAddedEvent events will be created with the URN that should be added to the current contact.
//
// ```
//   {
//     "type": "contact_urn_added",
//     "created_on": "2006-01-02T15:04:05Z",
//     "urn": "tel:+12345678900"
//   }
// ```
//
// @event contact_urn_added
type ContactURNAddedEvent struct {
	BaseEvent
	callerOrEngineEvent

	URN urns.URN `json:"urn" validate:"urn"`
}

// NewURNAddedEvent returns a new add URN event
func NewURNAddedEvent(urn urns.URN) *ContactURNAddedEvent {
	return &ContactURNAddedEvent{URN: urn}
}

// Type returns the type of this event
func (e *ContactURNAddedEvent) Type() string { return TypeContactURNAdded }

// Validate validates our event is valid and has all the assets it needs
func (e *ContactURNAddedEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

// Apply applies this event to the given run
func (e *ContactURNAddedEvent) Apply(run flows.FlowRun) error {
	run.Contact().AddURN(e.URN)
	return nil
}
