package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketNoteAdded, func() flows.Event { return &TicketNoteAddedEvent{} })
}

// TypeTicketNoteAdded is the type for our ticket note added events
const TypeTicketNoteAdded string = "ticket_note_added"

// TicketNoteAddedEvent events are created when a note is added to the currently open ticket.
//
//	{
//	  "type": "ticket_note_added",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "note": "this is weird"
//	}
//
// @event ticket_note_added
type TicketNoteAddedEvent struct {
	BaseEvent

	Note string `json:"note"`
}

// NewTicketNoteAdded returns a new ticket note added event
func NewTicketNoteAdded(note string) *TicketNoteAddedEvent {
	return &TicketNoteAddedEvent{
		BaseEvent: NewBaseEvent(TypeTicketNoteAdded),
		Note:      note,
	}
}
