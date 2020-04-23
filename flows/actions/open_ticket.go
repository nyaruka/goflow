package actions

import (
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
//     "subject": "Needs help",
//     "result_name": "Help Ticket"
//   }
//
// @action open_ticket
type OpenTicketAction struct {
	baseAction
	onlineAction

	Subject    string `json:"subject" validate:"required" engine:"evaluated"`
	ResultName string `json:"result_name" validate:"required"`
}

// NewOpenTicket creates a new open ticket action
func NewOpenTicket(uuid flows.ActionUUID, subject, resultName string) *OpenTicketAction {
	return &OpenTicketAction{
		baseAction: newBaseAction(TypeOpenTicket, uuid),
		Subject:    subject,
		ResultName: resultName,
	}
}

// Execute runs this action
func (a *OpenTicketAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	evaluatedSubject, err := run.EvaluateTemplate(a.Subject)
	if err != nil {
		logEvent(events.NewError(err))
	}

	ticket := a.open(run, step, evaluatedSubject, logEvent)
	if ticket != nil {
		a.saveResult(run, step, a.ResultName, ticket.ID, CategorySuccess, "", "", nil, logEvent)
	} else {
		a.saveResult(run, step, a.ResultName, "", CategoryFailure, "", "", nil, logEvent)
	}

	return nil
}

func (a *OpenTicketAction) open(run flows.FlowRun, step flows.Step, subject string, logEvent flows.EventCallback) *flows.Ticket {
	svc, err := run.Session().Engine().Services().Ticket(run.Session())
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	httpLogger := &flows.HTTPLogger{}

	ticket, err := svc.Open(run.Session(), subject, httpLogger.Log)
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	if len(httpLogger.Logs) > 0 {
		logEvent(events.NewTicketOpened(ticket, httpLogger.Logs))
	}

	return ticket
}
