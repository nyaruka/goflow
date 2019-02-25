package tools_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/tools"

	"github.com/stretchr/testify/assert"
)

func TestRefactorTemplate(t *testing.T) {
	testCases := []struct {
		old      string
		new      string
		hasError bool
	}{
		{`Hi @foo`, `Hi @foo`, false},
		{`Hi @(-1+2/3*4)`, `Hi @(-1 + 2 / 3 * 4)`, false},
		{`Hi @(foo[ 1 ] + foo[ "x" ])`, `Hi @(foo[1] + foo["x"])`, false},
		{`Hi @(1 / ) @(1+2)`, `Hi @(1 / ) @(1 + 2)`, true},
	}

	for _, tc := range testCases {
		actual, err := tools.RefactorTemplate(tc.old, []string{"foo"})

		if tc.hasError {
			assert.Error(t, err, "expected error for input: %s", tc.old)
		} else {
			assert.NoError(t, err, "unexpected error for input: %s, err: %s", tc.old, err)
		}

		assert.Equal(t, tc.new, actual, "refactor mismatch for input: %s", tc.old)
	}
}
