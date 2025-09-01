package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketAssigned, func() flows.Event { return &TicketAssigned{} })
}

// TypeTicketAssigned is the type for our ticket assigned events
const TypeTicketAssigned string = "ticket_assigned"

// TicketAssigned events are created when a ticket is assigned.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "ticket_assigned",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "ticket_uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2",
//	  "assignee": {
//	    "uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44",
//	    "name": "Bob"
//	  }
//	}
//
// @event ticket_assigned
type TicketAssigned struct {
	BaseEvent

	TicketUUID flows.TicketUUID      `json:"ticket_uuid" validate:"required,uuid"`
	Assignee   *assets.UserReference `json:"assignee"`
}

// NewTicketAssigned returns a new ticket assigned event
func NewTicketAssigned(ticketUUID flows.TicketUUID, assignee *assets.UserReference) *TicketAssigned {
	return &TicketAssigned{
		BaseEvent:  NewBaseEvent(TypeTicketAssigned),
		TicketUUID: ticketUUID,
		Assignee:   assignee,
	}
}
