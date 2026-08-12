package webhooks

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
)

type service struct {
	httpClient        *http.Client
	defaultHeaders    map[string]string
	restrictedDomains []string
	maxResponseBytes  int
}

// NewServiceFactory creates a new webhook service factory. The engine supplies the HTTP client and the maximum
// response size; the client's transport can be configured with tracing, mocking or access control as needed
// (see github.com/nyaruka/gocommon/httpx). The restricted domains are domains (e.g. messaging provider APIs)
// which flows shouldn't be calling directly.
func NewServiceFactory(defaultHeaders map[string]string, restrictedDomains []string) engine.WebhookServiceFactory {
	return func(eng flows.Engine, sa flows.SessionAssets) (flows.WebhookService, error) {
		return NewService(eng.HTTPClient(), defaultHeaders, restrictedDomains, eng.Options().MaxResponseBytes), nil
	}
}

// NewService creates a new default webhook service
func NewService(httpClient *http.Client, defaultHeaders map[string]string, restrictedDomains []string, maxResponseBytes int) flows.WebhookService {
	// build the client this service will call through, layering our concerns onto the transport we were given. The
	// read limit bounds how much we'll read from an untrusted endpoint and goes inside tracing so it applies before
	// the body is buffered into the trace; tracing is outermost so that a request denied by access control is still
	// captured as a trace. Tracing itself accumulates nothing, so this is done once here rather than per call - each
	// call gets its own traces from the collector it puts on the request context.
	inner := httpClient.Transport
	if maxResponseBytes > 0 {
		inner = httpx.WithReadLimit(inner, maxResponseBytes)
	}
	traced := *httpClient
	traced.Transport = httpx.WithTraces(inner)

	return &service{
		httpClient:        &traced,
		defaultHeaders:    defaultHeaders,
		restrictedDomains: restrictedDomains,
		maxResponseBytes:  maxResponseBytes,
	}
}

// IsRestricted returns whether the host of the given URL matches or is a subdomain of one of our
// configured restricted domains.
func (s *service) IsRestricted(u *url.URL) bool {
	host := strings.ToLower(u.Hostname())
	for _, domain := range s.restrictedDomains {
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return true
		}
	}
	return false
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

	// collect the traces of just this call - the tracing transport on our client is shared by every call through it,
	// and records into whichever collector the request's context carries
	ctx, traces := httpx.WithTraceCollector(request.Context())

	resp, err := s.httpClient.Do(request.WithContext(ctx))

	// tracing has already buffered the body into the trace; draining the handed-back body surfaces ErrResponseSize
	// if the response exceeded our limit
	var sizeErr error
	if resp != nil {
		if s.maxResponseBytes > 0 {
			_, sizeErr = io.ReadAll(resp.Body)
		}
		resp.Body.Close()
	}

	trace := traces.Last()
	if trace == nil {
		// the request couldn't even be dumped, so we have no trace to return
		return nil, err
	}

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
