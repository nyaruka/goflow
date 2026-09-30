package actions

import (
	"context"
	"strings"

	"github.com/nyaruka/goflow/core/events"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactEmail, func() flows.Action { return &SetContactEmail{} })
}

// TypeSetContactEmail is the type for the set contact email action
const TypeSetContactEmail string = "set_contact_email"

// SetContactEmail can be used to update the email address of the contact. The email is a template
// and white space is trimmed from the final value. An empty string clears the email.
// A [event:contact_email_changed] event will be created with the normalized value, or an
// [event:error] event if the value isn't a valid email address.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_contact_email",
//	  "email": "bob@example.com"
//	}
//
// @action set_contact_email
type SetContactEmail struct {
	baseAction
	universalAction

	Email string `json:"email" validate:"max=1000" engine:"evaluated"`
}

// NewSetContactEmail creates a new set email action
func NewSetContactEmail(uuid flows.ActionUUID, email string) *SetContactEmail {
	return &SetContactEmail{
		baseAction: newBaseAction(TypeSetContactEmail, uuid),
		Email:      email,
	}
}

// Execute runs this action
func (a *SetContactEmail) Execute(ctx context.Context, run flows.Run, step flows.Step, log events.EventLogger) error {
	email, ok := run.EvaluateTemplate(ctx, a.Email, log)
	email = strings.TrimSpace(email)

	if !ok {
		return nil
	}

	_, err := a.applyModifier(ctx, run, modifiers.NewEmail(email), log)
	return err
}
