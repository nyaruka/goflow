package definition

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// itemTranslations map a key for a node to a key - say "text" to "[je suis francais!]"
type itemTranslations map[string][]string

// languageTranslations map a node uuid to item_translations - say "node1-asdf" to { "text": "je suis francais!" }
type languageTranslations map[flows.UUID]itemTranslations

func (t *languageTranslations) GetTextArray(uuid flows.UUID, key string) ([]string, bool) {
	item, found := (*t)[uuid]
	if found {
		translation, found := item[key]
		if found {
			return translation, true
		}
	}
	return nil, false
}

// flowTranslations are our top level container for all the translations for a language
type flowTranslations map[utils.Language]*languageTranslations

func (t *flowTranslations) GetLanguageTranslations(lang utils.Language) (flows.Translations, bool) {
	translations, found := (*t)[lang]
	return translations, found
}
