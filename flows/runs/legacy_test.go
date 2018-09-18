package runs_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyExtra(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"bool": true, "number": 123.34, "text": "hello", "dict": {"foo": "bar"}, "list": [1, "x"]}`))
	}))
	server.Start()
	defer server.Close()

	session, err := test.CreateTestSession(server.URL, nil)
	require.NoError(t, err)

	run := session.Runs()[0]

	tests := []struct {
		template string
		output   string
	}{
		{"@legacy_extra.address.state", "WA"},
		{"@legacy_extra.ADDRESS.StaTE", "WA"},
		{"@legacy_extra.ADDRESS", `{"state":"WA"}`},
		{"@legacy_extra.bool", `true`},
		{"@legacy_extra.number", `123.34`},
		{"@legacy_extra.text", `hello`},
		{"@legacy_extra.list", `[1,"x"]`},
		{"@legacy_extra.list.0", `1`},
		{"@legacy_extra.list.1", `x`},
		{"@legacy_extra.dict.FOO", `bar`},
		{"@legacy_extra", `{"address":{"state":"WA"},"bool":true,"dict":{"foo":"bar"},"list":[1,"x"],"number":123.34,"source":"website","text":"hello"}`},
	}
	for _, tc := range tests {
		output, err := run.EvaluateTemplateAsString(tc.template, false)
		assert.NoError(t, err)
		assert.Equal(t, tc.output, output, "evaluate failed for %s", tc.template)
	}
}
