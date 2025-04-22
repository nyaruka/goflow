package definition

import (
	"fmt"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
)

// holds all property translations for a specific item, e.g.
//
//	{
//	  "text": ["Do you like cheese?"],
//	  "quick_replies": ["Yes", "No"],
//	  "_ui": {...}
//	}
type itemTranslation map[string]any

func (l itemTranslation) validate() error {
	for property := range l {
		if len(property) == 0 || len(property) > 64 {
			return fmt.Errorf("invalid property name '%s'", stringsx.TruncateEllipsis(property, 32))
		}
	}
	return nil
}

func (t itemTranslation) get(property string) []string {
	value, found := t[property]
	if !found {
		return nil
	}

	// flow editor is allowed to stuff _ui in here and it's not a string array
	asSlice, ok := value.([]any)
	if !ok {
		return nil
	}

	trans := make([]string, len(asSlice))
	for i, v := range asSlice {
		asString, ok := v.(string)
		if !ok {
			return nil
		}
		trans[i] = asString
	}

	// TODO editor sometimes saves empty rule translations as [""] which we should fix in a flow migration
	// but for now need to ignore
	if len(trans) == 0 || (len(trans) == 1 && trans[0] == "") {
		return nil
	}

	return trans
}

// holds all the item translations for a specific language, e.g.
//
//	{
//	  "f3368070-8db8-4549-872a-e69a9d060612": {
//	    "text": ["Do you like cheese?"],
//	    "quick_replies": ["Yes", "No"]
//	  },
//	  "7a1aec43-f3e1-42f0-b967-0ee75e725e3a": { ... }
//	}
type languageTranslation map[uuids.UUID]itemTranslation

func (l languageTranslation) validate() error {
	for uuid, item := range l {
		if !uuids.Is(string(uuid)) {
			return fmt.Errorf("invalid item uuid '%s'", uuid)
		}

		if err := item.validate(); err != nil {
			return fmt.Errorf("invalid item translation for '%s': %w", uuid, err)
		}
	}
	return nil
}

// returns the requested item translation
func (t languageTranslation) getTextArray(uuid uuids.UUID, property string) []string {
	item, found := t[uuid]
	if found {
		return item.get(property)
	}
	return nil
}

// creates/updates the requested item translation
func (t languageTranslation) setTextArray(uuid uuids.UUID, property string, translated []string) {
	_, found := t[uuid]
	if !found {
		t[uuid] = make(itemTranslation)
	}

	trans := make([]any, len(translated))
	for i, v := range translated {
		trans[i] = v
	}

	t[uuid][property] = trans
}

// our top level container for all the translations for all languages
type localization map[i18n.Language]languageTranslation

// NewLocalization creates a new empty localization
func NewLocalization() flows.Localization {
	return make(localization)
}

func (l localization) Validate() error {
	for lang, translation := range l {
		if len(lang) != 3 {
			return fmt.Errorf("invalid language code '%s'", lang)
		}
		if err := translation.validate(); err != nil {
			return fmt.Errorf("invalid translation for '%s': %w", lang, err)
		}
	}
	return nil
}

// Languages gets the list of languages included in this localization
func (l localization) Languages() []i18n.Language {
	languages := make([]i18n.Language, 0, len(l))
	for lang := range l {
		languages = append(languages, lang)
	}
	return languages
}

// GetItemTranslation gets an item translation
func (l localization) GetItemTranslation(lang i18n.Language, itemUUID uuids.UUID, property string) []string {
	translation, exists := l[lang]
	if exists {
		return translation.getTextArray(itemUUID, property)
	}
	return nil
}

// SetItemTranslation sets an item translation
func (l localization) SetItemTranslation(lang i18n.Language, itemUUID uuids.UUID, property string, translated []string) {
	_, found := l[lang]
	if !found {
		l[lang] = make(languageTranslation)
	}
	l[lang].setTextArray(itemUUID, property, translated)
}

// ReadLocalization reads entire localization flow segment
func ReadLocalization(data []byte) (flows.Localization, error) {
	translations := &localization{}
	if err := jsonx.Unmarshal(data, translations); err != nil {
		return nil, err
	}
	return translations, nil
}
