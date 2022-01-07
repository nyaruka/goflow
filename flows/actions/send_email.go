package actions

import (
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/pkg/errors"
)

func init() {
	registerType(TypeSendEmail, func() flows.Action { return &SendEmailAction{} })
}

// TypeSendEmail is the type for the send email action
const TypeSendEmail string = "send_email"

// SendEmailAction can be used to send an email to one or more recipients. The subject, body and addresses
// can all contain expressions.
//
// An [event:email_sent] event will be created if the email could be sent.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "send_email",
//     "addresses": ["@urns.mailto"],
//     "subject": "Here is your activation token",
//     "body": "Your activation token is @contact.fields.activation_token"
//   }
//
// @action send_email
type SendEmailAction struct {
	baseAction
	onlineAction

	Addresses []string `json:"addresses" validate:"required,min=1" engine:"evaluated"`
	Subject   string   `json:"subject" validate:"required" engine:"localized,evaluated"`
	Body      string   `json:"body" validate:"required" engine:"localized,evaluated"`
}

// NewSendEmail creates a new send email action
func NewSendEmail(uuid flows.ActionUUID, addresses []string, subject string, body string) *SendEmailAction {
	return &SendEmailAction{
		baseAction: newBaseAction(TypeSendEmail, uuid),
		Addresses:  addresses,
		Subject:    subject,
		Body:       body,
	}
}

// Execute creates the email events
func (a *SendEmailAction) Execute(run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	localizedSubject := run.GetText(uuids.UUID(a.UUID()), "subject", a.Subject)
	evaluatedSubject, err := run.EvaluateTemplate(localizedSubject)
	if err != nil {
		logEvent(events.NewError(err))
	}

	// make sure the subject is single line - replace '\t\n\r\f\v' to ' '
	evaluatedSubject = regexp.MustCompile(`\s+`).ReplaceAllString(evaluatedSubject, " ")
	evaluatedSubject = strings.TrimSpace(evaluatedSubject)

	if evaluatedSubject == "" {
		logEvent(events.NewErrorf("email subject evaluated to empty string, skipping"))
		return nil
	}

	localizedBody := run.GetText(uuids.UUID(a.UUID()), "body", a.Body)
	evaluatedBody, err := run.EvaluateTemplate(localizedBody)
	if err != nil {
		logEvent(events.NewError(err))
	}
	if evaluatedBody == "" {
		logEvent(events.NewErrorf("email body evaluated to empty string, skipping"))
		return nil
	}

	evaluatedAddresses := make([]string, 0)

	for _, address := range a.Addresses {
		evaluatedAddress, err := run.EvaluateTemplate(address)
		if err != nil {
			logEvent(events.NewError(err))
		}
		if evaluatedAddress == "" {
			logEvent(events.NewErrorf("email address evaluated to empty string, skipping"))
			continue
		}

		// strip mailto prefix if this is an email URN
		evaluatedAddress = strings.TrimPrefix(evaluatedAddress, "mailto:")

		evaluatedAddresses = append(evaluatedAddresses, evaluatedAddress)
	}

	// nothing to do if there are no addresses
	if len(evaluatedAddresses) == 0 {
		return nil
	}

	svc, err := run.Session().Engine().Services().Email(run.Session())
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	err = svc.Send(run.Session(), evaluatedAddresses, evaluatedSubject, evaluatedBody)
	if err != nil {
		logEvent(events.NewError(errors.Wrap(err, "unable to send email")))
	} else {
		logEvent(events.NewEmailSent(evaluatedAddresses, evaluatedSubject, evaluatedBody))
	}

	return nil
}
