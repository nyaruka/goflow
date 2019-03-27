package tests_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/routers/tests"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var xs = types.NewXText
var xn = types.RequireXNumberFromString
var xi = types.NewXNumberFromInt
var xd = types.NewXDateTime
var xt = types.NewXTime

type testResolvable struct{}

func (r *testResolvable) Resolve(env utils.Environment, key string) types.XValue {
	switch strings.ToLower(key) {
	case "foo":
		return types.NewXText("bar")
	case "zed":
		return types.NewXNumberFromInt(123)
	case "missing":
		return nil
	default:
		return types.NewXResolveError(r, key)
	}
}

func (r *testResolvable) Describe() string { return "test" }

// Reduce is called when this object needs to be reduced to a primitive
func (r *testResolvable) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText("hello")
}

// ToXJSON is called when this type is passed to @(json(...))
func (r *testResolvable) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, r, "foo", "zed").ToXJSON(env)
}

var kgl, _ = time.LoadLocation("Africa/Kigali")

var testTests = []struct {
	name     string
	args     []types.XValue
	match    types.XValue
	hasError bool
}{
	{"is_error", []types.XValue{xs("hello")}, nil, false},
	{"is_error", []types.XValue{nil}, nil, false},
	{"is_error", []types.XValue{types.NewXErrorf("I am error")}, types.NewXErrorf("I am error"), false},
	{"is_error", []types.XValue{}, nil, true},

	{"has_text", []types.XValue{xs("hello")}, xs("hello"), false},
	{"has_text", []types.XValue{xs("  ")}, nil, false},
	{"has_text", []types.XValue{nil}, nil, false},
	{"has_text", []types.XValue{xs("one"), xs("two")}, nil, true},

	{"has_beginning", []types.XValue{xs("hello"), xs("hell")}, xs("hell"), false},
	{"has_beginning", []types.XValue{xs("  HelloThere"), xs("hello")}, xs("Hello"), false},
	{"has_beginning", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},
	{"has_beginning", []types.XValue{nil, xs("hell")}, nil, false},
	{"has_beginning", []types.XValue{xs("hello"), nil}, nil, false},
	{"has_beginning", []types.XValue{xs(""), xs("hello")}, nil, false},
	{"has_beginning", []types.XValue{xs("hel"), xs("hello")}, nil, false},

	{"has_any_word", []types.XValue{xs("this.is.my.word"), xs("WORD word2 word")}, xs("word"), false},
	{"has_any_word", []types.XValue{xs("this.is.my.βήτα"), xs("βήτα")}, xs("βήτα"), false},
	{"has_any_word", []types.XValue{xs("I say to you📴"), xs("📴")}, xs("📴"), false},
	{"has_any_word", []types.XValue{xs("this World too"), xs("world")}, xs("World"), false},
	{"has_any_word", []types.XValue{xs("BUT not this one"), xs("world")}, nil, false},
	{"has_any_word", []types.XValue{xs(""), xs("world")}, nil, false},
	{"has_any_word", []types.XValue{xs("world"), xs("foo")}, nil, false},
	{"has_any_word", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},
	{"has_any_word", []types.XValue{xs("but foo"), nil}, nil, false},
	{"has_any_word", []types.XValue{nil, xs("but foo")}, nil, false},

	{"has_all_words", []types.XValue{xs("this.is.my.word"), xs("WORD word")}, xs("word"), false},
	{"has_all_words", []types.XValue{xs("this World too"), xs("world too")}, xs("World too"), false},
	{"has_all_words", []types.XValue{xs("BUT not this one"), xs("world")}, nil, false},
	{"has_all_words", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},

	{"has_phrase", []types.XValue{xs("you Must resist"), xs("must resist")}, xs("Must resist"), false},
	{"has_phrase", []types.XValue{xs("this world Too"), xs("world too")}, xs("world Too"), false},
	{"has_phrase", []types.XValue{xs("this world Too"), xs("")}, xs(""), false},
	{"has_phrase", []types.XValue{xs("this is not world"), xs("this world")}, nil, false},
	{"has_phrase", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},

	{"has_only_phrase", []types.XValue{xs("Must resist"), xs("must resist")}, xs("Must resist"), false},
	{"has_only_phrase", []types.XValue{xs(" world Too "), xs("world too")}, xs("world Too"), false},
	{"has_only_phrase", []types.XValue{xs("this world Too"), xs("")}, nil, false},
	{"has_only_phrase", []types.XValue{xs(""), xs("")}, xs(""), false},
	{"has_only_phrase", []types.XValue{xs("this world is my world"), xs("this world")}, nil, false},
	{"has_only_phrase", []types.XValue{xs("this world"), xs("this mighty")}, nil, false},
	{"has_only_phrase", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},

	{"has_beginning", []types.XValue{xs("Must resist"), xs("must resist")}, xs("Must resist"), false},
	{"has_beginning", []types.XValue{xs(" 2061212"), xs("206")}, xs("206"), false},
	{"has_beginning", []types.XValue{xs(" world Too foo"), xs("world too")}, xs("world Too"), false},
	{"has_beginning", []types.XValue{xs("but this world"), xs("this world")}, nil, false},
	{"has_beginning", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},

	{"has_pattern", []types.XValue{xs("<html>x</html>"), xs(`<\w+>`)}, xs("<html>"), false},
	{"has_pattern", []types.XValue{xs("<html>x</html>"), xs(`HTML`)}, xs("html"), false},
	{"has_pattern", []types.XValue{xs("<html>x</html>"), xs(`(?-i)HTML`)}, nil, false},
	{"has_pattern", []types.XValue{xs("12345"), xs(`\A\d{5}\z`)}, xs("12345"), false},
	{"has_pattern", []types.XValue{xs("12345 "), xs(`\A\d{5}\z`)}, nil, false},
	{"has_pattern", []types.XValue{xs(" 12345"), xs(`\A\d{5}\z`)}, nil, false},
	{"has_pattern", []types.XValue{xs("<html>x</html>"), xs(`[`)}, nil, true},

	{"has_number", []types.XValue{xs("the number 10")}, xn("10"), false},
	{"has_number", []types.XValue{xs("the number -10")}, xn("-10"), false},
	{"has_number", []types.XValue{xs("1-15")}, xn("1"), false},
	{"has_number", []types.XValue{xs("24ans")}, xn("24"), false},
	{"has_number", []types.XValue{xs("J'AI 20ANS")}, xn("20"), false},
	{"has_number", []types.XValue{xs("1,000,000")}, xn("1000000"), false},
	{"has_number", []types.XValue{xs("the number 10")}, xn("10"), false},
	{"has_number", []types.XValue{xs("O número é 500")}, xn("500"), false},
	{"has_number", []types.XValue{xs("another is -12.51")}, xn("-12.51"), false},
	{"has_number", []types.XValue{xs("hi.51")}, xn("51"), false},
	{"has_number", []types.XValue{xs("nothing here")}, nil, false},
	{"has_number", []types.XValue{xs("1OO l00")}, nil, false}, // no longer do substitutions
	{"has_number", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},

	{"has_number_lt", []types.XValue{xs("the number 10"), xs("11")}, xn("10"), false},
	{"has_number_lt", []types.XValue{xs("another is -12.51"), xs("12")}, xn("-12.51"), false},
	{"has_number_lt", []types.XValue{xs("nothing here"), xs("12")}, nil, false},
	{"has_number_lt", []types.XValue{xs("too big 15"), xs("12")}, nil, false},
	{"has_number_lt", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},
	{"has_number_lt", []types.XValue{xs("but foo"), nil}, nil, true},
	{"has_number_lt", []types.XValue{nil, xs("but foo")}, nil, true},

	{"has_number_lte", []types.XValue{xs("the number 10"), xs("11")}, xn("10"), false},
	{"has_number_lte", []types.XValue{xs("another is -12.51"), xs("12")}, xn("-12.51"), false},
	{"has_number_lte", []types.XValue{xs("nothing here"), xs("12")}, nil, false},
	{"has_number_lte", []types.XValue{xs("too big 15"), xs("12")}, nil, false},
	{"has_number_lte", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},

	{"has_number_eq", []types.XValue{xs("the number 10"), xs("10")}, xn("10"), false},
	{"has_number_eq", []types.XValue{xs("another is -12.51"), xs("-12.51")}, xn("-12.51"), false},
	{"has_number_eq", []types.XValue{xs("nothing here"), xs("12")}, nil, false},
	{"has_number_eq", []types.XValue{xs("wrong .51"), xs(".61")}, nil, false},
	{"has_number_eq", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},

	{"has_number_gte", []types.XValue{xs("the number 10"), xs("9")}, xn("10"), false},
	{"has_number_gte", []types.XValue{xs("another is -12.51"), xs("-13")}, xn("-12.51"), false},
	{"has_number_gte", []types.XValue{xs("nothing here"), xs("12")}, nil, false},
	{"has_number_gte", []types.XValue{xs("too small -12"), xs("-11")}, nil, false},
	{"has_number_gte", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},

	{"has_number_gt", []types.XValue{xs("the number 10"), xs("9")}, xn("10"), false},
	{"has_number_gt", []types.XValue{xs("another is -12.51"), xs("-13")}, xn("-12.51"), false},
	{"has_number_gt", []types.XValue{xs("nothing here"), xs("12")}, nil, false},
	{"has_number_gt", []types.XValue{xs("not great -12.51"), xs("-12.51")}, nil, false},
	{"has_number_gt", []types.XValue{xs("one"), xs("two"), xs("three")}, nil, true},

	{"has_number_between", []types.XValue{xs("the number 10"), xs("8"), xs("12")}, xn("10"), false},
	{"has_number_between", []types.XValue{xs("24ans"), xn("20"), xn("24")}, xn("24"), false},
	{"has_number_between", []types.XValue{xs("another is -12.51"), xs("-12.51"), xs("-10")}, xn("-12.51"), false},
	{"has_number_between", []types.XValue{xs("nothing here"), xs("10"), xs("15")}, nil, false},
	{"has_number_between", []types.XValue{xs("one"), xs("two")}, nil, true},
	{"has_number_between", []types.XValue{xs("but foo"), nil, xs("10")}, nil, true},
	{"has_number_between", []types.XValue{nil, xs("but foo"), xs("10")}, nil, true},
	{"has_number_between", []types.XValue{xs("a string"), xs("10"), xs("not number")}, nil, true},

	{"has_date", []types.XValue{xs("last date was 1.10.2017")}, xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)), false},
	{"has_date", []types.XValue{xs("last date was 1.10.99")}, xd(time.Date(1999, 10, 1, 15, 24, 30, 123456000, kgl)), false},
	{"has_date", []types.XValue{xs("this isn't a valid date 33.2.99")}, nil, false},
	{"has_date", []types.XValue{xs("no date at all")}, nil, false},
	{"has_date", []types.XValue{xs("too"), xs("many"), xs("args")}, nil, true},

	{"has_date_lt", []types.XValue{xs("last date was 1.10.2017"), xs("3.10.2017")}, xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)), false},
	{"has_date_lt", []types.XValue{xs("last date was 1.10.99"), xs("3.10.98")}, nil, false},
	{"has_date_lt", []types.XValue{xs("no date at all"), xs("3.10.98")}, nil, false},
	{"has_date_lt", []types.XValue{xs("too"), xs("many"), xs("args")}, nil, true},
	{"has_date_lt", []types.XValue{xs("last date was 1.10.2017"), nil}, nil, true},
	{"has_date_lt", []types.XValue{nil, xs("but foo")}, nil, true},

	{"has_date_eq", []types.XValue{xs("last date was 1.10.2017"), xs("1.10.2017")}, xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)), false},
	{"has_date_eq", []types.XValue{xs("last date was 1.10.99"), xs("3.10.98")}, nil, false},
	{"has_date_eq", []types.XValue{xs("2017-10-01T23:55:55.123456+02:00"), xs("1.10.2017")}, xd(time.Date(2017, 10, 1, 23, 55, 55, 123456000, kgl)), false},
	{"has_date_eq", []types.XValue{xs("2017-10-01T23:55:55.123456+01:00"), xs("1.10.2017")}, nil, false}, // would have been 2017-10-02 in env timezone
	{"has_date_eq", []types.XValue{xs("no date at all"), xs("3.10.98")}, nil, false},
	{"has_date_eq", []types.XValue{xs("too"), xs("many"), xs("args")}, nil, true},

	{"has_date_gt", []types.XValue{xs("last date was 1.10.2017"), xs("3.10.2016")}, xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)), false},
	{"has_date_gt", []types.XValue{xs("last date was 1.10.99"), xs("3.10.01")}, nil, false},
	{"has_date_gt", []types.XValue{xs("no date at all"), xs("3.10.98")}, nil, false},
	{"has_date_gt", []types.XValue{xs("too"), xs("many"), xs("args")}, nil, true},

	{"has_time", []types.XValue{xs("last time was 10:30")}, xt(utils.NewTimeOfDay(10, 30, 0, 0)), false},
	{"has_time", []types.XValue{xs("this isn't a valid time 59:77")}, nil, false},
	{"has_time", []types.XValue{xs("no time at all")}, nil, false},
	{"has_time", []types.XValue{xs("too"), xs("many"), xs("args")}, nil, true},

	{"has_email", []types.XValue{xs("my email is foo@bar.com.")}, xs("foo@bar.com"), false},
	{"has_email", []types.XValue{xs("my email is <foo1@bar-2.com>")}, xs("foo1@bar-2.com"), false},
	{"has_email", []types.XValue{xs("FOO@bar.whatzit")}, xs("FOO@bar.whatzit"), false},
	{"has_email", []types.XValue{xs("FOO@βήτα.whatzit")}, xs("FOO@βήτα.whatzit"), false},
	{"has_email", []types.XValue{xs("email is foo @ bar . com")}, nil, false},
	{"has_email", []types.XValue{xs("email is foo@bar")}, nil, false},
	{"has_email", []types.XValue{nil}, nil, false},
	{"has_email", []types.XValue{xs("too"), xs("many"), xs("args")}, nil, true},

	{"has_phone", []types.XValue{xs("my number is +250788123123")}, xs("+250788123123"), false},
	{"has_phone", []types.XValue{xs("my number is +593979111111")}, xs("+593979111111"), false},
	{"has_phone", []types.XValue{xs("my number is 0788123123")}, xs("+250788123123"), false}, // uses environment default
	{"has_phone", []types.XValue{xs("my number is 0788123123"), xs("RW")}, xs("+250788123123"), false},
	{"has_phone", []types.XValue{xs("my number is +250788123123"), xs("RW")}, xs("+250788123123"), false},
	{"has_phone", []types.XValue{xs("my number is +12065551212"), xs("RW")}, xs("+12065551212"), false},
	{"has_phone", []types.XValue{xs("my number is 12065551212"), xs("US")}, xs("+12065551212"), false},
	{"has_phone", []types.XValue{xs("my number is 206 555 1212"), xs("US")}, xs("+12065551212"), false},
	{"has_phone", []types.XValue{xs("my number is +10001112222"), xs("US")}, xs("+10001112222"), false},
	{"has_phone", []types.XValue{xs("my number is 10000"), xs("US")}, nil, false},
	{"has_phone", []types.XValue{xs("my number is 12067799294"), xs("BW")}, nil, false},
	{"has_phone", []types.XValue{xs("my number is none of your business"), xs("US")}, nil, false},
	{"has_phone", []types.XValue{}, nil, true},
	{"has_phone", []types.XValue{types.NewXErrorf("error")}, nil, true},
	{"has_phone", []types.XValue{xs("3245"), types.NewXErrorf("error")}, nil, true},
	{"has_phone", []types.XValue{xs("number"), nil}, nil, false},
	{"has_phone", []types.XValue{xs("too"), xs("many"), xs("args")}, nil, true},
}

