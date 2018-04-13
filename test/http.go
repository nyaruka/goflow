package test

import (
	"net"
	"net/http"
	"net/http/httptest"
)

// NewTestHTTPServer sets up a mock server for webhook actions
func NewTestHTTPServer() (*httptest.Server, error) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		defer r.Body.Close()
		w.Header().Set("Date", "Wed, 11 Apr 2018 18:24:30 GMT")

		switch cmd {
		case "success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{ "ok": "true" }`))
		case "echo":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(r.URL.Query().Get("content")))
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
	server.Start()
	return server, nil
}
