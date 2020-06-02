package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeLanguage, readLanguageModifier)
}

// TypeLanguage is the type of our language modifier
const TypeLanguage string = "language"

// LanguageModifier modifies the language of a contact
type LanguageModifier struct {
	baseModifier

	Language envs.Language `json:"language"`
}

// NewLanguage creates a new language modifier
func NewLanguage(language envs.Language) *LanguageModifier {
	return &LanguageModifier{
		baseModifier: newBaseModifier(TypeLanguage),
		Language:     language,
	}
}

// Apply applies this modification to the given contact
func (m *LanguageModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	if contact.Language() != m.Language {
		contact.SetLanguage(m.Language)
		log(events.NewContactLanguageChanged(m.Language))
		m.reevaluateGroups(env, assets, contact, log)
	}
}

var _ flows.Modifier = (*LanguageModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readLanguageModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &LanguageModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
