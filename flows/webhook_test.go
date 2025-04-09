package flows_test

import (
	"context"
	"testing"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookCall(t *testing.T) {
	ctx := context.Background()

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

func TestExtractJSON(t *testing.T) {
	tcs := []struct {
		body []byte
		json []byte
	}{
		{[]byte(`{`), nil}, // invalid JSON
		{[]byte(`"x"`), []byte(`"x"`)},
		{[]byte(`{"foo": ["x"]}`), []byte(`{"foo": ["x"]}`)},
		{[]byte("\"a\x80\x81b\""), []byte(`"ab"`)},                     // invalid UTF-8 sequences stripped
		{[]byte("\u0000{\"foo\": 123\u0000}"), []byte(`{"foo": 123}`)}, // null chars stripped
		{[]byte(`"a\u0000b"`), []byte(`"ab"`)},                         // escaped null chars stripped
		{[]byte(`"01\02\03"`), nil},                                    // \0 not valid JSON escape
		{[]byte(`"01\\02\\03"`), []byte(`"01\\02\\03"`)},
	}

	for _, tc := range tcs {
		actual := flows.ExtractJSON(tc.body)
		assert.Equal(t, string(tc.json), string(actual), "extracted JSON mismatch for %s", string(tc.body))
	}

	asXValue := types.JSONToXValue([]byte(`{"foo": "01\\02\\03"}`))
	asXObject := asXValue.(*types.XObject)
	foo, _ := asXObject.Get("foo")
	assert.Equal(t, types.NewXText(`01\02\03`), foo)
	assert.Equal(t, `"01\\02\\03"`, string(jsonx.MustMarshal(foo)))
}
