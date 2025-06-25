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
//	  "type": "ticket_closed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "ticket": {
//	    "uuid": "2e677ae6-9b57-423c-b022-7950503eef35",
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

	Ticket *flows.TicketEnvelope `json:"ticket"`
}

// NewTicketClosed returns a new ticket closed event
func NewTicketClosed(ticket *flows.Ticket) *TicketClosed {
	return &TicketClosed{
		BaseEvent: NewBaseEvent(TypeTicketClosed),
		Ticket:    ticket.Marshal(),
	}
}
