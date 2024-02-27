package excellent_test

import (
	"testing"

	. "github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/stretchr/testify/assert"
)

func TestParseTrees(t *testing.T) {
	tcs := []struct {
		expression string
		parsed     Expression
	}{
		{
			expression: `"hello\nworld"`,
			parsed:     &TextLiteral{Value: types.NewXText("hello\nworld")},
		},
		{
			expression: `"\w+"`,
			parsed:     &TextLiteral{Value: types.NewXText("\\w+")},
		},
		{
			expression: `"abc" & "cde"`,
			parsed: &Concatenation{
				Exp1: &TextLiteral{Value: types.NewXText("abc")},
				Exp2: &TextLiteral{Value: types.NewXText("cde")},
			},
		},
		{
			expression: `upper("abc")`,
			parsed: &FunctionCall{
				Func:   &ContextReference{Name: "upper"},
				Params: []Expression{&TextLiteral{Value: types.NewXText("abc")}},
			},
		},
		{
			expression: `(x) => upper(x)`,
			parsed: &AnonFunction{
				Args: []string{"x"},
				Body: &FunctionCall{
					Func:   &ContextReference{Name: "upper"},
					Params: []Expression{&ContextReference{Name: "x"}},
				},
			},
		},
		{
			expression: `((x) => upper(x))("abc")`,
			parsed: &FunctionCall{
				Func: &Parentheses{
					Exp: &AnonFunction{
						Args: []string{"x"},
						Body: &FunctionCall{
							Func:   &ContextReference{Name: "upper"},
							Params: []Expression{&ContextReference{Name: "x"}},
						},
					},
				},
				Params: []Expression{&TextLiteral{Value: types.NewXText("abc")}},
			},
		},
	}

	for _, tc := range tcs {
		exp, err := Parse(tc.expression, nil)
		assert.NoError(t, err)
		assert.Equal(t, tc.parsed, exp, "parsed mismatch for expression: %s", tc.expression)
	}
}

func TestExpressionVisitAndString(t *testing.T) {
	foo := &ContextReference{Name: "foo"}
	abc := &TextLiteral{types.NewXText("abc")}
	cde := &TextLiteral{types.NewXText("cde")}
	one := &NumberLiteral{Value: types.RequireXNumberFromString(`1`)}
	two := &NumberLiteral{Value: types.RequireXNumberFromString(`2`)}

	tcs := []struct {
		exp      Expression
		expected []string
	}{
		{foo, []string{"foo"}},

		{&DotLookup{Container: foo, Lookup: "bar"}, []string{`foo`, `foo.bar`}},
		{&DotLookup{Container: foo, Lookup: "1"}, []string{`foo`, `foo.1`}},

		{&ArrayLookup{Container: foo, Lookup: abc}, []string{"foo", `"abc"`, `foo["abc"]`}},
		{&ArrayLookup{Container: foo, Lookup: one}, []string{"foo", `1`, `foo[1]`}},

		{&FunctionCall{Func: foo, Params: []Expression{abc, one}}, []string{"foo", `"abc"`, `1`, `foo("abc", 1)`}},
		{&FunctionCall{Func: foo, Params: []Expression{}}, []string{`foo`, `foo()`}},

		{&AnonFunction{Args: []string{"x", "y"}, Body: abc}, []string{`"abc"`, `(x, y) => "abc"`}},

		{&Concatenation{Exp1: abc, Exp2: cde}, []string{`"abc"`, `"cde"`, `"abc" & "cde"`}},

		{&Addition{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 + 2`}},
		{&Subtraction{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 - 2`}},
		{&Multiplication{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 * 2`}},
		{&Division{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 / 2`}},
		{&Exponent{Expression: one, Exponent: two}, []string{`1`, `2`, `1 ^ 2`}},
		{&Negation{Exp: one}, []string{`1`, `-1`}},

		{&Equality{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 = 2`}},
		{&InEquality{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 != 2`}},
		{&LessThan{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 < 2`}},
		{&LessThanOrEqual{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 <= 2`}},
		{&GreaterThan{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 > 2`}},
		{&GreaterThanOrEqual{Exp1: one, Exp2: two}, []string{`1`, `2`, `1 >= 2`}},

		{&Parentheses{Exp: abc}, []string{`"abc"`, `("abc")`}},

		{&TextLiteral{Value: types.XTextEmpty}, []string{`""`}},
		{abc, []string{`"abc"`}},
		{&TextLiteral{Value: types.NewXText(`don't say "hello"`)}, []string{`"don't say \"hello\""`}},
		{&NumberLiteral{Value: types.RequireXNumberFromString(`123.5`)}, []string{`123.5`}},
		{&NumberLiteral{Value: types.RequireXNumberFromString(`123.0`)}, []string{`123`}},
		{&BooleanLiteral{Value: types.XBooleanTrue}, []string{`true`}},
		{&BooleanLiteral{Value: types.XBooleanFalse}, []string{`false`}},
		{&NullLiteral{}, []string{`null`}},
	}

	for _, tc := range tcs {
		var log []string
		visit := func(e Expression) { log = append(log, e.String()) } // log string version of each expression visited
		tc.exp.Visit(visit)
		assert.Equal(t, tc.expected, log)
		assert.Equal(t, tc.exp.String(), log[len(log)-1]) // last visit should be to the top level expression
	}
}
