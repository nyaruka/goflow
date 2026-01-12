package modifiers

import (
	"context"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeLanguage, readLanguage)
}

// TypeLanguage is the type of our language modifier
const TypeLanguage string = "language"

// Language modifies the language of a contact
type Language struct {
	baseModifier

	language i18n.Language
}

// NewLanguage creates a new language modifier
func NewLanguage(language i18n.Language) *Language {
	return &Language{
		baseModifier: newBaseModifier(TypeLanguage),
		language:     language,
	}
}

// Apply applies this modification to the given contact
func (m *Language) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventLogger) (bool, error) {
	if contact.Language() != m.language {
		contact.SetLanguage(m.language)
		log(events.NewContactLanguageChanged(m.language))
		return true, nil
	}
	return false, nil
}

var _ flows.Modifier = (*Language)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type languageEnvelope struct {
	utils.TypedEnvelope

	Language i18n.Language `json:"language"`
}

func readLanguage(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &languageEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewLanguage(e.Language), nil
}

func (m *Language) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&languageEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Language:      m.language,
	})
}
