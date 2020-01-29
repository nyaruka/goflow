package webhooks

import (
	"encoding/json"
	"net/http"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/utils/httpx"
)

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

	trace, err := httpx.DoTrace(s.httpClient, request, s.httpRetries, s.maxBodyBytes)
	if trace != nil {
		call := &flows.WebhookCall{Trace: trace}

		// for webhook calls, we're only interested in valid JSON response bodies
		if len(trace.ResponseBody) > 0 && !json.Valid(trace.ResponseBody) {
			call.ResponseBody = nil
			call.BodyIgnored = true
		}

		// throw away any error that happened prior to getting a response.. these will be surfaced to the user
		// as connection_error status on the response
		if trace.Response == nil {
			return call, nil
		}

		return call, err
	}

	return nil, err
}

var _ flows.WebhookService = (*service)(nil)
