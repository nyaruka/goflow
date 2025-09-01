package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketReopened, func() flows.Event { return &TicketReopened{} })
}

// TypeTicketReopened is the type for our ticket reopened events
const TypeTicketReopened string = "ticket_reopened"

// TicketReopened events are created when a ticket is reopened.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "ticket_reopened",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "ticket_uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2"
//	}
//
// @event ticket_reopened
type TicketReopened struct {
	BaseEvent

	TicketUUID flows.TicketUUID `json:"ticket_uuid" validate:"required,uuid"`
}

// NewTicketReopened returns a new ticket reopened event
func NewTicketReopened(ticketUUID flows.TicketUUID) *TicketReopened {
	return &TicketReopened{
		BaseEvent:  NewBaseEvent(TypeTicketReopened),
		TicketUUID: ticketUUID,
	}
}
