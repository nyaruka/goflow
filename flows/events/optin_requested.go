package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeOptInRequested, func() flows.Event { return &OptInRequested{} })
}

// TypeOptInRequested is our type for the optin event
const TypeOptInRequested string = "optin_requested"

// OptInRequested events are created when an action has created an optin to be sent.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
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
type OptInRequested struct {
	BaseEvent

	OptIn   *assets.OptInReference   `json:"optin" validate:"required"`
	Channel *assets.ChannelReference `json:"channel" validate:"required"`
	URN     urns.URN                 `json:"urn" validate:"required"`
}

// NewOptInRequested returns a new optin requested event
func NewOptInRequested(optIn *flows.OptIn, ch *assets.ChannelReference, urn urns.URN) *OptInRequested {
	return &OptInRequested{
		BaseEvent: NewBaseEvent(TypeOptInRequested),
		OptIn:     optIn.Reference(),
		Channel:   ch,
		URN:       urn,
	}
}
