package runs_test

import (
	"testing"

	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyExtra(t *testing.T) {
	server, err := test.NewTestHTTPServer(0)
	require.NoError(t, err)
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
		{"@legacy_extra.results.1", `{"state":"IN"}`},
		{"@legacy_extra.results.1.state", `IN`},
		{"@legacy_extra", `{"address":{"state":"WA"},"results":[{"state":"WA"},{"state":"IN"}],"source":"website"}`},
	}
	for _, tc := range tests {
		output, err := run.EvaluateTemplateAsString(tc.template, false)
		assert.NoError(t, err)
		assert.Equal(t, tc.output, output, "evaluate failed for %s", tc.template)
	}
}
