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

func (m *LanguageModifier) Apply(assets flows.SessionAssets, contact *flows.Contact) flows.Event {
	if contact.Language() != m.Language {
		contact.SetLanguage(m.Language)
		return events.NewContactLanguageChangedEvent(m.Language)
	}
	return nil
}

var _ Modifier = (*LanguageModifier)(nil)
