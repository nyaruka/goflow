package excellent

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type testResolvable struct{}

func (r *testResolvable) Resolve(key string) interface{} {
	switch key {
	case "foo":
		return "bar"
	case "zed":
		return 123
	case "missing":
		return nil
	default:
		return fmt.Errorf("no such thing")
	}
}

// Atomize is called when this object needs to be reduced to a primitive
func (r *testResolvable) Atomize() interface{} {
	return "hello"
}

var _ utils.Atomizable = (*testResolvable)(nil)
var _ utils.Resolvable = (*testResolvable)(nil)

func TestEvaluateTemplateAsString(t *testing.T) {

	varMap := map[string]interface{}{
		"string1": "foo",
		"string2": "bar",
		"汉字":      "simplified chinese",
		"int1":    1,
		"int2":    2,
		"dec1":    decimal.RequireFromString("1.5"),
		"dec2":    decimal.RequireFromString("2.5"),
		"words":   "one two three",
		"array":   utils.NewArray("one", "two", "three"),
		"thing":   &testResolvable{},
		"err":     fmt.Errorf("an error"),
	}
	vars := utils.NewMapResolver(varMap)

	keys := make([]string, 0, len(varMap))
	for key := range varMap {
		keys = append(keys, key)
	}

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

		{"@array", "one, two, three", false},
		{"@array[0]", "one, two, three[0]", false}, // [n] notation not supported outside expression
		{"@array.0", "one", false},                 // works as dot notation however
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
		eval, err := EvaluateTemplateAsString(env, vars, test.template, false, keys)

		if test.hasError {
			assert.Error(t, err, "expected error evaluating template '%s'", test.template)
		} else {
			assert.NoError(t, err, "unexpected error evaluating template '%s'", test.template)
		}
		if eval != test.expected {
			t.Errorf("Actual '%s' does not match expected '%s' evaluating template: '%s'", eval, test.expected, test.template)
		}
	}
}

func TestEvaluateTemplate(t *testing.T) {
	array1d := utils.NewArray("a", "b", "c")
	array2d := utils.NewArray(array1d, utils.NewArray("one", "two", "three"))

	varMap := map[string]interface{}{
		"string1": "foo",
		"string2": "bar",
		"key":     "four",
		"int1":    1,
		"int2":    2,
		"dec1":    decimal.RequireFromString("1.5"),
		"dec2":    decimal.RequireFromString("2.5"),
		"words":   "one two three",
		"array1d": array1d,
		"array2d": array2d,
	}

	vars := utils.NewMapResolver(varMap)

	keys := make([]string, 0, len(varMap))
	for key := range varMap {
		keys = append(keys, key)
	}

	env := utils.NewDefaultEnvironment()

	evaluateTests := []struct {
		template string
		expected interface{}
		hasError bool
	}{
		{"hello world", "hello world", false},
		{"@hello", "@hello", false},
		{"@(title(\"hello\"))", "Hello", false},

		{"@dec1", decimal.RequireFromString("1.5"), false},
		{"@(dec1 + dec2)", decimal.RequireFromString("4.0"), false},

		{"@array1d", array1d, false},
		{"@array1d.0", "a", false},
		{"@array1d.1", "b", false},
		{"@array2d.0.2", "c", false},
		{"@(array1d[0])", "a", false},
		{"@(array1d[1])", "b", false},
		{"@(array2d[0])", array1d, false},
		{"@(array2d[0][2])", "c", false},

		{"@string1 world", "foo world", false},

		{"@(-10)", -10, false},
		{"@(-asdf)", "", true},

		{"@(2^2)", 4, false},
		{"@(2^asdf)", "", true},
		{"@(asdf^2)", "", true},

		{"@(1+2)", 3, false},
		{"@(1-2.5)", decimal.RequireFromString("-1.5"), false},
		{"@(1-asdf)", "", true},
		{"@(asdf+1)", "", true},

		{"@(1*2)", 2, false},
		{"@(1/2)", decimal.RequireFromString("0.5"), false},
		{"@(1/0)", "", true},
		{"@(1*asdf)", "", true},
		{"@(asdf/1)", "", true},

		{"@(false)", false, false},
		{"@(TRUE)", true, false},

		{"@(1+1+1)", 3, false},
		{"@(5-2+1)", 4, false},
		{"@(2*3*4+2)", 26, false},
		{"@(4*3/4)", 3, false},
		{"@(4/2*4)", 8, false},
		{"@(2^2^2)", 16, false},
		{"@(11=11=11)", "", true},
		{"@(1<2<3)", "", true},
		{"@(\"a\" & \"b\" & \"c\")", "abc", false},
		{"@(1+3 <= 1+4)", true, false},

		{"@((1 = 1))", true, false},
		{"@((1 != 2))", true, false},
		{"@(2 > 1)", true, false},
		{"@(1 > 2)", false, false},
		{"@(2 >= 1)", true, false},
		{"@(1 >= 2)", false, false},
		{"@(1 <= 2)", true, false},
		{"@(2 <= 1)", false, false},
		{"@(1 < 2)", true, false},
		{"@(2 < 1)", false, false},
		{"@(1 = 1)", true, false},
		{"@(1 = 2)", false, false},
		{`@("asdf" = "basf")`, "", true},
		{"@(1 != 2)", true, false},
		{"@(1 != 1)", false, false},
		{"@(-1 = 1)", false, false},
		{"@(1 < asdf)", "", true},
		{`@("asdf" < "basf")`, "", true},

		{"@(\"foo\" & \"bar\")", "foobar", false},
		{"@(missing & \"bar\")", "", true},
		{"@(\"foo\" & missing)", "", true},

		{"@(TITLE(string1))", "Foo", false},
		{"@(MISSING(string1))", "", true},
		{"@(TITLE(string1, string2))", "", true},

		{"@(1 = asdf)", "", true},

		{"@(split(words, \" \").0)", "one", false},
		{"@(split(words, \" \")[1])", "two", false},
		{"@(split(words, \" \")[-1])", "three", false},
	}

	for _, test := range evaluateTests {
		eval, err := EvaluateTemplate(env, vars, test.template, keys)

		if test.hasError {
			assert.Error(t, err, "expected error evaluating template '%s'", test.template)
			continue
		} else {
			assert.NoError(t, err, "unexpected error evaluating template '%s'", test.template)
		}

		// first try reflect comparison
		equal := reflect.DeepEqual(eval, test.expected)

		// back down to our equality
		if !equal {
			cmp, err := utils.Compare(env, eval, test.expected)
			if err != nil {
				t.Errorf("Actual '%#v' does not match expected '%#v' evaluating template: '%s'", eval, test.expected, test.template)
			}
			equal = cmp == 0
		}

		if !equal {
			t.Errorf("Actual '%#v' does not match expected '%#v' evaluating template: '%s'", eval, test.expected, test.template)
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
