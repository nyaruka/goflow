package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeServiceCalled, func() flows.Event { return &ServiceCalled{} })
}

// TypeServiceCalled is our type for calling an external service
const TypeServiceCalled string = "service_called"

// ServiceCalled events are created when an engine service is called.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "service_called",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "service": "classifier",
//	  "classifier": {"uuid": "1c06c884-39dd-4ce4-ad9f-9a01cbe6c000", "name": "Booking"},
//	  "http_logs": [
//	    {
//	      "url": "https://api.wit.ai/message?v=20200513&q=hello",
//	      "status": "success",
//	      "request": "GET /message?v=20200513&q=hello HTTP/1.1",
//	      "response": "HTTP/1.1 200 OK\r\n\r\n{\"intents\":[]}",
//	      "created_on": "2006-01-02T15:04:05Z",
//	      "elapsed_ms": 123
//	    }
//	  ]
//	}
//
// @event service_called
type ServiceCalled struct {
	BaseEvent

	Service    string                      `json:"service"`
	Classifier *assets.ClassifierReference `json:"classifier,omitempty"`
	HTTPLogs   []*flows.HTTPLog            `json:"http_logs"`
}

// NewClassifierCalled returns a service called event for a classifier
func NewClassifierCalled(classifier *assets.ClassifierReference, httpLogs []*flows.HTTPLog) *ServiceCalled {
	return &ServiceCalled{
		BaseEvent:  NewBaseEvent(TypeServiceCalled),
		Service:    "classifier",
		Classifier: classifier,
		HTTPLogs:   httpLogs,
	}
}
