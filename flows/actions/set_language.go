package actions

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const SET_LANGUAGE string = "set_language"

type SetLanguageAction struct {
	BaseAction
	Language flows.Language `json:"language"    validate:"len=3"`
}

func (a *SetLanguageAction) Type() string { return SET_LANGUAGE }

func (a *SetLanguageAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *SetLanguageAction) Execute(run flows.FlowRun, step flows.Step) error {
	// update our current language
	run.SetLanguage(a.Language)

	// and update our contact as well
	contact := run.Contact()
	if contact != nil {
		contact.SetLanguage(a.Language)
	}

	run.AddEvent(step, &events.SetLanguageEvent{Language: a.Language})
	return nil
}
