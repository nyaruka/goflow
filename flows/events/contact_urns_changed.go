package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactURNsChanged, func() flows.Event { return &ContactURNsChanged{} })
}

// TypeContactURNsChanged is the type of our URNs changed event
const TypeContactURNsChanged string = "contact_urns_changed"

// ContactURNsChanged events are created when a contact's URNs have changed.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "contact_urns_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "urns": [
//	    "tel:+12345678900",
//	    "twitter:bob"
//	  ]
//	}
//
// @event contact_urns_changed
type ContactURNsChanged struct {
	BaseEvent

	URNs []urns.URN `json:"urns" validate:"dive,urn"`
}

// NewContactURNsChanged returns a new add URN event
func NewContactURNsChanged(urns []urns.URN) *ContactURNsChanged {
	return &ContactURNsChanged{
		BaseEvent: NewBaseEvent(TypeContactURNsChanged),
		URNs:      urns,
	}
}
