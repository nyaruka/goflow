package events

import (
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeWebhookCalled, func() flows.Event { return &WebhookCalled{} })
}

// TypeWebhookCalled is the type for our webhook events
const TypeWebhookCalled string = "webhook_called"

// WebhookCalled events are created when a webhook is called. The event contains
// the URL and the status of the response, as well as a full dump of the
// request and response.
//
//	{
//	  "type": "webhook_called",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "url": "http://localhost:49998/?cmd=success",
//	  "status": "success",
//	  "status_code": 200,
//	  "elapsed_ms": 123,
//	  "retries": 0,
//	  "request": "GET /?format=json HTTP/1.1",
//	  "response": "HTTP/1.1 200 OK\r\n\r\n{\"ip\":\"190.154.48.130\"}"
//	}
//
// @event webhook_called
type WebhookCalled struct {
	BaseEvent

	*flows.HTTPLogWithoutTime

	Resthook string `json:"resthook,omitempty"`
}

// NewWebhookCalled returns a new webhook called event
func NewWebhookCalled(trace *httpx.Trace, status flows.CallStatus, resthook string) *WebhookCalled {
	return &WebhookCalled{
		BaseEvent:          NewBaseEvent(TypeWebhookCalled),
		HTTPLogWithoutTime: flows.NewHTTPLogWithoutTime(trace, status, nil),
		Resthook:           resthook,
	}
}
