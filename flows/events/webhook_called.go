package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeWebhookCalled, func() flows.Event { return &WebhookCalledEvent{} })
}

// TypeWebhookCalled is the type for our webhook events
const TypeWebhookCalled string = "webhook_called"

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
//     "response": "HTTP/1.1 200 OK\r\n\r\n{\"ip\":\"190.154.48.130\"}"
//   }
//
// @event webhook_called
type WebhookCalledEvent struct {
	baseEvent

	*flows.HTTPTrace

	Resthook    string `json:"resthook,omitempty"`
	BodyIgnored bool   `json:"body_ignored,omitempty"`
}

// NewWebhookCalled returns a new webhook called event
func NewWebhookCalled(call *flows.WebhookCall, status flows.CallStatus, resthook string) *WebhookCalledEvent {
	return &WebhookCalledEvent{
		baseEvent:   newBaseEvent(TypeWebhookCalled),
		HTTPTrace:   flows.NewHTTPTrace(call.Trace, status),
		Resthook:    resthook,
		BodyIgnored: len(call.ResponseBody) > 0 && len(call.ResponseJSON) == 0, // i.e. there was a body but it couldn't be converted to JSON
	}
}
