package events

import (
	"fmt"
	"github.com/nyaruka/goflow/flows"
)

// TypeResthookCalled is the type for our resthook events
const TypeResthookCalled string = "resthook_called"

// ResthookSubscriberCall is call to a single subsccriber of a resthook
type ResthookSubscriberCall struct {
	URL        string              `json:"url" validate:"required"`
	Status     flows.WebhookStatus `json:"status" validate:"required"`
	StatusCode int                 `json:"status_code" validate:"required"`
}

// ResthookCalledEvent events are created when a resthook is called. The event contains the status and status code
// of each call to the resthook's subscribers. Applying this event updates @run.webhook in the context.
//
//   {
//     "type": "resthook_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "resthook": "new-registration",
//     "payload": "",
//     "calls": [
//       {
//         "url": "http://localhost:49998/?cmd=success",
//         "status": "success",
//         "status_code": 200
//       },{
//         "url": "https://api.ipify.org?format=json",
//         "status": "success",
//         "status_code": 410
//       }
//     ]
//   }
//
// @event resthook_called
type ResthookCalledEvent struct {
	baseEvent
	engineOnlyEvent

	Resthook string                    `json:"resthook" validate:"required"`
	Payload  string                    `json:"payload"`
	Calls    []*ResthookSubscriberCall `json:"calls" validate:"omitempty,dive"`
}

// NewResthookCalledEvent returns a new resthook called event
func NewResthookCalledEvent(resthook string, payload string, calls []*ResthookSubscriberCall) *ResthookCalledEvent {
	return &ResthookCalledEvent{
		baseEvent: newBaseEvent(),
		Resthook:  resthook,
		Payload:   payload,
		Calls:     calls,
	}
}

// Type returns the type of this event
func (e *ResthookCalledEvent) Type() string { return TypeResthookCalled }

// Apply applies this event to the given run
func (e *ResthookCalledEvent) Apply(run flows.FlowRun) error {
	var lastFailure, asWebhook *ResthookSubscriberCall

	for _, call := range e.Calls {
		if (call.StatusCode >= 200 && call.StatusCode < 300) || call.StatusCode == 410 {
			asWebhook = call
		} else {
			lastFailure = call
		}
	}

	if lastFailure != nil {
		asWebhook = lastFailure
	}

	if asWebhook != nil {
		// An HTTP 410 has a special meaning for resthook and should be considered a success within the run.
		// The onus is on the caller to remove the subscriber from the resthook.
		status := asWebhook.Status
		if asWebhook.StatusCode == 410 {
			status = flows.WebhookStatusSuccess
		}

		request := fmt.Sprintf("POST %s\r\n\r\n%s", asWebhook.URL, e.Payload)
		run.SetWebhook(flows.NewWebhookCall(asWebhook.URL, status, asWebhook.StatusCode, request, ""))
	} else {
		run.SetWebhook(nil)
	}
	return nil
}
