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
	registerType(TypeCallLLM, func() flows.Action { return &CallLLM{} })
}

// TypeCallLLM is the type for the call LLM action
const TypeCallLLM string = "call_llm"

// ModelErrorOutput is the output used when a model call fails
const ModelErrorOutput = "<ERROR>"

// CallLLM can be used to call an LLM. The instructions and input fields may be templates and will
// be evaluated at runtime.
//
// An [event:llm_called] event will be created if the model could be called. The action sets the local
// specified by `output_local` to the output of the LLM, or to `<ERROR>` if the call failed.
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
type CallLLM struct {
	baseAction
	onlineAction

	LLM          *assets.ModelReference `json:"llm"          validate:"required"`
	Instructions string                 `json:"instructions" validate:"required,max=10000"  engine:"evaluated"`
	Input        string                 `json:"input"        validate:"max=10000"           engine:"evaluated"`
	OutputLocal  string                 `json:"output_local" validate:"required,local_ref"`
}

// NewCallLLM creates a new call LLM action
func NewCallLLM(uuid flows.ActionUUID, llm *assets.ModelReference, instructions, input, outputLocal string) *CallLLM {
	return &CallLLM{
		baseAction:   newBaseAction(TypeCallLLM, uuid),
		LLM:          llm,
		Instructions: instructions,
		Input:        input,
		OutputLocal:  outputLocal,
	}
}

// Execute runs this action
func (a *CallLLM) Execute(ctx context.Context, run flows.Run, step flows.Step, log events.EventLogger) error {
	resp := a.call(ctx, run, log)
	if resp != nil {
		run.Locals().Set(a.OutputLocal, resp.Output)
	} else {
		run.Locals().Set(a.OutputLocal, ModelErrorOutput)
	}

	return nil
}

func (a *CallLLM) call(ctx context.Context, run flows.Run, log events.EventLogger) *core.ModelResponse {
	models := run.Session().Assets().Models()
	model := models.Get(a.LLM.UUID)
	if model == nil {
		log(events.NewDependencyError(a.LLM))
		return nil
	}
	if !model.HasRole(assets.ModelRoleEngine) {
		log(events.NewError(fmt.Sprintf("model %s does not have the engine role", a.LLM.UUID), ""))
		return nil
	}

	// substitute any variables in our instructions and input
	instructions, _ := run.EvaluateTemplate(ctx, a.Instructions, log)
	input, _ := run.EvaluateTemplate(ctx, a.Input, log)

	svc, err := run.Session().Engine().Services().Model(model)
	if err != nil {
		log(events.NewRawError(err))
		return nil
	}

	start := dates.Now()

	resp, err := svc.Response(ctx, instructions, input, 2500)
	if err != nil {
		log(events.NewRawError(err))
		return nil
	}

	log(events.NewLLMCalled(model.Reference(), instructions, input, resp, dates.Since(start)))

	return resp
}

func (a *CallLLM) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	dependency(a.LLM)
	local(a.OutputLocal)
}
