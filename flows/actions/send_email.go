package actions

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSendEmail, func() flows.Action { return &SendEmail{} })
}

// TypeSendEmail is the type for the send email action
const TypeSendEmail string = "send_email"

// SendEmail can be used to send an email to one or more recipients. The subject, body and addresses
// can all contain expressions.
//
// An [event:email_sent] event will be created if the email could be sent.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "send_email",
//	  "addresses": ["@urns.mailto"],
//	  "subject": "Here is your activation token",
//	  "body": "Your activation token is @contact.fields.activation_token"
//	}
//
// @action send_email
type SendEmail struct {
	baseAction
	onlineAction

	Addresses []string `json:"addresses" validate:"required,min=1,dive,max=1000" engine:"evaluated"`
	Subject   string   `json:"subject"   validate:"required,max=1000"            engine:"localized,evaluated"`
	Body      string   `json:"body"      validate:"required,max=10000"           engine:"localized,evaluated"`
}

// NewSendEmail creates a new send email action
func NewSendEmail(uuid flows.ActionUUID, addresses []string, subject string, body string) *SendEmail {
	return &SendEmail{
		baseAction: newBaseAction(TypeSendEmail, uuid),
		Addresses:  addresses,
		Subject:    subject,
		Body:       body,
	}
}

// Execute creates the email events
func (a *SendEmail) Execute(ctx context.Context, run flows.Run, step flows.Step, logEvent flows.EventCallback) error {
	localizedSubject, _ := run.GetText(uuids.UUID(a.UUID()), "subject", a.Subject)
	evaluatedSubject, _ := run.EvaluateTemplate(localizedSubject, logEvent)

	// make sure the subject is single line - replace '\t\n\r\f\v' to ' '
	evaluatedSubject = regexp.MustCompile(`\s+`).ReplaceAllString(evaluatedSubject, " ")
	evaluatedSubject = strings.TrimSpace(evaluatedSubject)

	if evaluatedSubject == "" {
		logEvent(events.NewError("email subject evaluated to empty string, skipping"))
		return nil
	}

	localizedBody, _ := run.GetText(uuids.UUID(a.UUID()), "body", a.Body)
	evaluatedBody, _ := run.EvaluateTemplate(localizedBody, logEvent)
	if evaluatedBody == "" {
		logEvent(events.NewError("email body evaluated to empty string, skipping"))
		return nil
	}

	evaluatedAddresses := make([]string, 0)

	for _, address := range a.Addresses {
		evaluatedAddress, _ := run.EvaluateTemplate(address, logEvent)
		if evaluatedAddress == "" {
			logEvent(events.NewError("email address evaluated to empty string, skipping"))
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

	svc, err := run.Session().Engine().Services().Email(run.Session().Assets())
	if err != nil {
		logEvent(events.NewError(err.Error()))
		return nil
	}

	err = svc.Send(evaluatedAddresses, evaluatedSubject, evaluatedBody)
	if err != nil {
		logEvent(events.NewError(fmt.Sprintf("unable to send email: %s", err.Error())))
	} else {
		logEvent(events.NewEmailSent(evaluatedAddresses, evaluatedSubject, evaluatedBody))
	}

	return nil
}
