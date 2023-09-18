package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeOptInCreated, func() flows.Event { return &OptInCreatedEvent{} })
}

// TypeOptInCreated is our type for the optin event
const TypeOptInCreated string = "optin_created"

// OptInCreatedEvent events are created when an action has created an optin to be sent.
//
//	{
//	  "type": "optin_created",
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
// @event optin_created
type OptInCreatedEvent struct {
	BaseEvent

	OptIn   *assets.OptInReference   `json:"optin" validate:"required,dive"`
	Channel *assets.ChannelReference `json:"channel" validate:"required,dive"`
	URN     urns.URN                 `json:"urn" validate:"required"`
}

// NewOptInCreated returns a new optin sent event
func NewOptInCreated(optIn *flows.OptIn, ch *flows.Channel, urn urns.URN) *OptInCreatedEvent {
	return &OptInCreatedEvent{
		BaseEvent: NewBaseEvent(TypeOptInCreated),
		OptIn:     optIn.Reference(),
		Channel:   ch.Reference(),
		URN:       urn,
	}
}
