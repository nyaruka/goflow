package events

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactLastSeenChanged, func() flows.Event { return &ContactLastSeenChanged{} })
}

// TypeContactLastSeenChanged is the type of our contact last seen changed event
const TypeContactLastSeenChanged string = "contact_last_seen_changed"

// ContactLastSeenChanged events are created when the last seen on of the contact has been changed.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "contact_last_seen_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "seen_on": "2026-02-03T15:04:05Z"
//	}
//
// @event contact_last_seen_changed
type ContactLastSeenChanged struct {
	BaseEvent

	LastSeenOn time.Time `json:"last_seen_on"`
}

// NewContactLastSeenChanged returns a new contact last seen changed event
func NewContactLastSeenChanged(seen time.Time) *ContactLastSeenChanged {
	return &ContactLastSeenChanged{
		BaseEvent:  NewBaseEvent(TypeContactLastSeenChanged),
		LastSeenOn: seen,
	}
}
