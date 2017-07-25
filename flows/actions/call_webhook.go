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
func (a *WebhookAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute runs this action
func (a *WebhookAction) Execute(run flows.FlowRun, step flows.Step) error {
	// substitute any variables in our url
	url, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.URL)
	if err != nil {
		run.AddError(step, err)
	}
	if url == "" {
		run.AddError(step, fmt.Errorf("call_webhook URL evaluated to empty string, skipping"))
		return nil
	}

	// substitute any body variables
	body := a.Body
	if body != "" {
		body, err = excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.Body)
		if err != nil {
			run.AddError(step, err)
		}
	}

	// build our request
	req, err := http.NewRequest(strings.ToUpper(a.Method), url, strings.NewReader(body))
	if err != nil {
		run.AddError(step, err)
		return nil
	}

	// add our headers, substituting any template vars
	for key, value := range a.Headers {
		headerValue, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), value)
		if err != nil {
			run.AddError(step, err)
		}

		req.Header.Add(key, headerValue)
	}

	rr, err := utils.MakeHTTPRequest(req)
	if err != nil {
		run.AddError(step, err)
	}
	run.SetWebhook(rr)

	run.ApplyEvent(step, a, events.NewWebhookCalledEvent(rr.URL(), rr.Status(), rr.StatusCode(), rr.Request(), rr.Response()))

	return nil
}
