package flows_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPLogs(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"http://temba.io/": []httpx.MockResponse{
			httpx.NewMockResponse(200, nil, "hello"),
			httpx.NewMockResponse(400, nil, "is error"),
			httpx.MockConnectionError,
		},
	}))

	trace1, err := httpx.NewTrace(http.DefaultClient, "GET", "http://temba.io/", nil, nil, nil, nil)
	require.NoError(t, err)

	trace2, err := httpx.NewTrace(http.DefaultClient, "GET", "http://temba.io/", nil, nil, nil, nil)
	require.NoError(t, err)

	trace3, err := httpx.NewTrace(http.DefaultClient, "GET", "http://temba.io/", nil, nil, nil, nil)
	require.EqualError(t, err, "unable to connect to server")

	log1 := flows.NewHTTPLog(trace1, flows.HTTPStatusFromCode)
	assert.Equal(t, flows.CallStatusSuccess, log1.Status)

	log2 := flows.NewHTTPLog(trace2, flows.HTTPStatusFromCode)
	assert.Equal(t, flows.CallStatusResponseError, log2.Status)

	log3 := flows.NewHTTPLog(trace3, flows.HTTPStatusFromCode)
	assert.Equal(t, flows.CallStatusConnectionError, log3.Status)
}
