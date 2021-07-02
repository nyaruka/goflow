package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
)

type service struct {
	httpClient     *http.Client
	httpRetries    *httpx.RetryConfig
	httpAccess     *httpx.AccessConfig
	defaultHeaders map[string]string
	maxBodyBytes   int
}

// NewServiceFactory creates a new webhook service factory
func NewServiceFactory(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, defaultHeaders map[string]string, maxBodyBytes int) engine.WebhookServiceFactory {
	return func(flows.Session) (flows.WebhookService, error) {
		return NewService(httpClient, httpRetries, httpAccess, defaultHeaders, maxBodyBytes), nil
	}
}

// NewService creates a new default webhook service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, defaultHeaders map[string]string, maxBodyBytes int) flows.WebhookService {
	return &service{
		httpClient:     httpClient,
		httpRetries:    httpRetries,
		httpAccess:     httpAccess,
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

	// If user has explicitly set Accept-Encoding: gzip, remove it as Transport will add this itself,
	// and it only does automatic decompression if its the one to set it.
	if request.Header.Get("Accept-Encoding") == "gzip" {
		request.Header.Del("Accept-Encoding")
	}

	trace, err := httpx.DoTrace(s.httpClient, request, s.httpRetries, s.httpAccess, s.maxBodyBytes)
	if trace != nil {
		call := &flows.WebhookCall{Trace: trace}

		// throw away any error that happened prior to getting a response.. these will be surfaced to the user
		// as connection_error status on the response
		if trace.Response == nil {
			return call, nil
		}

		if len(call.ResponseBody) > 0 {
			// strip out any invalid UTF-8
			bodyUTF8 := bytes.ToValidUTF8(call.ResponseBody, nil)
			if json.Valid(bodyUTF8) {
				call.ResponseJSON = bodyUTF8
			}
		}

		return call, err
	}

	return nil, err
}

var _ flows.WebhookService = (*service)(nil)
