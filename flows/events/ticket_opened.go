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
//       "body": "Where are my cookies?"
//     },
//     "http_logs": [
//       {
//         "url": "https://api.zendesk.com/new_ticket",
//         "status": "success",
//         "request": "POST /new_ticket HTTP/1.1",
//         "response": "HTTP/1.1 200 OK\r\n\r\n",
//         "created_on": "2020-04-23T15:04:05Z",
//         "elapsed_ms": 123
//       }
//     ]
//   }
//
// @event ticket_opened
type TicketOpenedEvent struct {
	baseEvent

	Ticket   *flows.Ticket    `json:"ticket"`
	HTTPLogs []*flows.HTTPLog `json:"http_logs,omitempty"`
}

// NewTicketOpened returns a new ticket opened event
func NewTicketOpened(ticket *flows.Ticket, httpLogs []*flows.HTTPLog) *TicketOpenedEvent {
	return &TicketOpenedEvent{
		baseEvent: newBaseEvent(TypeTicketOpened),
		Ticket:    ticket,
		HTTPLogs:  httpLogs,
	}
}
