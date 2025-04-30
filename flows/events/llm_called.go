package events

import (
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeLLMCalled, func() flows.Event { return &LLMCalledEvent{} })
}

// TypeLLMCalled is the type for our LLM calls events
const TypeLLMCalled string = "llm_called"

// LLMCalledEvent events are created when an LLM is called.
//
//	{
//	  "uuid": "019688A6-41d2-7366-958a-630e35c62431",
//	  "type": "llm_called",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "llm": {
//	    "uuid": "14115c03-b4c5-49e2-b9ac-390c43e9d7ce",
//	    "name": "GPT-4"
//	  },
//	  "instructions": "Categorize the following text as Positive or Negative",
//	  "input": "Please stop messaging me",
//	  "output": "Positive",
//	  "tokens_used": 567,
//	  "elapsed_ms": 123
//	}
//
// @event llm_called
type LLMCalledEvent struct {
	BaseEvent

	LLM          *assets.LLMReference `json:"llm" validate:"required"`
	Instructions string               `json:"instructions"`
	Input        string               `json:"input"`
	Output       string               `json:"output"`
	TokensUsed   int64                `json:"tokens_used"`
	ElapsedMS    int64                `json:"elapsed_ms"`
}

// NewLLMCalled returns a new LLM called event
func NewLLMCalled(llm *flows.LLM, instructions, input string, resp *flows.LLMResponse, elapsed time.Duration) *LLMCalledEvent {
	return &LLMCalledEvent{
		BaseEvent:    NewBaseEvent(TypeLLMCalled),
		LLM:          llm.Reference(),
		Instructions: instructions,
		Input:        input,
		Output:       resp.Output,
		TokensUsed:   resp.TokensUsed,
		ElapsedMS:    elapsed.Milliseconds(),
	}
}
