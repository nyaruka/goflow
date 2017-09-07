package actions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeCallWebhook is the type for our webhook action
const TypeCallWebhook string = "call_webhook"

// WebhookAction can be used to call an external service and insert the results in @run.webhook
// context variable. The body, header and url fields may be templates and will be evaluated at runtime.
//
// A `webhook_called` event will be created based on the results of the HTTP call.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "call_webhook",
//     "method": "get",
//     "url": "https://api.ipify.org?format=json",
//     "headers": {
//	      "Authorization": "Token AAFFZZHH"
//     }
//   }
// ```
//
// @action call_webhook
type WebhookAction struct {
	BaseAction
	Method  string            `json:"method"                validate:"required"`
	URL     string            `json:"url"                   validate:"required"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

// Type returns the type of this action
func (a *WebhookAction) Type() string { return TypeCallWebhook }

// Validate validates the fields on this action
func (a *WebhookAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *WebhookAction) Execute(run flows.FlowRun, step flows.Step) ([]flows.Event, error) {
	log := make([]flows.Event, 0)

	// substitute any variables in our url
	url, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.URL)
	if err != nil {
		log = append(log, events.NewErrorEvent(err))
	}
	if url == "" {
		log = append(log, events.NewErrorEvent(fmt.Errorf("call_webhook URL evaluated to empty string, skipping")))
		return log, nil
	}

	// substitute any body variables
	body := a.Body
	if body != "" {
		body, err = excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.Body)
		if err != nil {
			log = append(log, events.NewErrorEvent(err))
		}
	}

	// build our request
	req, err := http.NewRequest(strings.ToUpper(a.Method), url, strings.NewReader(body))
	if err != nil {
		log = append(log, events.NewErrorEvent(err))
		return log, nil
	}

	// add our headers, substituting any template vars
	for key, value := range a.Headers {
		headerValue, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), value)
		if err != nil {
			log = append(log, events.NewErrorEvent(err))
		}

		req.Header.Add(key, headerValue)
	}

	rr, err := utils.MakeHTTPRequest(req)
	if err != nil {
		log = append(log, events.NewErrorEvent(err))
	}
	run.SetWebhook(rr)

	log = append(log, events.NewWebhookCalledEvent(rr.URL(), rr.Status(), rr.StatusCode(), rr.Request(), rr.Response()))
	return log, nil
}
