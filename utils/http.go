package utils

import (
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
	RRSuccess RequestResponseStatus = "S"

	// RRConnectionFailure represents that the webhook had a connection failure
	RRConnectionFailure RequestResponseStatus = "C"

	// RRStatusFailure represents that the webhook had a non 2xx status code
	RRStatusFailure RequestResponseStatus = "F"
)

func (r RequestResponseStatus) String() string {
	return string(r)
}

// RequestResponse represents both the outgoing request and response for a particular URL/method/body
type RequestResponse interface {
	URL() string
	Status() RequestResponseStatus
	StatusCode() int
	Request() string
	Response() string
	Body() string
	JSON() JSONFragment
}

// MakeHTTPRequest fires the passed in http request, returning any errors encountered. RequestResponse is always set
// regardless of any errors being set
func MakeHTTPRequest(req *http.Request) (RequestResponse, error) {
	resp, err := getClient().Do(req)
	if err != nil {
		rr, _ := newRRFromRequestAndError(req, err)
		return rr, err
	}
	defer resp.Body.Close()

	rr, err := newRRFromResponse(resp)
	return rr, err
}

type requestResponse struct {
	url        string
	status     RequestResponseStatus
	statusCode int
	request    string
	response   string
	body       string
}

func (r *requestResponse) URL() string                   { return r.url }
func (r *requestResponse) Status() RequestResponseStatus { return r.status }
func (r *requestResponse) StatusCode() int               { return r.statusCode }
func (r *requestResponse) Request() string               { return r.request }
func (r *requestResponse) Response() string              { return r.response }
func (r *requestResponse) Body() string                  { return r.body }
func (r *requestResponse) JSON() JSONFragment            { return JSONFragment(r.body) }

func (r *requestResponse) Default() interface{} {
	return r.Body()
}

func (r *requestResponse) Resolve(key string) interface{} {
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

	return fmt.Errorf("No field '%s' on webhook", key)
}

// newRRFromResponse creates a new RequestResponse based on the passed in http request and error (when we received no response)
func newRRFromRequestAndError(r *http.Request, requestError error) (RequestResponse, error) {
	rr := requestResponse{}
	rr.url = r.URL.String()

	request, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return &rr, err
	}
	rr.request = string(request)
	rr.status = RRConnectionFailure
	rr.body = requestError.Error()

	return &rr, nil
}

// newRRFromResponse creates a new RequestResponse based on the passed in http Response
func newRRFromResponse(r *http.Response) (RequestResponse, error) {
	var err error
	rr := requestResponse{}
	rr.url = r.Request.URL.String()
	rr.statusCode = r.StatusCode

	request, err := httputil.DumpRequestOut(r.Request, true)
	if err != nil {
		return &rr, err
	}
	rr.request = string(request)

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

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type rrEnvelope struct {
	URL        string                `json:"url"`
	Status     RequestResponseStatus `json:"status"`
	StatusCode int                   `json:"status_code"`
	Request    string                `json:"request"`
	Response   string                `json:"response"`
}

func (r *requestResponse) UnmarshalJSON(data []byte) error {
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

	return nil
}

func (r *requestResponse) MarshalJSON() ([]byte, error) {
	var re rrEnvelope

	re.URL = r.url
	re.Status = r.status
	re.StatusCode = r.statusCode
	re.Request = r.request
	re.Response = r.response

	return json.Marshal(re)
}
