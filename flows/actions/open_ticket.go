package actions

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
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
//     "subject": "Needs help",
//     "body": "@input",
//     "result_name": "Help Ticket"
//   }
//
// @action open_ticket
type OpenTicketAction struct {
	baseAction
	onlineAction

	Ticketer   *assets.TicketerReference `json:"ticketer" validate:"required"`
	Subject    string                    `json:"subject" validate:"required" engine:"evaluated"`
	Body       string                    `json:"body" engine:"evaluated"`
	ResultName string                    `json:"result_name" validate:"required"`
}

// NewOpenTicket creates a new open ticket action
func NewOpenTicket(uuid flows.ActionUUID, ticketer *assets.TicketerReference, subject, body, resultName string) *OpenTicketAction {
	return &OpenTicketAction{
		baseAction: newBaseAction(TypeOpenTicket, uuid),
		Ticketer:   ticketer,
		Subject:    subject,
		Body:       body,
		ResultName: resultName,
	}
}

// Execute runs this action
func (a *OpenTicketAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	ticketers := run.Session().Assets().Ticketers()
	ticketer := ticketers.Get(a.Ticketer.UUID)

	evaluatedSubject, err := run.EvaluateTemplate(a.Subject)
	if err != nil {
		logEvent(events.NewError(err))
	}
	evaluatedBody, err := run.EvaluateTemplate(a.Body)
	if err != nil {
		logEvent(events.NewError(err))
	}

	ticket := a.open(run, step, ticketer, evaluatedSubject, evaluatedBody, logEvent)
	if ticket != nil {
		a.saveResult(run, step, a.ResultName, string(ticket.UUID), CategorySuccess, "", "", nil, logEvent)
	} else {
		a.saveResult(run, step, a.ResultName, "", CategoryFailure, "", "", nil, logEvent)
	}

	return nil
}

func (a *OpenTicketAction) open(run flows.FlowRun, step flows.Step, ticketer *flows.Ticketer, subject, body string, logEvent flows.EventCallback) *flows.Ticket {
	if run.Session().BatchStart() {
		logEvent(events.NewErrorf("can't open tickets during batch starts"))
		return nil
	}

	if ticketer == nil {
		logEvent(events.NewDependencyError(a.Ticketer))
		return nil
	}

	svc, err := run.Session().Engine().Services().Ticket(run.Session(), ticketer)
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	httpLogger := &flows.HTTPLogger{}

	ticket, err := svc.Open(run.Session(), subject, body, httpLogger.Log)
	if err != nil {
		logEvent(events.NewError(err))
	}
	if len(httpLogger.Logs) > 0 {
		logEvent(events.NewTicketerCalled(ticketer.Reference(), httpLogger.Logs))
	}
	if ticket != nil {
		logEvent(events.NewTicketOpened(ticket))
	}

	return ticket
}
