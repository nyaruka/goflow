package events

import (
	"fmt"
	"strconv"

	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeWebhookCalled, func() flows.Event { return &WebhookCalledEvent{} })
}

// TypeWebhookCalled is the type for our webhook events
const TypeWebhookCalled string = "webhook_called"

// WebhookCalledEvent events are created when a webhook is called. The event contains
// the status and status code of the response, as well as a full dump of the
// request and response. If this event has a `reult_name`, then applying this event creates
// a new result with that name. If the webhook returned valid JSON, that will be accessible
// through `extra` on the result.
//
//   {
//     "type": "webhook_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "url": "https://api.ipify.org?format=json",
//     "status": "success",
//     "status_code": 200,
//     "request": "GET https://api.ipify.org?format=json",
//     "response": "HTTP/1.1 200 OK {\"ip\":\"190.154.48.130\"}",
//     "result_name": "ip_check"
//   }
//
// @event webhook_called
type WebhookCalledEvent struct {
	BaseEvent
	engineOnlyEvent

	URL        string              `json:"url" validate:"required"`
	Status     flows.WebhookStatus `json:"status" validate:"required"`
	StatusCode int                 `json:"status_code"`
	Request    string              `json:"request" validate:"required"`
	Response   string              `json:"response"`
	ResultName string              `json:"result_name,omitempty"`
}

// NewWebhookCalledEvent returns a new webhook called event
func NewWebhookCalledEvent(webhook *flows.WebhookCall, resultName string) *WebhookCalledEvent {
	return &WebhookCalledEvent{
		BaseEvent:  NewBaseEvent(),
		URL:        webhook.URL(),
		Status:     webhook.Status(),
		StatusCode: webhook.StatusCode(),
		Request:    webhook.Request(),
		Response:   webhook.Response(),
		ResultName: resultName,
	}
}

// Type returns the type of this event
func (e *WebhookCalledEvent) Type() string { return TypeWebhookCalled }

// Apply applies this event to the given run
func (e *WebhookCalledEvent) Apply(run flows.FlowRun) error {
	// TODO remove
	run.SetWebhook(flows.NewWebhookCall(e.URL, e.Status, e.StatusCode, e.Request, e.Response))

	if e.ResultName != "" {
		nodeUUID := run.GetStep(e.StepUUID()).NodeUUID()
		e.saveWebhookResult(run, e.ResultName, flows.NewWebhookCall(e.URL, e.Status, e.StatusCode, e.Request, e.Response), nodeUUID)
	}
	return nil
}

func (e *BaseEvent) saveWebhookResult(run flows.FlowRun, resultName string, webhook *flows.WebhookCall, nodeUUID flows.NodeUUID) {
	input := fmt.Sprintf("%s %s", webhook.Method(), webhook.URL())
	value := strconv.Itoa(webhook.StatusCode())
	category := string(webhook.Status())
	extra := []byte(webhook.Body()) // TODO

	run.Results().Save(resultName, value, category, "", nodeUUID, &input, extra, e.CreatedOn())
}
