package test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/nyaruka/goflow/utils"
)

// TestHTTPClient a HTTP client instance for tests
var TestHTTPClient = utils.NewHTTPClient("goflow-testing")

// NewTestHTTPServer sets up a mock server for webhook actions
func NewTestHTTPServer(port int) (*httptest.Server, error) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(testHTTPHandler))

	if port > 0 {
		// manually create a listener for our test server so that our output is predictable
		l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			return nil, err
		}
		server.Listener = l
	}
	server.Start()
	return server, nil
}

func testHTTPHandler(w http.ResponseWriter, r *http.Request) {
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
	case "binary":
		typeParam := r.URL.Query().Get("type")
		if typeParam == "" {
			typeParam = "application/octet-stream"
		}

		sizeParam := r.URL.Query().Get("size")
		if sizeParam == "" {
			sizeParam = "10"
		}
		size, _ := strconv.Atoi(sizeParam)
		data := make([]byte, size)
		for i := 0; i < size; i++ {
			data[i] = byte(40 + i%10)
		}

		w.Header().Set("Content-Type", typeParam)
		w.Header().Set("Content-Length", sizeParam)

		w.WriteHeader(http.StatusOK)
		w.Write(data)

	case "unavailable":
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{ "errors": ["service unavailable"] }`))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "errors": ["bad_request"] }`))
	}
}
