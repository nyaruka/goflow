package actions

import (
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions/modifiers"
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

// NewSetContactLanguageAction creates a new set language action
func NewSetContactLanguageAction(uuid flows.ActionUUID, language string) *SetContactLanguageAction {
	return &SetContactLanguageAction{
		BaseAction: NewBaseAction(TypeSetContactLanguage, uuid),
		Language:   language,
	}
}

// Execute runs this action
func (a *SetContactLanguageAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	if run.Contact() == nil {
		logEvent(events.NewErrorEventf("can't execute action in session without a contact"))
		return nil
	}

	language, err := run.EvaluateTemplate(a.Language)
	language = strings.TrimSpace(language)

	// if we received an error, log it
	if err != nil {
		logEvent(events.NewErrorEvent(err))
		return nil
	}

	// language must be empty or valid language code
	lang := utils.NilLanguage
	if language != "" {
		lang, err = utils.ParseLanguage(language)
		if err != nil {
			logEvent(events.NewErrorEvent(err))
			return nil
		}
	}

	a.applyModifier(run, modifiers.NewLanguageModifier(lang), logModifier, logEvent)
	return nil
}

// Inspect inspects this object and any children
func (a *SetContactLanguageAction) Inspect(inspect func(flows.Inspectable)) {
	inspect(a)
}

// EnumerateTemplates enumerates all expressions on this object and its children
func (a *SetContactLanguageAction) EnumerateTemplates(localization flows.Localization, callback func(string)) {
	callback(a.Language)
}

// RewriteTemplates rewrites all templates on this object and its children
func (a *SetContactLanguageAction) RewriteTemplates(localization flows.Localization, rewrite func(string) string) {
	a.Language = rewrite(a.Language)
}
