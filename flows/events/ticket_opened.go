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
//     "ticket_id": "234562",
//     "subject": "Need help",
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

	TicketID string           `json:"ticket_id"`
	Subject  string           `json:"subject"`
	HTTPLogs []*flows.HTTPLog `json:"http_logs"`
}

// NewTicketOpened returns a new ticket opened event
func NewTicketOpened(ticket *flows.Ticket, httpLogs []*flows.HTTPLog) *TicketOpenedEvent {
	return &TicketOpenedEvent{
		baseEvent: newBaseEvent(TypeTicketOpened),
		TicketID:  ticket.ID,
		Subject:   ticket.Subject,
		HTTPLogs:  httpLogs,
	}
}