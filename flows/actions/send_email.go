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
// An [event:email_created] event will be created for each email address.
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
	onlineAction

	Addresses []string `json:"addresses" validate:"required,min=1"`
	Subject   string   `json:"subject" validate:"required"`
	Body      string   `json:"body" validate:"required"`
}

// NewSendEmailAction creates a new send email action
func NewSendEmailAction(uuid flows.ActionUUID, addresses []string, subject string, body string) *SendEmailAction {
	return &SendEmailAction{
		BaseAction: NewBaseAction(TypeSendEmail, uuid),
		Addresses:  addresses,
		Subject:    subject,
		Body:       body,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *SendEmailAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return nil
}

// Execute creates the email events
func (a *SendEmailAction) Execute(run flows.FlowRun, step flows.Step) error {
	subject, err := run.EvaluateTemplateAsString(a.Subject, false)
	if err != nil {
		a.logError(run, step, err)
	}

	// make sure the subject is single line - replace '\t\n\r\f\v' to ' '
	subject = regexp.MustCompile(`\s+`).ReplaceAllString(subject, " ")
	subject = strings.TrimSpace(subject)

	if subject == "" {
		a.logError(run, step, fmt.Errorf("email subject evaluated to empty string, skipping"))
		return nil
	}

	body, err := run.EvaluateTemplateAsString(a.Body, false)
	if err != nil {
		a.logError(run, step, err)
	}
	if body == "" {
		a.logError(run, step, fmt.Errorf("email body evaluated to empty string, skipping"))
		return nil
	}

	evaluatedAddresses := make([]string, 0)

	for _, address := range a.Addresses {
		evaluatedAddress, err := run.EvaluateTemplateAsString(address, false)
		if err != nil {
			a.logError(run, step, err)
		}
		if evaluatedAddress == "" {
			a.logError(run, step, fmt.Errorf("email address evaluated to empty string, skipping"))
			continue
		}

		// strip mailto prefix if this is an email URN
		if strings.HasPrefix(evaluatedAddress, "mailto:") {
			evaluatedAddress = evaluatedAddress[7:]
		}

		evaluatedAddresses = append(evaluatedAddresses, evaluatedAddress)
	}

	if len(evaluatedAddresses) > 0 {
		a.log(run, step, events.NewEmailCreatedEvent(evaluatedAddresses, subject, body))
	}

	return nil
}
