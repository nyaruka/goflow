package events

import (
	"github.com/nyaruka/goflow/flows"
)

// TypeResthookSubscriberCalled is the type for our resthook events
const TypeResthookSubscriberCalled string = "resthook_subscriber_called"

// ResthookSubscriberCalledEvent events are created a webhook call is made to a resthook subscriber.
// The event contains the status and status code of the response, as well as a full dump of the
// request and response. Applying this event updates @run.webhook in the context.
//
//   {
//     "type": "resthook_subscriber_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "resthook": "new-registration",
//     "url": "https://api.ipify.org?format=json",
//     "status": "success",
//     "status_code": 200,
//     "request": "POST https://api.ipify.org?format=json",
//     "response": ""
//   }
//
// @event resthook_subscriber_called
type ResthookSubscriberCalledEvent struct {
	baseEvent
	engineOnlyEvent

	Resthook   string              `json:"resthook" validate:"required"`
	URL        string              `json:"url" validate:"required"`
	Status     flows.WebhookStatus `json:"status" validate:"required"`
	StatusCode int                 `json:"status_code" validate:"required"`
	Request    string              `json:"request" validate:"required"`
	Response   string              `json:"response"`
}

// NewResthookSubscriberCalledEvent returns a new resthook called event
func NewResthookSubscriberCalledEvent(resthook string, url string, status flows.WebhookStatus, statusCode int, request string, response string) *ResthookSubscriberCalledEvent {
	return &ResthookSubscriberCalledEvent{
		baseEvent:  newBaseEvent(),
		Resthook:   resthook,
		URL:        url,
		Status:     status,
		StatusCode: statusCode,
		Request:    request,
		Response:   response,
	}
}

// Type returns the type of this event
func (e *ResthookSubscriberCalledEvent) Type() string { return TypeResthookSubscriberCalled }

// Apply applies this event to the given run
func (e *ResthookSubscriberCalledEvent) Apply(run flows.FlowRun) error {
	run.SetWebhook(flows.NewWebhookCall(e.URL, e.Status, e.StatusCode, e.Request, e.Response))
	return nil
}
