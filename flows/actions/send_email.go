package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeSendEmail is our type for the email action
const TypeSendEmail string = "send_email"

// EmailAction can be used to send an email to one or more recipients. The subject, body and emails can all contain templates.
//
// A `send_email` event will be created for each email address.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "send_email",
//     "subject": "Here is your activation token",
//     "body": "Your activation token is @contact.fields.activation_token",
//     "emails": ["@contact.urns.email"]
//   }
// ```
//
// @action send_email
type EmailAction struct {
	BaseAction
	Emails  []string `json:"emails"   validate:"required,min=1"`
	Subject string   `json:"subject"  validate:"required"`
	Body    string   `json:"body"     validate:"required"`
}

// Type returns the type of this action
func (a *EmailAction) Type() string { return TypeSendEmail }

// Validate validates the fields on this action
func (a *EmailAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute creates the email events
func (a *EmailAction) Execute(run flows.FlowRun, step flows.Step, log flows.ActionLog) error {
	for _, email := range a.Emails {
		email, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), email)
		if err != nil {
			log.Add(events.NewErrorEvent(err))
		}
		if email == "" {
			log.Add(events.NewErrorEvent(fmt.Errorf("send_email email evaluated to empty string, skipping")))
			continue
		}

		subject, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.Subject)
		if err != nil {
			log.Add(events.NewErrorEvent(err))
		}
		if subject == "" {
			log.Add(events.NewErrorEvent(fmt.Errorf("send_email subject evaluated to empty string, skipping")))
			continue
		}

		body, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.Body)
		if err != nil {
			log.Add(events.NewErrorEvent(err))
		}
		log.Add(events.NewSendEmailEvent(email, subject, body))
	}
	return nil
}
