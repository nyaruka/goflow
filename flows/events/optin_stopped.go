package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeOptInStopped, func() flows.Event { return &OptInStopped{} })
}

// TypeOptInStopped is our type for the optin stopped event
const TypeOptInStopped string = "optin_stopped"

// OptInStopped events are created when a contact has opted-out.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
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
type OptInStopped struct {
	BaseEvent

	OptIn   *assets.OptInReference   `json:"optin"   validate:"required"`
	Channel *assets.ChannelReference `json:"channel" validate:"required"`
}

// NewOptInStopped returns a new optin stopped event
func NewOptInStopped(optIn *flows.OptIn, ch *assets.ChannelReference) *OptInStopped {
	return &OptInStopped{
		BaseEvent: NewBaseEvent(TypeOptInStopped),
		OptIn:     optIn.Reference(),
		Channel:   ch,
	}
}
