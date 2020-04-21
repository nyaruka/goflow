package httpx

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nyaruka/goflow/utils/dates"

	"github.com/pkg/errors"
)

var debug = false

// Do makes the given HTTP request using the current requestor and retry config
func Do(client *http.Client, request *http.Request, retries *RetryConfig, access *AccessConfig) (*http.Response, error) {
	if access != nil {
		allowed, err := access.Allow(request)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, errors.Errorf("request to %s denied", request.URL.Hostname())
		}
	}

	var response *http.Response
	var err error
	retry := 0

	for {
		response, err = currentRequestor.Do(client, request)

		if retries != nil && retry < retries.MaxRetries() {
			backoff := retries.Backoff(retry)

			if retries.ShouldRetry(request, response, backoff) {
				time.Sleep(backoff)
				retry++
				continue
			}
		}

		break
	}

	return response, err
}

// Trace holds the complete trace of an HTTP request/response
type Trace struct {
	Request       *http.Request
	RequestTrace  []byte
	Response      *http.Response
	ResponseTrace []byte
	ResponseBody  []byte // response body stored separately
	StartTime     time.Time
	EndTime       time.Time
}

func (t *Trace) String() string {
	b := &strings.Builder{}
	b.WriteString(fmt.Sprintf(">>>>>>>> %s %s\n", t.Request.Method, t.Request.URL))
	b.WriteString(string(t.RequestTrace))
	b.WriteString("\n<<<<<<<<\n")
	b.WriteString(string(t.ResponseTrace))
	b.WriteString(string(t.ResponseBody))
	return b.String()
}

// ResponseTraceUTF8 returns a valid UTF-8 string version of trace, substituting the body with placeholder if it isn't valid UTF-8
func (t *Trace) ResponseTraceUTF8(placeholder string) string {
	// headers part assumed to be valid UTF-8
	s := string(t.ResponseTrace)

	// if body is valid UTF-8, include it
	if utf8.Valid(t.ResponseBody) {
		s += string(t.ResponseBody)
	} else {
		s += placeholder
	}

	return s
}

// DoTrace makes the given request saving traces of the complete request and response
func DoTrace(client *http.Client, request *http.Request, retries *RetryConfig, access *AccessConfig, maxBodyBytes int) (*Trace, error) {
	requestTrace, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, err
	}

	trace := &Trace{
		Request:      request,
		RequestTrace: requestTrace,
		StartTime:    dates.Now(),
	}

	response, err := Do(client, request, retries, access)
	trace.EndTime = dates.Now()

	if err != nil {
		return trace, err
	}

	trace.Response = response

	// save response trace without body which will be parsed separately
	responseTrace, err := httputil.DumpResponse(response, false)
	if err != nil {
		return trace, err
	}
	trace.ResponseTrace = responseTrace

	responseBody, err := readBody(response, maxBodyBytes)
	if err != nil {
		return trace, err
	}
	trace.ResponseBody = responseBody

	if debug {
		fmt.Println(trace.String())
	}

	return trace, nil
}

// NewRequest is a convenience method to create a request with the given headers
func NewRequest(method string, url string, body io.Reader, headers map[string]string) (*http.Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		r.Header.Set(key, value)
	}

	return r, nil
}

// attempts to read the body of an HTTP response
func readBody(response *http.Response, maxBodyBytes int) ([]byte, error) {
	defer response.Body.Close()

	if maxBodyBytes > 0 {
		// we will only read up to our max body bytes limit
		bodyReader := io.LimitReader(response.Body, int64(maxBodyBytes)+1)

		bodyBytes, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			return nil, err
		}

		// if we have no remaining bytes, error because the body was too big
		if bodyReader.(*io.LimitedReader).N <= 0 {
			return nil, errors.Errorf("webhook response body exceeds %d bytes limit", maxBodyBytes)
		}

		return bodyBytes, nil
	}

	// if there is no limit, read the entire body
	return ioutil.ReadAll(response.Body)
}

// Requestor is anything that can make an HTTP request with a client
type Requestor interface {
	Do(*http.Client, *http.Request) (*http.Response, error)
}

type defaultRequestor struct{}

func (r defaultRequestor) Do(client *http.Client, request *http.Request) (*http.Response, error) {
	return client.Do(request)
}

// DefaultRequestor is the default HTTP requestor
var DefaultRequestor Requestor = defaultRequestor{}
var currentRequestor = DefaultRequestor

// SetRequestor sets the requestor used by Request
func SetRequestor(requestor Requestor) {
	currentRequestor = requestor
}

// SetDebug enables debugging
func SetDebug(enabled bool) {
	debug = enabled
}
