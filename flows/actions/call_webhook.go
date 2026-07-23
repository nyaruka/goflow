package actions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/core/events"
	"github.com/nyaruka/goflow/flows"
	"golang.org/x/net/http/httpguts"
)

func isValidURL(u string) bool {
	if utf8.RuneCountInString(u) > 8192 {
		return false
	}
	_, err := url.Parse(u)
	return err == nil
}

// approximates the size in bytes of the request as it will be serialized on the wire
func requestSize(method, url string, headers map[string]string, body string) int {
	size := len(method) + len(url) + 12 // request line spaces, version and CRLF
	for key, value := range headers {
		size += len(key) + len(value) + 4 // colon, space and CRLF
	}
	return size + 2 + len(body) // blank line between headers and body
}

func init() {
	registerType(TypeCallWebhook, func() flows.Action { return &CallWebhook{} })
}

// TypeCallWebhook is the type for the call webhook action
const TypeCallWebhook string = "call_webhook"

// CallWebhook can be used to call an external service. The body, header and url fields may be
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
type CallWebhook struct {
	baseAction
	onlineAction

	Method     string            `json:"method"                                   validate:"required,http_method"`
	URL        string            `json:"url"                   engine:"evaluated" validate:"required,max=8192"`
	Headers    map[string]string `json:"headers,omitempty"     engine:"evaluated" validate:"max=10,dive,keys,max=100,endkeys,max=5000"`
	Body       string            `json:"body,omitempty"        engine:"evaluated" validate:"max=20000"`
	ResultName string            `json:"result_name,omitempty"                    validate:"omitempty,result_name"`
}

// NewCallWebhook creates a new call webhook action
func NewCallWebhook(uuid flows.ActionUUID, method string, url string, headers map[string]string, body string, resultName string) *CallWebhook {
	return &CallWebhook{
		baseAction: newBaseAction(TypeCallWebhook, uuid),
		Method:     method,
		URL:        url,
		Headers:    headers,
		Body:       body,
		ResultName: resultName,
	}
}

// Validate validates our action is valid
func (a *CallWebhook) Validate() error {
	for key := range a.Headers {
		if !httpguts.ValidHeaderFieldName(key) {
			return fmt.Errorf("header '%s' is not a valid HTTP header", key)
		}
	}

	return nil
}

// Execute runs this action
func (a *CallWebhook) Execute(ctx context.Context, run flows.Run, step flows.Step, log events.EventLogger) error {
	url, _ := run.EvaluateTemplate(ctx, a.URL, log)
	url = strings.TrimSpace(url)

	if url == "" {
		log(events.NewError("Webhook URL evaluated to empty string", ""))
		return nil
	}
	if !isValidURL(url) {
		log(events.NewError(fmt.Sprintf("Webhook URL evaluated to an invalid URL: '%s'", stringsx.TruncateEllipsis(url, 255)), ""))
		return nil
	}

	method := strings.ToUpper(a.Method)

	// substitute any header variables
	headers := make(map[string]string, len(a.Headers))
	for key, value := range a.Headers {
		headers[key], _ = run.EvaluateTemplate(ctx, value, log)
	}

	body := a.Body

	// substitute any body variables (bodies aren't truncated like other templates)
	if body != "" {
		body, _ = run.EvaluateTemplateText(ctx, body, nil, false, log)
	}

	// evaluated bodies aren't truncated like other templates because that could produce invalid JSON, but there
	// has to be an absolute cap on what we're prepared to send - e.g. a body template that embeds @webhook
	// multiple times could otherwise evaluate to something enormous - so we limit the overall request size
	maxRequestBytes := run.Session().Engine().Options().MaxRequestBytes
	if size := requestSize(method, url, headers, body); size > maxRequestBytes {
		log(events.NewError(fmt.Sprintf("Webhook request evaluated to %d bytes, exceeding the limit of %d", size, maxRequestBytes), events.ErrorCodeWebhookRequestSize))
		return nil
	}

	call := a.call(ctx, run, step, url, method, headers, body, log)
	run.SetWebhook(call)

	return nil
}

// Execute runs this action
func (a *CallWebhook) call(ctx context.Context, run flows.Run, step flows.Step, url, method string, headers map[string]string, body string, log events.EventLogger) *flows.WebhookCall {
	// build our request
	req, err := httpx.NewRequest(ctx, method, url, strings.NewReader(body), headers)
	if err != nil {
		// in theory this can't happen because we're already validating the method and the URL.. but just in case
		log(events.NewRawError(err))
		return nil
	}

	svc, err := run.Session().Engine().Services().Webhook(run.Session().Assets())
	if err != nil {
		log(events.NewRawError(err))
		return nil
	}

	trace, err := svc.Call(req)
	if err != nil {
		logCallError(err, log)
	}

	if trace != nil {
		call := flows.NewWebhookCall(trace)
		status := callStatus(trace, err, false)

		log(events.NewWebhookCalled(trace, status, ""))

		if a.ResultName != "" {
			a.saveLegacyWebhookResult(run, step, a.ResultName, call, status, log)
		}

		return call
	}

	return nil
}

func (a *CallWebhook) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	if a.ResultName != "" {
		result(flows.NewResultInfo(a.ResultName, webhookCategories))
	}
}

// logs an error from the webhook service, using a dedicated code where we have one
func logCallError(err error, log events.EventLogger) {
	if errors.Is(err, httpx.ErrResponseSize) {
		log(events.NewError(err.Error(), events.ErrorCodeWebhookResponseSize))
	} else {
		log(events.NewRawError(err))
	}
}

// determines the webhook status from the HTTP status code
func callStatus(t *httpx.Trace, err error, isResthook bool) core.CallStatus {
	if t.Response == nil || err != nil {
		return core.CallStatusConnectionError
	}
	if isResthook && t.Response.StatusCode == http.StatusGone {
		// https://zapier.com/developer/documentation/v2/rest-hooks/
		return core.CallStatusSubscriberGone
	}
	if t.Response.StatusCode/100 == 2 {
		return core.CallStatusSuccess
	}
	return core.CallStatusResponseError
}
