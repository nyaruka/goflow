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

// LLMNoMatchOutput is the output used when none of the categories fit the input
const LLMNoMatchOutput = "<CANT>"

// CallClassifier can be used to classify input into one of a set of categories using an LLM. The input field may be
// a template and will be evaluated at runtime.
//
// A [event:classifier_called] event will be created if the LLM could be called. The action sets the local specified
// by `output_local` to the name of the chosen category, to `<CANT>` if none of the categories fit the input, or to
// `<ERROR>` if the call failed.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "call_classifier",
//	  "llm": {
//	    "uuid": "14115c03-b4c5-49e2-b9ac-390c43e9d7ce",
//	    "name": "GPT-4"
//	  },
//	  "input": "@input.text",
//	  "categories": ["Flights", "Hotels"],
//	  "output_local": "_llm_output"
//	}
//
// @action call_classifier
type CallClassifier struct {
	baseAction
	onlineAction

	LLM         *assets.LLMReference `json:"llm"          validate:"required"`
	Input       string               `json:"input"        validate:"max=10000"                                   engine:"evaluated"`
	Categories  []string             `json:"categories"   validate:"required,min=1,max=255,unique,dive,result_category"`
	OutputLocal string               `json:"output_local" validate:"required,local_ref"`
}

// NewCallClassifier creates a new call classifier action
func NewCallClassifier(uuid flows.ActionUUID, llm *assets.LLMReference, input string, categories []string, outputLocal string) *CallClassifier {
	return &CallClassifier{
		baseAction:  newBaseAction(TypeCallClassifier, uuid),
		LLM:         llm,
		Input:       input,
		Categories:  categories,
		OutputLocal: outputLocal,
	}
}

// Execute runs this action
func (a *CallClassifier) Execute(ctx context.Context, run flows.Run, step flows.Step, log events.EventLogger) error {
	cls := a.call(ctx, run, log)
	if cls == nil {
		run.Locals().Set(a.OutputLocal, LLMErrorOutput)
	} else if cls.Category == "" {
		run.Locals().Set(a.OutputLocal, LLMNoMatchOutput)
	} else {
		run.Locals().Set(a.OutputLocal, cls.Category)
	}

	return nil
}

func (a *CallClassifier) call(ctx context.Context, run flows.Run, log events.EventLogger) *core.LLMClassification {
	llms := run.Session().Assets().LLMs()
	llm := llms.Get(a.LLM.UUID)
	if llm == nil {
		log(events.NewDependencyError(a.LLM))
		return nil
	}
	if !llm.HasRole(assets.LLMRoleEngine) {
		log(events.NewError(fmt.Sprintf("LLM %s does not have the engine role", a.LLM.UUID), ""))
		return nil
	}

	input, _ := run.EvaluateTemplate(ctx, a.Input, log)

	svc, err := run.Session().Engine().Services().LLM(llm)
	if err != nil {
		log(events.NewRawError(err))
		return nil
	}

	start := dates.Now()

	cls, err := svc.Classify(ctx, input, a.Categories)
	if err != nil {
		log(events.NewRawError(err))
		return nil
	}

	log(events.NewClassifierCalled(llm.Reference(), input, a.Categories, cls, dates.Since(start)))

	return cls
}

func (a *CallClassifier) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	dependency(a.LLM)
	local(a.OutputLocal)
}
