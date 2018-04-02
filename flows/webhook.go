package flows

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

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

// WebhookCall represents both the outgoing request and response for a webhook call
type WebhookCall struct {
	url        string
	method     string
	status     WebhookStatus
	statusCode int
	request    string
	response   string
	body       string
}

// MakeWebhookCall fires the passed in http request, returning any errors encountered. RequestResponse is always set
// regardless of any errors being set
func MakeWebhookCall(req *http.Request) (*WebhookCall, error) {
	requestTrace, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		rr, _ := newWebhookCallFromError(req, string(requestTrace), err)
		return rr, err
	}

	resp, err := utils.NewHTTPClient().Do(req)
	if err != nil {
		w, _ := newWebhookCallFromError(req, string(requestTrace), err)
		return w, err
	}
	defer resp.Body.Close()

	w, err := newWebhookCallFromResponse(string(requestTrace), resp)
	return w, err
}

// URL returns the full URL
func (w *WebhookCall) URL() string { return w.url }

// Status returns the response status message
func (w *WebhookCall) Status() WebhookStatus { return w.status }

// StatusCode returns the response status code
func (w *WebhookCall) StatusCode() int { return w.statusCode }

// Request returns the request trace
func (w *WebhookCall) Request() string { return w.request }

// Response returns the response trace
func (w *WebhookCall) Response() string { return w.response }

// Body returns the response body
func (w *WebhookCall) Body() string { return w.body }

// JSON returns the response as a JSON fragment
func (w *WebhookCall) JSON() utils.JSONFragment { return utils.JSONFragment([]byte(w.body)) }

// Resolve resolves the given key when this webhook is referenced in an expression
func (w *WebhookCall) Resolve(key string) interface{} {
	switch key {
	case "body":
		return w.Body()
	case "json":
		return w.JSON()
	case "url":
		return w.URL()
	case "request":
		return w.Request()
	case "response":
		return w.Response()
	case "status":
		return string(w.Status())
	case "status_code":
		return w.StatusCode()
	}

	return fmt.Errorf("no field '%s' on webhook", key)
}

// Atomize is called when this object needs to be reduced to a primitive
func (w *WebhookCall) Atomize() interface{} {
	return w.body
}

var _ utils.Atomizable = (*WebhookCall)(nil)
var _ utils.Resolvable = (*WebhookCall)(nil)

// newRRFromResponse creates a new RequestResponse based on the passed in http request and error (when we received no response)
func newWebhookCallFromError(r *http.Request, requestTrace string, requestError error) (*WebhookCall, error) {
	rr := WebhookCall{}
	rr.url = r.URL.String()

	rr.request = requestTrace
	rr.status = WebhookStatusConnectionError
	rr.body = requestError.Error()

	return &rr, nil
}

// newWebhookCallFromResponse creates a new RequestResponse based on the passed in http Response
func newWebhookCallFromResponse(requestTrace string, r *http.Response) (*WebhookCall, error) {
	var err error
	rr := WebhookCall{}
	rr.url = r.Request.URL.String()
	rr.statusCode = r.StatusCode

	// set our status based on our status code
	if rr.statusCode/100 == 2 {
		rr.status = WebhookStatusSuccess
	} else {
		rr.status = WebhookStatusResponseError
	}

	rr.request = requestTrace

	// figure out if our Response is something that looks like text from our headers
	isText := false
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "text") ||
		strings.Contains(contentType, "json") ||
		strings.Contains(contentType, "utf") ||
		strings.Contains(contentType, "javascript") ||
		strings.Contains(contentType, "xml") {

		isText = true
	}

	// only dump the whole body if this looks like text
	response, err := httputil.DumpResponse(r, isText)
	if err != nil {
		return &rr, err
	}
	rr.response = string(response)

	if isText {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &rr, err
		}
		rr.body = strings.TrimSpace(string(bodyBytes))
	} else {
		// no body for non-text responses but add it to our Response log so users know why
		rr.response = rr.response + "\nNon-text body, ignoring"
	}

	return &rr, nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type webhookCallEnvelope struct {
	URL        string        `json:"url"`
	Status     WebhookStatus `json:"status"`
	StatusCode int           `json:"status_code"`
	Body       string        `json:"body"`
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
	w.body = envelope.Body
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
		Body:       r.body,
	})
}
