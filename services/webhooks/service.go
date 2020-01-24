package webhooks

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/pkg/errors"
)

// response content-types that we'll fetch
var fetchResponseContentTypes = map[string]bool{
	"application/json":       true,
	"application/javascript": true,
	"application/xml":        true,
	"text/html":              true,
	"text/plain":             true,
	"text/xml":               true,
	"text/javascript":        true,
}

type service struct {
	httpClient     *http.Client
	httpRetries    *httpx.RetryConfig
	defaultHeaders map[string]string
	maxBodyBytes   int
}

// NewServiceFactory creates a new webhook service factory
func NewServiceFactory(httpClient *http.Client, httpRetries *httpx.RetryConfig, defaultHeaders map[string]string, maxBodyBytes int) engine.WebhookServiceFactory {
	return func(flows.Session) (flows.WebhookService, error) {
		return NewService(httpClient, httpRetries, defaultHeaders, maxBodyBytes), nil
	}
}

// NewService creates a new default webhook service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, defaultHeaders map[string]string, maxBodyBytes int) flows.WebhookService {
	return &service{
		httpClient:     httpClient,
		httpRetries:    httpRetries,
		defaultHeaders: defaultHeaders,
		maxBodyBytes:   maxBodyBytes,
	}
}

func (s *service) Call(session flows.Session, request *http.Request) (*flows.WebhookCall, error) {
	// set any headers with defaults
	for k, v := range s.defaultHeaders {
		if request.Header.Get(k) == "" {
			request.Header.Set(k, v)
		}
	}

	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, err
	}

	start := dates.Now()
	response, err := httpx.Do(s.httpClient, request, s.httpRetries)
	timeTaken := dates.Now().Sub(start)

	if err != nil {
		return &flows.WebhookCall{
			URL:        request.URL.String(),
			Method:     request.Method,
			StatusCode: 0,
			Request:    dump,
			Response:   nil,
		}, nil
	}

	return s.newCallFromResponse(dump, response, s.maxBodyBytes, timeTaken)
}

// creates a new call based on the passed in http response
func (s *service) newCallFromResponse(requestTrace []byte, response *http.Response, maxBodyBytes int, timeTaken time.Duration) (*flows.WebhookCall, error) {
	defer response.Body.Close()

	// save response trace without body which will be parsed separately
	responseTrace, err := httputil.DumpResponse(response, false)
	if err != nil {
		return nil, err
	}

	w := &flows.WebhookCall{
		URL:        response.Request.URL.String(),
		Method:     response.Request.Method,
		StatusCode: response.StatusCode,
		Request:    requestTrace,
		Response:   responseTrace,
		TimeTaken:  timeTaken,
	}

	body, err := readBody(response, maxBodyBytes)
	if err != nil {
		return nil, err
	}

	w.Response = append(w.Response, body...)

	if json.Valid(body) {
		w.ResponseJSON = body
	} else {
		w.BodyIgnored = true
	}

	return w, nil
}

// attempts to read the body of an HTTP response
func readBody(response *http.Response, maxBodyBytes int) ([]byte, error) {
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

var _ flows.WebhookService = (*service)(nil)
