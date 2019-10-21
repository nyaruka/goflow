package events

import (
	"time"

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
//     "request": "GET /?format=json HTTP/1.1",
//     "response": "HTTP/1.1 200 OK\r\n\r\n{\"ip\":\"190.154.48.130\"}"
//   }
//
// @event webhook_called
type WebhookCalledEvent struct {
	baseEvent

	URL         string           `json:"url" validate:"required"`
	Status      flows.CallStatus `json:"status" validate:"required"`
	Request     string           `json:"request" validate:"required"`
	Response    string           `json:"response"`
	ElapsedMS   int              `json:"elapsed_ms"`
	Resthook    string           `json:"resthook,omitempty"`
	StatusCode  int              `json:"status_code,omitempty"`
	BodyIgnored bool             `json:"body_ignored,omitempty"`
}

// NewWebhookCalled returns a new webhook called event
func NewWebhookCalled(webhook *flows.WebhookCall) *WebhookCalledEvent {
	return &WebhookCalledEvent{
		baseEvent:   newBaseEvent(TypeWebhookCalled),
		URL:         webhook.URL,
		Status:      webhook.Status,
		Request:     string(webhook.Request),
		Response:    string(webhook.Response),
		ElapsedMS:   int(webhook.TimeTaken / time.Millisecond),
		Resthook:    webhook.Resthook,
		StatusCode:  webhook.StatusCode,
		BodyIgnored: webhook.BodyIgnored,
	}
}
