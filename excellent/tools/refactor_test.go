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
	tcs := []struct {
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

	eval := excellent.NewEvaluator()
	env := envs.NewBuilder().Build()
	ctx := types.NewXObject(map[string]types.XValue{
		"foo": types.NewXObject(map[string]types.XValue{
			"bar": types.NewXNumberFromInt(123),
		}),
	})
	topLevels := []string{"foo"}

	tx := func(excellent.Expression) bool { return true } // always refactor

	for _, tc := range tcs {
		actual, err := tools.RefactorTemplate(tc.template, topLevels, tx)

		assert.Equal(t, tc.refactored, actual, "refactor mismatch for template: %s", tc.template)

		if tc.hasError {
			assert.Error(t, err, "expected error for template: %s", tc.template)
		} else {
			assert.NoError(t, err, "unexpected error for template: %s, err: %s", tc.template, err)

			// test that the original and the refactored template evaluate equally
			originalValue, _, _ := eval.Template(env, ctx, tc.template, nil)
			refactoredValue, _, _ := eval.Template(env, ctx, actual, nil)

			assert.Equal(t, originalValue, refactoredValue, "refactoring of template %s gives different value: %s", tc.template, refactoredValue)
		}
	}
}

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
		actual, err := tools.RefactorTemplate(tc.template, topLevels, tools.ContextRefRename(tc.from, tc.to))
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual, "refactor mismatch for template: %s", tc.template)
	}
}
