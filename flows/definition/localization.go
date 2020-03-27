package definition

import (
	"encoding/json"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"
)

// holds all property translations for a specific item, e.g.
// {
//   "text": "Do you like cheese?"
//	 "quick_replies": ["Yes", "No"]
// }
type itemTranslation map[string][]string

// holds all the item translations for a specific language, e.g.
// {
//   "f3368070-8db8-4549-872a-e69a9d060612": {
//	   "text": "Do you like cheese?"
//	   "quick_replies": ["Yes", "No"]
//   },
//   "7a1aec43-f3e1-42f0-b967-0ee75e725e3a": { ... }
// }
type languageTranslation map[uuids.UUID]itemTranslation

// returns the requested item translation
func (t languageTranslation) getTextArray(uuid uuids.UUID, property string) []string {
	item, found := t[uuid]
	if found {
		translation, found := item[property]
		if found {
			return translation
		}
	}
	return nil
}

// creates/updates the requested item translation
func (t languageTranslation) setTextArray(uuid uuids.UUID, property string, translated []string) {
	_, found := t[uuid]
	if !found {
		t[uuid] = make(itemTranslation)
	}

	t[uuid][property] = translated
}

func (t languageTranslation) Enumerate(callback func(uuids.UUID, string, []string)) {
	for uuid, it := range t {
		for property, texts := range it {
			callback(uuid, property, texts)
		}
	}
}

// our top level container for all the translations for all languages
type localization map[envs.Language]languageTranslation

// NewLocalization creates a new empty localization
func NewLocalization() flows.Localization {
	return make(localization)
}

// Languages gets the list of languages included in this localization
func (l localization) Languages() []envs.Language {
	languages := make([]envs.Language, 0, len(l))
	for lang := range l {
		languages = append(languages, lang)
	}
	return languages
}

// GetItemTranslation gets an item translation
func (l localization) GetItemTranslation(lang envs.Language, itemUUID uuids.UUID, property string) []string {
	translation, exists := l[lang]
	if exists {
		return translation.getTextArray(itemUUID, property)
	}
	return nil
}

// SetItemTranslation sets an item translation
func (l localization) SetItemTranslation(lang envs.Language, itemUUID uuids.UUID, property string, translated []string) {
	_, found := l[lang]
	if !found {
		l[lang] = make(languageTranslation)
	}
	l[lang].setTextArray(itemUUID, property, translated)
}

// ReadLocalization reads entire localization flow segment
func ReadLocalization(data json.RawMessage) (flows.Localization, error) {
	translations := &localization{}
	if err := jsonx.Unmarshal(data, translations); err != nil {
		return nil, err
	}
	return translations, nil
}
