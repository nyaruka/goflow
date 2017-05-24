package events

// TypeEmail is our type for the email event
const TypeEmail string = "email"

// EmailEvent events will be created for each email in an `email` action.
//
// ```
//   {
//     "step": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "created_on": "2006-01-02T15:04:05Z",
//     "type": "email",
//     "email": "foo@bar.com",
//     "subject": "Your activation token",
//     "body": "Your activation token is AAFFKKEE"
//   }
// ```
//
// @event email
type EmailEvent struct {
	BaseEvent
	Email   string `json:"email"   validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body"    validate:"required"`
}

// NewEmailEvent returns a new email event witht he passed in subject, body and emails
func NewEmailEvent(email string, subject string, body string) *EmailEvent {
	return &EmailEvent{
		Subject: subject,
		Body:    body,
		Email:   email,
	}
}

// Type returns the type of this event
func (a *EmailEvent) Type() string { return TypeEmail }
