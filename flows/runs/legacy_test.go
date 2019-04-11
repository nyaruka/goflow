package runs_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyExtra(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"bool": true, "number": 123.34, "text": "hello", "dict": {"foo": "bar", "1": "xx"}, "list": [1, "x"]}`))
	}))
	server.Start()
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, nil)
	require.NoError(t, err)

	run := session.Runs()[0]

	tests := []struct {
		template string
		output   string
	}{
		{"@legacy_extra.address.state", "WA"},
		{"@legacy_extra.ADDRESS.StaTE", "WA"},
		{"@legacy_extra.ADDRESS", `{state: WA}`},
		{"@legacy_extra.bool", `true`},
		{"@legacy_extra.number", `123.34`},
		{"@legacy_extra.text", `hello`},
		{"@legacy_extra.list", `[1, x]`},
		{"@(legacy_extra.list[0])", `1`},
		{"@(legacy_extra.list[1])", `x`},
		{"@legacy_extra.dict.FOO", `bar`},
		{`@(legacy_extra.dict["1"])`, `xx`},
		{"@legacy_extra", `{address: {state: WA}, bool: true, dict: {1: xx, foo: bar}, list: [1, x], number: 123.34, source: website, text: hello, webhook: {"bool": true, "number": 123.34, "text": "hello", "dict": {"foo": "bar", "1": "xx"}, "list": [1, "x"]}}`},
	}
	for _, tc := range tests {
		output, err := run.EvaluateTemplate(tc.template)
		assert.NoError(t, err)
		assert.Equal(t, tc.output, output, "evaluate failed for %s", tc.template)
	}

	// can also add something which is an array
	result := flows.NewResult("webhook", "200", "Success", "", flows.NodeUUID(""), "", []byte(`[{"foo": 123}, {"foo": 345}]`), utils.Now())
	run.SaveResult(result)

	output, err := run.EvaluateTemplate(`@(legacy_extra[0])`)
	assert.NoError(t, err)
	assert.Equal(t, `{foo: 123}`, output)
}
