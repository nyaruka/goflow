package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/utils"
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
			call.ResponseJSON, call.ResponseCleaned = ExtractJSON(call.ResponseBody)
		}

		return call, err
	}

	return nil, err
}

func ExtractJSON(body []byte) ([]byte, bool) {
	// we make a best effort to turn the body into JSON, so we strip out:
	//  1. any invalid UTF-8 sequences
	//  2. null chars
	//  3. escaped null chars (\u0000)
	cleaned := bytes.ToValidUTF8(body, nil)
	cleaned = bytes.ReplaceAll(cleaned, []byte{0}, nil)
	cleaned = utils.ReplaceEscapedNulls(cleaned, nil)

	if json.Valid(cleaned) {
		changed := !bytes.Equal(body, cleaned)
		return cleaned, changed
	}
	return nil, false
}

var _ flows.WebhookService = (*service)(nil)
