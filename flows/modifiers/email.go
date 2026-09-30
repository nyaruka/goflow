package modifiers

import (
	"context"
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/core/events"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeEmail, readEmail)
}

// TypeEmail is the type of our email modifier
const TypeEmail string = "email"

// Email modifies the email address of a contact, creating a contact_email_changed event if it changed.
// Addresses are normalized and invalid ones are rejected with an error event. An empty value clears it.
type Email struct {
	baseModifier

	email string
}

// NewEmail creates a new email modifier
func NewEmail(email string) *Email {
	return &Email{
		baseModifier: newBaseModifier(TypeEmail),
		email:        email,
	}
}

// Apply applies this modification to the given contact
func (m *Email) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *core.Contact, log events.EventLogger) (bool, error) {
	email := ""
	if m.email != "" {
		var valid bool
		email, valid = core.NormalizeEmail(m.email)
		if !valid {
			log(events.NewError(fmt.Sprintf("'%s' is not a valid email address", m.email), events.ErrorCodeEmailInvalid, "email", m.email))
			return false, nil
		}
	}

	if contact.SetEmail(email) {
		log(events.NewContactEmailChanged(email))
		return true, nil
	}
	return false, nil
}

var _ flows.Modifier = (*Email)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type emailEnvelope struct {
	utils.TypedEnvelope

	Email string `json:"email"`
}

func readEmail(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &emailEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewEmail(e.Email), nil
}

func (m *Email) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&emailEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Email:         m.email,
	})
}
