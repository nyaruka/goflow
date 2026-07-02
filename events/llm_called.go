package events

import (
	"time"

	"github.com/nyaruka/goflow/assets"
)

func init() {
	registerType(TypeLLMCalled, func() Event { return &LLMCalled{} })
}

// TypeLLMCalled is the type for our LLM calls events
const TypeLLMCalled string = "llm_called"

// LLMCalled events are created when an LLM is called.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "llm_called",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "llm": {
//	    "uuid": "14115c03-b4c5-49e2-b9ac-390c43e9d7ce",
//	    "name": "GPT-4"
//	  },
//	  "instructions": "Categorize the following text as Positive or Negative",
//	  "input": "Please stop messaging me",
//	  "output": "Positive",
//	  "tokens": {"input": 234, "output": 333},
//	  "elapsed_ms": 123
//	}
//
// @event llm_called
type LLMCalled struct {
	BaseEvent

	LLM          *assets.LLMReference `json:"llm" validate:"required"`
	Instructions string               `json:"instructions"`
	Input        string               `json:"input"`
	Output       string               `json:"output"`
	Tokens       LLMTokens            `json:"tokens"`
	ElapsedMS    int64                `json:"elapsed_ms"`
}

type LLMTokens struct {
	Input  int64 `json:"input"`
	Output int64 `json:"output"`
}

// NewLLMCalled returns a new LLM called event
func NewLLMCalled(llm *assets.LLMReference, instructions, input string, resp *LLMResponse, elapsed time.Duration) *LLMCalled {
	return &LLMCalled{
		BaseEvent:    NewBaseEvent(TypeLLMCalled),
		LLM:          llm,
		Instructions: instructions,
		Input:        input,
		Output:       resp.Output,
		Tokens:       LLMTokens{Input: resp.TokensInput, Output: resp.TokensOutput},
		ElapsedMS:    elapsed.Milliseconds(),
	}
}
