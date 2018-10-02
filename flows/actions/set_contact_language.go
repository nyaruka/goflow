package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeSetContactLanguage, func() flows.Action { return &SetContactLanguageAction{} })
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
	BaseAction
	universalAction

	Language string `json:"language"`
}

// Type returns the type of this action
func (a *SetContactLanguageAction) Type() string { return TypeSetContactLanguage }

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactLanguageAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *SetContactLanguageAction) Execute(run flows.FlowRun, step flows.Step) error {
	if run.Contact() == nil {
		a.logError(run, step, fmt.Errorf("can't execute action in session without a contact"))
		return nil
	}

	language, err := a.evaluateLocalizableTemplate(run, "language", a.Language)
	language = strings.TrimSpace(language)

	// if we received an error, log it
	if err != nil {
		a.logError(run, step, err)
		return nil
	}

	// language must be empty or valid language code
	lang := utils.NilLanguage
	if language != "" {
		lang, err = utils.ParseLanguage(language)
		if err != nil {
			a.logError(run, step, err)
			return nil
		}
	}

	if run.Contact().Language() != lang {
		run.Contact().SetLanguage(lang)
		a.log(run, step, events.NewContactLanguageChangedEvent(language))

		a.reevaluateDynamicGroups(run, step)
	}

	return nil
}
