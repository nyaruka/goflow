package actions

import (
	"context"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/events"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeCallClassifier, func() flows.Action { return &CallClassifier{} })
}

// Execute only ever saves a Failure result now, but all three categories are still advertised
// from Inspect so existing flows with router cases wired to Success or Skipped continue to validate.
var classificationCategories = []string{CategorySuccess, CategorySkipped, CategoryFailure}

// TypeCallClassifier is the type for the call classifier action
const TypeCallClassifier string = "call_classifier"

// CallClassifier is a deprecated action that previously called an NLU classifier. It is retained
// only so that existing flow definitions remain parseable, and now always saves a result with
// category Failure.
type CallClassifier struct {
	baseAction
	onlineAction

	Input      string `json:"input" validate:"required" engine:"evaluated"`
	ResultName string `json:"result_name" validate:"required,result_name"`
}

// NewCallClassifier creates a new call classifier action
func NewCallClassifier(uuid flows.ActionUUID, input string, resultName string) *CallClassifier {
	return &CallClassifier{
		baseAction: newBaseAction(TypeCallClassifier, uuid),
		Input:      input,
		ResultName: resultName,
	}
}

// Execute runs this action
func (a *CallClassifier) Execute(ctx context.Context, run flows.Run, step flows.Step, log events.EventLogger) error {
	input, _ := run.EvaluateTemplate(a.Input, log)

	log(events.NewWarning("NLU classifiers are no longer supported"))
	a.saveResult(run, step, a.ResultName, "0", CategoryFailure, "", input, nil, log)
	return nil
}

func (a *CallClassifier) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	result(flows.NewResultInfo(a.ResultName, classificationCategories))
}
