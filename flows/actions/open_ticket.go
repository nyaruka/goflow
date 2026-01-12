package actions

import (
	"context"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeOpenTicket, func() flows.Action { return &OpenTicket{} })
}

const (
	// TypeOpenTicket is the type for the open ticket action
	TypeOpenTicket string = "open_ticket"

	OpenTicketOutputLocal = "_new_ticket"
)

// OpenTicket is used to open a ticket for the contact if they don't already have an open ticket.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "open_ticket",
//	  "topic": {
//	    "uuid": "472a7a73-96cb-4736-b567-056d987cc5b4",
//	    "name": "Weather"
//	  },
//	  "assignee": {"uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44", "name": "Bob McTickets"},
//	  "note": "@input"
//	}
//
// @action open_ticket
type OpenTicket struct {
	baseAction
	onlineAction

	Topic    *assets.TopicReference `json:"topic"`
	Assignee *assets.UserReference  `json:"assignee,omitempty"`
	Note     string                 `json:"note,omitempty"     engine:"evaluated"`
}

// NewOpenTicket creates a new open ticket action
func NewOpenTicket(uuid flows.ActionUUID, topic *assets.TopicReference, assignee *assets.UserReference, note string) *OpenTicket {
	return &OpenTicket{
		baseAction: newBaseAction(TypeOpenTicket, uuid),
		Topic:      topic,
		Assignee:   assignee,
		Note:       note,
	}
}

// Execute runs this action
func (a *OpenTicket) Execute(ctx context.Context, run flows.Run, step flows.Step, log flows.EventLogger) error {
	sa := run.Session().Assets()

	// get topic or fallback to default
	var topic *flows.Topic
	if a.Topic != nil {
		topic = sa.Topics().Get(a.Topic.UUID)
	} else {
		topic = sa.Topics().FindByName("General")
	}

	var assignee *flows.User
	if a.Assignee != nil {
		assignee = resolveUser(run, a.Assignee, log)
	}

	evaluatedNote, _ := run.EvaluateTemplate(a.Note, log)
	evaluatedNote = strings.TrimSpace(evaluatedNote)

	ticket, err := a.open(ctx, run, topic, assignee, evaluatedNote, log)
	if err != nil {
		return err
	} else if ticket != nil {
		run.Locals().Set(OpenTicketOutputLocal, string(ticket.UUID()))
	} else {
		run.Locals().Set(OpenTicketOutputLocal, "")
	}

	return nil
}

func (a *OpenTicket) open(ctx context.Context, run flows.Run, topic *flows.Topic, assignee *flows.User, note string, log flows.EventLogger) (*flows.Ticket, error) {
	if run.Session().BatchStart() {
		log(events.NewError("can't open tickets during batch starts"))
		return nil, nil
	}

	if a.Topic != nil && topic == nil {
		log(events.NewDependencyError(a.Topic))
		return nil, nil
	}

	mod := modifiers.NewTicketOpen(topic, assignee, note)
	modified, err := a.applyModifier(ctx, run, mod, log)
	if err != nil {
		return nil, err
	}
	if modified {
		// if we were able to open a ticket, return it
		if lastOpen := run.Session().Contact().Tickets().LastOpen(); lastOpen != nil {
			return lastOpen, nil
		}
	}
	return nil, nil
}

func (a *OpenTicket) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	if a.Topic != nil {
		dependency(a.Topic)
	}
	if a.Assignee != nil {
		dependency(a.Assignee)
	}
	local(OpenTicketOutputLocal)
}
