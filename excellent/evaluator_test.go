package excellent

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/utils"
)

type TestVars struct {
	vars map[string]interface{}
}

func (v *TestVars) Resolve(key string) interface{} {
	val, present := v.vars[key]
	if !present {
		return fmt.Errorf("No such key: %s", key)
	}
	return val
}

func (v *TestVars) Default() interface{} {
	return v.vars
}

func TestEvaluateTemplateAsString(t *testing.T) {
	varMap := make(map[string]interface{})
	varMap["string1"] = "foo"
	varMap["string2"] = "bar"
	varMap["汉字"] = "simplified chinese"
	varMap["int1"] = 1
	varMap["int2"] = 2
	varMap["dec1"] = 1.5
	varMap["dec2"] = 2.5
	varMap["words"] = "one two three"
	varMap["array"] = []string{"one", "two", "three"}
	vars := &TestVars{varMap}

	evaluateAsStringTests := []struct {
		template string
		expected string
		hasError bool
	}{

		{"hello world", "hello world", false},
		{"@(\"hello\\nworld\")", "hello\nworld", false},
		{"@(\"hello😁world\")", "hello😁world", false},
		{"@(\"hello\\U0001F601world\")", "hello😁world", false},
		{"@hello", "@hello", false},
		{"@hello.bar", "@hello.bar", false},
		{"@(title(\"hello\"))", "Hello", false},
		{"@(title(hello))", "", true},
		{"Hello @(title(string1))", "Hello Foo", false},
		{"Hello @@string1", "Hello @string1", false},
		{"My email is foo@bar.com", "My email is foo@bar.com", false},

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

		{"@missing", "@missing", false},
		{"@(TITLE(missing))", "", true},

		{"@array", "one, two, three", false},
		{"@array[0]", "one, two, three[0]", false}, // [n] notation not supported outside expression
		{"@array.0", "one", false},                 // works as dot notation however
		{"@(array [0])", "one", false},
		{"@(array[0])", "one", false},
		{"@(array.0)", "one", false},

		{"@(array[-1])", "three", false},
		{"@(array.-1)", "", true},

		{"@(split(words, \" \").0)", "one", false},
		{"@(split(words, \" \")[1])", "two", false},
		{"@(split(words, \" \")[-1])", "three", false},
	}

	env := utils.NewDefaultEnvironment()
	for _, test := range evaluateAsStringTests {
		eval, err := EvaluateTemplateAsString(env, vars, test.template)
		if err != nil {
			if !test.hasError {
				t.Errorf("Received error evaluating '%s': %s", test.template, err)
			}
		} else {
			if test.hasError {
				t.Errorf("Did not receive error evaluating '%s'", test.template)
			}
		}

		if eval != test.expected {
			t.Errorf("Actual '%s' does not match expected '%s' evaluating template: '%s'", eval, test.expected, test.template)
		}
	}
}

func TestEvaluateTemplate(t *testing.T) {
	arr := []string{"a", "b", "c"}

	strMap := make(map[string]string)
	strMap["1"] = "one"
	strMap["2"] = "two"
	strMap["3"] = "three"
	strMap["four"] = "four"
	strMap["with space"] = "spacy"
	strMap["with-dash"] = "dashy"
	strMap["汉字"] = "simplified chinese"

	intMap := make(map[int]string)
	intMap[1] = "one"
	intMap[2] = "two"
	intMap[3] = "three"

	innerMap := make(map[string]interface{})
	innerMap["int_map"] = intMap

	innerArr := []map[string]string{strMap}

	varMap := make(map[string]interface{})
	varMap["string1"] = "foo"
	varMap["string2"] = "bar"
	varMap["key"] = "four"
	varMap["int1"] = 1
	varMap["int2"] = 2
	varMap["dec1"] = 1.5
	varMap["dec2"] = 2.5
	varMap["words"] = "one two three"
	varMap["array1"] = arr
	varMap["str_map"] = strMap
	varMap["int_map"] = intMap
	varMap["inner_map"] = innerMap
	varMap["inner_arr"] = innerArr
	vars := &TestVars{varMap}

	env := utils.NewDefaultEnvironment()

	evaluateTests := []struct {
		template string
		expected interface{}
		hasError bool
	}{
		{"hello world", "hello world", false},
		{"@hello", "@hello", true},
		{"@(title(\"hello\"))", "Hello", false},

		{"@dec1", 1.5, false},
		{"@(dec1 + dec2)", 4, false},
		{"@array1", arr, false},
		{"@str_map", strMap, false},
		{"@int_map", intMap, false},
		{"@int_map.1", "one", false},
		{"@str_map.1", "one", false},
		{"@(str_map[1])", "one", false},
		{"@(str_map[10])", nil, false},
		{"@(str_map.汉字)", "simplified chinese", false},
		{"@(int_map[1])", "one", false},
		{"@(int_map[10])", nil, false},
		{"@(str_map[\"four\"])", "four", false},
		{"@(str_map[key])", "four", false},
		{"@(str_map[lower(key)])", "four", false},
		{"@(title(missing))", "", true},
		{`@(str_map["with-dash"])`, "dashy", false},
		{`@(str_map["with space"])`, "spacy", false},
		{`@(inner_map["int_map"].1)`, `one`, false},
		{`@(inner_map.int_map.1)`, `one`, false},
		{`@(inner_arr[0].four)`, `four`, false},
		{`@(inner_arr[0].0)`, nil, false},
		{`@(inner_arr[0].1)`, `one`, false},

		{"@string1 world", "foo world", false},

		{"@(-10)", -10, false},
		{"@(-asdf)", "", true},

		{"@(2^2)", 4, false},
		{"@(2^asdf)", "", true},
		{"@(asdf^2)", "", true},

		{"@(1+2)", 3, false},
		{"@(1-2.5)", -1.5, false},
		{"@(1-asdf)", "", true},
		{"@(asdf+1)", "", true},

		{"@(1*2)", 2, false},
		{"@(1/2)", .5, false},
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
		{"@(1 != 2)", true, false},
		{"@(1 != 1)", false, false},
		{"@(-1 = 1)", false, false},
		{"@(1 < asdf)", "", true},

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
		eval, err := EvaluateTemplate(env, vars, test.template)
		if err != nil {
			if !test.hasError {
				t.Errorf("Received error evaluating '%s': %s", test.template, err)
			}
		} else {
			if test.hasError {
				t.Errorf("Did not receive error evaluating '%s'", test.template)
			}
		}

		if test.hasError {
			continue
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
	scanner := NewXScanner(strings.NewReader("12"))

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
