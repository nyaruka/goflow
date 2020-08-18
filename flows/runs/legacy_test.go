package runs_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyExtra(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"bool": true, "number": 123.34, "text": "hello", "object": {"foo": "bar", "1": "xx"}, "array": [1, "x"]}`))
	}))
	server.Start()
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
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
		{"@legacy_extra.array", `[1, x]`},
		{"@(legacy_extra.array[0])", `1`},
		{"@(legacy_extra.array[1])", `x`},
		{"@legacy_extra.object.FOO", `bar`},
		{`@(legacy_extra.object["1"])`, `xx`},
		{"@legacy_extra", `{address: {state: WA}, array: [1, x], bool: true, entities: {location: [{confidence: 1, value: Quito}]}, intent: {"intents":[{"name":"book_flight","confidence":0.5},{"name":"book_hotel","confidence":0.25}],"entities":{"location":[{"value":"Quito","confidence":1}]}}, intents: [{confidence: 0.5, name: book_flight}, {confidence: 0.25, name: book_hotel}], number: 123.34, object: {1: xx, foo: bar}, source: website, text: hello, webhook: {"bool": true, "number": 123.34, "text": "hello", "object": {"foo": "bar", "1": "xx"}, "array": [1, "x"]}}`},
	}
	for _, tc := range tests {
		output, err := run.EvaluateTemplate(tc.template)
		assert.NoError(t, err)
		assert.Equal(t, tc.output, output, "evaluate failed for %s", tc.template)
	}

	// can also add something which is an array
	result := flows.NewResult("webhook", "200", "Success", "", flows.NodeUUID(""), "", []byte(`[{"foo": 123}, {"foo": 345}]`), dates.Now())
	run.SaveResult(result)

	output, err := run.EvaluateTemplate(`@(legacy_extra[0])`)
	assert.NoError(t, err)
	assert.Equal(t, `{foo: 123}`, output)
}
