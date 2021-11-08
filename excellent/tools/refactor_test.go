package tools_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/tools"
	"github.com/nyaruka/goflow/excellent/types"
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

	env := envs.NewBuilder().Build()
	ctx := types.NewXObject(map[string]types.XValue{
		"foo": types.NewXObject(map[string]types.XValue{
			"bar": types.NewXNumberFromInt(123),
		}),
	})
	topLevels := []string{"foo"}

	for _, tc := range testCases {
		actual, err := tools.RefactorTemplate(tc.template, topLevels)

		assert.Equal(t, tc.refactored, actual, "refactor mismatch for template: %s", tc.template)

		if tc.hasError {
			assert.Error(t, err, "expected error for template: %s", tc.template)
		} else {
			assert.NoError(t, err, "unexpected error for template: %s, err: %s", tc.template, err)

			// test that the original and the refactored template evaluate equally
			originalValue, _ := excellent.EvaluateTemplate(env, ctx, tc.template, nil)
			refactoredValue, _ := excellent.EvaluateTemplate(env, ctx, actual, nil)

			assert.Equal(t, originalValue, refactoredValue, "refactoring of template %s gives different value: %s", tc.template, refactoredValue)
		}
	}
}
