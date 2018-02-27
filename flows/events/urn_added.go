package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

// TypeURNAdded is the type of our add URN event
const TypeURNAdded string = "urn_added"

// URNAddedEvent events will be created with the URN that should be added to the current contact.
//
// ```
//   {
//     "type": "urn_added",
//     "created_on": "2006-01-02T15:04:05Z",
//     "urn": "tel:+12345678900"
//   }
// ```
//
// @event urn_added
type URNAddedEvent struct {
	BaseEvent
	URN urns.URN `json:"urn" validate:"urn"`
}

// NewURNAddedEvent returns a new add URN event
func NewURNAddedEvent(urn urns.URN) *URNAddedEvent {
	return &URNAddedEvent{URN: urn}
}

// Type returns the type of this event
func (e *URNAddedEvent) Type() string { return TypeURNAdded }

// Apply applies this event to the given run
func (e *URNAddedEvent) Apply(run flows.FlowRun) error {
	run.Contact().AddURN(e.URN)
	return nil
}
