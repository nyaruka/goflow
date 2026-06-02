package webhooks

import (
	"net/http"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
)

type service struct {
	httpClient     *http.Client
	defaultHeaders map[string]string
	maxBodyBytes   int
}

// NewServiceFactory creates a new webhook service factory. The engine supplies the HTTP client; its transport can be
// configured with tracing, mocking or access control as needed (see github.com/nyaruka/gocommon/httpx).
func NewServiceFactory(defaultHeaders map[string]string, maxBodyBytes int) engine.WebhookServiceFactory {
	return func(httpClient *http.Client, sa flows.SessionAssets) (flows.WebhookService, error) {
		return NewService(httpClient, defaultHeaders, maxBodyBytes), nil
	}
}

// NewService creates a new default webhook service
func NewService(httpClient *http.Client, defaultHeaders map[string]string, maxBodyBytes int) flows.WebhookService {
	return &service{
		httpClient:     httpClient,
		defaultHeaders: defaultHeaders,
		maxBodyBytes:   maxBodyBytes,
	}
}

func (s *service) Call(request *http.Request) (*httpx.Trace, error) {
	// set any headers with defaults
	for k, v := range s.defaultHeaders {
		if request.Header.Get(k) == "" {
			request.Header.Set(k, v)
		}
	}

	// If user has explicitly set Accept-Encoding: gzip, remove it as Transport will add this itself,
	// and it only does automatic decompression if its the one to set it.
	if request.Header.Get("Accept-Encoding") == "gzip" {
		request.Header.Del("Accept-Encoding")
	}

	// wrap the configured transport with a per-call tracer so we capture this request's trace. Tracing is the
	// outermost wrapper so that a request denied by access control is still captured as a trace.
	captureBytes := -1 // a non-positive limit means capture the entire body
	if s.maxBodyBytes > 0 {
		captureBytes = s.maxBodyBytes + 1 // capture one byte beyond our limit so we can detect responses that exceed it
	}
	tracing := httpx.WithTracing(s.httpClient.Transport, captureBytes)
	client := *s.httpClient
	client.Transport = tracing

	_, err := client.Do(request)

	traces := tracing.Traces()
	if len(traces) == 0 {
		// the request couldn't even be dumped, so we have no trace to return
		return nil, err
	}
	trace := traces[len(traces)-1]

	if s.maxBodyBytes > 0 && len(trace.ResponseBody) > s.maxBodyBytes {
		// response body exceeded our limit
		trace.ResponseBody = nil
		return trace, httpx.ErrResponseSize
	}

	if err != nil && trace.Response == nil {
		// throw away any error that happened prior to getting a response.. these will be surfaced to the user
		// as connection_error status on the response
		return trace, nil
	}

	return trace, err
}

var _ flows.WebhookService = (*service)(nil)
