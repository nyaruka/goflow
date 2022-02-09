package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeEmailCreated, func() flows.Event { return &EmailCreatedEvent{} })
}

// TypeEmailCreated is our type for the email event
const TypeEmailCreated string = "email_created"

// EmailCreatedEvent is no longer used but old sessions might include these
type EmailCreatedEvent struct {
	BaseEvent

	Addresses []string `json:"addresses" validate:"required,min=1"`
	Subject   string   `json:"subject" validate:"required"`
	Body      string   `json:"body"`
}
