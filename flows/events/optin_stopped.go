package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeOptInStopped, func() flows.Event { return &OptInStoppedEvent{} })
}

// TypeOptInStopped is our type for the optin stopped event
const TypeOptInStopped string = "optin_stopped"

// OptInStoppedEvent events are created when a contact has opted-out.
//
//	{
//	  "type": "optin_stopped",
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
// @event optin_stopped
type OptInStoppedEvent struct {
	BaseEvent

	OptIn   *assets.OptInReference   `json:"optin" validate:"required"`
	Channel *assets.ChannelReference `json:"channel,omitempty"` // TODO make required
}

// NewOptInStopped returns a new optin stopped event
func NewOptInStopped(optIn *flows.OptIn, ch *assets.ChannelReference) *OptInStoppedEvent {
	return &OptInStoppedEvent{
		BaseEvent: NewBaseEvent(TypeOptInStopped),
		OptIn:     optIn.Reference(),
		Channel:   ch,
	}
}
