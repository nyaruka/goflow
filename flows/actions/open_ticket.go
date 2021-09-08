package actions

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/pkg/errors"
)

func init() {
	registerType(TypeOpenTicket, func() flows.Action { return &OpenTicketAction{} })
}

// TypeOpenTicket is the type for the open ticket action
const TypeOpenTicket string = "open_ticket"

// OpenTicketAction is used to open a ticket for the contact.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "open_ticket",
//     "ticketer": {
//       "uuid": "19dc6346-9623-4fe4-be80-538d493ecdf5",
//       "name": "Support Tickets"
//     },
//     "topic": {
//       "uuid": "472a7a73-96cb-4736-b567-056d987cc5b4",
//       "name": "Weather"
//     },
//     "body": "@input",
//     "assignee": {"email": "bob@nyaruka.com", "name": "Bob McTickets"},
//     "result_name": "Help Ticket"
//   }
//
// @action open_ticket
type OpenTicketAction struct {
	baseAction
	onlineAction

	Ticketer   *assets.TicketerReference `json:"ticketer" validate:"required,dive"`
	Topic      *assets.TopicReference    `json:"topic" validate:"omitempty,dive"`
	Subject    string                    `json:"subject" engine:"evaluated"`
	Body       string                    `json:"body" engine:"evaluated"`
	Assignee   *assets.UserReference     `json:"assignee" validate:"omitempty,dive"`
	ResultName string                    `json:"result_name" validate:"required"`
}

// NewOpenTicket creates a new open ticket action
func NewOpenTicket(uuid flows.ActionUUID, ticketer *assets.TicketerReference, topic *assets.TopicReference, subject, body string, assignee *assets.UserReference, resultName string) *OpenTicketAction {
	return &OpenTicketAction{
		baseAction: newBaseAction(TypeOpenTicket, uuid),
		Ticketer:   ticketer,
		Topic:      topic,
		Subject:    subject,
		Body:       body,
		Assignee:   assignee,
		ResultName: resultName,
	}
}

// Validate validates our action is valid
func (a *OpenTicketAction) Validate() error {
	if (a.Subject == "" && a.Topic == nil) || (a.Subject != "" && a.Topic != nil) {
		return errors.New("must have one of 'subject' or 'topic'")
	}

	return nil
}

// Execute runs this action
func (a *OpenTicketAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	sa := run.Session().Assets()

	ticketer := sa.Ticketers().Get(a.Ticketer.UUID)

	var topic *flows.Topic
	if a.Topic != nil {
		topic = sa.Topics().Get(a.Topic.UUID)
	} else {
		topic = sa.Topics().FindByName("General") // TODO remove when editor adds topics
	}
	var assignee *flows.User
	if a.Assignee != nil {
		assignee = resolveUser(run, a.Assignee, logEvent)
	}

	evaluatedSubject, err := run.EvaluateTemplate(a.Subject)
	if err != nil {
		logEvent(events.NewError(err))
	}
	evaluatedBody, err := run.EvaluateTemplate(a.Body)
	if err != nil {
		logEvent(events.NewError(err))
	}

	ticket := a.open(run, step, ticketer, topic, evaluatedSubject, evaluatedBody, assignee, logEvent)
	if ticket != nil {
		a.saveResult(run, step, a.ResultName, string(ticket.UUID()), CategorySuccess, "", "", nil, logEvent)
	} else {
		a.saveResult(run, step, a.ResultName, "", CategoryFailure, "", "", nil, logEvent)
	}

	return nil
}

func (a *OpenTicketAction) open(run flows.FlowRun, step flows.Step, ticketer *flows.Ticketer, topic *flows.Topic, subject, body string, assignee *flows.User, logEvent flows.EventCallback) *flows.Ticket {
	if run.Session().BatchStart() {
		logEvent(events.NewErrorf("can't open tickets during batch starts"))
		return nil
	}

	if ticketer == nil {
		logEvent(events.NewDependencyError(a.Ticketer))
		return nil
	}
	if a.Topic != nil && topic == nil {
		logEvent(events.NewDependencyError(a.Topic))
		return nil
	}

	svc, err := run.Session().Engine().Services().Ticket(run.Session(), ticketer)
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	httpLogger := &flows.HTTPLogger{}

	ticket, err := svc.Open(run.Session(), topic, subject, body, assignee, httpLogger.Log)
	if err != nil {
		logEvent(events.NewError(err))
	}
	if len(httpLogger.Logs) > 0 {
		logEvent(events.NewTicketerCalled(ticketer.Reference(), httpLogger.Logs))
	}
	if ticket != nil {
		logEvent(events.NewTicketOpened(ticket))

		run.Contact().Tickets().Add(ticket)
	}

	return ticket
}
