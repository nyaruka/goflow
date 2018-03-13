package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

// RequestResponseStatus represents the status of a WebhookRequeset
type RequestResponseStatus string

const (
	// RRSuccess represents that the webhook was successful
	RRSuccess RequestResponseStatus = "success"

	// RRConnectionError represents that the webhook had a connection error
	RRConnectionError RequestResponseStatus = "connection_error"

	// RRResponseError represents that the webhook response had a non 2xx status code
	RRResponseError RequestResponseStatus = "response_error"
)

func init() {
	Validator.RegisterAlias("http_method", "eq=GET|eq=HEAD|eq=POST|eq=PUT|eq=PATCH|eq=DELETE")
}

func (r RequestResponseStatus) String() string {
	return string(r)
}

// RequestResponse represents both the outgoing request and response for a particular URL/method/body
type RequestResponse struct {
	url        string
	method     string
	status     RequestResponseStatus
	statusCode int
	request    string
	response   string
	body       string
}

// MakeHTTPRequest fires the passed in http request, returning any errors encountered. RequestResponse is always set
// regardless of any errors being set
func MakeHTTPRequest(req *http.Request) (*RequestResponse, error) {
	requestTrace, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		rr, _ := newRRFromRequestAndError(req, string(requestTrace), err)
		return rr, err
	}

	resp, err := getClient().Do(req)
	if err != nil {
		rr, _ := newRRFromRequestAndError(req, string(requestTrace), err)
		return rr, err
	}
	defer resp.Body.Close()

	rr, err := newRRFromResponse(string(requestTrace), resp)
	return rr, err
}

// URL returns the full URL
func (r *RequestResponse) URL() string { return r.url }

// Status returns the response status message
func (r *RequestResponse) Status() RequestResponseStatus { return r.status }

// StatusCode returns the response status code
func (r *RequestResponse) StatusCode() int { return r.statusCode }

// Request returns the request trace
func (r *RequestResponse) Request() string { return r.request }

// Response returns the response trace
func (r *RequestResponse) Response() string { return r.response }

// Body returns the response body
func (r *RequestResponse) Body() string { return r.body }

// JSON returns the response as a JSON fragment
func (r *RequestResponse) JSON() JSONFragment { return JSONFragment([]byte(r.body)) }

// Resolve resolves the given key when this webhook is referenced in an expression
func (r *RequestResponse) Resolve(key string) interface{} {
	switch key {
	case "body":
		return r.Body()
	case "json":
		return r.JSON()
	case "url":
		return r.URL()
	case "request":
		return r.Request()
	case "response":
		return r.Response()
	case "status":
		return r.Status()
	case "status_code":
		return r.StatusCode()
	}

	return fmt.Errorf("no field '%s' on webhook", key)
}

// Default returns the value of this webhook when it is the result of an expression
func (r *RequestResponse) Default() interface{} {
	return r
}

func (r *RequestResponse) String() string {
	return r.body
}

var _ VariableResolver = (*RequestResponse)(nil)

// newRRFromResponse creates a new RequestResponse based on the passed in http request and error (when we received no response)
func newRRFromRequestAndError(r *http.Request, requestTrace string, requestError error) (*RequestResponse, error) {
	rr := RequestResponse{}
	rr.url = r.URL.String()

	rr.request = requestTrace
	rr.status = RRConnectionError
	rr.body = requestError.Error()

	return &rr, nil
}

// newRRFromResponse creates a new RequestResponse based on the passed in http Response
func newRRFromResponse(requestTrace string, r *http.Response) (*RequestResponse, error) {
	var err error
	rr := RequestResponse{}
	rr.url = r.Request.URL.String()
	rr.statusCode = r.StatusCode

	// set our status based on our status code
	if rr.statusCode/100 == 2 {
		rr.status = RRSuccess
	} else {
		rr.status = RRResponseError
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

var (
	transport *http.Transport
	client    *http.Client
	once      sync.Once
)

func getClient() *http.Client {
	once.Do(func() {
		timeout := time.Duration(15 * time.Second)
		transport = &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		}
		client = &http.Client{Transport: transport, Timeout: timeout}
	})

	return client
}

func getInsecureClient() *http.Client {
	once.Do(func() {
		timeout := time.Duration(15 * time.Second)
		transport = &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: transport, Timeout: timeout}
	})

	return client
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type rrEnvelope struct {
	URL        string                `json:"url"`
	Status     RequestResponseStatus `json:"status"`
	StatusCode int                   `json:"status_code"`
	Body       string                `json:"body"`
	Request    string                `json:"request"`
	Response   string                `json:"response"`
}

// UnmarshalJSON unmarshals a request response from the given JSON
func (r *RequestResponse) UnmarshalJSON(data []byte) error {
	var envelope rrEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return err
	}

	r.url = envelope.URL
	r.status = envelope.Status
	r.statusCode = envelope.StatusCode
	r.request = envelope.Request
	r.response = envelope.Response
	r.body = envelope.Body

	return nil
}

// MarshalJSON marshals this request reponse into JSON
func (r *RequestResponse) MarshalJSON() ([]byte, error) {
	var re rrEnvelope

	re.URL = r.url
	re.Status = r.status
	re.StatusCode = r.statusCode
	re.Request = r.request
	re.Response = r.response
	re.Body = r.body

	return json.Marshal(re)
}
