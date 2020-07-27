package actions

import (
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactLanguage, func() flows.Action { return &SetContactLanguageAction{} })
}

// TypeSetContactLanguage is the type for the set contact Language action
const TypeSetContactLanguage string = "set_contact_language"

// SetContactLanguageAction can be used to update the name of the contact. The language is a localizable
// template and white space is trimmed from the final value. An empty string clears the language.
// A [event:contact_language_changed] event will be created with the corresponding value.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "set_contact_language",
//     "language": "eng"
//   }
//
// @action set_contact_language
type SetContactLanguageAction struct {
	baseAction
	universalAction

	Language string `json:"language" engine:"evaluated"`
}

// NewSetContactLanguage creates a new set language action
func NewSetContactLanguage(uuid flows.ActionUUID, language string) *SetContactLanguageAction {
	return &SetContactLanguageAction{
		baseAction: newBaseAction(TypeSetContactLanguage, uuid),
		Language:   language,
	}
}

// Execute runs this action
func (a *SetContactLanguageAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	if run.Contact() == nil {
		logEvent(events.NewErrorf("can't execute action in session without a contact"))
		return nil
	}

	language, err := run.EvaluateTemplate(a.Language)
	language = strings.TrimSpace(language)

	// if we received an error, log it
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	// language must be empty or valid language code
	lang := envs.NilLanguage
	if language != "" {
		lang, err = envs.ParseLanguage(language)
		if err != nil {
			logEvent(events.NewError(err))
			return nil
		}
	}

	a.applyModifier(run, modifiers.NewLanguage(lang), logModifier, logEvent)
	return nil
}
