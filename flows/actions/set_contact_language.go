package actions

import (
	"context"
	"strings"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactLanguage, func() flows.Action { return &SetContactLanguage{} })
}

// TypeSetContactLanguage is the type for the set contact Language action
const TypeSetContactLanguage string = "set_contact_language"

// SetContactLanguage can be used to update the name of the contact. The language is a localizable
// template and white space is trimmed from the final value. An empty string clears the language.
// A [event:contact_language_changed] event will be created with the corresponding value.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_contact_language",
//	  "language": "eng"
//	}
//
// @action set_contact_language
type SetContactLanguage struct {
	baseAction
	universalAction

	Language string `json:"language" engine:"evaluated"`
}

// NewSetContactLanguage creates a new set language action
func NewSetContactLanguage(uuid flows.ActionUUID, language string) *SetContactLanguage {
	return &SetContactLanguage{
		baseAction: newBaseAction(TypeSetContactLanguage, uuid),
		Language:   language,
	}
}

// Execute runs this action
func (a *SetContactLanguage) Execute(ctx context.Context, run flows.Run, step flows.Step, log flows.EventLogger) error {
	language, ok := run.EvaluateTemplate(a.Language, log)
	language = strings.TrimSpace(language)

	if !ok {
		return nil
	}

	// language must be empty or valid language code
	lang := i18n.NilLanguage
	var err error
	if language != "" {
		lang, err = i18n.ParseLanguage(language)
		if err != nil {
			log(events.NewRawError(err))
			return nil
		}
	}

	_, err = a.applyModifier(ctx, run, modifiers.NewLanguage(lang), log)
	return err
}
