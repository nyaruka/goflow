package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/require"
)

func TestWebhookCall(t *testing.T) {
	ctx := t.Context()

	eng := test.NewMockedEngine(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, map[string]string{"Content-Type": "application/json"}, []byte(`{"foo":123}`)),
			httpx.NewMockResponse(400, nil, []byte("is error")),
			httpx.MockConnectionError,
			httpx.MockConnectionError,
			httpx.MockConnectionError,
		},
	})
	env := envs.NewBuilder().Build()
	svc, err := eng.Services().Webhook(nil)
	require.NoError(t, err)

	request := func(method string) *flows.WebhookCall {
		req1, err := httpx.NewRequest(ctx, method, "http://temba.io/", nil, nil)
		require.NoError(t, err)

		trace, err := svc.Call(req1)
		require.NoError(t, err)
		require.NotNil(t, trace)

		return flows.NewWebhookCall(trace)
	}

	call1 := request("GET")

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("GET http://temba.io/"),
		"status":      types.NewXNumberFromInt(200),
		"headers":     types.NewXObject(map[string]types.XValue{"Content-Type": types.NewXText("application/json")}),
		"json":        types.NewXObject(map[string]types.XValue{"foo": types.NewXNumberFromInt(123)}),
	}), core.Context(env, call1))

	call2 := request("POST")

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("POST http://temba.io/"),
		"status":      types.NewXNumberFromInt(400),
		"headers":     types.XObjectEmpty,
		"json":        nil,
	}), core.Context(env, call2))

	call3 := request("GET")

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("GET http://temba.io/"),
		"status":      types.NewXNumberFromInt(0),
		"headers":     types.XObjectEmpty,
		"json":        nil,
	}), core.Context(env, call3))
}
