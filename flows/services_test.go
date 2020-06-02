package flows_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPLogs(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, nil, "hello"),
			httpx.NewMockResponse(400, nil, "is error"),
			httpx.MockConnectionError,
		},
	}))

	req1, err := httpx.NewRequest("GET", "http://temba.io/", nil, nil)
	require.NoError(t, err)
	trace1, err := httpx.DoTrace(http.DefaultClient, req1, nil, nil, -1)
	require.NoError(t, err)

	req2, err := httpx.NewRequest("GET", "http://temba.io/", nil, nil)
	require.NoError(t, err)
	trace2, err := httpx.DoTrace(http.DefaultClient, req2, nil, nil, -1)
	require.NoError(t, err)

	req3, err := httpx.NewRequest("GET", "http://temba.io/", nil, nil)
	require.NoError(t, err)
	trace3, err := httpx.DoTrace(http.DefaultClient, req3, nil, nil, -1)
	require.EqualError(t, err, "unable to connect to server")

	log1 := flows.NewHTTPLog(trace1, flows.HTTPStatusFromCode, nil)
	assert.Equal(t, flows.CallStatusSuccess, log1.Status)

	log2 := flows.NewHTTPLog(trace2, flows.HTTPStatusFromCode, nil)
	assert.Equal(t, flows.CallStatusResponseError, log2.Status)

	log3 := flows.NewHTTPLog(trace3, flows.HTTPStatusFromCode, nil)
	assert.Equal(t, flows.CallStatusConnectionError, log3.Status)
}

func TestHTTPLogsRedaction(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"http://temba.io/code/987654321/": {
			httpx.NewMockResponse(200, nil, `{"value": "987654321", "secret": "43t34wf#@f3"}`),
			httpx.NewMockResponse(400, nil, "The code is 987654321, I said 987654321"),
		},
	}))

	// code in URL and a header
	req1, err := httpx.NewRequest("GET", "http://temba.io/code/987654321/", nil, map[string]string{"X-Code": "987654321"})
	require.NoError(t, err)
	trace1, err := httpx.DoTrace(http.DefaultClient, req1, nil, nil, -1)
	require.NoError(t, err)

	// code in URL the request body
	req2, err := httpx.NewRequest("GET", "http://temba.io/code/987654321/", strings.NewReader("My code is 987654321"), nil)
	require.NoError(t, err)
	trace2, err := httpx.DoTrace(http.DefaultClient, req2, nil, nil, -1)
	require.NoError(t, err)

	redactor := utils.NewRedactor(flows.RedactionMask, "987654321", "43t34wf#@f3")

	log1 := flows.NewHTTPLog(trace1, flows.HTTPStatusFromCode, redactor)
	assert.Equal(t, "GET /code/****************/ HTTP/1.1\r\nHost: temba.io\r\nUser-Agent: Go-http-client/1.1\r\nX-Code: ****************\r\nAccept-Encoding: gzip\r\n\r\n", log1.Request)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 47\r\n\r\n{\"value\": \"****************\", \"secret\": \"****************\"}", log1.Response)

	log2 := flows.NewHTTPLog(trace2, flows.HTTPStatusFromCode, redactor)
	assert.Equal(t, "GET /code/****************/ HTTP/1.1\r\nHost: temba.io\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 20\r\nAccept-Encoding: gzip\r\n\r\nMy code is ****************", log2.Request)
	assert.Equal(t, "HTTP/1.0 400 Bad Request\r\nContent-Length: 39\r\n\r\nThe code is ****************, I said ****************", log2.Response)
}
