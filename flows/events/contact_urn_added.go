package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeContactURNAdded, func() flows.Event { return &ContactURNAddedEvent{} })
}

// TypeContactURNAdded is the type of our add URN event
const TypeContactURNAdded string = "contact_urn_added"

// ContactURNAddedEvent events are created when a URN is added to a contact.
//
//   {
//     "type": "contact_urn_added",
//     "created_on": "2006-01-02T15:04:05Z",
//     "urn": "tel:+12345678900"
//   }
//
// @event contact_urn_added
type ContactURNAddedEvent struct {
	BaseEvent

	URN urns.URN `json:"urn" validate:"urn"`
}

// NewURNAddedEvent returns a new add URN event
func NewURNAddedEvent(urn urns.URN) *ContactURNAddedEvent {
	return &ContactURNAddedEvent{BaseEvent: NewBaseEvent(), URN: urn}
}

// Type returns the type of this event
func (e *ContactURNAddedEvent) Type() string { return TypeContactURNAdded }
