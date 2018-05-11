package excellent

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

var xs = types.NewXText
var xn = types.RequireXNumberFromString
var xi = types.NewXNumberFromInt
var xd = types.NewXDateTime

type testXObject struct {
	foo string
	bar int
}

func NewTestXObject(foo string, bar int) *testXObject {
	return &testXObject{foo: foo, bar: bar}
}

// Describe returns a representation of this type for error messages
func (v *testXObject) Describe() string { return "test" }

func (v *testXObject) Reduce() types.XPrimitive { return types.NewXText(v.foo) }

func (v *testXObject) Resolve(key string) types.XValue {
	switch key {
	case "foo":
		return types.NewXText("bar")
	case "zed":
		return types.NewXNumberFromInt(123)
	case "missing":
		return nil
	default:
		return types.NewXResolveError(v, key)
	}
}

// ToXJSON is called when this type is passed to @(json(...))
func (v *testXObject) ToXJSON() types.XText {
	return types.ResolveKeys(v, "foo", "bar").ToXJSON()
}

var _ types.XValue = &testXObject{}
var _ types.XResolvable = &testXObject{}

func TestEvaluateTemplateAsString(t *testing.T) {

	vars := types.NewXMap(map[string]types.XValue{
		"string1": types.NewXText("foo"),
		"string2": types.NewXText("bar"),
		"汉字":      types.NewXText("simplified chinese"),
		"int1":    types.NewXNumberFromInt(1),
		"int2":    types.NewXNumberFromInt(2),
		"dec1":    types.RequireXNumberFromString("1.5"),
		"dec2":    types.RequireXNumberFromString("2.5"),
		"words":   types.NewXText("one two three"),
		"array":   types.NewXArray(types.NewXText("one"), types.NewXText("two"), types.NewXText("three")),
		"thing":   NewTestXObject("hello", 123),
		"err":     types.NewXError(fmt.Errorf("an error")),
	})

	evaluateAsStringTests := []struct {
		template string
		expected string
		hasError bool
	}{
		{"hello world", "hello world", false},
		{"@(\"hello\\nworld\")", "hello\nworld", false},
		{"@(\"hello😁world\")", "hello😁world", false},
		{"@(\"hello\\U0001F601world\")", "hello😁world", false},
		{"@(title(\"hello\"))", "Hello", false},
		{"@(title(hello))", "", true},
		{"Hello @(title(string1))", "Hello Foo", false},
		{"Hello @@string1", "Hello @string1", false},

		// an identifier which isn't valid top-level is ignored completely
		{"@hello", "@hello", false},
		{"@hello.bar", "@hello.bar", false},
		{"My email is foo@bar.com", "My email is foo@bar.com", false},

		// identifier which is valid top-level, errors and isn't echo'ed back
		{"@string1.xxx", "", true},

		{"1 + 2", "1 + 2", false},
		{"@(1 + 2)", "3", false},
		{"@@string1", "@string1", false},
		{"@@@string1", "@foo", false},

		{"@string1@string2", "foobar", false},
		{"@(string1 & string2)", "foobar", false},
		{"@string1.@string2", "foo.bar", false},
		{"@string1.@string2.@string3", "foo.bar.@string3", false},

		{"@(汉字)", "simplified chinese", false},
		{"@(string1", "@(string1", false},
		{"@ (string1", "@ (string1", false},
		{"@ (string1)", "@ (string1)", false},

		{"@(int1 + int2)", "3", false},
		{"@(1 + \"asdf\")", "", true},

		{"@(int1 + string1)", "", true},

		{"@(dec1 + dec2)", "4", false},

		{"@(TITLE(missing))", "", true},
		{"@(TITLE(string1.xxx))", "", true},

		{"@array", `["one","two","three"]`, false},
		{"@array[0]", `["one","two","three"][0]`, false}, // [n] notation not supported outside expression
		{"@array.0", "one", false},                       // works as dot notation however
		{"@(array [0])", "one", false},
		{"@(array[0])", "one", false},
		{"@(array.0)", "one", false},
		{"@(array[-1])", "three", false}, // negative index
		{"@(array.-1)", "", true},        // invalid negative index

		{"@(split(words, \" \").0)", "one", false},
		{"@(split(words, \" \")[1])", "two", false},
		{"@(split(words, \" \")[-1])", "three", false},

		{"@(thing.foo)", "bar", false},
		{"@(thing.zed)", "123", false},
		{"@(thing.missing)", "", false},    // missing is nil which becomes empty string
		{"@(thing.missing.xxx)", "", true}, // but can't look up a property on nil
		{"@(thing.xxx)", "", true},
	}

	env := utils.NewDefaultEnvironment()
	for _, test := range evaluateAsStringTests {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic evaluating template %s", test.template)
			}
		}()

		eval, err := EvaluateTemplateAsString(env, vars, test.template, false, vars.Keys())

		if test.hasError {
			assert.Error(t, err, "expected error evaluating template '%s'", test.template)
		} else {
			assert.NoError(t, err, "unexpected error evaluating template '%s'", test.template)

			if eval != test.expected {
				t.Errorf("Actual '%s' does not match expected '%s' evaluating template: '%s'", eval, test.expected, test.template)
			}
		}
	}
}

