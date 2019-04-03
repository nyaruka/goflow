package tests_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/routers/tests"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var xs = types.NewXText
var xn = types.RequireXNumberFromString
var xi = types.NewXNumberFromInt
var xd = types.NewXDateTime
var xt = types.NewXTime
var xa = types.NewXArray
var result = tests.NewTrueResult
var resultWithExtra = tests.NewTrueResultWithExtra
var ERROR = types.NewXErrorf("any error")

var kgl, _ = time.LoadLocation("Africa/Kigali")

var testTests = []struct {
	name     string
	args     []types.XValue
	expected types.XValue
}{
	{"is_error", []types.XValue{xs("hello")}, nil},
	{"is_error", []types.XValue{nil}, nil},
	{"is_error", []types.XValue{types.NewXErrorf("I am error")}, result(types.NewXErrorf("I am error"))},
	{"is_error", []types.XValue{}, ERROR},

	{"has_text", []types.XValue{xs("hello")}, result(xs("hello"))},
	{"has_text", []types.XValue{xs("  ")}, nil},
	{"has_text", []types.XValue{nil}, nil},
	{"has_text", []types.XValue{xs("one"), xs("two")}, ERROR},
	{"has_text", []types.XValue{ERROR}, ERROR},

	{"has_beginning", []types.XValue{xs("hello"), xs("hell")}, result(xs("hell"))},
	{"has_beginning", []types.XValue{xs("  HelloThere"), xs("hello")}, result(xs("Hello"))},
	{"has_beginning", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_beginning", []types.XValue{nil, xs("hell")}, nil},
	{"has_beginning", []types.XValue{xs("hello"), nil}, nil},
	{"has_beginning", []types.XValue{xs(""), xs("hello")}, nil},
	{"has_beginning", []types.XValue{xs("hel"), xs("hello")}, nil},
	{"has_beginning", []types.XValue{ERROR, ERROR}, ERROR},
	{"has_beginning", []types.XValue{}, ERROR},

	{"has_any_word", []types.XValue{xs("this.is.my.word"), xs("WORD word2 word")}, result(xs("word"))},
	{"has_any_word", []types.XValue{xs("this.is.my.βήτα"), xs("βήτα")}, result(xs("βήτα"))},
	{"has_any_word", []types.XValue{xs("I say to you📴"), xs("📴")}, result(xs("📴"))},
	{"has_any_word", []types.XValue{xs("this World too"), xs("world")}, result(xs("World"))},
	{"has_any_word", []types.XValue{xs("BUT not this one"), xs("world")}, nil},
	{"has_any_word", []types.XValue{xs(""), xs("world")}, nil},
	{"has_any_word", []types.XValue{xs("world"), xs("foo")}, nil},
	{"has_any_word", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_any_word", []types.XValue{xs("but foo"), nil}, nil},
	{"has_any_word", []types.XValue{nil, xs("but foo")}, nil},

	{"has_all_words", []types.XValue{xs("this.is.my.word"), xs("WORD word")}, result(xs("word"))},
	{"has_all_words", []types.XValue{xs("this World too"), xs("world too")}, result(xs("World too"))},
	{"has_all_words", []types.XValue{xs("BUT not this one"), xs("world")}, nil},
	{"has_all_words", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},

	{"has_phrase", []types.XValue{xs("you Must resist"), xs("must resist")}, result(xs("Must resist"))},
	{"has_phrase", []types.XValue{xs("this world Too"), xs("world too")}, result(xs("world Too"))},
	{"has_phrase", []types.XValue{xs("this world Too"), xs("")}, result(xs(""))},
	{"has_phrase", []types.XValue{xs("this is not world"), xs("this world")}, nil},
	{"has_phrase", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},

	{"has_only_phrase", []types.XValue{xs("Must resist"), xs("must resist")}, result(xs("Must resist"))},
	{"has_only_phrase", []types.XValue{xs(" world Too "), xs("world too")}, result(xs("world Too"))},
	{"has_only_phrase", []types.XValue{xs("this world Too"), xs("")}, nil},
	{"has_only_phrase", []types.XValue{xs(""), xs("")}, result(xs(""))},
	{"has_only_phrase", []types.XValue{xs("this world is my world"), xs("this world")}, nil},
	{"has_only_phrase", []types.XValue{xs("this world"), xs("this mighty")}, nil},
	{"has_only_phrase", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},

	{"has_beginning", []types.XValue{xs("Must resist"), xs("must resist")}, result(xs("Must resist"))},
	{"has_beginning", []types.XValue{xs(" 2061212"), xs("206")}, result(xs("206"))},
	{"has_beginning", []types.XValue{xs(" world Too foo"), xs("world too")}, result(xs("world Too"))},
	{"has_beginning", []types.XValue{xs("but this world"), xs("this world")}, nil},
	{"has_beginning", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},

	{"has_pattern", []types.XValue{xs("<html>x</html>"), xs(`<\w+>`)}, resultWithExtra(xs("<html>"), types.NewXDict(map[string]types.XValue{"0": xs("<html>")}))},
	{"has_pattern", []types.XValue{xs("<html>x</html>"), xs(`HTML`)}, resultWithExtra(xs("html"), types.NewXDict(map[string]types.XValue{"0": xs("html")}))},
	{"has_pattern", []types.XValue{xs("<html>x</html>"), xs(`(?-i)HTML`)}, nil},
	{"has_pattern", []types.XValue{xs("12345"), xs(`\A\d{5}\z`)}, resultWithExtra(xs("12345"), types.NewXDict(map[string]types.XValue{"0": xs("12345")}))},
	{"has_pattern", []types.XValue{xs("12345 "), xs(`\A\d{5}\z`)}, nil},
	{"has_pattern", []types.XValue{xs(" 12345"), xs(`\A\d{5}\z`)}, nil},
	{"has_pattern", []types.XValue{xs("<html>x</html>"), xs(`[`)}, ERROR},

	{"has_number", []types.XValue{xs("the number 10")}, result(xn("10"))},
	{"has_number", []types.XValue{xs("the number -10")}, result(xn("-10"))},
	{"has_number", []types.XValue{xs("1-15")}, result(xn("1"))},
	{"has_number", []types.XValue{xs("24ans")}, result(xn("24"))},
	{"has_number", []types.XValue{xs("J'AI 20ANS")}, result(xn("20"))},
	{"has_number", []types.XValue{xs("1,000,000")}, result(xn("1000000"))},
	{"has_number", []types.XValue{xs("the number 10")}, result(xn("10"))},
	{"has_number", []types.XValue{xs("O número é 500")}, result(xn("500"))},
	{"has_number", []types.XValue{xs("another is -12.51")}, result(xn("-12.51"))},
	{"has_number", []types.XValue{xs("hi.51")}, result(xn("51"))},
	{"has_number", []types.XValue{xs("nothing here")}, nil},
	{"has_number", []types.XValue{xs("1OO l00")}, nil}, // no longer do substitutions
	{"has_number", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},

	{"has_number_lt", []types.XValue{xs("the number 10"), xs("11")}, result(xn("10"))},
	{"has_number_lt", []types.XValue{xs("another is -12.51"), xs("12")}, result(xn("-12.51"))},
	{"has_number_lt", []types.XValue{xs("nothing here"), xs("12")}, nil},
	{"has_number_lt", []types.XValue{xs("too big 15"), xs("12")}, nil},
	{"has_number_lt", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_number_lt", []types.XValue{xs("but foo"), nil}, ERROR},
	{"has_number_lt", []types.XValue{nil, xs("but foo")}, ERROR},

	{"has_number_lte", []types.XValue{xs("the number 10"), xs("11")}, result(xn("10"))},
	{"has_number_lte", []types.XValue{xs("another is -12.51"), xs("12")}, result(xn("-12.51"))},
	{"has_number_lte", []types.XValue{xs("nothing here"), xs("12")}, nil},
	{"has_number_lte", []types.XValue{xs("too big 15"), xs("12")}, nil},
	{"has_number_lte", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},

	{"has_number_eq", []types.XValue{xs("the number 10"), xs("10")}, result(xn("10"))},
	{"has_number_eq", []types.XValue{xs("another is -12.51"), xs("-12.51")}, result(xn("-12.51"))},
	{"has_number_eq", []types.XValue{xs("nothing here"), xs("12")}, nil},
	{"has_number_eq", []types.XValue{xs("wrong .51"), xs(".61")}, nil},
	{"has_number_eq", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},

	{"has_number_gte", []types.XValue{xs("the number 10"), xs("9")}, result(xn("10"))},
	{"has_number_gte", []types.XValue{xs("another is -12.51"), xs("-13")}, result(xn("-12.51"))},
	{"has_number_gte", []types.XValue{xs("nothing here"), xs("12")}, nil},
	{"has_number_gte", []types.XValue{xs("too small -12"), xs("-11")}, nil},
	{"has_number_gte", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},

	{"has_number_gt", []types.XValue{xs("the number 10"), xs("9")}, result(xn("10"))},
	{"has_number_gt", []types.XValue{xs("another is -12.51"), xs("-13")}, result(xn("-12.51"))},
	{"has_number_gt", []types.XValue{xs("nothing here"), xs("12")}, nil},
	{"has_number_gt", []types.XValue{xs("not great -12.51"), xs("-12.51")}, nil},
	{"has_number_gt", []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},

	{"has_number_between", []types.XValue{xs("the number 10"), xs("8"), xs("12")}, result(xn("10"))},
	{"has_number_between", []types.XValue{xs("24ans"), xn("20"), xn("24")}, result(xn("24"))},
	{"has_number_between", []types.XValue{xs("another is -12.51"), xs("-12.51"), xs("-10")}, result(xn("-12.51"))},
	{"has_number_between", []types.XValue{xs("nothing here"), xs("10"), xs("15")}, nil},
	{"has_number_between", []types.XValue{xs("one"), xs("two")}, ERROR},
	{"has_number_between", []types.XValue{xs("but foo"), nil, xs("10")}, ERROR},
	{"has_number_between", []types.XValue{nil, xs("but foo"), xs("10")}, ERROR},
	{"has_number_between", []types.XValue{xs("a string"), xs("10"), xs("not number")}, ERROR},

	{"has_date", []types.XValue{xs("last date was 1.10.2017")}, result(xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date", []types.XValue{xs("last date was 1.10.99")}, result(xd(time.Date(1999, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date", []types.XValue{xs("this isn't a valid date 33.2.99")}, nil},
	{"has_date", []types.XValue{xs("no date at all")}, nil},
	{"has_date", []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},

	{"has_date_lt", []types.XValue{xs("last date was 1.10.2017"), xs("3.10.2017")}, result(xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date_lt", []types.XValue{xs("last date was 1.10.99"), xs("3.10.98")}, nil},
	{"has_date_lt", []types.XValue{xs("no date at all"), xs("3.10.98")}, nil},
	{"has_date_lt", []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},
	{"has_date_lt", []types.XValue{xs("last date was 1.10.2017"), nil}, ERROR},
	{"has_date_lt", []types.XValue{nil, xs("but foo")}, ERROR},

	{"has_date_eq", []types.XValue{xs("last date was 1.10.2017"), xs("1.10.2017")}, result(xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date_eq", []types.XValue{xs("last date was 1.10.99"), xs("3.10.98")}, nil},
	{"has_date_eq", []types.XValue{xs("2017-10-01T23:55:55.123456+02:00"), xs("1.10.2017")}, result(xd(time.Date(2017, 10, 1, 23, 55, 55, 123456000, kgl)))},
	{"has_date_eq", []types.XValue{xs("2017-10-01T23:55:55.123456+01:00"), xs("1.10.2017")}, nil}, // would have been 2017-10-02 in env timezone
	{"has_date_eq", []types.XValue{xs("no date at all"), xs("3.10.98")}, nil},
	{"has_date_eq", []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},

	{"has_date_gt", []types.XValue{xs("last date was 1.10.2017"), xs("3.10.2016")}, result(xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date_gt", []types.XValue{xs("last date was 1.10.99"), xs("3.10.01")}, nil},
	{"has_date_gt", []types.XValue{xs("no date at all"), xs("3.10.98")}, nil},
	{"has_date_gt", []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},

	{"has_time", []types.XValue{xs("last time was 10:30")}, result(xt(utils.NewTimeOfDay(10, 30, 0, 0)))},
	{"has_time", []types.XValue{xs("this isn't a valid time 59:77")}, nil},
	{"has_time", []types.XValue{xs("no time at all")}, nil},
	{"has_time", []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},

	{"has_email", []types.XValue{xs("my email is foo@bar.com.")}, result(xs("foo@bar.com"))},
	{"has_email", []types.XValue{xs("my email is <foo1@bar-2.com>")}, result(xs("foo1@bar-2.com"))},
	{"has_email", []types.XValue{xs("FOO@bar.whatzit")}, result(xs("FOO@bar.whatzit"))},
	{"has_email", []types.XValue{xs("FOO@βήτα.whatzit")}, result(xs("FOO@βήτα.whatzit"))},
	{"has_email", []types.XValue{xs("email is foo @ bar . com")}, nil},
	{"has_email", []types.XValue{xs("email is foo@bar")}, nil},
	{"has_email", []types.XValue{nil}, nil},
	{"has_email", []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},

	{"has_phone", []types.XValue{xs("my number is +250788123123")}, result(xs("+250788123123"))},
	{"has_phone", []types.XValue{xs("my number is +593979111111")}, result(xs("+593979111111"))},
	{"has_phone", []types.XValue{xs("my number is 0788123123")}, result(xs("+250788123123"))}, // uses environment default
	{"has_phone", []types.XValue{xs("my number is 0788123123"), xs("RW")}, result(xs("+250788123123"))},
	{"has_phone", []types.XValue{xs("my number is +250788123123"), xs("RW")}, result(xs("+250788123123"))},
	{"has_phone", []types.XValue{xs("my number is +12065551212"), xs("RW")}, result(xs("+12065551212"))},
	{"has_phone", []types.XValue{xs("my number is 12065551212"), xs("US")}, result(xs("+12065551212"))},
	{"has_phone", []types.XValue{xs("my number is 206 555 1212"), xs("US")}, result(xs("+12065551212"))},
	{"has_phone", []types.XValue{xs("my number is +10001112222"), xs("US")}, result(xs("+10001112222"))},
	{"has_phone", []types.XValue{xs("my number is 10000"), xs("US")}, nil},
	{"has_phone", []types.XValue{xs("my number is 12067799294"), xs("BW")}, nil},
	{"has_phone", []types.XValue{xs("my number is none of your business"), xs("US")}, nil},
	{"has_phone", []types.XValue{}, ERROR},
	{"has_phone", []types.XValue{ERROR}, ERROR},
	{"has_phone", []types.XValue{xs("3245"), ERROR}, ERROR},
	{"has_phone", []types.XValue{xs("number"), nil}, nil},
	{"has_phone", []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},

	{
		"has_group",
		[]types.XValue{
			xa(
				types.NewXDict(map[string]types.XValue{"uuid": xs("group-uuid-1"), "name": xs("Testers")}),
				types.NewXDict(map[string]types.XValue{"uuid": xs("group-uuid-2"), "name": xs("Customers")}),
			),
			xs("group-uuid-2"),
		},
		types.NewXDict(map[string]types.XValue{
			"match": types.NewXDict(map[string]types.XValue{"uuid": xs("group-uuid-2"), "name": xs("Customers")}),
		}),
	},
	{"has_group", []types.XValue{xa(), xs("group-uuid-1")}, nil},
	{"has_group", []types.XValue{ERROR, ERROR}, ERROR},
	{"has_group", []types.XValue{}, ERROR},
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

	for _, tc := range testTests {
		testID := fmt.Sprintf("%s(%#v)", tc.name, tc.args)

		testFunc, exists := tests.XTESTS[tc.name]
		require.True(t, exists, "no such registered function: %s", tc.name)

		result := testFunc(env, tc.args...)

		// don't check error equality - just check that we got an error if we expected one
		if tc.expected == ERROR {
			assert.True(t, types.IsXError(result), "expecting error, got %T{%s} for ", result, result, testID)
		} else {
			if !types.Equals(env, result, tc.expected) {
				assert.Fail(t, "", "unexpected value, expected %T{%s}, got %T{%s} for ", tc.expected, tc.expected, result, result, testID)
			}
		}
	}
}

func TestEvaluateTemplate(t *testing.T) {
	vars := types.NewXDict(map[string]types.XValue{
		"int1":   types.NewXNumberFromInt(1),
		"int2":   types.NewXNumberFromInt(2),
		"array1": types.NewXArray(xs("one"), xs("two"), xs("three")),
		"thing": types.NewXDict(map[string]types.XValue{
			"foo":     types.NewXText("bar"),
			"zed":     types.NewXNumberFromInt(123),
			"missing": nil,
		}),
		"err": types.NewXErrorf("an error"),
	})

	evalTests := []struct {
		template string
		expected string
		hasError bool
	}{
		{"@(is_error(array1[100]))", "{match: index 100 out of range for 3 items}", false}, // errors are like any other value
		{"@(is_error(array1.100))", "{match: array has no property '100'}", false},
		{`@(is_error(round("foo", "bar")))`, "{match: error calling ROUND: unable to convert \"foo\" to a number}", false},
		{`@(is_error(err))`, "{match: an error}", false},
		{"@(is_error(thing.foo))", "", false},
		{"@(is_error(thing.xxx))", "{match: dict has no property 'xxx'}", false},
		{"@(is_error(1 / 0))", "{match: division by zero}", false},
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
