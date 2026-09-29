package actions

import (
	"context"
	"fmt"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/core/events"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeCallClassifier, func() flows.Action { return &CallClassifier{} })
}

// TypeCallClassifier is the type for the call classifier action
const TypeCallClassifier string = "call_classifier"

// ClassifierNoneOutput is the output used when the model isn't confident enough in any option
const ClassifierNoneOutput = "<NONE>"

// CallClassifier can be used to classify input as one of a set of options using a model. The input field may be a
// template and will be evaluated at runtime. Each option can have a description to help the model decide whether it
// fits.
//
// A [event:classifier_called] event will be created if the model could be called. The action sets the local specified
// by `output_local` to the name of the chosen option, or to `<ERROR>` if the call failed. The model always chooses one
// of the options, so if its confidence in that option is below `min_confidence` (between 0 and 1), the local is set
// to `<NONE>` instead.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "call_classifier",
//	  "model": {
//	    "uuid": "14115c03-b4c5-49e2-b9ac-390c43e9d7ce",
//	    "name": "GPT-4"
//	  },
//	  "input": "@input.text",
//	  "options": [
//	    {"name": "Flights", "description": "Booking or changing flights"},
//	    {"name": "Hotels", "description": "Booking or changing hotel rooms"}
//	  ],
//	  "min_confidence": 0.5,
//	  "output_local": "_classification"
//	}
//
// @action call_classifier
type CallClassifier struct {
	baseAction
	onlineAction

	Model         *assets.ModelReference   `json:"model"                    validate:"required"`
	Input         string                   `json:"input"                    validate:"max=10000"                          engine:"evaluated"`
	Options       []*core.ClassifierOption `json:"options"                  validate:"required,min=1,max=10,unique=Name,dive"`
	MinConfidence float64                  `json:"min_confidence,omitempty" validate:"min=0,max=1"`
	OutputLocal   string                   `json:"output_local"             validate:"required,local_ref"`
}

// NewCallClassifier creates a new call classifier action
func NewCallClassifier(uuid flows.ActionUUID, model *assets.ModelReference, input string, options []*core.ClassifierOption, minConfidence float64, outputLocal string) *CallClassifier {
	return &CallClassifier{
		baseAction:    newBaseAction(TypeCallClassifier, uuid),
		Model:         model,
		Input:         input,
		Options:       options,
		MinConfidence: minConfidence,
		OutputLocal:   outputLocal,
	}
}

// Validate validates our action is valid
func (a *CallClassifier) Validate() error {
	for _, o := range a.Options {
		if o.Name == ModelErrorOutput || o.Name == ClassifierNoneOutput {
			return fmt.Errorf("options can't include %s", o.Name)
		}
	}
	return nil
}

// Execute runs this action
func (a *CallClassifier) Execute(ctx context.Context, run flows.Run, step flows.Step, log events.EventLogger) error {
	cls := a.call(ctx, run, log)
	if cls == nil {
		run.Locals().Set(a.OutputLocal, ModelErrorOutput)
	} else if cls.Confidence < a.MinConfidence {
		run.Locals().Set(a.OutputLocal, ClassifierNoneOutput)
	} else {
		run.Locals().Set(a.OutputLocal, cls.Option)
	}

	return nil
}

func (a *CallClassifier) call(ctx context.Context, run flows.Run, log events.EventLogger) *core.Classification {
	models := run.Session().Assets().Models()
	model := models.Get(a.Model.UUID)
	if model == nil {
		log(events.NewDependencyError(a.Model))
		return nil
	}
	if !model.HasRole(assets.ModelRoleClassify) {
		log(events.NewError(fmt.Sprintf("model %s does not have the classify role", a.Model.UUID), ""))
		return nil
	}

	input, _ := run.EvaluateTemplate(ctx, a.Input, log)

	svc, err := run.Session().Engine().Services().Model(model)
	if err != nil {
		log(events.NewRawError(err))
		return nil
	}

	start := dates.Now()

	cls, err := svc.Classify(ctx, input, a.Options)
	if err != nil {
		log(events.NewRawError(err))
		return nil
	}

	log(events.NewClassifierCalled(model.Reference(), input, a.Options, cls, dates.Since(start)))

	return cls
}

func (a *CallClassifier) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	dependency(a.Model)
	local(a.OutputLocal)
}
