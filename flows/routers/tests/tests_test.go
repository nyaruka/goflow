package tests_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows/routers/tests"
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

func newDecimal(val string) decimal.Decimal {
	dec, _ := decimal.NewFromString(val)
	return dec
}

// noStr is used to blow up our type conversions in the tests below
type noStr struct {
}

var testTests = []struct {
	name     string
	args     []interface{}
	matched  bool
	match    interface{}
	hasError bool
}{
	{"is_error", []interface{}{"hello"}, false, nil, false},
	{"is_error", []interface{}{nil}, false, nil, false},
	{"is_error", []interface{}{fmt.Errorf("I am error")}, true, fmt.Errorf("I am error"), false},
	{"is_error", []interface{}{}, false, nil, true},

	{"has_text", []interface{}{"hello"}, true, "hello", false},
	{"has_text", []interface{}{"  "}, false, nil, false},
	{"has_text", []interface{}{"one", "two"}, false, nil, true},
	{"has_text", []interface{}{noStr{}}, false, nil, true},

	{"has_beginning", []interface{}{"hello", "hell"}, true, "hell", false},
	{"has_beginning", []interface{}{"  HelloThere", "hello"}, true, "Hello", false},
	{"has_beginning", []interface{}{"one", "two", "three"}, false, nil, true},
	{"has_beginning", []interface{}{noStr{}, "hell"}, false, nil, true},
	{"has_beginning", []interface{}{"hello", noStr{}}, false, nil, true},
	{"has_beginning", []interface{}{"", "hello"}, false, nil, false},
	{"has_beginning", []interface{}{"hel", "hello"}, false, nil, false},

	{"has_any_word", []interface{}{"this.is.my.word", "WORD word2 word"}, true, "word", false},
	{"has_any_word", []interface{}{"this.is.my.βήτα", "βήτα"}, true, "βήτα", false},
	{"has_any_word", []interface{}{"I say to you📴", "📴"}, true, "📴", false},
	{"has_any_word", []interface{}{"this World too", "world"}, true, "World", false},
	{"has_any_word", []interface{}{"BUT not this one", "world"}, false, nil, false},
	{"has_any_word", []interface{}{"", "world"}, false, nil, false},
	{"has_any_word", []interface{}{"world", "foo"}, false, nil, false},
	{"has_any_word", []interface{}{"one", "two", "three"}, false, nil, true},
	{"has_any_word", []interface{}{"but foo", noStr{}}, false, nil, true},
	{"has_any_word", []interface{}{noStr{}, "but foo"}, false, nil, true},

	{"has_all_words", []interface{}{"this.is.my.word", "WORD word"}, true, "word", false},
	{"has_all_words", []interface{}{"this World too", "world too"}, true, "World too", false},
	{"has_all_words", []interface{}{"BUT not this one", "world"}, false, nil, false},
	{"has_all_words", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_phrase", []interface{}{"you Must resist", "must resist"}, true, "Must resist", false},
	{"has_phrase", []interface{}{"this world Too", "world too"}, true, "world Too", false},
	{"has_phrase", []interface{}{"this world Too", ""}, true, "", false},
	{"has_phrase", []interface{}{"this is not world", "this world"}, false, nil, false},
	{"has_phrase", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_only_phrase", []interface{}{"Must resist", "must resist"}, true, "Must resist", false},
	{"has_only_phrase", []interface{}{" world Too ", "world too"}, true, "world Too", false},
	{"has_only_phrase", []interface{}{"this world Too", ""}, false, nil, false},
	{"has_only_phrase", []interface{}{"", ""}, true, "", false},
	{"has_only_phrase", []interface{}{"this world is my world", "this world"}, false, nil, false},
	{"has_only_phrase", []interface{}{"this world", "this mighty"}, false, nil, false},
	{"has_only_phrase", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_beginning", []interface{}{"Must resist", "must resist"}, true, "Must resist", false},
	{"has_beginning", []interface{}{" 2061212", "206"}, true, "206", false},
	{"has_beginning", []interface{}{" world Too foo", "world too"}, true, "world Too", false},
	{"has_beginning", []interface{}{"but this world", "this world"}, false, nil, false},
	{"has_beginning", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_number", []interface{}{"the number 10"}, true, newDecimal("10"), false},
	{"has_number", []interface{}{"the number 1o"}, true, newDecimal("10"), false},
	{"has_number", []interface{}{"the number lo"}, true, newDecimal("10"), false},
	{"has_number", []interface{}{"another is -12.51"}, true, newDecimal("-12.51"), false},
	{"has_number", []interface{}{".51"}, true, newDecimal(".51"), false},
	{"has_number", []interface{}{"nothing here"}, false, nil, false},
	{"has_number", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_number_lt", []interface{}{"the number 10", "11"}, true, newDecimal("10"), false},
	{"has_number_lt", []interface{}{"another is -12.51", "12"}, true, newDecimal("-12.51"), false},
	{"has_number_lt", []interface{}{"nothing here", "12"}, false, nil, false},
	{"has_number_lt", []interface{}{"too big 15", "12"}, false, nil, false},
	{"has_number_lt", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_number_lt", []interface{}{"but foo", noStr{}}, false, nil, true},
	{"has_number_lt", []interface{}{noStr{}, "but foo"}, false, nil, true},

	{"has_number_lte", []interface{}{"the number 10", "11"}, true, newDecimal("10"), false},
	{"has_number_lte", []interface{}{"another is -12.51", "12"}, true, newDecimal("-12.51"), false},
	{"has_number_lte", []interface{}{"nothing here", "12"}, false, nil, false},
	{"has_number_lte", []interface{}{"too big 15", "12"}, false, nil, false},
	{"has_number_lte", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_number_eq", []interface{}{"the number 10", "10"}, true, newDecimal("10"), false},
	{"has_number_eq", []interface{}{"another is -12.51", "-12.51"}, true, newDecimal("-12.51"), false},
	{"has_number_eq", []interface{}{"nothing here", "12"}, false, nil, false},
	{"has_number_eq", []interface{}{"wrong .51", ".61"}, false, nil, false},
	{"has_number_eq", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_number_gte", []interface{}{"the number 10", "9"}, true, newDecimal("10"), false},
	{"has_number_gte", []interface{}{"another is -12.51", "-13"}, true, newDecimal("-12.51"), false},
	{"has_number_gte", []interface{}{"nothing here", "12"}, false, nil, false},
	{"has_number_gte", []interface{}{"too small -12", "-11"}, false, nil, false},
	{"has_number_gte", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_number_gt", []interface{}{"the number 10", "9"}, true, newDecimal("10"), false},
	{"has_number_gt", []interface{}{"another is -12.51", "-13"}, true, newDecimal("-12.51"), false},
	{"has_number_gt", []interface{}{"nothing here", "12"}, false, nil, false},
	{"has_number_gt", []interface{}{"not great -12.51", "-12.51"}, false, nil, false},
	{"has_number_gt", []interface{}{"one", "two", "three"}, false, nil, true},

	{"has_number_between", []interface{}{"the number 10", "8", "12"}, true, newDecimal("10"), false},
	{"has_number_between", []interface{}{"another is -12.51", "-12.51", "-10"}, true, newDecimal("-12.51"), false},
	{"has_number_between", []interface{}{"nothing here", "10", "15"}, false, nil, false},
	{"has_number_between", []interface{}{"one", "two"}, false, nil, true},

	{"has_number_between", []interface{}{"but foo", noStr{}, "10"}, false, nil, true},
	{"has_number_between", []interface{}{noStr{}, "but foo", "10"}, false, nil, true},
	{"has_number_between", []interface{}{"a string", "10", "not number"}, false, nil, true},

	{"has_date", []interface{}{"last date was 1.10.2017"}, true, time.Date(2017, 10, 1, 0, 0, 0, 0, time.UTC), false},
	{"has_date", []interface{}{"last date was 1.10.99"}, true, time.Date(1999, 10, 1, 0, 0, 0, 0, time.UTC), false},
	{"has_date", []interface{}{"this isn't a valid date 33.2.99"}, false, nil, false},
	{"has_date", []interface{}{"no date at all"}, false, nil, false},
	{"has_date", []interface{}{"too many", "args"}, false, nil, true},

	{"has_date_lt", []interface{}{"last date was 1.10.2017", "3.10.2017"}, true, time.Date(2017, 10, 1, 0, 0, 0, 0, time.UTC), false},
	{"has_date_lt", []interface{}{"last date was 1.10.99", "3.10.98"}, false, nil, false},
	{"has_date_lt", []interface{}{"no date at all", "3.10.98"}, false, nil, false},
	{"has_date_lt", []interface{}{"too", "many", "args"}, false, nil, true},

	{"has_date_lt", []interface{}{"last date was 1.10.2017", noStr{}}, false, nil, true},
	{"has_date_lt", []interface{}{noStr{}, "but foo"}, false, nil, true},

	{"has_date_eq", []interface{}{"last date was 1.10.2017", "1.10.2017"}, true, time.Date(2017, 10, 1, 0, 0, 0, 0, time.UTC), false},
	{"has_date_eq", []interface{}{"last date was 1.10.99", "3.10.98"}, false, nil, false},
	{"has_date_eq", []interface{}{"no date at all", "3.10.98"}, false, nil, false},
	{"has_date_eq", []interface{}{"too", "many", "args"}, false, nil, true},

	{"has_date_gt", []interface{}{"last date was 1.10.2017", "3.10.2016"}, true, time.Date(2017, 10, 1, 0, 0, 0, 0, time.UTC), false},
	{"has_date_gt", []interface{}{"last date was 1.10.99", "3.10.01"}, false, nil, false},
	{"has_date_gt", []interface{}{"no date at all", "3.10.98"}, false, nil, false},
	{"has_date_gt", []interface{}{"too", "many", "args"}, false, nil, true},

	{"has_email", []interface{}{"my email is foo@bar.com."}, true, "foo@bar.com", false},
	{"has_email", []interface{}{"my email is <foo1@bar-2.com>"}, true, "foo1@bar-2.com", false},
	{"has_email", []interface{}{"FOO@bar.whatzit"}, true, "FOO@bar.whatzit", false},
	{"has_email", []interface{}{"FOO@βήτα.whatzit"}, true, "FOO@βήτα.whatzit", false},
	{"has_email", []interface{}{"email is foo @ bar . com"}, false, nil, false},
	{"has_email", []interface{}{"email is foo@bar"}, false, nil, false},
	{"has_email", []interface{}{noStr{}}, false, nil, true},
	{"has_email", []interface{}{"too", "many", "args"}, false, nil, true},

	{"has_phone", []interface{}{"my number is 0788123123", "RW"}, true, "+250788123123", false},
	{"has_phone", []interface{}{"my number is +250788123123", "RW"}, true, "+250788123123", false},
	{"has_phone", []interface{}{"my number is +12065551212", "RW"}, true, "+12065551212", false},
	{"has_phone", []interface{}{"my number is 12065551212", "US"}, true, "+12065551212", false},
	{"has_phone", []interface{}{"my number is 206 555 1212", "US"}, true, "+12065551212", false},
	{"has_phone", []interface{}{"my number is none of your business", "US"}, false, nil, false},
	{"has_phone", []interface{}{noStr{}}, false, nil, true},
	{"has_phone", []interface{}{"number", noStr{}}, false, nil, true},
	{"has_phone", []interface{}{"too", "many", "args"}, false, nil, true},
}

func TestTests(t *testing.T) {
	env := utils.NewEnvironment(utils.DateFormatDayMonthYear, utils.TimeFormatHourMinuteSecond, time.UTC, utils.LanguageList{})

	for _, test := range testTests {
		testFunc := tests.XTESTS[test.name]

		result := testFunc(env, test.args...)
		err, _ := result.(error)

		if test.hasError {
			assert.Error(t, err, "expected error running test %s(%#v)", test.name, test.args)
		} else {
			assert.NoError(t, err, "unexpected error running test %s(%#v): %v", test.name, test.args, err)

			// otherwise, cast to our result
			testResult := result.(tests.XTestResult)

			// check our expected match
			assert.Equal(t, test.matched, testResult.Matched(), "unexpected matched value: %v for test %s(%#v)", testResult.Matched(), test.name, test.args)

			// and the match itself
			if !reflect.DeepEqual(testResult.Match(), test.match) {
				assert.Fail(t, "Unexpected match value, expected '%s', got '%s' for test %s(%#v)", test.match, testResult.Match(), test.name, test.args)
			}
		}
	}
}

func TestEvaluateTemplateAsString(t *testing.T) {
	varMap := map[string]interface{}{
		"int1":  1,
		"int2":  2,
		"array": utils.NewArray("one", "two", "three"),
		"thing": &testResolvable{},
		"err":   fmt.Errorf("an error"),
	}
	vars := utils.NewMapResolver(varMap)

	keys := make([]string, 0, len(varMap))
	for key := range varMap {
		keys = append(keys, key)
	}

	evalTests := []struct {
		template string
		expected string
		hasError bool
	}{
		{"@(is_error(array[100]))", "true", false}, // errors are like any other value
		{"@(is_error(array.100))", "true", false},
		{`@(is_error(round("foo", "bar")))`, "true", false},
		{`@(is_error(err))`, "true", false},
		{"@(is_error(thing.foo))", "false", false},
		{"@(is_error(thing.xxx))", "true", false},
		{"@(is_error(1 / 0))", "true", false},
	}

	env := utils.NewDefaultEnvironment()
	for _, test := range evalTests {
		eval, err := excellent.EvaluateTemplateAsString(env, vars, test.template, false, keys)

		if test.hasError {
			assert.Error(t, err, "expected error evaluating template '%s'", test.template)
		} else {
			assert.NoError(t, err, "unexpected error evaluating template '%s'", test.template)

			assert.Equal(t, test.expected, eval, "actual '%s' does not match expected '%s' evaluating template: '%s'", eval, test.expected, test.template)
		}
	}
}
