package events

import "github.com/nyaruka/goflow/flows"

func init() {
	registerType(TypeContactStatusChanged, func() flows.Event { return &ContactStatusChanged{} })
}

// TypeContactStatusChanged is the type of our contact status changed event
const TypeContactStatusChanged string = "contact_status_changed"

// ContactStatusChanged events are created when the status of the contact has been changed.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "contact_timezone_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "status": "blocked"
//	}
//
// @event contact_status_changed
type ContactStatusChanged struct {
	BaseEvent

	Status flows.ContactStatus `json:"status"`
}

// NewContactStatusChanged returns a new contact_status_changed event
func NewContactStatusChanged(status flows.ContactStatus) *ContactStatusChanged {
	return &ContactStatusChanged{
		BaseEvent: NewBaseEvent(TypeContactStatusChanged),
		Status:    status,
	}
}
