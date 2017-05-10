package actions

import (
	"net/http"
	"strings"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

const WEBHOOK string = "webhook"

type WebhookAction struct {
	BaseAction
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

func (a *WebhookAction) Type() string { return WEBHOOK }

func (a *WebhookAction) Validate() error {
	return utils.ValidateAll(a)
}

func (a *WebhookAction) Execute(run flows.FlowRun, step flows.Step) error {
	// substitute any variables in our url
	url, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.URL)
	if err != nil {
		run.AddError(step, err)
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
	req, err := http.NewRequest(strings.ToUpper(a.Method), url, nil)

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

	event := events.WebhookEvent{URL: rr.URL(), Status: rr.Status(), StatusCode: rr.StatusCode(), Request: rr.Request(), Response: rr.Response()}
	run.AddEvent(step, &event)

	return nil
}
