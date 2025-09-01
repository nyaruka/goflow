package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketNoteAdded, func() flows.Event { return &TicketNoteAdded{} })
}

// TypeTicketNoteAdded is the type for our ticket note added events
const TypeTicketNoteAdded string = "ticket_note_added"

// TicketNoteAdded events are created when a note is added to a ticket.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "ticket_note_added",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "ticket_uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2",
//	  "note": "This looks important!"
//	}
//
// @event ticket_note_added
type TicketNoteAdded struct {
	BaseEvent

	TicketUUID flows.TicketUUID `json:"ticket_uuid" validate:"required,uuid"`
	Note       string           `json:"note"        validate:"required"`
}

// NewTicketNoteAdded returns a new ticket note added event
func NewTicketNoteAdded(ticketUUID flows.TicketUUID, note string) *TicketNoteAdded {
	return &TicketNoteAdded{
		BaseEvent:  NewBaseEvent(TypeTicketNoteAdded),
		TicketUUID: ticketUUID,
		Note:       note,
	}
}
