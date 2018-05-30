package flows

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// response content-types that we'll save as @run.webhook.body
var saveResponseContentTypes = map[string]bool{
	"application/json":       true,
	"application/javascript": true,
	"application/xml":        true,
	"text/html":              true,
	"text/plain":             true,
	"text/xml":               true,
}

// WebhookStatus represents the status of a WebhookRequest
type WebhookStatus string

const (
	// WebhookStatusSuccess represents that the webhook was successful
	WebhookStatusSuccess WebhookStatus = "success"

	// WebhookStatusConnectionError represents that the webhook had a connection error
	WebhookStatusConnectionError WebhookStatus = "connection_error"

	// WebhookStatusResponseError represents that the webhook response had a non 2xx status code
	WebhookStatusResponseError WebhookStatus = "response_error"
)

func (r WebhookStatus) String() string {
	return string(r)
}

// WebhookCall describes a call made to an external service. It has several properties which can be accessed in expressions:
//
//  * `status` the status of the webhook - one of "success", "connection_error" or "response_error"
//  * `status_code` the status code of the response
//  * `body` the body of the response
//  * `json` the parsed JSON response (if response body was JSON)
//  * `json.[key]` sub-elements of the parsed JSON response
//  * `request` the raw request made, including headers
//  * `response` the raw response received, including headers
//
// Examples:
//
//   @run.webhook.status_code -> 200
//   @run.webhook.json.results.0.state -> WA
//
// @context webhook
type WebhookCall struct {
	url        string
	status     WebhookStatus
	statusCode int
	request    string
	response   string
}

// MakeWebhookCall fires the passed in http request, returning any errors encountered. RequestResponse is always set
// regardless of any errors being set
func MakeWebhookCall(session Session, request *http.Request) (*WebhookCall, error) {
	var response *http.Response
	var requestDump string
	var err error

	// if our config has mocks, look for a matching one
	mock := findMockedRequest(session, request)
	if mock != nil {
		response, requestDump, err = session.HTTPClient().MockWithDump(request, mock.Status, mock.Body)
	} else {
		if session.EngineConfig().DisableWebhooks() {
			response, requestDump, err = session.HTTPClient().MockWithDump(request, 200, "DISABLED")
		} else {
			response, requestDump, err = session.HTTPClient().DoWithDump(request)
		}
	}

	if err != nil {
		return newWebhookCallFromError(request, requestDump, err), err
	}

	return newWebhookCallFromResponse(requestDump, response, session.EngineConfig().MaxWebhookResponseBytes())
}

// URL returns the full URL
func (w *WebhookCall) URL() string { return w.url }

// Method returns the full HTTP method
func (w *WebhookCall) Method() string { return w.request[:strings.IndexRune(w.request, ' ')] }

// Status returns the response status message
func (w *WebhookCall) Status() WebhookStatus { return w.status }

// StatusCode returns the response status code
func (w *WebhookCall) StatusCode() int { return w.statusCode }

// Request returns the request trace
func (w *WebhookCall) Request() string { return w.request }

// Response returns the response trace
func (w *WebhookCall) Response() string { return w.response }

// Body returns the response body
func (w *WebhookCall) Body() string {
	parts := strings.SplitN(w.response, "\r\n\r\n", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// JSON returns the response as a JSON fragment
func (w *WebhookCall) JSON() types.XValue { return types.JSONToXValue([]byte(w.Body())) }

// Resolve resolves the given key when this webhook is referenced in an expression
func (w *WebhookCall) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "url":
		return types.NewXText(w.URL())
	case "request":
		return types.NewXText(w.Request())
	case "response":
		return types.NewXText(w.Response())
	case "status":
		return types.NewXText(string(w.Status()))
	case "status_code":
		return types.NewXNumberFromInt(w.StatusCode())
	case "json":
		return w.JSON()
	}

	return types.NewXResolveError(w, key)
}

