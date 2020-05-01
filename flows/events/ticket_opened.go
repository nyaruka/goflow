package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketOpened, func() flows.Event { return &TicketOpenedEvent{} })
}

// TypeTicketOpened is the type for our ticket opened events
const TypeTicketOpened string = "ticket_opened"

// TicketOpenedEvent events are created when a new ticket is opened.
//
//   {
//     "type": "ticket_opened",
//     "created_on": "2006-01-02T15:04:05Z",
//     "ticket": {
//       "uuid": "2e677ae6-9b57-423c-b022-7950503eef35",
//       "ticketer": {
//         "uuid": "d605bb96-258d-4097-ad0a-080937db2212",
//         "name": "Support Tickets"
//       },
//       "subject": "Need help",
//       "body": "Where are my cookies?",
//       "external_id": "32526523"
//     }
//   }
//
// @event ticket_opened
type TicketOpenedEvent struct {
	baseEvent

	Ticket *flows.Ticket `json:"ticket"`
}

// NewTicketOpened returns a new ticket opened event
func NewTicketOpened(ticket *flows.Ticket) *TicketOpenedEvent {
	return &TicketOpenedEvent{
		baseEvent: newBaseEvent(TypeTicketOpened),
		Ticket:    ticket,
	}
}
