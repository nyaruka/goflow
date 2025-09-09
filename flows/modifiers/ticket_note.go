package modifiers

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeTicketNote, readTicketNote)
}

// TypeTicketNote is the type of our note modifier
const TypeTicketNote string = "ticket_note"

// TicketNote adds a note to a ticket
type TicketNote struct {
	baseModifier

	note string
}

// NewTicketNote creates a new note modifier
func NewTicketNote(note string) *TicketNote {
	return &TicketNote{
		baseModifier: newBaseModifier(TypeTicketNote),
		note:         note,
	}
}

// Apply applies this modification to the given ticket
func (m *TicketNote) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, ticket *flows.Ticket, log flows.EventCallback) bool {
	if ticket != nil {
		log(events.NewTicketNoteAdded(ticket.UUID(), m.note))
	}
	return true
}

var _ flows.Modifier = (*TicketNote)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketNoteEnvelope struct {
	utils.TypedEnvelope

	Note string `json:"note"`
}

func readTicketNote(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketNoteEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewTicketNote(e.Note), nil
}

func (m *TicketNote) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketNoteEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Note:          m.note,
	})
}
