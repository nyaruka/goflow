package events

func init() {
	registerType(TypeContactEmailChanged, func() Event { return &ContactEmailChanged{} })
}

// TypeContactEmailChanged is the type of our contact email changed event
const TypeContactEmailChanged string = "contact_email_changed"

// ContactEmailChanged events are created when the email address of the contact has been changed.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "contact_email_changed",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "email": "bob@example.com"
//	}
//
// @event contact_email_changed
type ContactEmailChanged struct {
	BaseEvent

	Email string `json:"email"`
}

// NewContactEmailChanged returns a new contact email changed event
func NewContactEmailChanged(email string) *ContactEmailChanged {
	return &ContactEmailChanged{
		BaseEvent: NewBaseEvent(TypeContactEmailChanged),
		Email:     email,
	}
}
