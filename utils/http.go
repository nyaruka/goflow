package utils

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"time"
)

var httpHeaderUserAgent = "User-Agent"

func init() {
	Validator.RegisterAlias("http_method", "eq=GET|eq=HEAD|eq=POST|eq=PUT|eq=PATCH|eq=DELETE")
}

// HTTPClient is a client for HTTP requests
type HTTPClient struct {
	client           *http.Client
	defaultUserAgent string
}

// NewHTTPClient creates a new HTTP client with our default options
func NewHTTPClient(defaultUserAgent string) *HTTPClient {
	// support single tls renegotiation
	tlsConfig := &tls.Config{
		Renegotiation: tls.RenegotiateOnceAsClient,
	}

	return &HTTPClient{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
				TLSClientConfig: tlsConfig,
			},
			Timeout: time.Duration(15 * time.Second),
		},
		defaultUserAgent: defaultUserAgent,
	}
}

func (c *HTTPClient) prepareRequest(request *http.Request) {
	// if user-agent isn't set, use our default
	if request.Header.Get(httpHeaderUserAgent) == "" {
		request.Header.Set(httpHeaderUserAgent, c.defaultUserAgent)
	}
}

// Do does the given HTTP request
func (c *HTTPClient) Do(request *http.Request) (*http.Response, error) {
	c.prepareRequest(request)

	return c.client.Do(request)
}

// DoWithDump does the given HTTP request and returns a dump of the entire request
func (c *HTTPClient) DoWithDump(request *http.Request) (*http.Response, string, error) {
	c.prepareRequest(request)

	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, "", err
	}

	response, err := c.client.Do(request)

	return response, string(dump), err
}

// MockWithDump mocks the given HTTP request and returns a dump of the entire request
func (c *HTTPClient) MockWithDump(request *http.Request, mockStatus int, mockResponse string) (*http.Response, string, error) {
	c.prepareRequest(request)

	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		return nil, "", err
	}

	recorder := httptest.NewRecorder()
	recorder.WriteString(mockResponse)
	recorder.Code = mockStatus

	response := recorder.Result()
	response.Request = request

	return response, string(dump), nil
}
