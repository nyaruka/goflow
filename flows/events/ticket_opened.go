package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketOpened, func() flows.Event { return &TicketOpenedEvent{} })
}

// TypeTicketOpened is the type for our ticket opened events
const TypeTicketOpened string = "ticket_opened"

type Ticket struct {
	UUID       flows.TicketUUID          `json:"uuid"                   validate:"required,uuid4"`
	Ticketer   *assets.TicketerReference `json:"ticketer"               validate:"required,dive"`
	Subject    string                    `json:"subject"`
	Body       string                    `json:"body"`
	ExternalID string                    `json:"external_id,omitempty"`
	Assignee   *assets.UserReference     `json:"assignee,omitempty"     validate:"omitempty,dive"`
}

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
//       "external_id": "32526523",
//       "assignee": {"email": "bob@nyaruka.com", "name": "Bob"}
//     }
//   }
//
// @event ticket_opened
type TicketOpenedEvent struct {
	baseEvent

	Ticket *Ticket `json:"ticket"`
}

// NewTicketOpened returns a new ticket opened event
func NewTicketOpened(ticket *flows.Ticket) *TicketOpenedEvent {
	return &TicketOpenedEvent{
		baseEvent: newBaseEvent(TypeTicketOpened),
		Ticket: &Ticket{
			UUID:       ticket.UUID(),
			Ticketer:   ticket.Ticketer().Reference(),
			Subject:    ticket.Subject(),
			Body:       ticket.Body(),
			ExternalID: ticket.ExternalID(),
			Assignee:   ticket.Assignee().Reference(),
		},
	}
}
