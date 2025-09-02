package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeTicketAssigneeChanged, func() flows.Event { return &TicketAssigneeChanged{} })
}

// TypeTicketAssigneeChanged is the type for our ticket assigned events
const TypeTicketAssigneeChanged string = "ticket_assignee_changed"

// TicketAssigneeChanged events are created when a ticket is assigned or unassigned.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "ticket_assignee_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "ticket_uuid": "019905d4-5f7b-71b8-bcb8-6a68de2d91d2",
//	  "assignee": {
//	    "uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44",
//	    "name": "Bob"
//	  },
//	  "initial": true
//	}
//
// @event ticket_assignee_changed
type TicketAssigneeChanged struct {
	BaseEvent

	TicketUUID flows.TicketUUID      `json:"ticket_uuid" validate:"required,uuid"`
	Assignee   *assets.UserReference `json:"assignee"`
	Initial    bool                  `json:"initial,omitempty"`
}

// NewTicketAssigneeChanged returns a new ticket assignee changed event
func NewTicketAssigneeChanged(ticketUUID flows.TicketUUID, assignee *assets.UserReference, initial bool) *TicketAssigneeChanged {
	return &TicketAssigneeChanged{
		BaseEvent:  NewBaseEvent(TypeTicketAssigneeChanged),
		TicketUUID: ticketUUID,
		Assignee:   assignee,
		Initial:    initial,
	}
}
