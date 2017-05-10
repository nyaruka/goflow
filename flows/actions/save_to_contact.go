package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const SAVE_TO_CONTACT string = "save_to_contact"

type SaveToContactAction struct {
	BaseAction
	Field flows.FieldUUID `json:"field"    validate:"nonzero"`
	Name  string          `json:"name"     validate:"nonzero"`
	Value string          `json:"value"    validate:"nonzero"`
}

func (a *SaveToContactAction) Type() string { return SAVE_TO_CONTACT }

func (a *SaveToContactAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *SaveToContactAction) Execute(run flows.FlowRun, step flows.Step) error {
	// get our localized value if any
	template := run.GetText(flows.UUID(a.Uuid), "value", a.Value)
	value, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), template)

	// if we received an error, log it
	if err != nil {
		run.AddError(step, err)
	}

	// log our event
	event := events.SaveToContactEvent{Field: a.Field, Name: a.Name, Value: value}
	run.AddEvent(step, &event)

	// save to our field dictionary
	run.Contact().Fields().Save(a.Field, a.Name, value, *event.CreatedOn())

	return nil
}
