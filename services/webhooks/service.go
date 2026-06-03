package webhooks

import (
	"errors"
	"io"
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

	// bound how many bytes we'll read from an untrusted endpoint, wrapped inside tracing so the limit applies before
	// the body is buffered into the trace
	inner := s.httpClient.Transport
	if s.maxBodyBytes > 0 {
		inner = httpx.WithReadLimit(inner, s.maxBodyBytes)
	}

	// wrap the configured transport with a per-call tracer so we capture this request's trace. Tracing is the
	// outermost wrapper so that a request denied by access control is still captured as a trace.
	tracing := httpx.WithTraces(inner)
	client := *s.httpClient
	client.Transport = tracing

	resp, err := client.Do(request)

	// tracing has already buffered the body into the trace; draining the handed-back body surfaces ErrResponseSize
	// if the response exceeded our limit
	var sizeErr error
	if resp != nil {
		if s.maxBodyBytes > 0 {
			_, sizeErr = io.ReadAll(resp.Body)
		}
		resp.Body.Close()
	}

	traces := tracing.Traces()
	if len(traces) == 0 {
		// the request couldn't even be dumped, so we have no trace to return
		return nil, err
	}
	trace := traces[len(traces)-1]

	if errors.Is(sizeErr, httpx.ErrResponseSize) {
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
