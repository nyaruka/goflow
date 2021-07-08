package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"

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
			// strip out any invalid UTF-8 and replace any escaped nulls
			cleaned := replaceEscapedNulls(bytes.ToValidUTF8(call.ResponseBody, nil))

			if json.Valid(cleaned) {
				call.ResponseJSON = cleaned
			}
		}

		return call, err
	}

	return nil, err
}

var _ flows.WebhookService = (*service)(nil)

// replaces any `\u0000` sequences with the unicode replacement char `\ufffd`
// a sequence such as `\\u0000` is preserved as it is an escaped slash followed by the sequence `u0000`
func replaceEscapedNulls(data []byte) []byte {
	return nullEscapeRegex.ReplaceAllFunc(data, func(m []byte) []byte {
		slashes := bytes.Count(m, []byte(`\`))
		if slashes%2 == 0 {
			return m
		}

		return append(bytes.Repeat([]byte(`\`), slashes-1), replacementChar...)
	})
}

var nullEscapeRegex = regexp.MustCompile(`\\+u0{4}`)
var replacementChar = []byte(`\ufffd`)
