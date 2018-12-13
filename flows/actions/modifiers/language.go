package modifiers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeLanguage, func() Modifier { return &LanguageModifier{} })
}

// TypeLanguage is the type of our language modifier
const TypeLanguage string = "language"

// LanguageModifier modifies the language of a contact
type LanguageModifier struct {
	baseModifier

	Language utils.Language `json:"language"`
}

// NewLanguageModifier creates a new language modifier
func NewLanguageModifier(language utils.Language) *LanguageModifier {
	return &LanguageModifier{
		baseModifier: newBaseModifier(TypeLanguage),
		Language:     language,
	}
}

// Apply applies this modification to the given contact
func (m *LanguageModifier) Apply(env utils.Environment, assets flows.SessionAssets, contact *flows.Contact, log func(flows.Event)) {
	if contact.Language() != m.Language {
		contact.SetLanguage(m.Language)
		log(events.NewContactLanguageChangedEvent(m.Language))
		m.reevaluateDynamicGroups(env, assets, contact, log)
	}
}

var _ Modifier = (*LanguageModifier)(nil)
