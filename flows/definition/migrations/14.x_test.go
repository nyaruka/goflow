package migrations_test

import (
	"strings"
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tests that escaping literal {{ doesn't change the output of templates
func TestMigrate14_7_0PreservesOutput(t *testing.T) {
	tcs := []string{
		`{{x}}`,
		`{{{x}}}`,
		`{{{{`,
		`{{ @contact.name }}`,
		`\{{x}}`,
		`\\{{x}}`,
		`\\\{{x}}`,
		`a\b \{x} \`,
		`@{{x}}`,
		`@@{{x}}`,
		`@\{{x}}`,
		`@@\{{x}}`,
		`@@@\{{x}}`,
		`@ {{x}}`,
		`foo@bar.com {{x}}`,
		`\{{\w+}}`,
	}

	env := envs.NewBuilder().Build()
	ctx := types.NewXObject(map[string]types.XValue{})
	eval := excellent.NewEvaluator(excellent.DefaultEvaluationBudget)

	for _, tc := range tcs {
		// before braces, templates with no expressions in an empty context output their text with @@ unescaped
		var expected strings.Builder
		scanner := excellent.NewXScanner(strings.NewReader(tc), nil)
		scanner.SetBraces(false)
		for tokenType, token := scanner.Scan(); tokenType != excellent.EOF; tokenType, token = scanner.Scan() {
			expected.WriteString(excellent.Source(tokenType, token))
		}

		flow := migrations.Flow{
			"spec_version": "14.6.0",
			"nodes": []any{map[string]any{"actions": []any{map[string]any{"type": "send_msg", "text": tc}}}},
		}
		migrated, err := migrations.Migrate14_7_0(flow, nil)
		require.NoError(t, err)

		text := migrated.Nodes()[0].Actions()[0]["text"].(string)
		actual, _, err := eval.Template(t.Context(), env, ctx, text, nil)

		assert.NoError(t, err, "unexpected error for %s migrated to %s", tc, text)
		assert.Equal(t, expected.String(), actual, "output mismatch for %s migrated to %s", tc, text)
	}
}
