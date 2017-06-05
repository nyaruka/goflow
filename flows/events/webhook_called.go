package events

import "github.com/nyaruka/goflow/utils"

// TypeWebhookCalled is the type for our webhook events
const TypeWebhookCalled string = "webhook_called"

// WebhookCalledEvent events are created when a webhook is called. The event contains
// the status and status code of the response, as well as a full dump of the
// request and response.
//
// ```
//   {
//    "step_uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//    "created_on": "2006-01-02T15:04:05Z",
//    "type": "webhook_called",
//    "url": "https://api.ipify.org?format=json",
//    "status": "S",
//    "status_code": 200,
//    "request": "GET https://api.ipify.org?format=json",
//    "response": "HTTP/1.1 200 OK {\"ip\":\"190.154.48.130\"}"
//   }
// ```
//
// @event webhook_called
type WebhookCalledEvent struct {
	BaseEvent
	URL        string                      `json:"url"         validate:"required"`
	Status     utils.RequestResponseStatus `json:"status"      validate:"required"`
	StatusCode int                         `json:"status_code" validate:"required"`
	Request    string                      `json:"request"     validate:"required"`
	Response   string                      `json:"response"`
}

// Type returns the type of this event
func (e *WebhookCalledEvent) Type() string { return TypeWebhookCalled }
