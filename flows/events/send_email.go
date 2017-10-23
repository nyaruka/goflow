package events

import "github.com/nyaruka/goflow/flows"

// TypeSendEmail is our type for the email event
const TypeSendEmail string = "send_email"

// SendEmailEvent events are created for each recipient which should receive an email.
//
// ```
//   {
//     "type": "send_email",
//     "created_on": "2006-01-02T15:04:05Z",
//     "addresses": ["foo@bar.com"],
//     "subject": "Your activation token",
//     "body": "Your activation token is AAFFKKEE"
//   }
// ```
//
// @event send_email
type SendEmailEvent struct {
	BaseEvent
	Addresses []string `json:"addresses" validate:"required,min=1"`
	Subject   string   `json:"subject" validate:"required"`
	Body      string   `json:"body"`
}

// NewSendEmailEvent returns a new email event with the passed in subject, body and emails
func NewSendEmailEvent(addresses []string, subject string, body string) *SendEmailEvent {
	return &SendEmailEvent{
		BaseEvent: NewBaseEvent(),
		Addresses: addresses,
		Subject:   subject,
		Body:      body,
	}
}

// Type returns the type of this event
func (a *SendEmailEvent) Type() string { return TypeSendEmail }

// Apply applies this event to the given run
func (e *SendEmailEvent) Apply(run flows.FlowRun) error {
	return nil
}
