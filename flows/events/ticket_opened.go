package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketOpened, func() flows.Event { return &TicketOpened{} })
}

// TypeTicketOpened is the type for our ticket opened events
const TypeTicketOpened string = "ticket_opened"

// TicketOpened events are created when a new ticket is opened.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "ticket_opened",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "ticket": {
//	    "uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2",
//	    "status": "open",
//	    "topic": {
//	      "uuid": "add17edf-0b6e-4311-bcd7-a64b2a459157",
//	      "name": "Weather"
//	    },
//	    "assignee": {"uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44", "name": "Bob"}
//	  },
//	  "note": "this is weird"
//	}
//
// @event ticket_opened
type TicketOpened struct {
	BaseEvent

	Ticket *flows.TicketEnvelope `json:"ticket"`
	Note   string                `json:"note,omitempty"`
}

// NewTicketOpened returns a new ticket opened event
func NewTicketOpened(ticket *flows.Ticket, note string) *TicketOpened {
	return &TicketOpened{
		BaseEvent: NewBaseEvent(TypeTicketOpened),
		Ticket:    ticket.Marshal(),
		Note:      note,
	}
}
