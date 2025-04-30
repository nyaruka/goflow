package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactURNsChanged, func() flows.Event { return &ContactURNsChangedEvent{} })
}

// TypeContactURNsChanged is the type of our URNs changed event
const TypeContactURNsChanged string = "contact_urns_changed"

// ContactURNsChangedEvent events are created when a contact's URNs have changed.
//
//	{
//	  "uuid": "019688A6-41d2-7366-958a-630e35c62431",
//	  "type": "contact_urns_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "urns": [
//	    "tel:+12345678900",
//	    "twitter:bob"
//	  ]
//	}
//
// @event contact_urns_changed
type ContactURNsChangedEvent struct {
	BaseEvent

	URNs []urns.URN `json:"urns" validate:"dive,urn"`
}

// NewContactURNsChanged returns a new add URN event
func NewContactURNsChanged(urns []urns.URN) *ContactURNsChangedEvent {
	return &ContactURNsChangedEvent{
		BaseEvent: NewBaseEvent(TypeContactURNsChanged),
		URNs:      urns,
	}
}
