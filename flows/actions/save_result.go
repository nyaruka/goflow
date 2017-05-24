package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeSaveResult is our type for the save result action
const TypeSaveResult string = "save_result"

// SaveResultAction can be used to save a result for a flow. The result will be available in the context
// for the run as @run.results.[name]. The optional category can be used as a way of categorizing results,
// this can be useful for reporting or analytics.
//
// Both the value and category fields may be templates. A `save_result` event will be created with the
// final values.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "save_result",
//     "name": "gender",
//     "value": "m",
//     "category": "Male"
//   }
//
// @action save_result
type SaveResultAction struct {
	BaseAction
	Name     string `json:"name"        validate:"required"`
	Value    string `json:"value"       validate:"required"`
	Category string `json:"category"`
}

// Type returns the type of this action
func (a *SaveResultAction) Type() string { return TypeSaveResult }

// Validate validates the fields on this action
func (a *SaveResultAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute runs this action
func (a *SaveResultAction) Execute(run flows.FlowRun, step flows.Step) error {
	// get our localized value if any
	template := run.GetText(flows.UUID(a.UUID), "value", a.Value)
	value, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), template)

	// log any error received
	if err != nil {
		run.AddError(step, err)
	}

	// log our event
	event := events.NewSaveResult(step.Node(), a.Name, value, a.Category)
	run.AddEvent(step, event)

	// and save our result
	run.Results().Save(step.Node(), a.Name, value, a.Category, *event.CreatedOn())

	return nil
}
