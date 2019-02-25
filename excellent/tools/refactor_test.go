package tools_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/tools"

	"github.com/stretchr/testify/assert"
)

func TestRefactorTemplate(t *testing.T) {
	testCases := []struct {
		template   string
		refactored string
		hasError   bool
	}{
		{``, ``, false},
		{`Hi @foo`, `Hi @foo`, false},
		{`@(foo)`, `@(foo)`, false},
		{`@( "Hello"+12345.123 )`, `@("Hello" + 12345.123)`, false},
		{`@foo.bar`, `@foo.bar`, false},
		{`@(foo . bar)`, `@(foo.bar)`, false},
		{`@(OR(TRUE, False, Null))`, `@(or(true, false, null))`, false},
		{`@(foo[ 1 ] + foo[ "x" ])`, `@(foo[1] + foo["x"])`, false},
		{`@(-1+( 2/3 )*4^5)`, `@(-1 + (2 / 3) * 4 ^ 5)`, false},
		{`@("x"&"y")`, `@("x" & "y")`, false},
		{`@(AND("x"="y", "x"!="y"))`, `@(and("x" = "y", "x" != "y"))`, false},
		{`@(AND(1>2, 3<4, 5>=6, 7<=8))`, `@(and(1 > 2, 3 < 4, 5 >= 6, 7 <= 8))`, false},
		{`@(FOO_Func(x, y))`, `@(foo_func(x, y))`, false},
		{`@(1 / ) @(1+2)`, `@(1 / ) @(1 + 2)`, true},
	}

	for _, tc := range testCases {
		actual, err := tools.RefactorTemplate(tc.template, []string{"foo"})

		if tc.hasError {
			assert.Error(t, err, "expected error for template: %s", tc.template)
		} else {
			assert.NoError(t, err, "unexpected error for template: %s, err: %s", tc.template, err)
		}

		assert.Equal(t, tc.refactored, actual, "refactor mismatch for template: %s", tc.template)
	}
}
