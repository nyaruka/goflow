package flows_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPLogs(t *testing.T) {
	ctx := context.Background()

	tracing := httpx.WithTraces(httpx.WithMocks(http.DefaultTransport, map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, nil, []byte("hello \\u0000")),
			httpx.NewMockResponse(400, nil, []byte("is error")),
			httpx.MockConnectionError,
		},
		"http://temba.io/?x=" + strings.Repeat("x", 3000): {
			httpx.NewMockResponse(200, nil, []byte("hello")),
		},
	}))
	client := &http.Client{Transport: tracing}

	req1, err := httpx.NewRequest(ctx, "GET", "http://temba.io/", nil, nil)
	require.NoError(t, err)
	_, err = client.Do(req1)
	require.NoError(t, err)

	req2, err := httpx.NewRequest(ctx, "GET", "http://temba.io/", nil, nil)
	require.NoError(t, err)
	_, err = client.Do(req2)
	require.NoError(t, err)

	req3, err := httpx.NewRequest(ctx, "GET", "http://temba.io/", nil, nil)
	require.NoError(t, err)
	_, err = client.Do(req3)
	require.ErrorContains(t, err, "unable to connect to server")

	req4, err := httpx.NewRequest(ctx, "GET", "http://temba.io/?x="+strings.Repeat("x", 3000), nil, nil) // URL exceeds limit
	require.NoError(t, err)
	_, err = client.Do(req4)
	require.NoError(t, err)

	traces := tracing.Traces()
	require.Len(t, traces, 4)
	trace1, trace2, trace3, trace4 := traces[0], traces[1], traces[2], traces[3]

	log1 := flows.NewHTTPLog(trace1, flows.HTTPStatusFromCode, nil)
	assert.Equal(t, flows.CallStatusSuccess, log1.Status)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 12\r\n\r\nhello �", log1.Response) // escaped null should have been replaced

	log2 := flows.NewHTTPLog(trace2, flows.HTTPStatusFromCode, nil)
	assert.Equal(t, flows.CallStatusResponseError, log2.Status)

	log3 := flows.NewHTTPLog(trace3, flows.HTTPStatusFromCode, nil)
	assert.Equal(t, flows.CallStatusConnectionError, log3.Status)

	log4 := flows.NewHTTPLog(trace4, flows.HTTPStatusFromCode, nil)
	assert.Equal(t, "http://temba.io/?x="+strings.Repeat("x", 2026)+"...", log4.URL) // trimmed
}
