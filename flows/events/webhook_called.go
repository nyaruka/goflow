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
// the URL and the status of the response, as well as a full dump of the
// request and response.
//
//   {
//     "type": "webhook_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "url": "http://localhost:49998/?cmd=success",
//     "status": "success",
//     "elapsed_ms": 123,
//     "request": "GET /?format=json HTTP/1.1",
//     "response": "HTTP/1.1 200 OK\r\n\r\n{\"ip\":\"190.154.48.130\"}"
//   }
//
// @event webhook_called
type WebhookCalledEvent struct {
	BaseEvent

	URL       string              `json:"url" validate:"required"`
	Resthook  string              `json:"resthook,omitempty"`
	Status    flows.WebhookStatus `json:"status" validate:"required"`
	ElapsedMS int                 `json:"elapsed_ms"`
	Request   string              `json:"request" validate:"required"`
	Response  string              `json:"response"`
}

// NewWebhookCalledEvent returns a new webhook called event
func NewWebhookCalledEvent(webhook *flows.WebhookCall) *WebhookCalledEvent {
	return &WebhookCalledEvent{
		BaseEvent: NewBaseEvent(TypeWebhookCalled),
		URL:       webhook.URL(),
		Resthook:  webhook.Resthook(),
		Status:    webhook.Status(),
		ElapsedMS: int(webhook.TimeTaken().Nanoseconds() / 1e6),
		Request:   webhook.Request(),
		Response:  webhook.Response(),
	}
}
