package events

import (
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/core"
)

func init() {
	registerType(TypeClassifierCalled, func() Event { return &ClassifierCalled{} })
}

// TypeClassifierCalled is the type for our classifier calls events
const TypeClassifierCalled string = "classifier_called"

// ClassifierCalled events are created when a model is called to classify some input as one of a set of options.
// Confidence is a level of how reliable the choice is. Any per-option probabilities the model provides are included for
// diagnostics, but they aren't comparable across models and may not cover every option.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "classifier_called",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "model": {
//	    "uuid": "14115c03-b4c5-49e2-b9ac-390c43e9d7ce",
//	    "name": "GPT-4"
//	  },
//	  "input": "I'd like to book a room for two nights",
//	  "options": ["Flights", "Hotels"],
//	  "option": "Hotels",
//	  "confidence": "high",
//	  "probabilities": {"Flights": 0.29, "Hotels": 0.71},
//	  "tokens": {"input": 123, "output": 5},
//	  "elapsed_ms": 123
//	}
//
// @event classifier_called
type ClassifierCalled struct {
	BaseEvent

	Model         *assets.ModelReference    `json:"model" validate:"required"`
	Input         string                    `json:"input"`
	Options       []string                  `json:"options"`
	Option        string                    `json:"option"`
	Confidence    core.ClassifierConfidence `json:"confidence" validate:"required,classifier_confidence"`
	Probabilities map[string]float64        `json:"probabilities,omitempty"`
	Tokens        core.ModelTokens          `json:"tokens"`
	ElapsedMS     int64                     `json:"elapsed_ms"`
}

// NewClassifierCalled returns a new classifier called event
func NewClassifierCalled(model *assets.ModelReference, input string, options []*core.ClassifierOption, cls *core.Classification, elapsed time.Duration) *ClassifierCalled {
	return &ClassifierCalled{
		BaseEvent:     NewBaseEvent(TypeClassifierCalled),
		Model:         model,
		Input:         input,
		Options:       optionNames(options),
		Option:        cls.Option,
		Confidence:    cls.Confidence,
		Probabilities: cls.Probabilities,
		Tokens:        cls.Tokens,
		ElapsedMS:     elapsed.Milliseconds(),
	}
}

func optionNames(options []*core.ClassifierOption) []string {
	names := make([]string, len(options))
	for i, o := range options {
		names[i] = o.Name
	}
	return names
}
