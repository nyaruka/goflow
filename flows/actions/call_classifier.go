package actions

import (
	"context"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeCallClassifier, func() flows.Action { return &CallClassifier{} })
}

var classificationCategories = []string{CategorySuccess, CategorySkipped, CategoryFailure}

// TypeCallClassifier is the type for the call classifier action
const TypeCallClassifier string = "call_classifier"

// CallClassifier can be used to classify the intent and entities from a given input using an NLU classifier. It always
// saves a result indicating whether the classification was successful, skipped or failed, and what the extracted intents
// and entities were.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "call_classifier",
//	  "classifier": {
//	    "uuid": "1c06c884-39dd-4ce4-ad9f-9a01cbe6c000",
//	    "name": "Booking"
//	  },
//	  "input": "@input.text",
//	  "result_name": "Intent"
//	}
//
// @action call_classifier
type CallClassifier struct {
	baseAction
	onlineAction

	Classifier *assets.ClassifierReference `json:"classifier" validate:"required"`
	Input      string                      `json:"input" validate:"required" engine:"evaluated"`
	ResultName string                      `json:"result_name" validate:"required,result_name"`
}

// NewCallClassifier creates a new call classifier action
func NewCallClassifier(uuid flows.ActionUUID, classifier *assets.ClassifierReference, input string, resultName string) *CallClassifier {
	return &CallClassifier{
		baseAction: newBaseAction(TypeCallClassifier, uuid),
		Classifier: classifier,
		Input:      input,
		ResultName: resultName,
	}
}

// Execute runs this action
func (a *CallClassifier) Execute(ctx context.Context, run flows.Run, step flows.Step, log flows.EventLogger) error {
	classifiers := run.Session().Assets().Classifiers()
	classifier := classifiers.Get(a.Classifier.UUID)

	// substitute any variables in our input
	input, _ := run.EvaluateTemplate(a.Input, log)

	classification, skipped := a.classify(ctx, run, input, classifier, log)
	if classification != nil {
		a.saveSuccess(run, step, input, classification, log)
	} else if skipped {
		a.saveSkipped(run, step, input, log)
	} else {
		a.saveFailure(run, step, input, log)
	}

	return nil
}

func (a *CallClassifier) classify(ctx context.Context, run flows.Run, input string, classifier *flows.Classifier, log flows.EventLogger) (*flows.Classification, bool) {
	if input == "" {
		log(events.NewError("can't classify empty input, skipping classification", ""))
		return nil, true
	}
	if classifier == nil {
		log(events.NewDependencyError(a.Classifier))
		return nil, false
	}

	svc, err := run.Session().Engine().Services().Classification(classifier)
	if err != nil {
		log(events.NewError(err.Error(), ""))
		return nil, false
	}

	httpLogger := &flows.HTTPLogger{}

	classification, err := svc.Classify(ctx, run.Session().MergedEnvironment(), input, httpLogger.Log)

	if len(httpLogger.Logs) > 0 {
		log(events.NewClassifierCalled(classifier.Reference(), httpLogger.Logs))
	}

	if err != nil {
		log(events.NewError(err.Error(), ""))
		return nil, false
	}

	return classification, false
}

func (a *CallClassifier) saveSuccess(run flows.Run, step flows.Step, input string, classification *flows.Classification, log flows.EventLogger) {
	// result value is name of top ranked intent if there is one
	value := ""
	if len(classification.Intents) > 0 {
		value = classification.Intents[0].Name
	}
	extra, _ := jsonx.Marshal(classification)

	a.saveResult(run, step, a.ResultName, value, CategorySuccess, "", input, extra, log)
}

func (a *CallClassifier) saveSkipped(run flows.Run, step flows.Step, input string, log flows.EventLogger) {
	a.saveResult(run, step, a.ResultName, "0", CategorySkipped, "", input, nil, log)
}

func (a *CallClassifier) saveFailure(run flows.Run, step flows.Step, input string, log flows.EventLogger) {
	a.saveResult(run, step, a.ResultName, "0", CategoryFailure, "", input, nil, log)
}

func (a *CallClassifier) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	dependency(a.Classifier)

	if a.ResultName != "" {
		result(flows.NewResultInfo(a.ResultName, classificationCategories))
	}
}
