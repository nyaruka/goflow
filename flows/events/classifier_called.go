package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeClassifierCalled, func() flows.Event { return &ClassifierCalledEvent{} })
}

// TypeClassifierCalled is our type for the classification event
const TypeClassifierCalled string = "classifier_called"

// ClassifierCalledEvent events are created when a NLU classifier is called.
//
//   {
//     "type": "classifier_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "classifier": {"uuid": "1c06c884-39dd-4ce4-ad9f-9a01cbe6c000", "name": "Booking"},
//     "http_logs": [
//       {
//         "url": "https://api.wit.ai/message?v=20170307&q=hello",
//         "status": "success",
//         "request": "GET /message?v=20170307&q=hello HTTP/1.1",
//         "response": "HTTP/1.1 200 OK\r\n\r\n{\"intents\":[]}",
//         "created_on": "2006-01-02T15:04:05Z",
//         "elapsed_ms": 123
//       }
//     ]
//   }
//
// @event classifier_called
type ClassifierCalledEvent struct {
	baseEvent

	Classifier *assets.ClassifierReference `json:"classifier" validate:"required"`
	HTTPLogs   []*flows.HTTPLog            `json:"http_logs"`
}

// NewClassifierCalled returns a classifier called event
func NewClassifierCalled(classifier *assets.ClassifierReference, httpLogs []*flows.HTTPLog) *ClassifierCalledEvent {
	return &ClassifierCalledEvent{
		baseEvent:  newBaseEvent(TypeClassifierCalled),
		Classifier: classifier,
		HTTPLogs:   httpLogs,
	}
}
