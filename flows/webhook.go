package flows

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// WebhookCall holds the details of a webhook call
type WebhookCall struct {
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	ResponseStatus  int               `json:"status"`
	ResponseHeaders map[string]string `json:"headers"`
	ResponseJSON    json.RawMessage   `json:"json"`
}

// NewWebhookCall creates a new webhook call from a trace
func NewWebhookCall(t *httpx.Trace) *WebhookCall {
	respStatus := 0
	var respHeaders map[string]string

	if t.Response != nil {
		respStatus = t.Response.StatusCode

		respHeaders = make(map[string]string, len(t.Response.Header))
		for k := range t.Response.Header {
			respHeaders[k] = t.Response.Header.Get(k)
		}
	}

	c := &WebhookCall{
		Method:          t.Request.Method,
		URL:             t.Request.URL.String(),
		ResponseStatus:  respStatus,
		ResponseHeaders: respHeaders,
	}

	if len(t.ResponseBody) > 0 {
		c.ResponseJSON = ExtractJSON(t.ResponseBody)
	}

	return c
}

// Context returns the properties available in expressions
//
//	__default__:text -> the method and URL
//	status:number -> the response status code
//	headers:any -> the response headers
//	json:any -> the response body if valid JSON
//
// @context webhook
func (w *WebhookCall) Context(env envs.Environment) map[string]types.XValue {
	headers := types.NewXLazyObject(func() map[string]types.XValue {
		values := make(map[string]types.XValue, len(w.ResponseHeaders))
		for k, v := range w.ResponseHeaders {
			values[k] = types.NewXText(v)
		}
		return values
	})

	json := types.JSONToXValue(w.ResponseJSON)
	if types.IsXError(json) {
		json = nil
	}

	return map[string]types.XValue{
		"__default__": types.NewXText(fmt.Sprintf("%s %s", w.Method, w.URL)),
		"status":      types.NewXNumberFromInt(w.ResponseStatus),
		"headers":     headers,
		"json":        json,
	}
}

func ExtractJSON(body []byte) []byte {
	// we make a best effort to turn the body into JSON, so we strip out:
	//  1. any invalid UTF-8 sequences
	//  2. null chars
	//  3. escaped null chars (\u0000)
	cleaned := bytes.ToValidUTF8(body, nil)
	cleaned = bytes.ReplaceAll(cleaned, []byte{0}, nil)
	cleaned = []byte(httpx.ReplaceEscapedNulls(string(cleaned), ""))

	if json.Valid(cleaned) {
		return cleaned
	}
	return nil
}
