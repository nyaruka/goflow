package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const SAVE_RESULT string = "save_result"

type SaveResultAction struct {
	BaseAction
	Name     string `json:"name"        validate:"nonzero"`
	Value    string `json:"value"       validate:"nonzero"`
	Category string `json:"category"    validate:"nonzero"`
}

func (a *SaveResultAction) Type() string { return SAVE_RESULT }

func (a *SaveResultAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *SaveResultAction) Execute(run flows.FlowRun, step flows.Step) error {
	// get our localized value if any
	template := run.GetText(flows.UUID(a.Uuid), "value", a.Value)
	value, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), template)

	// log any error received
	if err != nil {
		run.AddError(step, err)
	}

	// log our event
	event := events.NewResultEvent(step.Node(), a.Name, value, a.Category)
	run.AddEvent(step, event)

	// and save our result
	run.Results().Save(step.Node(), a.Name, value, a.Category, *event.CreatedOn())

	return nil
}
