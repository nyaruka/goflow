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
	registerType(TypeOpenTicket, func() flows.Action { return &OpenTicketAction{} })
}

// TypeOpenTicket is the type for the open ticket action
const TypeOpenTicket string = "open_ticket"

// OpenTicketAction is used to open a ticket for the contact if they don't already have an open ticket.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "open_ticket",
//	  "topic": {
//	    "uuid": "472a7a73-96cb-4736-b567-056d987cc5b4",
//	    "name": "Weather"
//	  },
//	  "note": "@input",
//	  "assignee": {"email": "bob@nyaruka.com", "name": "Bob McTickets"}
//	}
//
// @action open_ticket
type OpenTicketAction struct {
	baseAction
	onlineAction

	Topic      *assets.TopicReference `json:"topic"`
	Note       string                 `json:"note"                  engine:"evaluated"`
	Assignee   *assets.UserReference  `json:"assignee,omitempty"`
	ResultName string                 `json:"result_name,omitempty"             validate:"omitempty,result_name"`
}

// NewOpenTicket creates a new open ticket action
func NewOpenTicket(uuid flows.ActionUUID, topic *assets.TopicReference, note string, assignee *assets.UserReference, resultName string) *OpenTicketAction {
	return &OpenTicketAction{
		baseAction: newBaseAction(TypeOpenTicket, uuid),
		Topic:      topic,
		Note:       note,
		Assignee:   assignee,
		ResultName: resultName,
	}
}

// Execute runs this action
func (a *OpenTicketAction) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
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
		assignee = resolveUser(run, a.Assignee, logEvent)
	}

	evaluatedNote, _ := run.EvaluateTemplate(a.Note, logEvent)
	evaluatedNote = strings.TrimSpace(evaluatedNote)

	ticket := a.open(run, topic, assignee, evaluatedNote, logModifier, logEvent)
	if ticket != nil {
		run.Locals().Set("_new_ticket", string(ticket.UUID()))
	} else {
		run.Locals().Set("_new_ticket", "")
	}

	if a.ResultName != "" {
		if ticket != nil {
			a.saveResult(run, step, a.ResultName, string(ticket.UUID()), CategorySuccess, "", "", nil, logEvent)
		} else {
			a.saveResult(run, step, a.ResultName, "", CategoryFailure, "", "", nil, logEvent)
		}
	}

	return nil
}

func (a *OpenTicketAction) open(run flows.Run, topic *flows.Topic, assignee *flows.User, note string, logModifier flows.ModifierCallback, logEvent flows.EventCallback) *flows.Ticket {
	if run.Session().BatchStart() {
		logEvent(events.NewError("can't open tickets during batch starts"))
		return nil
	}

	if a.Topic != nil && topic == nil {
		logEvent(events.NewDependencyError(a.Topic))
		return nil
	}

	mod := modifiers.NewTicket(topic, assignee, note)

	if a.applyModifier(run, mod, logModifier, logEvent) {
		// if we were able to open a ticket, return it
		return run.Session().Contact().Ticket()
	}
	return nil
}

func (a *OpenTicketAction) Inspect(result func(*flows.ResultInfo), dependency func(assets.Reference)) {
	if a.Topic != nil {
		dependency(a.Topic)
	}
	if a.Assignee != nil {
		dependency(a.Assignee)
	}
}
