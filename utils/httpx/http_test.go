package httpx_test

import (
	"net/http"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoTrace(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)

	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))

	server := test.NewTestHTTPServer(52025)

	// test with a text response
	request, err := httpx.NewRequest("GET", server.URL+"?cmd=success", nil, nil)
	require.NoError(t, err)

	trace, err := httpx.DoTrace(http.DefaultClient, request, nil, nil, -1)
	assert.NoError(t, err)
	assert.Equal(t, "GET /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:52025\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n", string(trace.RequestTrace))
	assert.Equal(t, `{ "ok": "true" }`, string(trace.ResponseBody))
	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, "{ \"ok\": \"true\" }", string(trace.ResponseBody))
	assert.Equal(t, time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC), trace.StartTime)
	assert.Equal(t, time.Date(2019, 10, 7, 15, 21, 31, 123456789, time.UTC), trace.EndTime)

	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }", string(trace.ResponseTraceUTF8("...")))
	assert.Equal(t, ">>>>>>>> GET http://127.0.0.1:52025?cmd=success\nGET /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:52025\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n\n<<<<<<<<\nHTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }", trace.String())

	// test with a binary response
	request, err = httpx.NewRequest("GET", server.URL+"?cmd=binary&size=1000", nil, nil)
	require.NoError(t, err)

	trace, err = httpx.DoTrace(http.DefaultClient, request, nil, nil, -1)
	assert.NoError(t, err)
	assert.Equal(t, "GET /?cmd=binary&size=1000 HTTP/1.1\r\nHost: 127.0.0.1:52025\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Length: 1000\r\nContent-Type: application/octet-stream\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, 1000, len(trace.ResponseBody))
	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Length: 1000\r\nContent-Type: application/octet-stream\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n...", string(trace.ResponseTraceUTF8("...")))
}

func TestMaxBodyBytes(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	testBody := `abcdefghijklmnopqrstuvwxyz`

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://temba.io": {
			httpx.NewMockResponse(200, nil, testBody),
			httpx.NewMockResponse(200, nil, testBody),
			httpx.NewMockResponse(200, nil, testBody),
			httpx.NewMockResponse(200, nil, testBody),
		},
	}))

	call := func(maxBodyBytes int) (*httpx.Trace, error) {
		request, _ := http.NewRequest("GET", "https://temba.io", nil)
		return httpx.DoTrace(http.DefaultClient, request, nil, nil, maxBodyBytes)
	}

	trace, err := call(-1) // no body limit
	assert.NoError(t, err)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 26\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, testBody, string(trace.ResponseBody))

	trace, err = call(1000) // limit bigger than body
	assert.NoError(t, err)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 26\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, testBody, string(trace.ResponseBody))

	trace, err = call(len(testBody)) // limit same as body
	assert.NoError(t, err)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 26\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, testBody, string(trace.ResponseBody))

	trace, err = call(10) // limit smaller than body
	assert.EqualError(t, err, `webhook response body exceeds 10 bytes limit`)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 26\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, ``, string(trace.ResponseBody))
}

func TestNonUTF8(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://temba.io": {
			httpx.MockResponse{Status: 200, Headers: nil, Body: []byte{'\xc3', '\x28'}},
		},
	}))

	request, err := httpx.NewRequest("GET", "https://temba.io", nil, nil)
	require.NoError(t, err)

	trace, err := httpx.DoTrace(http.DefaultClient, request, nil, nil, -1)
	assert.NoError(t, err)
	assert.Equal(t, "GET / HTTP/1.1\r\nHost: temba.io\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 2\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, []byte{'\xc3', '\x28'}, trace.ResponseBody)
	assert.True(t, utf8.Valid(trace.ResponseTrace))
	assert.False(t, utf8.Valid(trace.ResponseBody))
}
