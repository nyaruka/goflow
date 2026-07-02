package modifiers

import (
	"context"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/events"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeTicketClose, readTicketClose)
}

// TypeTicketClose is the type of our close modifier
const TypeTicketClose string = "ticket_close"

// TicketClose closes an open ticket
type TicketClose struct {
	baseModifier

	ticketUUID events.TicketUUID
}

// NewTicketClose creates a new close modifier
func NewTicketClose(ticketUUID events.TicketUUID) *TicketClose {
	return &TicketClose{
		baseModifier: newBaseModifier(TypeTicketClose),
		ticketUUID:   ticketUUID,
	}
}

// Apply applies this modification to the given contact
func (m *TicketClose) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log events.EventLogger) (bool, error) {
	ticket := contact.Tickets().Find(m.ticketUUID)

	if ticket != nil && ticket.Status() != events.TicketStatusClosed {
		ticket.SetStatus(events.TicketStatusClosed)
		log(events.NewTicketClosed(ticket.UUID()))
		return true, nil
	}
	return false, nil
}

var _ flows.Modifier = (*TicketClose)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketCloseEnvelope struct {
	utils.TypedEnvelope

	TicketUUID events.TicketUUID `json:"ticket_uuid" validate:"required,uuid"`
}

func readTicketClose(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketCloseEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewTicketClose(e.TicketUUID), nil
}

func (m *TicketClose) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketCloseEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		TicketUUID:    m.ticketUUID,
	})
}
