package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeContactNameChanged, func() flows.Event { return &ContactNameChanged{} })
}

// TypeContactNameChanged is the type of our contact name changed event
const TypeContactNameChanged string = "contact_name_changed"

// ContactNameChanged events are created when the name of the contact has been changed.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "contact_name_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "name": "Bob Smith"
//	}
//
// @event contact_name_changed
type ContactNameChanged struct {
	BaseEvent

	Name string `json:"name"`
}

// NewContactNameChanged returns a new contact name changed event
func NewContactNameChanged(name string) *ContactNameChanged {
	return &ContactNameChanged{
		BaseEvent: NewBaseEvent(TypeContactNameChanged),
		Name:      name,
	}
}
