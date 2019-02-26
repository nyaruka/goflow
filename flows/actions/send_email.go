package actions

import (
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
//     "addresses": ["@(contact.urns.mailto[0])"],
//     "subject": "Here is your activation token",
//     "body": "Your activation token is @contact.fields.activation_token"
//   }
//
// @action send_email
type SendEmailAction struct {
	BaseAction
	onlineAction

	Addresses []string `json:"addresses" validate:"required,min=1" engine:"evaluate"`
	Subject   string   `json:"subject" validate:"required" engine:"evaluate"`
	Body      string   `json:"body" validate:"required" engine:"evaluate"`
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
func (a *SendEmailAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	subject, err := run.EvaluateTemplate(a.Subject)
	if err != nil {
		logEvent(events.NewErrorEvent(err))
	}

	// make sure the subject is single line - replace '\t\n\r\f\v' to ' '
	subject = regexp.MustCompile(`\s+`).ReplaceAllString(subject, " ")
	subject = strings.TrimSpace(subject)

	if subject == "" {
		logEvent(events.NewErrorEventf("email subject evaluated to empty string, skipping"))
		return nil
	}

	body, err := run.EvaluateTemplate(a.Body)
	if err != nil {
		logEvent(events.NewErrorEvent(err))
	}
	if body == "" {
		logEvent(events.NewErrorEventf("email body evaluated to empty string, skipping"))
		return nil
	}

	evaluatedAddresses := make([]string, 0)

	for _, address := range a.Addresses {
		evaluatedAddress, err := run.EvaluateTemplate(address)
		if err != nil {
			logEvent(events.NewErrorEvent(err))
		}
		if evaluatedAddress == "" {
			logEvent(events.NewErrorEventf("email address evaluated to empty string, skipping"))
			continue
		}

		// strip mailto prefix if this is an email URN
		if strings.HasPrefix(evaluatedAddress, "mailto:") {
			evaluatedAddress = evaluatedAddress[7:]
		}

		evaluatedAddresses = append(evaluatedAddresses, evaluatedAddress)
	}

	if len(evaluatedAddresses) > 0 {
		logEvent(events.NewEmailCreatedEvent(evaluatedAddresses, subject, body))
	}

	return nil
}
