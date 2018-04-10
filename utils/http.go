package utils

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

func init() {
	Validator.RegisterAlias("http_method", "eq=GET|eq=HEAD|eq=POST|eq=PUT|eq=PATCH|eq=DELETE")
}

var (
	transport *http.Transport
	client    *http.Client
	once      sync.Once
)

func NewHTTPClient() *http.Client {
	once.Do(func() {
		timeout := time.Duration(15 * time.Second)
		transport = &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		}
		client = &http.Client{Transport: transport, Timeout: timeout}
	})

	return client
}

func getInsecureClient() *http.Client {
	once.Do(func() {
		timeout := time.Duration(15 * time.Second)
		transport = &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: transport, Timeout: timeout}
	})

	return client
}

// NewTestHTTPServer sets up a mock server for webhook actions
func NewTestHTTPServer() (*httptest.Server, error) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		defer r.Body.Close()
		w.Header().Set("Date", "")

		switch cmd {
		case "success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{ "ok": "true" }`))
		case "unavailable":
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{ "errors": ["service unavailable"] }`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{ "errors": ["bad_request"] }`))
		}
	}))
	// manually create a listener for our test server so that our output is predictable
	l, err := net.Listen("tcp", "127.0.0.1:49999")
	if err != nil {
		return nil, err
	}
	server.Listener = l
	return server, nil
}
