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
	registerType(TypeTicketAssignee, readTicketAssignee)
}

// TypeTicketAssignee is the type of our assignee modifier
const TypeTicketAssignee string = "ticket_assignee"

// TicketAssignee modifies the assignee of a ticket
type TicketAssignee struct {
	baseModifier

	assignee *flows.User
}

// NewTicketAssignee creates a new assignee modifier
func NewTicketAssignee(assignee *flows.User) *TicketAssignee {
	return &TicketAssignee{
		baseModifier: newBaseModifier(TypeTicketAssignee),
		assignee:     assignee,
	}
}

// Apply applies this modification to the given ticket
func (m *TicketAssignee) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, ticket *flows.Ticket, log flows.EventCallback) bool {
	if ticket != nil && ticket.Assignee() != m.assignee {
		prevAssignee := ticket.Assignee().Reference()
		thisAssignee := m.assignee.Reference()

		ticket.SetAssignee(m.assignee)
		log(events.NewTicketAssigneeChanged(ticket.UUID(), thisAssignee, prevAssignee))
		return true
	}
	return false
}

var _ flows.Modifier = (*TicketAssignee)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketAssigneeEnvelope struct {
	utils.TypedEnvelope

	Assignee *assets.UserReference `json:"assignee"`
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

	return NewTicketAssignee(assignee), nil
}

func (m *TicketAssignee) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketAssigneeEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Assignee:      m.assignee.Reference(),
	})
}
