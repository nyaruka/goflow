package tools_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/tools"
	"github.com/stretchr/testify/assert"
)

func TestFindContextRefsInTemplate(t *testing.T) {
	testCases := []struct {
		template string
		paths    [][]string
		hasError bool
	}{
		{``, [][]string{}, false},
		{`Hi @foo @foo.bar`, [][]string{{`foo`}, {`foo`}, {`foo`, `bar`}}, false},
		{`@foo.bar.x.y`, [][]string{{`foo`}, {`foo`, `bar`}, {`foo`, `bar`, `x`}, {`foo`, `bar`, `x`, `y`}}, false},
		{`@(foo) @(foo.bar)`, [][]string{{`foo`}, {`foo`}, {`foo`, `bar`}}, false},
		{`@((foo))`, [][]string{{`foo`}}, false},
		{`@((FOO))`, [][]string{{`foo`}}, false},
		{`@(lower(foo.bar))`, [][]string{{`foo`}, {`foo`, `bar`}}, false},
		{`@(foo["bar"])`, [][]string{{`foo`}, {`foo`, `bar`}}, false},
		{`@(3 * (foo.bar + 1) / 2)`, [][]string{{`foo`}, {`foo`, `bar`}}, false},
		{`@("foo.bar")`, [][]string{}, false},
		{`@(webhook.0.kd_prov)`, [][]string{{"webhook"}, {"webhook", "0"}, {"webhook", "0", "kd_prov"}}, false},
	}

	for _, tc := range testCases {
		actual := make([][]string, 0)

		err := tools.FindContextRefsInTemplate(tc.template, []string{"foo"}, func(path []string) {
			actual = append(actual, path)
		})

		if tc.hasError {
			assert.Error(t, err, "expected error for template: %s", tc.template)
		} else {
			assert.NoError(t, err, "unexpected error for template: %s, err: %s", tc.template, err)
		}

		assert.Equal(t, tc.paths, actual, "audit context mismatch for input: %s", tc.template)
	}
}
