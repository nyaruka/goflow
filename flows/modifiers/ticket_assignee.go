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
	registerType(TypeTicketAssignee, readTicketAssignee)
}

// TypeTicketAssignee is the type of our assignee modifier
const TypeTicketAssignee string = "ticket_assignee"

// TicketAssignee modifies the assignee of tickets
type TicketAssignee struct {
	baseModifier

	ticketUUIDs []flows.TicketUUID
	assignee    *flows.User
}

// NewTicketAssignee creates a new assignee modifier
func NewTicketAssignee(ticketUUIDs []flows.TicketUUID, assignee *flows.User) *TicketAssignee {
	return &TicketAssignee{
		baseModifier: newBaseModifier(TypeTicketAssignee),
		ticketUUIDs:  ticketUUIDs,
		assignee:     assignee,
	}
}

// Apply applies this modification to the given contact
func (m *TicketAssignee) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	modified := false

	for _, ticket := range contact.Tickets().All() {
		if slices.Contains(m.ticketUUIDs, ticket.UUID()) && ticket.Assignee() != m.assignee {
			prevAssignee := ticket.Assignee().Reference()
			thisAssignee := m.assignee.Reference()

			ticket.SetAssignee(m.assignee)
			log(events.NewTicketAssigneeChanged(ticket.UUID(), thisAssignee, prevAssignee))
			modified = true
		}
	}
	return modified
}

var _ flows.Modifier = (*TicketAssignee)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketAssigneeEnvelope struct {
	utils.TypedEnvelope

	TicketUUIDs []flows.TicketUUID    `json:"ticket_uuids" validate:"required,dive,uuid"`
	Assignee    *assets.UserReference `json:"assignee"`
}

func readTicketAssignee(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketAssigneeEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var assignee *flows.User
	if e.Assignee != nil {
		assignee = sa.Users().Get(e.Assignee.UUID)
		if assignee == nil {
			missing(e.Assignee, nil)
		}
	}

	return NewTicketAssignee(e.TicketUUIDs, assignee), nil
}

func (m *TicketAssignee) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketAssigneeEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		TicketUUIDs:   m.ticketUUIDs,
		Assignee:      m.assignee.Reference(),
	})
}
