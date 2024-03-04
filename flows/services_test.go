package flows_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPLogs(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, nil, []byte("hello \\u0000")),
			httpx.NewMockResponse(400, nil, []byte("is error")),
			httpx.MockConnectionError,
		},
		"http://temba.io/?x=" + strings.Repeat("x", 3000): {
			httpx.NewMockResponse(200, nil, []byte("hello")),
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

	req4, err := httpx.NewRequest("GET", "http://temba.io/?x="+strings.Repeat("x", 3000), nil, nil) // URL exceeds limit
	require.NoError(t, err)
	trace4, err := httpx.DoTrace(http.DefaultClient, req4, nil, nil, -1)
	require.NoError(t, err)

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

func TestWebhookCall(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, map[string]string{"Content-Type": "application/json"}, []byte(`{"foo":123}`)),
			httpx.NewMockResponse(400, nil, []byte("is error")),
			httpx.MockConnectionError,
			httpx.MockConnectionError,
			httpx.MockConnectionError,
		},
	}))

	eng := test.NewEngine()
	env := envs.NewBuilder().Build()
	svc, err := eng.Services().Webhook(nil)
	require.NoError(t, err)

	request := func(method string) *flows.WebhookCall {
		req1, err := httpx.NewRequest(method, "http://temba.io/", nil, nil)
		require.NoError(t, err)

		call, err := svc.Call(req1)
		require.NoError(t, err)
		require.NotNil(t, call)
		return call
	}

	call1 := request("GET")

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("GET http://temba.io/"),
		"status":      types.NewXNumberFromInt(200),
		"headers":     types.NewXObject(map[string]types.XValue{"Content-Type": types.NewXText("application/json")}),
		"json":        types.NewXObject(map[string]types.XValue{"foo": types.NewXNumberFromInt(123)}),
	}), flows.Context(env, call1))

	call2 := request("POST")

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("POST http://temba.io/"),
		"status":      types.NewXNumberFromInt(400),
		"headers":     types.XObjectEmpty,
		"json":        nil,
	}), flows.Context(env, call2))

	call3 := request("GET")

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("GET http://temba.io/"),
		"status":      types.NewXNumberFromInt(0),
		"headers":     types.XObjectEmpty,
		"json":        nil,
	}), flows.Context(env, call3))
}
