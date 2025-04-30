package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeOptInRequested, func() flows.Event { return &OptInRequestedEvent{} })
}

// TypeOptInRequested is our type for the optin event
const TypeOptInRequested string = "optin_requested"

// OptInRequestedEvent events are created when an action has created an optin to be sent.
//
//	{
//	  "uuid": "019688A6-41d2-7366-958a-630e35c62431",
//	  "type": "optin_requested",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "optin": {
//	    "uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
//	    "name": "Joke Of The Day"
//	  },
//	  "channel": {
//	    "uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648",
//	    "name": "Facebook"
//	  },
//	  "urn": "tel:+12065551212"
//	}
//
// @event optin_requested
type OptInRequestedEvent struct {
	BaseEvent

	OptIn   *assets.OptInReference   `json:"optin" validate:"required"`
	Channel *assets.ChannelReference `json:"channel" validate:"required"`
	URN     urns.URN                 `json:"urn" validate:"required"`
}

// NewOptInRequested returns a new optin requested event
func NewOptInRequested(optIn *flows.OptIn, ch *flows.Channel, urn urns.URN) *OptInRequestedEvent {
	return &OptInRequestedEvent{
		BaseEvent: NewBaseEvent(TypeOptInRequested),
		OptIn:     optIn.Reference(),
		Channel:   ch.Reference(),
		URN:       urn,
	}
}