func TestEvaluateTemplate(t *testing.T) {
	array1d := types.NewXArray(types.NewXText("a"), types.NewXText("b"), types.NewXText("c"))
	array2d := types.NewXArray(array1d, types.NewXArray(types.NewXText("one"), types.NewXText("two"), types.NewXText("three")))

	vars := types.NewXMap(map[string]types.XValue{
		"string1": types.NewXText("foo"),
		"string2": types.NewXText("bar"),
		"key":     types.NewXText("four"),
		"int1":    types.NewXNumberFromInt(1),
		"int2":    types.NewXNumberFromInt(2),
		"dec1":    types.RequireXNumberFromString("1.5"),
		"dec2":    types.RequireXNumberFromString("2.5"),
		"words":   types.NewXText("one two three"),
		"array1d": array1d,
		"array2d": array2d,
	})

	env := utils.NewDefaultEnvironment()

	evaluateTests := []struct {
		template string
		expected types.XValue
		hasError bool
	}{
		{"hello world", xs("hello world"), false},
		{"@hello", xs("@hello"), false},
		{"@(title(\"hello\"))", xs("Hello"), false},

		{"@dec1", xn("1.5"), false},
		{"@(dec1 + dec2)", xn("4.0"), false},

		{"@array1d", array1d, false},
		{"@array1d.0", xs("a"), false},
		{"@array1d.1", xs("b"), false},
		{"@array2d.0.2", xs("c"), false},
		{"@(array1d[0])", xs("a"), false},
		{"@(array1d[1])", xs("b"), false},
		{"@(array2d[0])", array1d, false},
		{"@(array2d[0][2])", xs("c"), false},

		{"@string1 world", xs("foo world"), false},

		{"@(-10)", xi(-10), false},
		{"@(-asdf)", nil, true},

		{"@(2^2)", xi(4), false},
		{"@(2^asdf)", nil, true},
		{"@(asdf^2)", nil, true},

		{"@(1+2)", xi(3), false},
		{"@(1-2.5)", xn("-1.5"), false},
		{"@(1-asdf)", nil, true},
		{"@(asdf+1)", nil, true},

		{"@(1*2)", xi(2), false},
		{"@(1/2)", xn("0.5"), false},
		{"@(1/0)", nil, true},
		{"@(1*asdf)", nil, true},
		{"@(asdf/1)", nil, true},

		{"@(false)", types.XBooleanFalse, false},
		{"@(TRUE)", types.XBooleanTrue, false},

		{"@(1+1+1)", xi(3), false},
		{"@(5-2+1)", xi(4), false},
		{"@(2*3*4+2)", xi(26), false},
		{"@(4*3/4)", xi(3), false},
		{"@(4/2*4)", xi(8), false},
		{"@(2^2^2)", xi(16), false},
		{"@(\"a\" & \"b\" & \"c\")", xs("abc"), false},
		{"@(1+3 <= 1+4)", types.XBooleanTrue, false},

		// string equality
		{`@("asdf" = "asdf")`, types.XBooleanTrue, false},
		{`@("asdf" = "basf")`, types.XBooleanFalse, false},
		{`@("asdf" = "ASDF")`, types.XBooleanFalse, false}, // case-sensitive
		{`@("asdf" != "asdf")`, types.XBooleanFalse, false},
		{`@("asdf" != "basf")`, types.XBooleanTrue, false},

		// bool equality
		{"@(true = true)", types.XBooleanTrue, false},
		{"@(true = false)", types.XBooleanFalse, false},
		{"@(true = TRUE)", types.XBooleanTrue, false},

		// numerical equality
		{"@((1 = 1))", types.XBooleanTrue, false},
		{"@((1 != 2))", types.XBooleanTrue, false},
		{"@(1 = 1)", types.XBooleanTrue, false},
		{"@(1 = 2)", types.XBooleanFalse, false},
		{"@(1 != 2)", types.XBooleanTrue, false},
		{"@(1 != 1)", types.XBooleanFalse, false},
		{"@(-1 = 1)", types.XBooleanFalse, false},
		{"@(1.0 = 1)", types.XBooleanTrue, false},
		{"@(1.1 = 1.10)", types.XBooleanTrue, false},
		{"@(1.1234 = 1.10)", types.XBooleanFalse, false},
		{`@(1 = number("1.0"))`, types.XBooleanTrue, false},
		{"@(11=11=11)", types.XBooleanFalse, false}, // 11=11 -> TRUE, then TRUE != 11

		// date equality
		{`@(datetime("2018-04-16") = datetime("2018-04-16"))`, types.XBooleanTrue, false},
		{`@(datetime("2018-04-16") != datetime("2018-04-16"))`, types.XBooleanFalse, false},
		{`@(datetime("2018-04-16") = datetime("2017-03-20"))`, types.XBooleanFalse, false},
		{`@(datetime("2018-04-16") != datetime("2017-03-20"))`, types.XBooleanTrue, false},
		{`@(datetime("xxx") == datetime("2017-03-20"))`, nil, true},

		// other comparsions must be numerical
		{"@(2 > 1)", types.XBooleanTrue, false},
		{"@(1 > 2)", types.XBooleanFalse, false},
		{"@(2 >= 1)", types.XBooleanTrue, false},
		{"@(1 >= 2)", types.XBooleanFalse, false},
		{"@(1 <= 2)", types.XBooleanTrue, false},
		{"@(2 <= 1)", types.XBooleanFalse, false},
		{"@(1 < 2)", types.XBooleanTrue, false},
		{"@(2 < 1)", types.XBooleanFalse, false},
		{`@(1 < "asdf")`, nil, true}, // can't use with strings
		{`@("asdf" < "basf")`, nil, true},
		{"@(1<2<3)", nil, true}, // can't chain

		// nulls
		{"@(null)", nil, false},
		{"@(NULL)", nil, false},
		{"@(null = NULL)", types.XBooleanTrue, false},
		{"@(null != NULL)", types.XBooleanFalse, false},

		{"@(\"foo\" & \"bar\")", xs("foobar"), false},
		{"@(missing & \"bar\")", nil, true},
		{"@(\"foo\" & missing)", nil, true},

		{"@(TITLE(string1))", xs("Foo"), false},
		{"@(MISSING(string1))", nil, true},
		{"@(TITLE(string1, string2))", nil, true},

		{"@(1 = asdf)", nil, true},       // asdf isn't a valid context item
		{"@(asdf = 1)", nil, true},       // asdf isn't a valid context item
		{"@((1 / 0).field)", nil, true},  // can't resolve a property on an error value
		{"@((1 / 0)[0])", nil, true},     // can't index into an error value
		{"@(array1d[1 / 0])", nil, true}, // index expression can't be an error

		{"@(split(words, \" \").0)", xs("one"), false},
		{"@(split(words, \" \")[1])", xs("two"), false},
		{"@(split(words, \" \")[-1])", xs("three"), false},
	}

	for _, test := range evaluateTests {
		result, err := EvaluateTemplate(env, vars, test.template, vars.Keys())

		if test.hasError {
			assert.Error(t, err, "expected error evaluating template '%s'", test.template)
			continue
		} else {
			assert.NoError(t, err, "unexpected error evaluating template '%s'", test.template)
		}

		if !types.Equals(result, test.expected) {
			assert.Fail(t, "", "unexpected value, expected %T{%s}, got %T{%s} for function %s(%#v)", test.expected, test.expected, result, result, test.template)
		}
	}
}

