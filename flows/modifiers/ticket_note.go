package modifiers

import (
	"slices"

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

// TicketNote adds a note to tickets
type TicketNote struct {
	baseModifier

	ticketUUIDs []flows.TicketUUID
	note        string
}

// NewTicketNote creates a new note modifier
func NewTicketNote(ticketUUIDs []flows.TicketUUID, note string) *TicketNote {
	return &TicketNote{
		baseModifier: newBaseModifier(TypeTicketNote),
		ticketUUIDs:  ticketUUIDs,
		note:         note,
	}
}

// Apply applies this modification to the given contact
func (m *TicketNote) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	modified := false

	for _, ticket := range contact.Tickets().All() {
		if slices.Contains(m.ticketUUIDs, ticket.UUID()) {
			log(events.NewTicketNoteAdded(ticket.UUID(), m.note))
			modified = true
		}
	}
	return modified
}

var _ flows.Modifier = (*TicketNote)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketNoteEnvelope struct {
	utils.TypedEnvelope

	TicketUUIDs []flows.TicketUUID `json:"ticket_uuids" validate:"required,dive,uuid"`
	Note        string             `json:"note"`
}

func readTicketNote(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketNoteEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewTicketNote(e.TicketUUIDs, e.Note), nil
}

func (m *TicketNote) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketNoteEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		TicketUUIDs:   m.ticketUUIDs,
		Note:          m.note,
	})
}
