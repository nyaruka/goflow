package tickets

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeAssignee is the type of our assignee modifier
const TypeAssignee string = "assignee"

// Assignee modifies the assignee of a ticket
type Assignee struct {
	baseModifier

	assignee *flows.User
}

// NewAssignee creates a new assignee modifier
func NewAssignee(assignee *flows.User) *Assignee {
	return &Assignee{
		baseModifier: newBaseModifier(TypeAssignee),
		assignee:     assignee,
	}
}

// Apply applies this modification to the given ticket
func (m *Assignee) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, ticket *flows.Ticket, log flows.EventCallback) bool {
	if ticket.Assignee() != m.assignee {
		var prevAssignee, thisAssignee *assets.UserReference
		if ticket.Assignee() != nil {
			prevAssignee = ticket.Assignee().Reference()
		}
		if m.assignee != nil {
			thisAssignee = m.assignee.Reference()
		}

		ticket.SetAssignee(m.assignee)
		log(events.NewTicketAssigneeChanged(ticket.UUID(), prevAssignee, thisAssignee))
		return true
	}
	return false
}

var _ flows.TicketModifier = (*Assignee)(nil)
