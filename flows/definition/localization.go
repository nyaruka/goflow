package definition

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// the translations for a specific item, e.g.
// {
//   "text": "Do you like cheese?"
//	 "quick_replies": ["Yes", "No"]
// }
type itemTranslations map[string][]string

// the translations for a specific language, e.g.
// {
//   "f3368070-8db8-4549-872a-e69a9d060612": {
//	   "text": "Do you like cheese?"
//	   "quick_replies": ["Yes", "No"]
//   },
//   "7a1aec43-f3e1-42f0-b967-0ee75e725e3a": { ... }
// }
type languageTranslations map[utils.UUID]itemTranslations

// GetTextArray returns the requested item translation
func (t languageTranslations) GetTextArray(uuid utils.UUID, property string) []string {
	item, found := t[uuid]
	if found {
		translation, found := item[property]
		if found {
			return translation
		}
	}
	return nil
}

// SetTextArray updates the requested item translation
func (t languageTranslations) SetTextArray(uuid utils.UUID, property string, translated []string) {
	_, found := t[uuid]
	if !found {
		t[uuid] = make(itemTranslations)
	}

	t[uuid][property] = translated
}

// our top level container for all the translations for all languages
type localization map[utils.Language]languageTranslations

func NewLocalization() flows.Localization {
	return make(localization)
}

// Languages gets the list of languages included in this localization
func (l localization) Languages() []utils.Language {
	languages := make([]utils.Language, 0, len(l))
	for lang := range l {
		languages = append(languages, lang)
	}
	return languages
}

// AddItemTranslation adds a new item translation
func (l localization) AddItemTranslation(lang utils.Language, itemUUID utils.UUID, property string, translated []string) {
	_, found := l[lang]
	if !found {
		l[lang] = make(languageTranslations)
	}
	l[lang].SetTextArray(itemUUID, property, translated)
}

// GetTranslations returns the translations for the given language
func (l localization) GetTranslations(lang utils.Language) flows.Translations {
	return l[lang]
}

// ReadLocalization reads entire localization flow segment
func ReadLocalization(data json.RawMessage) (flows.Localization, error) {
	translations := &localization{}
	if err := json.Unmarshal(data, translations); err != nil {
		return nil, err
	}
	return translations, nil
}
