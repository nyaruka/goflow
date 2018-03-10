package events

import "github.com/nyaruka/goflow/flows"

// TypeEmailCreated is our type for the email event
const TypeEmailCreated string = "email_created"

// EmailCreatedEvent events are created for each recipient which should receive an email.
//
// ```
//   {
//     "type": "email_created",
//     "created_on": "2006-01-02T15:04:05Z",
//     "addresses": ["foo@bar.com"],
//     "subject": "Your activation token",
//     "body": "Your activation token is AAFFKKEE"
//   }
// ```
//
// @event email_created
type EmailCreatedEvent struct {
	BaseEvent
	Addresses []string `json:"addresses" validate:"required,min=1"`
	Subject   string   `json:"subject" validate:"required"`
	Body      string   `json:"body"`
}

// NewEmailCreatedEvent returns a new email event with the passed in subject, body and emails
func NewEmailCreatedEvent(addresses []string, subject string, body string) *EmailCreatedEvent {
	return &EmailCreatedEvent{
		BaseEvent: NewBaseEvent(),
		Addresses: addresses,
		Subject:   subject,
		Body:      body,
	}
}

// Type returns the type of this event
func (e *EmailCreatedEvent) Type() string { return TypeEmailCreated }

// Apply applies this event to the given run
func (e *EmailCreatedEvent) Apply(run flows.FlowRun) error {
	return nil
}
