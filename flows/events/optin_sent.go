package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeOptInSent, func() flows.Event { return &OptInSentEvent{} })
}

// TypeOptInSent is our type for the optin event
const TypeOptInSent string = "optin_sent"

// OptInSentEvent events are created when an action has sent an optin.
//
//	{
//	  "type": "optin_sent",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "optin": {
//	    "uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
//	    "name": "Joke Of The Day"
//	  }
//	}
//
// @event optin_sent
type OptInSentEvent struct {
	BaseEvent

	OptIn *assets.OptInReference `json:"optin" validate:"required,dive"`
}

// NewOptInSent returns a new optin sent event
func NewOptInSent(optIn *flows.OptIn) *OptInSentEvent {
	return &OptInSentEvent{
		BaseEvent: NewBaseEvent(TypeOptInSent),
		OptIn:     optIn.Reference(),
	}
}
