package events

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
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
//     "response": "HTTP/1.1 200 OK\r\n\r\n{\"ip\":\"190.154.48.130\"}",
//     "result_name": "IP Check"
//   }
//
// @event webhook_called
type WebhookCalledEvent struct {
	BaseEvent
	engineOnlyEvent

	URL        string              `json:"url" validate:"required"`
	Status     flows.WebhookStatus `json:"status" validate:"required"`
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
		Request:    webhook.Request(),
		Response:   webhook.Response(),
		ResultName: resultName,
	}
}

// Type returns the type of this event
func (e *WebhookCalledEvent) Type() string { return TypeWebhookCalled }

// Apply applies this event to the given run
func (e *WebhookCalledEvent) Apply(run flows.FlowRun) error {
	if e.ResultName != "" {
		nodeUUID := run.GetStep(e.StepUUID()).NodeUUID()
		return e.saveWebhookResult(run, e.ResultName, e.URL, e.Request, e.Response, nodeUUID)
	}
	return nil
}

var webhookStatusCategories = map[flows.WebhookStatus]string{
	flows.WebhookStatusSuccess:         "Success",
	flows.WebhookStatusResponseError:   "Failure",
	flows.WebhookStatusConnectionError: "Failure",
}

func (e *BaseEvent) saveWebhookResult(run flows.FlowRun, resultName, url, requestTrace string, responseTrace string, nodeUUID flows.NodeUUID) error {
	webhook, err := flows.ReconstructWebhookCall(url, requestTrace, responseTrace)
	if err != nil {
		return err
	}

	input := fmt.Sprintf("%s %s", webhook.Method(), webhook.URL())
	value := strconv.Itoa(webhook.StatusCode())
	category := webhookStatusCategories[webhook.Status()]

	body := []byte(webhook.Body())
	var extra json.RawMessage

	// try to parse body as JSON
	if utils.IsValidJSON(body) {
		// if that was successful, the body is valid JSON and extra is the body
		extra = body
	} else {
		// if not, treat body as text and encode as a JSON string
		extra, _ = json.Marshal(string(body))
	}

	run.Results().Save(resultName, value, category, "", nodeUUID, &input, extra, e.CreatedOn())
	return nil
}
