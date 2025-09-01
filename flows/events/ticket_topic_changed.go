package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketTopicChanged, func() flows.Event { return &TicketTopicChanged{} })
}

// TypeTicketTopicChanged is the type for our ticket topic changed events
const TypeTicketTopicChanged string = "ticket_topic_changed"

// TicketTopicChanged events are created when a ticket's topic is changed.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "ticket_topic_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "ticket_uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2",
//	  "topic": {
//	    "uuid": "add17edf-0b6e-4311-bcd7-a64b2a459157",
//	    "name": "Weather"
//	  }
//	}
//
// @event ticket_topic_changed
type TicketTopicChanged struct {
	BaseEvent

	TicketUUID flows.TicketUUID       `json:"ticket_uuid" validate:"required,uuid"`
	Topic      *assets.TopicReference `json:"topic"       validate:"required"`
}

// NewTicketTopicChanged returns a new ticket topic changed event
func NewTicketTopicChanged(ticketUUID flows.TicketUUID, topic *assets.TopicReference) *TicketTopicChanged {
	return &TicketTopicChanged{
		BaseEvent:  NewBaseEvent(TypeTicketTopicChanged),
		TicketUUID: ticketUUID,
		Topic:      topic,
	}
}