func TestTests(t *testing.T) {
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)))
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	env := utils.NewEnvironmentBuilder().
		WithDateFormat(utils.DateFormatDayMonthYear).
		WithTimeFormat(utils.TimeFormatHourMinuteSecond).
		WithTimezone(kgl).
		WithDefaultCountry(utils.Country("RW")).
		Build()

	for _, test := range testTests {
		testFunc := tests.XTESTS[test.name]

		result := testFunc(env, test.args...)
		err, _ := result.(error)

		if test.hasError {
			assert.Error(t, err, "expected error running test %s(%#v)", test.name, test.args)
		} else {
			assert.NoError(t, err, "unexpected error running test %s(%#v): %v", test.name, test.args, err)

			// otherwise, cast to our result
			testResult, _ := result.(tests.XTestResult)

			if test.match == nil {
				assert.Nil(t, testResult.Match(), "unexpected match value for test %s(%#v)", test.name, test.args)
			} else {
				if !reflect.DeepEqual(testResult.Match(), test.match) {
					assert.Fail(t, "", "unexpected match value, expected '%s', got '%s' for test %s(%#v)", test.match, testResult.Match(), test.name, test.args)
				}
			}
		}
	}
}

func TestEvaluateTemplate(t *testing.T) {
	vars := types.NewXDict(map[string]types.XValue{
		"int1":  types.NewXNumberFromInt(1),
		"int2":  types.NewXNumberFromInt(2),
		"array": types.NewXArray(xs("one"), xs("two"), xs("three")),
		"thing": &testResolvable{},
		"err":   types.NewXErrorf("an error"),
	})

	evalTests := []struct {
		template string
		expected string
		hasError bool
	}{
		{"@(is_error(array[100]))", "true", false}, // errors are like any other value
		{"@(is_error(array.100))", "true", false},
		{`@(is_error(round("foo", "bar")))`, "true", false},
		{`@(is_error(err))`, "true", false},
		{"@(is_error(thing.foo))", "", false},
		{"@(is_error(thing.xxx))", "true", false},
		{"@(is_error(1 / 0))", "true", false},
	}

	env := utils.NewEnvironmentBuilder().Build()
	for _, test := range evalTests {
		eval, err := excellent.EvaluateTemplate(env, vars, test.template, vars.Keys())

		if test.hasError {
			assert.Error(t, err, "expected error evaluating template '%s'", test.template)
		} else {
			assert.NoError(t, err, "unexpected error evaluating template '%s'", test.template)

			assert.Equal(t, test.expected, eval, "actual '%s' does not match expected '%s' evaluating template: '%s'", eval, test.expected, test.template)
		}
	}
}

func TestParseDecimalFuzzy(t *testing.T) {
	parseTests := []struct {
		input    string
		expected decimal.Decimal
		format   *utils.NumberFormat
	}{
		{"1234", decimal.RequireFromString("1234"), utils.DefaultNumberFormat},
		{"1,234.567", decimal.RequireFromString("1234.567"), utils.DefaultNumberFormat},
		{"1.234,567", decimal.RequireFromString("1234.567"), &utils.NumberFormat{DecimalSymbol: ",", DigitGroupingSymbol: "."}},
		{"100.00", decimal.RequireFromString("100.00"), utils.DefaultNumberFormat},
	}

	for _, test := range parseTests {
		val, err := tests.ParseDecimalFuzzy(test.input, test.format)

		assert.NoError(t, err)
		assert.Equal(t, test.expected, val, "parse decimal failed for input '%s'", test.input)
	}
}