func TestEvaluationErrors(t *testing.T) {

	vars := types.NewXMap(map[string]types.XValue{
		"foo": types.NewXText("bar"),
	})
	env := utils.NewDefaultEnvironment()

	errorTests := []struct {
		template string
		errorMsg string
	}{
		// parser errors
		{`@(")`, `error evaluating @("): syntax error at "`},
		{`@('x')`, `error evaluating @('x'): syntax error at 'x'`},
		{`@(")@('x')`, `error evaluating @("): syntax error at ", error evaluating @('x'): syntax error at 'x'`},
		{`@(0 / )`, `error evaluating @(0 / ): syntax error at `},

		// resolver errors
		{`@(NULL.x)`, `error evaluating @(NULL.x): null has no property 'x'`},
		{`@("abc".v)`, `error evaluating @("abc".v): "abc" has no property 'v'`},
		{`@(False.g)`, `error evaluating @(False.g): false has no property 'g'`},
		{`@(1.1.0)`, `error evaluating @(1.1.0): 1.1 has no property '0'`},
		{`@(hello)`, `error evaluating @(hello): map has no property 'hello'`}, // this context is a map
		{`@(foo.x)`, `error evaluating @(foo.x): "bar" has no property 'x'`},
		{`@foo.x`, `error evaluating @foo.x: "bar" has no property 'x'`},
		{`@(array(1, 2)[5])`, `error evaluating @(array(1, 2)[5]): index 5 out of range for 2 items`},

		// conversion errors
		{`@(1 + null)`, `error evaluating @(1 + null): unable to convert null to a number`},
		{`@(1 + true)`, `error evaluating @(1 + true): unable to convert true to a number`},
		{`@("a" + 2)`, `error evaluating @("a" + 2): unable to convert "a" to a number`},
		{`@(format_datetime("x"))`, `error evaluating @(format_datetime("x")): error calling FORMAT_DATETIME: unable to convert "x" to a datetime`},
		{`@(format_datetime(3))`, `error evaluating @(format_datetime(3)): error calling FORMAT_DATETIME: unable to convert 3 to a datetime`},

		// function call errors
		{`@(FOO())`, `error evaluating @(FOO()): no function with name 'foo'`},
		{`@(length(1))`, `error evaluating @(length(1)): error calling LENGTH: value doesn't have length`},
		{`@(word_count())`, `error evaluating @(word_count()): error calling WORD_COUNT: need 1 argument(s), got 0`},
		{`@(word_count("a", "b", "c"))`, `error evaluating @(word_count("a", "b", "c")): error calling WORD_COUNT: need 1 argument(s), got 3`},
	}
	for _, tc := range errorTests {
		result, err := EvaluateTemplateAsString(env, vars, tc.template, false, vars.Keys())
		assert.Equal(t, "", result)
		assert.NotNil(t, err)

		if err != nil {
			assert.Equal(t, tc.errorMsg, err.Error(), "error message mismatch for template '%s'", tc.template)
		}
	}
}

func TestScanner(t *testing.T) {
	scanner := NewXScanner(strings.NewReader("12"), []string{})

	if scanner.read() != '1' {
		t.Errorf("Expected '1'")
	}
	scanner.unread('1')
	if scanner.read() != '1' {
		t.Errorf("Expected '1'")
	}
	if scanner.read() != '2' {
		t.Errorf("Expected '2'")
	}
	if scanner.read() != eof {
		t.Errorf("Expected eof")
	}
}
