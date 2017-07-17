package events

import "github.com/nyaruka/goflow/flows"

// TypeSendEmail is our type for the email event
const TypeSendEmail string = "send_email"

// SendEmailEvent events are created for each recipient which should receive an email.
//
// ```
//   {
//     "type": "send_email",
//     "step_uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "created_on": "2006-01-02T15:04:05Z",
//     "email": "foo@bar.com",
//     "subject": "Your activation token",
//     "body": "Your activation token is AAFFKKEE"
//   }
// ```
//
// @event send_email
type SendEmailEvent struct {
	BaseEvent
	Email   string `json:"email"   validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body"`
}

// NewSendEmailEvent returns a new email event with the passed in subject, body and emails
func NewSendEmailEvent(email string, subject string, body string) *SendEmailEvent {
	return &SendEmailEvent{
		BaseEvent: NewBaseEvent(),
		Subject:   subject,
		Body:      body,
		Email:     email,
	}
}

// Type returns the type of this event
func (a *SendEmailEvent) Type() string { return TypeSendEmail }

// Apply applies this event to the given run
func (e *SendEmailEvent) Apply(run flows.FlowRun) {}
