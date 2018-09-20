package events

import (
	"fmt"
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
// of each call to the resthook's subscribers, as well as the payload sent to each subscriber. If this event has a
// `result_name`, then applying this event creates a new result with that name based on one of the calls. The call
// used will the last one unless one has failed, in which case it is used instead. If the call returned valid JSON,
// that will be accessible through `extra` on the result.
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
//     ],
//     "result_name": "ip_check"
//   }
//
// @event resthook_called
type ResthookCalledEvent struct {
	BaseEvent
	engineOnlyEvent

	Resthook   string                    `json:"resthook" validate:"required"`
	Payload    string                    `json:"payload"`
	Calls      []*ResthookSubscriberCall `json:"calls" validate:"omitempty,dive"`
	ResultName string                    `json:"result_name,omitempty"`
}

// NewResthookCalledEvent returns a new resthook called event
func NewResthookCalledEvent(resthook string, payload string, calls []*ResthookSubscriberCall, resultName string) *ResthookCalledEvent {
	return &ResthookCalledEvent{
		BaseEvent:  NewBaseEvent(),
		Resthook:   resthook,
		Payload:    payload,
		Calls:      calls,
		ResultName: resultName,
	}
}

// Type returns the type of this event
func (e *ResthookCalledEvent) Type() string { return TypeResthookCalled }

// Apply applies this event to the given run
func (e *ResthookCalledEvent) Apply(run flows.FlowRun) error {
	// no result namem then nothing to do
	if e.ResultName == "" {
		return nil
	}

	// select one of our calls to become the webhook result
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

		nodeUUID := run.GetStep(e.StepUUID()).NodeUUID()

		e.saveWebhookResult(run, e.ResultName, flows.NewWebhookCall(asWebhook.URL, asWebhook.Status, asWebhook.StatusCode, request, asWebhook.Response), nodeUUID)

		// TODO remove
		run.SetWebhook(flows.NewWebhookCall(asWebhook.URL, asWebhook.Status, asWebhook.StatusCode, request, asWebhook.Response))
	} else {
		run.SetWebhook(nil)
	}
	return nil
}
