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
	registerType(TypeTicketReopen, readTicketReopen)
}

// TypeTicketReopen is the type of our reopen modifier
const TypeTicketReopen string = "ticket_reopen"

// TicketReopen reopens a closed ticket
type TicketReopen struct {
	baseModifier

	ticketUUID flows.TicketUUID
}

// NewTicketReopen creates a new reopen modifier
func NewTicketReopen(ticketUUID flows.TicketUUID) *TicketReopen {
	return &TicketReopen{
		baseModifier: newBaseModifier(TypeTicketReopen),
		ticketUUID:   ticketUUID,
	}
}

// Apply applies this modification to the given contact
func (m *TicketReopen) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	// if there's already an open ticket, nothing to do
	if contact.Tickets().Open().Count() > 0 {
		return false
	}

	ticket := contact.Tickets().Find(m.ticketUUID)

	if ticket != nil && ticket.Status() != flows.TicketStatusOpen {
		ticket.SetStatus(flows.TicketStatusOpen)
		log(events.NewTicketReopened(ticket.UUID()))
		return true
	}
	return false
}

var _ flows.Modifier = (*TicketReopen)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketReopenEnvelope struct {
	utils.TypedEnvelope

	TicketUUID flows.TicketUUID `json:"ticket_uuid" validate:"required,uuid"`
}

func readTicketReopen(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketReopenEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewTicketReopen(e.TicketUUID), nil
}

func (m *TicketReopen) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketReopenEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		TicketUUID:    m.ticketUUID,
	})
}
