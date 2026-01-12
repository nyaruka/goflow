package modifiers

import (
	"context"

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

	ticketUUID flows.TicketUUID
	note       string
}

// NewTicketNote creates a new note modifier
func NewTicketNote(ticketUUID flows.TicketUUID, note string) *TicketNote {
	return &TicketNote{
		baseModifier: newBaseModifier(TypeTicketNote),
		ticketUUID:   ticketUUID,
		note:         note,
	}
}

// Apply applies this modification to the given contact
func (m *TicketNote) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventLogger) (bool, error) {
	ticket := contact.Tickets().Find(m.ticketUUID)

	if ticket != nil {
		log(events.NewTicketNoteAdded(ticket.UUID(), m.note))
		return true, nil
	}
	return false, nil
}

var _ flows.Modifier = (*TicketNote)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketNoteEnvelope struct {
	utils.TypedEnvelope

	TicketUUID flows.TicketUUID `json:"ticket_uuid" validate:"required,uuid"`
	Note       string           `json:"note"`
}

func readTicketNote(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketNoteEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewTicketNote(e.TicketUUID, e.Note), nil
}

func (m *TicketNote) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketNoteEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		TicketUUID:    m.ticketUUID,
		Note:          m.note,
	})
}
