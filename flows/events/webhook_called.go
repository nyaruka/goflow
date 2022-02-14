package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeWebhookCalled, func() flows.Event { return &WebhookCalledEvent{} })
}

// TypeWebhookCalled is the type for our webhook events
const TypeWebhookCalled string = "webhook_called"

type Extraction string

const (
	ExtractionNone    Extraction = "none"    // no response or body was empty
	ExtractionValid   Extraction = "valid"   // body was valid JSON
	ExtractionCleaned Extraction = "cleaned" // body could be made into JSON with some cleaning
	ExtractionIgnored Extraction = "ignored" // body couldn't be made into JSON and was ignored
)

// WebhookCalledEvent events are created when a webhook is called. The event contains
// the URL and the status of the response, as well as a full dump of the
// request and response.
//
//   {
//     "type": "webhook_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "url": "http://localhost:49998/?cmd=success",
//     "status": "success",
//     "status_code": 200,
//     "elapsed_ms": 123,
//     "retries": 0,
//     "request": "GET /?format=json HTTP/1.1",
//     "response": "HTTP/1.1 200 OK\r\n\r\n{\"ip\":\"190.154.48.130\"}",
//     "extraction": "valid"
//   }
//
// @event webhook_called
type WebhookCalledEvent struct {
	BaseEvent

	*flows.HTTPTrace

	Resthook   string     `json:"resthook,omitempty"`
	Extraction Extraction `json:"extraction"`
}

// NewWebhookCalled returns a new webhook called event
func NewWebhookCalled(call *flows.WebhookCall, status flows.CallStatus, resthook string) *WebhookCalledEvent {
	extraction := ExtractionNone
	if len(call.ResponseBody) > 0 {
		if len(call.ResponseJSON) > 0 {
			if call.ResponseCleaned {
				extraction = ExtractionCleaned
			} else {
				extraction = ExtractionValid
			}
		} else {
			extraction = ExtractionIgnored
		}
	}

	return &WebhookCalledEvent{
		BaseEvent:  NewBaseEvent(TypeWebhookCalled),
		HTTPTrace:  flows.NewHTTPTrace(call.Trace, status),
		Resthook:   resthook,
		Extraction: extraction,
	}
}
