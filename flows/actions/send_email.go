package actions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeSendEmail, func() flows.Action { return &SendEmailAction{} })
}

// TypeSendEmail is the type for the send email action
const TypeSendEmail string = "send_email"

// SendEmailAction can be used to send an email to one or more recipients. The subject, body and addresses
// can all contain expressions.
//
// An `email_created` event will be created for each email address.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "send_email",
//     "addresses": ["@contact.urns.mailto.0"],
//     "subject": "Here is your activation token",
//     "body": "Your activation token is @contact.fields.activation_token"
//   }
//
// @action send_email
type SendEmailAction struct {
	BaseAction
	Addresses []string `json:"addresses" validate:"required,min=1"`
	Subject   string   `json:"subject" validate:"required"`
	Body      string   `json:"body" validate:"required"`
}

// Type returns the type of this action
func (a *SendEmailAction) Type() string { return TypeSendEmail }

// Validate validates our action is valid and has all the assets it needs
func (a *SendEmailAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute creates the email events
func (a *SendEmailAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	subject, err := run.EvaluateTemplateAsString(a.Subject, false)
	if err != nil {
		log.Add(events.NewErrorEvent(err))
	}
	if subject == "" {
		log.Add(events.NewErrorEvent(fmt.Errorf("send_email subject evaluated to empty string, skipping")))
		return nil
	}

	// make sure the subject is single line - replace '\t\n\r\f\v' to ' '
	subject = regexp.MustCompile(`\s+`).ReplaceAllString(subject, " ")

	body, err := run.EvaluateTemplateAsString(a.Body, false)
	if err != nil {
		log.Add(events.NewErrorEvent(err))
	}
	if body == "" {
		log.Add(events.NewErrorEvent(fmt.Errorf("send_email body evaluated to empty string, skipping")))
		return nil
	}

	evaluatedAddresses := make([]string, 0)

	for _, address := range a.Addresses {
		evaluatedAddress, err := run.EvaluateTemplateAsString(address, false)
		if err != nil {
			log.Add(events.NewErrorEvent(err))
		}
		if evaluatedAddress == "" {
			log.Add(events.NewErrorEvent(fmt.Errorf("send_email address evaluated to empty string, skipping")))
			continue
		}

		// strip mailto prefix if this is an email URN
		if strings.HasPrefix(evaluatedAddress, "mailto:") {
			evaluatedAddress = evaluatedAddress[7:]
		}

		evaluatedAddresses = append(evaluatedAddresses, evaluatedAddress)
	}

	if len(evaluatedAddresses) > 0 {
		log.Add(events.NewEmailCreatedEvent(evaluatedAddresses, subject, body))
	}

	return nil
}
