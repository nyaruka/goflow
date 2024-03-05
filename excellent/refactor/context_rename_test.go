package refactor_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/refactor"
	"github.com/stretchr/testify/assert"
)

func TestContextRefRename(t *testing.T) {
	tcs := []struct {
		template string
		from     string
		to       string
		expected string
	}{
		{"@foo", "foo", "bar", "@bar"},
		{" @foo @foo ", "foo", "bar", " @bar @bar "},
		{"@(foo.uuid + 1)", "foo", "bar", "@(bar.uuid + 1)"},
		{"@(Upper(Foo))", "foo", "bar", "@(upper(bar))"},
		{"@webhook", "webhook", "webhook.json", "@webhook.json"},
		{"@( webhook[0] )", "webhook", "webhook.json", "@(webhook.json[0])"},
		{"@( 1 +  2)", "webhook", "webhook.json", "@( 1 +  2)"}, // unchanged because no change needed
	}

	topLevels := []string{"foo", "webhook"}

	for _, tc := range tcs {
		actual, err := refactor.Template(tc.template, topLevels, refactor.ContextRefRename(tc.from, tc.to))
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual, "refactor mismatch for template: %s", tc.template)
	}
}
