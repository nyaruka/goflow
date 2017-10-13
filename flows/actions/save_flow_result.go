package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeSaveFlowResult is our type for the save result action
const TypeSaveFlowResult string = "save_flow_result"

// SaveFlowResultAction can be used to save a result for a flow. The result will be available in the context
// for the run as @run.results.[name]. The optional category can be used as a way of categorizing results,
// this can be useful for reporting or analytics.
//
// Both the value and category fields may be templates. A `save_flow_result` event will be created with the
// final values.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "save_flow_result",
//     "result_name": "gender",
//     "value": "m",
//     "category": "Male"
//   }
// ```
//
// @action save_flow_result
type SaveFlowResultAction struct {
	BaseAction
	ResultName string `json:"result_name"        validate:"required"`
	Value      string `json:"value"`
	Category   string `json:"category"`
}

// Type returns the type of this action
func (a *SaveFlowResultAction) Type() string { return TypeSaveFlowResult }

// Validate validates the fields on this action
func (a *SaveFlowResultAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *SaveFlowResultAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	// get our localized value if any
	template := run.GetText(flows.UUID(a.UUID()), "value", a.Value)
	value, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), template)

	// log any error received
	if err != nil {
		log.Add(events.NewErrorEvent(err))
	}

	template = run.GetText(flows.UUID(a.UUID()), "category", a.Category)
	categoryLocalized, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), template)
	if err != nil {
		log.Add(events.NewErrorEvent(err))
	}

	if a.Category == categoryLocalized {
		categoryLocalized = ""
	}

	log.Add(events.NewSaveFlowResult(step.NodeUUID(), a.ResultName, value, a.Category, categoryLocalized))
	return nil
}
