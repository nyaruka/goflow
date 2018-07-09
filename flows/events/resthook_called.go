package events

import (
	"fmt"
	"net/http"

	"github.com/nyaruka/goflow/flows"
)

// TypeResthookCalled is the type for our resthook events
const TypeResthookCalled string = "resthook_called"

// ResthookSubscriberCall is call to a single subsccriber of a resthook
type ResthookSubscriberCall struct {
	URL        string              `json:"url" validate:"required"`
	Status     flows.WebhookStatus `json:"status" validate:"required"`
	StatusCode int                 `json:"status_code" validate:"required"`
	Response   string              `json:"response"`
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
		URL:        webhook.URL(),
		Status:     status,
		StatusCode: webhook.StatusCode(),
		Response:   webhook.Response(),
	}
}

// ResthookCalledEvent events are created when a resthook is called. The event contains the status and status code
// of each call to the resthook's subscribers, as well as the payload sent to each subscriber. Applying this event
// updates @run.webhook in the context to the results of the last subscriber call. However if one of the subscriber
// calls fails, then it is used to update @run.webhook instead.
//
//   {
//     "type": "resthook_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "resthook": "new-registration",
//     "payload": "{...}",
//     "calls": [
//       {
//         "url": "http://localhost:49998/?cmd=success",
//         "status": "success",
//         "status_code": 200,
//         "response": "{\"errors\":[]}"
//       },{
//         "url": "https://api.ipify.org?format=json",
//         "status": "success",
//         "status_code": 410,
//         "response": "{\"errors\":[\"Unsubscribe\"]}"
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
		if call.Status == flows.WebhookStatusSuccess {
			asWebhook = call
		} else {
			lastFailure = call
		}
	}

	if lastFailure != nil {
		asWebhook = lastFailure
	}

	if asWebhook != nil {
		// reconstruct the request dump
		request := fmt.Sprintf("POST %s\r\n\r\n%s", asWebhook.URL, e.Payload)

		run.SetWebhook(flows.NewWebhookCall(asWebhook.URL, asWebhook.Status, asWebhook.StatusCode, request, asWebhook.Response))
	} else {
		run.SetWebhook(nil)
	}
	return nil
}
