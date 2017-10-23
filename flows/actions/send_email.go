package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeSendEmail is our type for the email action
const TypeSendEmail string = "send_email"

// EmailAction can be used to send an email to one or more recipients. The subject, body and addresses
// can all contain expressions.
//
// A `send_email` event will be created for each email address.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "send_email",
//     "addresses": ["@contact.urns.email"],
//     "subject": "Here is your activation token",
//     "body": "Your activation token is @contact.fields.activation_token"
//   }
// ```
//
// @action send_email
type EmailAction struct {
	BaseAction
	Addresses []string `json:"addresses" validate:"required,min=1"`
	Subject   string   `json:"subject" validate:"required"`
	Body      string   `json:"body" validate:"required"`
}

// Type returns the type of this action
func (a *EmailAction) Type() string { return TypeSendEmail }

// Validate validates our action is valid and has all the assets it needs
func (a *EmailAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute creates the email events
func (a *EmailAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	subject, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.Subject)
	if err != nil {
		log.Add(events.NewErrorEvent(err))
	}
	if subject == "" {
		log.Add(events.NewErrorEvent(fmt.Errorf("send_email subject evaluated to empty string, skipping")))
		return nil
	}

	body, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.Body)
	if err != nil {
		log.Add(events.NewErrorEvent(err))
	}
	if body == "" {
		log.Add(events.NewErrorEvent(fmt.Errorf("send_email body evaluated to empty string, skipping")))
		return nil
	}

	evaluatedAddresses := make([]string, 0)

	for _, address := range a.Addresses {
		evaluatedAddress, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), address)
		if err != nil {
			log.Add(events.NewErrorEvent(err))
		}
		if evaluatedAddress == "" {
			log.Add(events.NewErrorEvent(fmt.Errorf("send_email address evaluated to empty string, skipping")))
			continue
		}
		evaluatedAddresses = append(evaluatedAddresses, evaluatedAddress)
	}

	if len(evaluatedAddresses) > 0 {
		log.Add(events.NewSendEmailEvent(evaluatedAddresses, subject, body))
	}

	return nil
}
