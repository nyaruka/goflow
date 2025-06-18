package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeOptInStarted, func() flows.Event { return &OptInStartedEvent{} })
}

// TypeOptInStarted is our type for the optin started event
const TypeOptInStarted string = "optin_started"

// OptInStartedEvent events are created when a contact has opted-in.
//
//	{
//	  "type": "optin_started",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "optin": {
//	    "uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
//	    "name": "Joke Of The Day"
//	  },
//	  "channel": {
//	    "uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648",
//	    "name": "Facebook"
//	  }
//	}
//
// @event optin_started
type OptInStartedEvent struct {
	BaseEvent

	OptIn   *assets.OptInReference   `json:"optin" validate:"required"`
	Channel *assets.ChannelReference `json:"channel,omitempty"` // TODO make required
}

// NewOptInStarted returns a new optin started event
func NewOptInStarted(optIn *flows.OptIn, ch *assets.ChannelReference) *OptInStartedEvent {
	return &OptInStartedEvent{
		BaseEvent: NewBaseEvent(TypeOptInStarted),
		OptIn:     optIn.Reference(),
		Channel:   ch,
	}
}
