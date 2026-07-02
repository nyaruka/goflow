package flows

import (
	"slices"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/envs"
)

// TranslationsForContact is a utility to help callers get the message content for a contact
func TranslationsForContact(e envs.Environment, b core.BroadcastTranslations, c *Contact, baseLanguage i18n.Language) (*core.MsgContent, i18n.Locale) {
	// get the set of languages to merge translations from
	languages := make([]i18n.Language, 0, 3)

	// highest priority is the contact language if it is valid
	if c.Language() != i18n.NilLanguage && slices.Contains(e.AllowedLanguages(), c.Language()) {
		languages = append(languages, c.Language())
	}

	// then the default workspace language, then the base language
	languages = append(languages, e.DefaultLanguage(), baseLanguage)

	content := &core.MsgContent{}
	language := i18n.NilLanguage
	country := e.DefaultCountry()
	if c.Country() != i18n.NilCountry {
		country = c.Country()
	}

	for _, lang := range languages {
		trans := b[lang]
		if trans != nil {
			if content.Text == "" && trans.Text != "" {
				content.Text = trans.Text
				language = lang
			}
			if len(content.Attachments) == 0 && len(trans.Attachments) > 0 {
				content.Attachments = trans.Attachments
			}
			if len(content.QuickReplies) == 0 && len(trans.QuickReplies) > 0 {
				content.QuickReplies = trans.QuickReplies
			}
		}
	}

	return content, i18n.NewLocale(language, country)
}
