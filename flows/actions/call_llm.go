package actions

import (
	"context"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeCallLLM, func() flows.Action { return &CallLLMAction{} })
}

// TypeCallLLM is the type for the call LLM action
const TypeCallLLM string = "call_llm"

// LLMErrorOutput is the output used when the LLM call fails
const LLMErrorOutput = "<ERROR>"

// CallLLMAction can be used to call an LLM.
//
// An [event:llm_called] event will be created if the LLM could be called.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "call_llm",
//	  "llm": {
//	    "uuid": "14115c03-b4c5-49e2-b9ac-390c43e9d7ce",
//	    "name": "GPT-4"
//	  },
//	  "instructions": "Categorize the following text as positive or negative",
//	  "input": "@input.text",
//	  "output_local": "_llm_output"
//	}
//
// @action call_llm
type CallLLMAction struct {
	baseAction
	onlineAction

	LLM          *assets.LLMReference `json:"llm"          validate:"required"`
	Instructions string               `json:"instructions" validate:"required"            engine:"evaluated"`
	Input        string               `json:"input"                                       engine:"evaluated"`
	OutputLocal  string               `json:"output_local" validate:"required,local_ref"`
}

// NewCallLLM creates a new call LLM action
func NewCallLLM(uuid flows.ActionUUID, llm *assets.LLMReference, instructions, input, outputLocal string) *CallLLMAction {
	return &CallLLMAction{
		baseAction:   newBaseAction(TypeCallLLM, uuid),
		LLM:          llm,
		Instructions: instructions,
		Input:        input,
		OutputLocal:  outputLocal,
	}
}

// Execute runs this action
func (a *CallLLMAction) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	resp := a.call(ctx, run, logEvent)
	if resp != nil {
		run.Locals().Set(a.OutputLocal, resp.Output)
	} else {
		run.Locals().Set(a.OutputLocal, LLMErrorOutput)
	}

	return nil
}

func (a *CallLLMAction) call(ctx context.Context, run flows.Run, logEvent flows.EventCallback) *flows.LLMResponse {
	llms := run.Session().Assets().LLMs()
	llm := llms.Get(a.LLM.UUID)
	if llm == nil {
		logEvent(events.NewDependencyError(a.LLM))
		return nil
	}

	// substitute any variables in our instructions and input
	instructions, _ := run.EvaluateTemplate(a.Instructions, logEvent)
	input, _ := run.EvaluateTemplate(a.Input, logEvent)

	svc, err := run.Session().Engine().Services().LLM(llm)
	if err != nil {
		logEvent(events.NewError(err.Error()))
		return nil
	}

	start := dates.Now()

	resp, err := svc.Response(ctx, instructions, input, 2500)
	if err != nil {
		logEvent(events.NewError(err.Error()))
		return nil
	}

	logEvent(events.NewLLMCalled(llm, instructions, input, resp, dates.Since(start)))

	return resp
}
