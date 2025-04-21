package actions

import (
	"context"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeSetRunResult, func() flows.Action { return &SetRunResultAction{} })
}

// TypeSetRunResult is the type for the set run result action
const TypeSetRunResult string = "set_run_result"

// SetRunResultAction can be used to save a result for a flow. The result will be available in the context
// for the run as @results.[name]. The optional category can be used as a way of categorizing results,
// this can be useful for reporting or analytics.
//
// Both the value and category fields may be templates. A [event:run_result_changed] event will be created with the
// final values.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_run_result",
//	  "name": "Gender",
//	  "value": "m",
//	  "category": "Male"
//	}
//
// @action set_run_result
type SetRunResultAction struct {
	baseAction
	universalAction

	Name     string `json:"name"                                  validate:"required,result_name"`
	Value    string `json:"value"              engine:"evaluated"`
	Category string `json:"category,omitempty" engine:"localized" validate:"omitempty,result_category"`
}

// NewSetRunResult creates a new set run result action
func NewSetRunResult(uuid flows.ActionUUID, name string, value string, category string) *SetRunResultAction {
	return &SetRunResultAction{
		baseAction: newBaseAction(TypeSetRunResult, uuid),
		Name:       name,
		Value:      value,
		Category:   category,
	}
}

// Execute runs this action
func (a *SetRunResultAction) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	value, ok := run.EvaluateTemplate(a.Value, logEvent)
	if !ok {
		return nil
	}

	categoryLocalized, _ := run.GetText(uuids.UUID(a.UUID()), "category", a.Category)
	if a.Category == categoryLocalized {
		categoryLocalized = ""
	}

	a.saveResult(run, step, a.Name, value, a.Category, categoryLocalized, "", nil, logEvent)
	return nil
}

func (a *SetRunResultAction) Inspect(result func(*flows.ResultInfo)) {
	if a.Category != "" {
		result(flows.NewResultInfo(a.Name, []string{a.Category}))
	} else {
		result(flows.NewResultInfo(a.Name, []string{}))
	}
}
