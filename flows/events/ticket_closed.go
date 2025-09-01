package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketClosed, func() flows.Event { return &TicketClosed{} })
}

// TypeTicketClosed is the type for our ticket closed events
const TypeTicketClosed string = "ticket_closed"

// TicketClosed events are created when a ticket is closed.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "ticket_closed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "ticket_uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2",
//	  "ticket": {
//	    "uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2",
//	    "topic": {
//	      "uuid": "add17edf-0b6e-4311-bcd7-a64b2a459157",
//	      "name": "Weather"
//	    },
//	    "assignee": {"uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44", "name": "Bob"}
//	  }
//	}
//
// @event ticket_closed
type TicketClosed struct {
	BaseEvent

	TicketUUID flows.TicketUUID      `json:"ticket_uuid" validate:"omitempty,uuid"`
	Ticket     *flows.TicketEnvelope `json:"ticket,omitempty"` // deprecated
}

// NewTicketClosed returns a new ticket closed event
func NewTicketClosed(ticket *flows.Ticket) *TicketClosed {
	return &TicketClosed{
		BaseEvent:  NewBaseEvent(TypeTicketClosed),
		TicketUUID: ticket.UUID(),
		Ticket:     ticket.Marshal(),
	}
}
