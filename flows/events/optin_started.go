package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeOptInStarted, func() flows.Event { return &OptInStarted{} })
}

// TypeOptInStarted is our type for the optin started event
const TypeOptInStarted string = "optin_started"

// OptInStarted events are created when a contact has opted-in.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
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
type OptInStarted struct {
	BaseEvent

	OptIn   *assets.OptInReference   `json:"optin"   validate:"required"`
	Channel *assets.ChannelReference `json:"channel" validate:"required"`
}

// NewOptInStarted returns a new optin started event
func NewOptInStarted(optIn *assets.OptInReference, ch *assets.ChannelReference) *OptInStarted {
	return &OptInStarted{
		BaseEvent: NewBaseEvent(TypeOptInStarted),
		OptIn:     optIn,
		Channel:   ch,
	}
}
