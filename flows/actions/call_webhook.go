package actions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"

	"golang.org/x/net/http/httpguts"
)

func isValidURL(u string) bool {
	if utf8.RuneCountInString(u) > 2048 {
		return false
	}
	_, err := url.Parse(u)
	return err == nil
}

func init() {
	registerType(TypeCallWebhook, func() flows.Action { return &CallWebhookAction{} })
}

// TypeCallWebhook is the type for the call webhook action
const TypeCallWebhook string = "call_webhook"

// CallWebhookAction can be used to call an external service. The body, header and url fields may be
// templates and will be evaluated at runtime. A [event:webhook_called] event will be created based on
// the results of the HTTP call. If this action has a `result_name`, then additionally it will create
// a new result with that name. The value of the result will be the status code and the category will be
// `Success` or `Failed`. If the webhook returned valid JSON which is less than 10000 bytes, that will be
// accessible through `extra` on the result. The last JSON response from a webhook call in the current
// sprint will additionally be accessible in expressions as `@webhook` regardless of size.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "call_webhook",
//	  "method": "GET",
//	  "url": "http://localhost:49998/?cmd=success",
//	  "headers": {
//	    "Authorization": "Token AAFFZZHH"
//	  },
//	  "result_name": "webhook"
//	}
//
// @action call_webhook
type CallWebhookAction struct {
	baseAction
	onlineAction

	Method     string            `json:"method"                                   validate:"required,http_method"`
	URL        string            `json:"url"                   engine:"evaluated" validate:"required"`
	Headers    map[string]string `json:"headers,omitempty"     engine:"evaluated"`
	Body       string            `json:"body,omitempty"        engine:"evaluated"`
	ResultName string            `json:"result_name,omitempty"                    validate:"omitempty,result_name"`
}

// NewCallWebhook creates a new call webhook action
func NewCallWebhook(uuid flows.ActionUUID, method string, url string, headers map[string]string, body string, resultName string) *CallWebhookAction {
	return &CallWebhookAction{
		baseAction: newBaseAction(TypeCallWebhook, uuid),
		Method:     method,
		URL:        url,
		Headers:    headers,
		Body:       body,
		ResultName: resultName,
	}
}

// Validate validates our action is valid
func (a *CallWebhookAction) Validate() error {
	for key := range a.Headers {
		if !httpguts.ValidHeaderFieldName(key) {
			return fmt.Errorf("header '%s' is not a valid HTTP header", key)
		}
	}

	return nil
}

// Execute runs this action
func (a *CallWebhookAction) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	url, _ := run.EvaluateTemplate(a.URL, logEvent)
	url = strings.TrimSpace(url)

	if url == "" {
		logEvent(events.NewError("webhook URL evaluated to empty string"))
		return nil
	}
	if !isValidURL(url) {
		logEvent(events.NewError(fmt.Sprintf("webhook URL evaluated to an invalid URL: '%s'", url)))
		return nil
	}

	method := strings.ToUpper(a.Method)
	body := a.Body

	// substitute any body variables
	if body != "" {
		// webhook bodies aren't truncated like other templates
		body, _ = run.EvaluateTemplateText(body, nil, false, logEvent)
	}

	call := a.call(ctx, run, step, url, method, body, logEvent)
	run.SetWebhook(call)

	return nil
}

// Execute runs this action
func (a *CallWebhookAction) call(ctx context.Context, run flows.Run, step flows.Step, url, method, body string, logEvent flows.EventCallback) *flows.WebhookCall {
	// build our request
	req, err := httpx.NewRequest(ctx, method, url, strings.NewReader(body), nil)
	if err != nil {
		// in theory this can't happen because we're already validating the method and the URL.. but just in case
		logEvent(events.NewError(err.Error()))
		return nil
	}

	// add the custom headers, substituting any template vars
	for key, value := range a.Headers {
		headerValue, _ := run.EvaluateTemplate(value, logEvent)

		req.Header.Add(key, headerValue)
	}

	svc, err := run.Session().Engine().Services().Webhook(run.Session().Assets())
	if err != nil {
		logEvent(events.NewError(err.Error()))
		return nil
	}

	trace, err := svc.Call(req)
	if err != nil {
		logEvent(events.NewError(err.Error()))
	}

	if trace != nil {
		call := flows.NewWebhookCall(trace)
		status := callStatus(trace, err, false)

		logEvent(events.NewWebhookCalled(trace, status, ""))

		if a.ResultName != "" {
			a.saveWebhookResult(run, step, a.ResultName, call, status, logEvent)
		}

		return call
	}

	return nil
}

func (a *CallWebhookAction) Inspect(result func(*flows.ResultInfo)) {
	if a.ResultName != "" {
		result(flows.NewResultInfo(a.ResultName, webhookCategories))
	}
}

// determines the webhook status from the HTTP status code
func callStatus(t *httpx.Trace, err error, isResthook bool) flows.CallStatus {
	if t.Response == nil || err != nil {
		return flows.CallStatusConnectionError
	}
	if isResthook && t.Response.StatusCode == http.StatusGone {
		// https://zapier.com/developer/documentation/v2/rest-hooks/
		return flows.CallStatusSubscriberGone
	}
	if t.Response.StatusCode/100 == 2 {
		return flows.CallStatusSuccess
	}
	return flows.CallStatusResponseError
}
