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
	registerType(TypeTicketClose, readTicketClose)
}

// TypeTicketClose is the type of our close modifier
const TypeTicketClose string = "ticket_close"

// TicketClose closes tickets
type TicketClose struct {
	baseModifier

	ticketUUIDs []flows.TicketUUID
}

// NewTicketClose creates a new close modifier
func NewTicketClose(ticketUUIDs []flows.TicketUUID) *TicketClose {
	return &TicketClose{
		baseModifier: newBaseModifier(TypeTicketClose),
		ticketUUIDs:  ticketUUIDs,
	}
}

// Apply applies this modification to the given contact
func (m *TicketClose) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	modified := false

	for _, ticket := range contact.Tickets().All() {
		if slices.Contains(m.ticketUUIDs, ticket.UUID()) && ticket.Status() != flows.TicketStatusClosed {
			ticket.SetStatus(flows.TicketStatusClosed)
			log(events.NewTicketClosed(ticket))
			modified = true
		}
	}
	return modified
}

var _ flows.Modifier = (*TicketClose)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketCloseEnvelope struct {
	utils.TypedEnvelope

	TicketUUIDs []flows.TicketUUID `json:"ticket_uuids" validate:"required,dive,uuid"`
}

func readTicketClose(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketCloseEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewTicketClose(e.TicketUUIDs), nil
}

func (m *TicketClose) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketCloseEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		TicketUUIDs:   m.ticketUUIDs,
	})
}