// Describe returns a representation of this type for error messages
func (w *WebhookCall) Describe() string { return "webhook" }

// Reduce reduces this to a string of method and URL, e.g. "GET http://example.com/hook.php"
func (w *WebhookCall) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(fmt.Sprintf("%s %s", w.Method(), w.URL()))
}

// ToXJSON is called when this type is passed to @(json(...))
func (w *WebhookCall) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, w, "body", "json", "url", "request", "response", "status", "status_code").ToXJSON(env)
}

var _ types.XValue = (*WebhookCall)(nil)
var _ types.XResolvable = (*WebhookCall)(nil)

// newWebhookCallFromError creates a new webhook call based on the passed in http request and error (when we received no response)
func newWebhookCallFromError(request *http.Request, requestTrace string, requestError error) *WebhookCall {
	return &WebhookCall{
		url:        request.URL.String(),
		status:     WebhookStatusConnectionError,
		statusCode: 0,
		request:    requestTrace,
		response:   requestError.Error(),
	}
}

// newWebhookCallFromResponse creates a new RequestResponse based on the passed in http Response
func newWebhookCallFromResponse(requestTrace string, response *http.Response, maxBodyBytes int) (*WebhookCall, error) {
	defer response.Body.Close()

	w := &WebhookCall{
		url:        response.Request.URL.String(),
		statusCode: response.StatusCode,
		request:    requestTrace,
	}

	// set our status based on our status code
	if w.statusCode/100 == 2 {
		w.status = WebhookStatusSuccess
	} else {
		w.status = WebhookStatusResponseError
	}

	// save response dump without body which will be parsed separately
	responseDump, err := httputil.DumpResponse(response, false)
	if err != nil {
		return nil, err
	}
	w.response = string(responseDump)

	// only save response body's if we have a supported content-type
	contentType := response.Header.Get("Content-Type")
	mediaType, _, _ := mime.ParseMediaType(contentType)
	saveBody := saveResponseContentTypes[mediaType]

	if saveBody {
		// only read up to our max body bytes limit
		bodyReader := io.LimitReader(response.Body, int64(maxBodyBytes)+1)

		bodyBytes, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			return nil, err
		}

		// if we have no remaining bytes, error because the body was too big
		if bodyReader.(*io.LimitedReader).N <= 0 {
			return nil, fmt.Errorf("webhook response body exceeds %d bytes limit", maxBodyBytes)
		}

		w.response += string(bodyBytes)
	} else {
		// no body for non-text responses but add it to our Response log so users know why
		w.response += "Non-text body, ignoring"
	}

	return w, nil
}

//------------------------------------------------------------------------------------------
// Request Mocking
//------------------------------------------------------------------------------------------

type WebhookMock struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func findMockedRequest(session Session, request *http.Request) *WebhookMock {
	for _, mock := range session.EngineConfig().WebhookMocks() {
		if strings.EqualFold(mock.Method, request.Method) && strings.EqualFold(mock.URL, request.URL.String()) {
			return mock
		}
	}
	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type webhookCallEnvelope struct {
	URL        string        `json:"url"`
	Status     WebhookStatus `json:"status"`
	StatusCode int           `json:"status_code"`
	Request    string        `json:"request"`
	Response   string        `json:"response"`
}

// UnmarshalJSON unmarshals a request response from the given JSON
func (w *WebhookCall) UnmarshalJSON(data []byte) error {
	var envelope webhookCallEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return err
	}

	w.url = envelope.URL
	w.status = envelope.Status
	w.statusCode = envelope.StatusCode
	w.request = envelope.Request
	w.response = envelope.Response
	return nil
}

// MarshalJSON marshals this request response into JSON
func (r *WebhookCall) MarshalJSON() ([]byte, error) {
	return json.Marshal(&webhookCallEnvelope{
		URL:        r.url,
		Status:     r.status,
		StatusCode: r.statusCode,
		Request:    r.request,
		Response:   r.response,
	})
}
