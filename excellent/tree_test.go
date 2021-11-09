package excellent

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tcs := []struct {
		expression string
		parsed     Expression
	}{
		{
			expression: `"hello\nworld"`,
			parsed:     &TextLiteral{val: types.NewXText("hello\nworld")},
		},
		{
			expression: `"\w+"`,
			parsed:     &TextLiteral{val: types.NewXText("\\w+")},
		},
		{
			expression: `"abc" & "cde"`,
			parsed: &Concatenation{
				exp1: &TextLiteral{val: types.NewXText("abc")},
				exp2: &TextLiteral{val: types.NewXText("cde")},
			},
		},
		{
			expression: `upper("abc")`,
			parsed: &FunctionCall{
				function: &ContextReference{name: "upper"},
				params:   []Expression{&TextLiteral{val: types.NewXText("abc")}},
			},
		},
		{
			expression: `(x) => upper(x)`,
			parsed: &AnonFunction{
				args: []string{"x"},
				body: &FunctionCall{
					function: &ContextReference{name: "upper"},
					params:   []Expression{&ContextReference{name: "x"}},
				},
			},
		},
		{
			expression: `((x) => upper(x))("abc")`,
			parsed: &FunctionCall{
				function: &Parentheses{
					exp: &AnonFunction{
						args: []string{"x"},
						body: &FunctionCall{
							function: &ContextReference{name: "upper"},
							params:   []Expression{&ContextReference{name: "x"}},
						},
					},
				},
				params: []Expression{&TextLiteral{val: types.NewXText("abc")}},
			},
		},
	}

	for _, tc := range tcs {
		exp, err := Parse(tc.expression, nil)
		assert.NoError(t, err)
		assert.Equal(t, tc.parsed, exp, "parsed mismatch for expression: %s", tc.expression)
	}
}

func TestExpressionsToString(t *testing.T) {
	foo := &ContextReference{name: "foo"}
	abc := &TextLiteral{types.NewXText("abc")}
	cde := &TextLiteral{types.NewXText("cde")}
	one := &NumberLiteral{val: types.RequireXNumberFromString(`1`)}
	two := &NumberLiteral{val: types.RequireXNumberFromString(`2`)}

	assert.Equal(t, `foo`, foo.String())

	assert.Equal(t, `foo.bar`, (&DotLookup{container: foo, lookup: "bar"}).String())
	assert.Equal(t, `foo.1`, (&DotLookup{container: foo, lookup: "1"}).String())

	assert.Equal(t, `foo["abc"]`, (&ArrayLookup{container: foo, lookup: abc}).String())
	assert.Equal(t, `foo[1]`, (&ArrayLookup{container: foo, lookup: one}).String())

	assert.Equal(t, `foo("abc", 1)`, (&FunctionCall{function: foo, params: []Expression{abc, one}}).String())
	assert.Equal(t, `foo()`, (&FunctionCall{function: foo, params: []Expression{}}).String())

	assert.Equal(t, `(x, y) => "abc"`, (&AnonFunction{args: []string{"x", "y"}, body: abc}).String())

	assert.Equal(t, `"abc" & "cde"`, (&Concatenation{exp1: abc, exp2: cde}).String())

	assert.Equal(t, `1 + 2`, (&Addition{exp1: one, exp2: two}).String())
	assert.Equal(t, `1 - 2`, (&Subtraction{exp1: one, exp2: two}).String())
	assert.Equal(t, `1 * 2`, (&Multiplication{exp1: one, exp2: two}).String())
	assert.Equal(t, `1 / 2`, (&Division{exp1: one, exp2: two}).String())
	assert.Equal(t, `1 ^ 2`, (&Exponent{expression: one, exponent: two}).String())
	assert.Equal(t, `-1`, (&Negation{exp: one}).String())

	assert.Equal(t, `1 = 2`, (&Equality{exp1: one, exp2: two}).String())
	assert.Equal(t, `1 != 2`, (&InEquality{exp1: one, exp2: two}).String())

	assert.Equal(t, `1 < 2`, (&LessThan{exp1: one, exp2: two}).String())
	assert.Equal(t, `1 <= 2`, (&LessThanOrEqual{exp1: one, exp2: two}).String())
	assert.Equal(t, `1 > 2`, (&GreaterThan{exp1: one, exp2: two}).String())
	assert.Equal(t, `1 >= 2`, (&GreaterThanOrEqual{exp1: one, exp2: two}).String())

	assert.Equal(t, `("abc")`, (&Parentheses{exp: abc}).String())

	assert.Equal(t, `""`, (&TextLiteral{val: types.XTextEmpty}).String())
	assert.Equal(t, `"abc"`, abc.String())
	assert.Equal(t, `"don't say \"hello\""`, (&TextLiteral{val: types.NewXText(`don't say "hello"`)}).String())
	assert.Equal(t, `123.5`, (&NumberLiteral{val: types.RequireXNumberFromString(`123.5`)}).String())
	assert.Equal(t, `123`, (&NumberLiteral{val: types.RequireXNumberFromString(`123.0`)}).String())
	assert.Equal(t, `true`, (&BooleanLiteral{val: types.XBooleanTrue}).String())
	assert.Equal(t, `false`, (&BooleanLiteral{val: types.XBooleanFalse}).String())
	assert.Equal(t, `null`, (&NullLiteral{}).String())
}
