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

// ClassifierCalled events are created when a model is called to classify some input into one of a set of categories.
// Confidence is a measure of how reliable the choice is and isn't necessarily the chosen category's probability. The
// per-category probabilities are only included if the model provides them.
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
//	  "categories": ["Flights", "Hotels"],
//	  "category": "Hotels",
//	  "confidence": 0.86,
//	  "probabilities": {"Flights": 0.29, "Hotels": 0.71},
//	  "tokens": {"input": 123, "output": 5},
//	  "elapsed_ms": 123
//	}
//
// @event classifier_called
type ClassifierCalled struct {
	BaseEvent

	Model         *assets.ModelReference `json:"model" validate:"required"`
	Input         string                 `json:"input"`
	Categories    []string               `json:"categories"`
	Category      string                 `json:"category"`
	Confidence    float64                `json:"confidence"`
	Probabilities map[string]float64     `json:"probabilities,omitempty"`
	Tokens        ModelTokens            `json:"tokens"`
	ElapsedMS     int64                  `json:"elapsed_ms"`
}

// NewClassifierCalled returns a new classifier called event
func NewClassifierCalled(model *assets.ModelReference, input string, categories []string, cls *core.Classification, elapsed time.Duration) *ClassifierCalled {
	return &ClassifierCalled{
		BaseEvent:     NewBaseEvent(TypeClassifierCalled),
		Model:         model,
		Input:         input,
		Categories:    categories,
		Category:      cls.Category,
		Confidence:    cls.Confidence,
		Probabilities: cls.Probabilities,
		Tokens:        ModelTokens{Input: cls.TokensInput, Output: cls.TokensOutput},
		ElapsedMS:     elapsed.Milliseconds(),
	}
}
