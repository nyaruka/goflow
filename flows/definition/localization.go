package definition

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// itemTranslations map a key for a node to a key - say "text" to "je suis francais!"
type itemTranslations map[string]string

// languageTranslations map a node uuid to item_translations - say "node1-asdf" to { "text": "je suis francais!" }
type languageTranslations map[flows.UUID]itemTranslations

func (t *languageTranslations) GetText(uuid flows.UUID, key string, backdown string) string {
	item, found := (*t)[uuid]
	if found {
		translation, found := item[key]
		if found {
			return translation
		}
	}
	return backdown
}

// flowTranslations are our top level container for all the translations for a language
type flowTranslations map[utils.Language]*languageTranslations

func (t *flowTranslations) GetTranslations(lang utils.Language) flows.Translations {
	translations, found := (*t)[lang]
	if found {
		return translations
	}
	return nil
}
