package definition

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// itemTranslations map a key for a node to a key - say "text" to "[je suis francais!]"
type itemTranslations map[string][]string

// languageTranslations map a node uuid to item_translations - say "node1-asdf" to { "text": "je suis francais!" }
type languageTranslations map[utils.UUID]itemTranslations

func (t *languageTranslations) GetTextArray(uuid utils.UUID, key string) []string {
	item, found := (*t)[uuid]
	if found {
		translation, found := item[key]
		if found {
			return translation
		}
	}
	return nil
}

// our top level container for all the translations for all languages
type localization map[utils.Language]*languageTranslations

func (t localization) Languages() utils.LanguageList {
	languages := make(utils.LanguageList, 0, len(t))
	for lang := range t {
		languages = append(languages, lang)
	}
	return languages
}

func (t localization) GetTranslations(lang utils.Language) flows.Translations {
	translations, found := t[lang]
	if found {
		return translations
	}
	return nil
}

// ReadLocalization reads entire localization flow segment
func ReadLocalization(data json.RawMessage) (flows.Localization, error) {
	translations := &localization{}
	if err := json.Unmarshal(data, translations); err != nil {
		return nil, err
	}
	return translations, nil
}
