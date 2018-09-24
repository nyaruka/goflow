package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeWebhookCalled, func() flows.Event { return &WebhookCalledEvent{} })
}

// TypeWebhookCalled is the type for our webhook events
const TypeWebhookCalled string = "webhook_called"

// WebhookCalledEvent events are created when a webhook is called. The event contains
// the status and status code of the response, as well as a full dump of the
// request and response. If this event has a `result_name`, then applying this event creates
// a new result with that name. If the webhook returned valid JSON, that will be accessible
// through `extra` on the result.
//
//   {
//     "type": "webhook_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "url": "https://api.ipify.org/?format=json",
//     "status": "success",
//     "request": "GET /?format=json HTTP/1.1",
//     "response": "HTTP/1.1 200 OK\r\n\r\n{\"ip\":\"190.154.48.130\"}"
//   }
//
// @event webhook_called
type WebhookCalledEvent struct {
	BaseEvent

	URL      string              `json:"url" validate:"required"`
	Status   flows.WebhookStatus `json:"status" validate:"required"`
	Request  string              `json:"request" validate:"required"`
	Response string              `json:"response"`
}

// NewWebhookCalledEvent returns a new webhook called event
func NewWebhookCalledEvent(webhook *flows.WebhookCall) *WebhookCalledEvent {
	return &WebhookCalledEvent{
		BaseEvent: NewBaseEvent(),
		URL:       webhook.URL(),
		Status:    webhook.Status(),
		Request:   webhook.Request(),
		Response:  webhook.Response(),
	}
}

// Type returns the type of this event
func (e *WebhookCalledEvent) Type() string { return TypeWebhookCalled }
