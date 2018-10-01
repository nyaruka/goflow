package flows

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"
)

var DefaultWebhookPayload = `{
	"contact": {"uuid": "@contact.uuid", "name": @(json(contact.name)), "urn": @(json(if(default(run.input.urn, default(contact.urns.0, null)), text(default(run.input.urn, default(contact.urns.0, null))), null)))},
	"flow": @(json(run.flow)),
	"path": @(json(run.path)),
	"results": @(json(run.results)),
	"run": {"uuid": "@run.uuid", "created_on": "@run.created_on"},
	"input": @(json(run.input)),
	"channel": @(json(if(run.input, run.input.channel, null)))
}`

// response content-types that we'll fetch
var fetchResponseContentTypes = map[string]bool{
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

func WebhookStatusFromCode(code int) WebhookStatus {
	if code/100 == 2 {
		return WebhookStatusSuccess
	}
	return WebhookStatusResponseError
}

func (r WebhookStatus) String() string {
	return string(r)
}

// WebhookCall is a call made to an external service
type WebhookCall struct {
	url           string
	request       *http.Request
	response      *http.Response
	status        WebhookStatus
	timeTaken     time.Duration
	requestTrace  string
	responseTrace string
}

// MakeWebhookCall fires the passed in http request, returning any errors encountered. RequestResponse is always set
// regardless of any errors being set
func MakeWebhookCall(session Session, request *http.Request) (*WebhookCall, error) {
	var response *http.Response
	var requestDump string
	var err error
	var timeTaken time.Duration

	// if our config has mocks, look for a matching one
	mock := findMockedRequest(session, request)
	if mock != nil {
		response, requestDump, err = session.HTTPClient().MockWithDump(request, mock.Status, mock.Body)
	} else {
		if session.EngineConfig().DisableWebhooks() {
			response, requestDump, err = session.HTTPClient().MockWithDump(request, 200, "DISABLED")
		} else {
			start := utils.Now()
			response, requestDump, err = session.HTTPClient().DoWithDump(request)
			timeTaken = utils.Now().Sub(start)
		}
	}

	if err != nil {
		return newWebhookCallFromError(request, requestDump, err), err
	}

	return newWebhookCallFromResponse(requestDump, response, session.EngineConfig().MaxWebhookResponseBytes(), timeTaken)
}

// URL returns the full URL
func (w *WebhookCall) URL() string { return w.url }

// Method returns the full HTTP method
func (w *WebhookCall) Method() string { return w.request.Method }

// Status returns the response status message
func (w *WebhookCall) Status() WebhookStatus { return w.status }

// StatusCode returns the response status code
func (w *WebhookCall) StatusCode() int {
	if w.response != nil {
		return w.response.StatusCode
	}
	return 0
}

// TimeTaken returns the time taken to make the request
func (w *WebhookCall) TimeTaken() time.Duration { return w.timeTaken }

// Request returns the request trace
func (w *WebhookCall) Request() string { return w.requestTrace }

// Response returns the response trace
func (w *WebhookCall) Response() string { return w.responseTrace }

// Body returns the response body
func (w *WebhookCall) Body() string {
	parts := strings.SplitN(w.responseTrace, "\r\n\r\n", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// newWebhookCallFromError creates a new webhook call based on the passed in http request and error (when we received no response)
func newWebhookCallFromError(request *http.Request, requestTrace string, requestError error) *WebhookCall {
	return &WebhookCall{
		url:           request.URL.String(),
		request:       request,
		response:      nil,
		status:        WebhookStatusConnectionError,
		requestTrace:  requestTrace,
		responseTrace: requestError.Error(),
	}
}

// newWebhookCallFromResponse creates a new RequestResponse based on the passed in http Response
func newWebhookCallFromResponse(requestTrace string, response *http.Response, maxBodyBytes int, timeTaken time.Duration) (*WebhookCall, error) {
	defer response.Body.Close()

	// save response trace without body which will be parsed separately
	responseTrace, err := httputil.DumpResponse(response, false)
	if err != nil {
		return nil, err
	}

	w := &WebhookCall{
		url:           response.Request.URL.String(),
		request:       response.Request,
		response:      response,
		status:        WebhookStatusFromCode(response.StatusCode),
		requestTrace:  requestTrace,
		responseTrace: string(responseTrace),
		timeTaken:     timeTaken,
	}

	// only save response body's if we have a supported content-type
	contentType := response.Header.Get("Content-Type")
	mediaType, _, _ := mime.ParseMediaType(contentType)
	saveBody := fetchResponseContentTypes[mediaType]

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

		w.responseTrace += string(bodyBytes)
	} else {
		// no body for non-text responses but add it to our Response log so users know why
		w.responseTrace += "Non-text body, ignoring"
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
