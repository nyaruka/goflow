package tickets

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeNote is the type of our note modifier
const TypeNote string = "note"

// Note adds a note to a ticket
type Note struct {
	baseModifier

	note string
}

// NewNote creates a new note modifier
func NewNote(note string) *Note {
	return &Note{
		baseModifier: newBaseModifier(TypeNote),
		note:         note,
	}
}

// Apply applies this modification to the given ticket
func (m *Note) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, ticket *flows.Ticket, log flows.EventCallback) bool {
	log(events.NewTicketNoteAdded(ticket.UUID(), m.note))
	return true
}

var _ flows.TicketModifier = (*Note)(nil)
