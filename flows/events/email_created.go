package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeEmailCreated, func() flows.Event { return &EmailCreatedEvent{} })
}

// TypeEmailCreated is our type for the email event
const TypeEmailCreated string = "email_created"

// EmailCreatedEvent events are created when an action wants to send an email.
//
//   {
//     "type": "email_created",
//     "created_on": "2006-01-02T15:04:05Z",
//     "addresses": ["foo@bar.com"],
//     "subject": "Your activation token",
//     "body": "Your activation token is AAFFKKEE"
//   }
//
// @event email_created
type EmailCreatedEvent struct {
	baseEvent

	Addresses []string `json:"addresses" validate:"required,min=1"`
	Subject   string   `json:"subject" validate:"required"`
	Body      string   `json:"body"`
}

// NewEmailCreated returns a new email event with the passed in subject, body and emails
func NewEmailCreated(addresses []string, subject string, body string) *EmailCreatedEvent {
	return &EmailCreatedEvent{
		baseEvent: newBaseEvent(TypeEmailCreated),
		Addresses: addresses,
		Subject:   subject,
		Body:      body,
	}
}
