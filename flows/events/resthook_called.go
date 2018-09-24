package events

import (
	"net/http"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeResthookCalled, func() flows.Event { return &ResthookCalledEvent{} })
}

// TypeResthookCalled is the type for our resthook events
const TypeResthookCalled string = "resthook_called"

// ResthookSubscriberCall is call to a single subsccriber of a resthook
type ResthookSubscriberCall struct {
	URL      string              `json:"url" validate:"required"`
	Status   flows.WebhookStatus `json:"status" validate:"required"`
	Request  string              `json:"request" validate:"required"`
	Response string              `json:"response"`
}

// NewResthookSubscriberCall creates a new subscriber call from the given webhook call
func NewResthookSubscriberCall(webhook *flows.WebhookCall) *ResthookSubscriberCall {
	// An HTTP 410 has a special meaning for resthook and should be considered a success within the run.
	// The onus is on the caller to remove the subscriber from the resthook.
	status := webhook.Status()
	if webhook.StatusCode() == http.StatusGone {
		status = flows.WebhookStatusSuccess
	}

	return &ResthookSubscriberCall{
		URL:      webhook.URL(),
		Status:   status,
		Request:  webhook.Request(),
		Response: webhook.Response(),
	}
}

// ResthookCalledEvent events are created when a resthook is called. The event contains the status and status code
// of each call to the resthook's subscribers, as well as the payload sent to each subscriber. If this event has a
// `result_name`, then applying this event creates a new result with that name based on one of the calls. The call
// used will the last one unless one has failed, in which case it is used instead. If the call returned valid JSON,
// that will be accessible through `extra` on the result.
//
//   {
//     "type": "resthook_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "resthook": "new-registration",
//     "calls": [
//       {
//         "url": "http://localhost:49998/?cmd=success",
//         "status": "success",
//         "request": "POST /?cmd=success HTTP/1.1",
//         "response": "HTTP/1.1 200 OK\r\n\r\n{\"errors\":[]}"
//       },{
//         "url": "https://api.ipify.org/?format=json",
//         "status": "success",
//         "request": "POST /?format=json HTTP/1.1",
//         "response": "HTTP/1.1 410 Gone\r\n\r\n{\"errors\":[\"Unsubscribe\"]}"
//       }
//     ]
//   }
//
// @event resthook_called
type ResthookCalledEvent struct {
	BaseEvent

	Resthook string                    `json:"resthook" validate:"required"`
	Calls    []*ResthookSubscriberCall `json:"calls" validate:"omitempty,dive"`
}

// NewResthookCalledEvent returns a new resthook called event
func NewResthookCalledEvent(resthook string, webhooks []*flows.WebhookCall) *ResthookCalledEvent {
	calls := make([]*ResthookSubscriberCall, len(webhooks))
	for w := range webhooks {
		calls[w] = NewResthookSubscriberCall(webhooks[w])
	}

	return &ResthookCalledEvent{
		BaseEvent: NewBaseEvent(),
		Resthook:  resthook,
		Calls:     calls,
	}
}

// Type returns the type of this event
func (e *ResthookCalledEvent) Type() string { return TypeResthookCalled }
