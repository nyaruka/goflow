package cases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/routers/cases"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var xs = types.NewXText
var xn = types.RequireXNumberFromString
var xd = types.NewXDateTime
var xt = types.NewXTime
var xa = types.NewXArray
var xj = func(s string) types.XValue { return types.JSONToXValue([]byte(s)) }
var result = cases.NewTrueResult
var resultWithExtra = cases.NewTrueResultWithExtra
var falseResult = cases.FalseResult
var ERROR = types.NewXErrorf("any error")

var kgl, _ = time.LoadLocation("Africa/Kigali")
var dmy = envs.NewBuilder().
	WithDateFormat(envs.DateFormatDayMonthYear).
	WithTimeFormat(envs.TimeFormatHourMinuteSecond).
	WithTimezone(kgl).
	WithDefaultCountry(i18n.Country("RW")).
	Build()
var ara = envs.NewBuilder().
	WithInputCollation(envs.CollationArabicVariants).
	Build()

var assetsJSON = `{
	"flows": [
        {
            "uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
            "name": "Test",
            "spec_version": "13.1.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        }
	],
    "locations": [
		{
			"name": "Rwanda",
			"aliases": ["Ruanda"],		
			"children": [
				{
					"name": "Kigali City",
					"aliases": ["Kigali", "Kigari"],
					"children": [
						{
							"name": "Gasabo",
							"children": [
								{
									"name": "Gisozi"
								},
								{
									"name": "Ndera"
								}
							]
						},
						{
							"name": "Nyarugenge",
							"children": []
						}
					]
				},
				{
					"name": "Québec",
					"aliases": ["Q.C", "Le Québec", "Quebec", "Que,", "Que", "Qc", "Québec"],
					"children": []
				},
				{
					"name": "Île-de-France",
					"aliases": [],
					"children": []
				},
				{
					"name": "Paktika",
					"aliases": ["Janikhel", "Terwa", "Yahyakhel", "Yusufkhel", "\u067e\u06a9\u062a\u06cc\u06a9\u0627", "\u062a\u0631\u0648\u0648", "\u06cc\u062d\u06cc\u06cc \u062e\u06cc\u0644", "\u06cc\u0648\u0633\u0641 \u062e\u06cc\u0644"],
					"children": []
				}
			]
		}
	]
}`

var testTests = []struct {
	name     string
	env      envs.Environment
	args     []types.XValue
	expected types.XValue
}{

	{"has_error", dmy, []types.XValue{xs("hello")}, falseResult},
	{"has_error", dmy, []types.XValue{nil}, falseResult},
	{"has_error", dmy, []types.XValue{types.NewXErrorf("I am error")}, result(xs("I am error"))},
	{"has_error", dmy, []types.XValue{}, ERROR},

	{"has_text", dmy, []types.XValue{xs("hello")}, result(xs("hello"))},
	{"has_text", dmy, []types.XValue{xs("  ")}, falseResult},
	{"has_text", dmy, []types.XValue{nil}, falseResult},
	{"has_text", dmy, []types.XValue{xs("one"), xs("two")}, ERROR},
	{"has_text", dmy, []types.XValue{ERROR}, ERROR},

	{"has_only_text", dmy, []types.XValue{xs("hello"), xs("hello")}, result(xs("hello"))},
	{"has_only_text", dmy, []types.XValue{xs("hello-world"), xs("hello-world")}, result(xs("hello-world"))},
	{"has_only_text", dmy, []types.XValue{xs("HELLO"), xs("hello")}, falseResult}, // case sensitive
	{"has_only_text", dmy, []types.XValue{xs("hello"), ERROR}, ERROR},
	{"has_only_text", dmy, []types.XValue{ERROR, xs("hello")}, ERROR},

	{"has_beginning", dmy, []types.XValue{xs("hello"), xs("hell")}, result(xs("hell"))},
	{"has_beginning", dmy, []types.XValue{xs("  HelloThere"), xs("hello")}, result(xs("Hello"))},
	{"has_beginning", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_beginning", dmy, []types.XValue{nil, xs("hell")}, falseResult},
	{"has_beginning", dmy, []types.XValue{xs("hello"), nil}, falseResult},
	{"has_beginning", dmy, []types.XValue{xs(""), xs("hello")}, falseResult},
	{"has_beginning", dmy, []types.XValue{xs("hel"), xs("hello")}, falseResult},
	{"has_beginning", dmy, []types.XValue{ERROR, ERROR}, ERROR},
	{"has_beginning", dmy, []types.XValue{}, ERROR},

	{"has_any_word", dmy, []types.XValue{xs("this.is.my.word"), xs("WORD word2 word")}, result(xs("word"))},
	{"has_any_word", dmy, []types.XValue{xs("this.is.my.βήτα"), xs("βήτα")}, result(xs("βήτα"))},
	{"has_any_word", dmy, []types.XValue{xs("I say to you📴"), xs("📴")}, result(xs("📴"))},
	{"has_any_word", dmy, []types.XValue{xs("this World too"), xs("world")}, result(xs("World"))},
	{"has_any_word", dmy, []types.XValue{xs("I don't like it"), xs("don't dont")}, result(xs("don't"))},
	{"has_any_word", dmy, []types.XValue{xs("BUT not this one"), xs("world")}, falseResult},
	{"has_any_word", dmy, []types.XValue{xs(""), xs("world")}, falseResult},
	{"has_any_word", dmy, []types.XValue{xs("world"), xs("foo")}, falseResult},
	{"has_any_word", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_any_word", dmy, []types.XValue{xs("but foo"), nil}, falseResult},
	{"has_any_word", dmy, []types.XValue{nil, xs("but foo")}, falseResult},
	{"has_any_word", dmy, []types.XValue{xs("بلی"), xs("بلی")}, result(xs("بلی"))}, // using regular collation
	{"has_any_word", dmy, []types.XValue{xs("بلي"), xs("بلی")}, falseResult},       // not quite the same yeh
	{"has_any_word", dmy, []types.XValue{xs("بلى"), xs("بلی")}, falseResult},       // not quite the same yeh
	{"has_any_word", ara, []types.XValue{xs("بلی"), xs("بلی")}, result(xs("بلی"))}, // using ara-far collation
	{"has_any_word", ara, []types.XValue{xs("بلي"), xs("بلی")}, result(xs("بلي"))}, // using ara-far collation
	{"has_any_word", ara, []types.XValue{xs("بلى"), xs("بلی")}, result(xs("بلى"))}, // using ara-far collation
	{"has_any_word", dmy, []types.XValue{}, ERROR},

	{"has_all_words", dmy, []types.XValue{xs("this.is.my.word"), xs("WORD word")}, result(xs("word"))},
	{"has_all_words", dmy, []types.XValue{xs("this World too"), xs("world too")}, result(xs("World too"))},
	{"has_all_words", dmy, []types.XValue{xs("BUT not this one"), xs("world")}, falseResult},
	{"has_all_words", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_all_words", dmy, []types.XValue{}, ERROR},

	{"has_phrase", dmy, []types.XValue{xs("you Must resist"), xs("must resist")}, result(xs("Must resist"))},
	{"has_phrase", dmy, []types.XValue{xs("this world Too"), xs("world too")}, result(xs("world Too"))},
	{"has_phrase", dmy, []types.XValue{xs("this world Too"), xs("")}, result(xs(""))},
	{"has_phrase", dmy, []types.XValue{xs("this is not world"), xs("this world")}, falseResult},
	{"has_phrase", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_phrase", dmy, []types.XValue{}, ERROR},

	{"has_only_phrase", dmy, []types.XValue{xs("Must resist"), xs("must resist")}, result(xs("Must resist"))},
	{"has_only_phrase", dmy, []types.XValue{xs(" world Too "), xs("world too")}, result(xs("world Too"))},
	{"has_only_phrase", dmy, []types.XValue{xs("this world Too"), xs("")}, falseResult},
	{"has_only_phrase", dmy, []types.XValue{xs(""), xs("")}, result(xs(""))},
	{"has_only_phrase", dmy, []types.XValue{xs("this world is my world"), xs("this world")}, falseResult},
	{"has_only_phrase", dmy, []types.XValue{xs("this world"), xs("this mighty")}, falseResult},
	{"has_only_phrase", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_only_phrase", dmy, []types.XValue{}, ERROR},

	{"has_beginning", dmy, []types.XValue{xs("Must resist"), xs("must resist")}, result(xs("Must resist"))},
	{"has_beginning", dmy, []types.XValue{xs(" 2061212"), xs("206")}, result(xs("206"))},
	{"has_beginning", dmy, []types.XValue{xs(" world Too foo"), xs("world too")}, result(xs("world Too"))},
	{"has_beginning", dmy, []types.XValue{xs("but this world"), xs("this world")}, falseResult},
	{"has_beginning", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_beginning", dmy, []types.XValue{}, ERROR},

	{"has_pattern", dmy, []types.XValue{xs("<html>x</html>"), xs(`<\w+>`)}, resultWithExtra(xs("<html>"), types.NewXObject(map[string]types.XValue{"0": xs("<html>")}))},
	{"has_pattern", dmy, []types.XValue{xs("<html>x</html>"), xs(`HTML`)}, resultWithExtra(xs("html"), types.NewXObject(map[string]types.XValue{"0": xs("html")}))},
	{"has_pattern", dmy, []types.XValue{xs("<html>x</html>"), xs(`(?-i)HTML`)}, falseResult},
	{"has_pattern", dmy, []types.XValue{xs("12345"), xs(`\A\d{5}\z`)}, resultWithExtra(xs("12345"), types.NewXObject(map[string]types.XValue{"0": xs("12345")}))},
	{"has_pattern", dmy, []types.XValue{xs("12345 "), xs(`\A\d{5}\z`)}, falseResult},
	{"has_pattern", dmy, []types.XValue{xs(" 12345"), xs(`\A\d{5}\z`)}, falseResult},
	{"has_pattern", dmy, []types.XValue{xs(`hi there 😀`), xs("[\U0001F600-\U0001F64F]")}, resultWithExtra(xs("😀"), types.NewXObject(map[string]types.XValue{"0": xs("😀")}))},
	{"has_pattern", dmy, []types.XValue{xs(`hi there`), xs("[\U0001F600-\U0001F64F]")}, falseResult},
	{"has_pattern", dmy, []types.XValue{xs(`hi there 😂`), xs("[😀-🙏]")}, resultWithExtra(xs("😂"), types.NewXObject(map[string]types.XValue{"0": xs("😂")}))},
	{"has_pattern", dmy, []types.XValue{xs("<html>x</html>"), xs(`[`)}, ERROR},
	{"has_pattern", dmy, []types.XValue{}, ERROR},

	{"has_number", dmy, []types.XValue{xs("the number 10")}, result(xn("10"))},
	{"has_number", dmy, []types.XValue{xs("the number -10")}, result(xn("-10"))},
	{"has_number", dmy, []types.XValue{xs("1-15")}, result(xn("1"))},
	{"has_number", dmy, []types.XValue{xs("24ans")}, result(xn("24"))},
	{"has_number", dmy, []types.XValue{xs("J'AI 20ANS")}, result(xn("20"))},
	{"has_number", dmy, []types.XValue{xs("1,000,000")}, result(xn("1000000"))},
	{"has_number", dmy, []types.XValue{xs("the number 10")}, result(xn("10"))},
	{"has_number", dmy, []types.XValue{xs("O número é 500")}, result(xn("500"))},
	{"has_number", dmy, []types.XValue{xs("another is -12.51")}, result(xn("-12.51"))},
	{"has_number", dmy, []types.XValue{xs("hi.51")}, result(xn("51"))},
	{"has_number", dmy, []types.XValue{xs("hi .51")}, result(xn("0.51"))},
	{"has_number", dmy, []types.XValue{xs(".51")}, result(xn("0.51"))},
	{"has_number", dmy, []types.XValue{xs("١٢٣٤")}, result(xn("1234"))}, // Arabic
	{"has_number", dmy, []types.XValue{xs("۱۲۳۴")}, result(xn("1234"))}, // Eastern Arabic
	{"has_number", dmy, []types.XValue{xs("٠.٥")}, result(xn("0.5"))},
	{"has_number", dmy, []types.XValue{xs("nothing here")}, falseResult},
	{"has_number", dmy, []types.XValue{xs("lOO")}, falseResult}, // no longer do substitutions
	{"has_number", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_number", dmy, []types.XValue{}, ERROR},

	{"has_number_lt", dmy, []types.XValue{xs("the number 10"), xs("11")}, result(xn("10"))},
	{"has_number_lt", dmy, []types.XValue{xs("another is -12.51"), xs("12")}, result(xn("-12.51"))},
	{"has_number_lt", dmy, []types.XValue{xs("١٠"), xs("11")}, result(xn("10"))},
	{"has_number_lt", dmy, []types.XValue{xs("nothing here"), xs("12")}, falseResult},
	{"has_number_lt", dmy, []types.XValue{xs("too big 15"), xs("12")}, falseResult},
	{"has_number_lt", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_number_lt", dmy, []types.XValue{xs("but foo"), falseResult}, ERROR},
	{"has_number_lt", dmy, []types.XValue{nil, xs("but foo")}, ERROR},
	{"has_number_lt", dmy, []types.XValue{}, ERROR},

	{"has_number_lte", dmy, []types.XValue{xs("the number 10"), xs("11")}, result(xn("10"))},
	{"has_number_lte", dmy, []types.XValue{xs("another is -12.51"), xs("12")}, result(xn("-12.51"))},
	{"has_number_lte", dmy, []types.XValue{xs("١٠"), xs("11")}, result(xn("10"))},
	{"has_number_lte", dmy, []types.XValue{xs("nothing here"), xs("12")}, falseResult},
	{"has_number_lte", dmy, []types.XValue{xs("too big 15"), xs("12")}, falseResult},
	{"has_number_lte", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_number_lte", dmy, []types.XValue{}, ERROR},

	{"has_number_eq", dmy, []types.XValue{xs("the number 10"), xs("10")}, result(xn("10"))},
	{"has_number_eq", dmy, []types.XValue{xs("another is -12.51"), xs("-12.51")}, result(xn("-12.51"))},
	{"has_number_eq", dmy, []types.XValue{xs("١٠"), xs("10")}, result(xn("10"))},
	{"has_number_eq", dmy, []types.XValue{xs("nothing here"), xs("12")}, falseResult},
	{"has_number_eq", dmy, []types.XValue{xs("wrong .51"), xs(".61")}, falseResult},
	{"has_number_eq", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_number_eq", dmy, []types.XValue{}, ERROR},

	{"has_number_gte", dmy, []types.XValue{xs("the number 10"), xs("9")}, result(xn("10"))},
	{"has_number_gte", dmy, []types.XValue{xs("another is -12.51"), xs("-13")}, result(xn("-12.51"))},
	{"has_number_gte", dmy, []types.XValue{xs("١٠"), xs("9")}, result(xn("10"))},
	{"has_number_gte", dmy, []types.XValue{xs("nothing here"), xs("12")}, falseResult},
	{"has_number_gte", dmy, []types.XValue{xs("too small -12"), xs("-11")}, falseResult},
	{"has_number_gte", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_number_gte", dmy, []types.XValue{}, ERROR},

	{"has_number_gt", dmy, []types.XValue{xs("the number 10"), xs("9")}, result(xn("10"))},
	{"has_number_gt", dmy, []types.XValue{xs("another is -12.51"), xs("-13")}, result(xn("-12.51"))},
	{"has_number_gt", dmy, []types.XValue{xs("١٠"), xs("9")}, result(xn("10"))},
	{"has_number_gt", dmy, []types.XValue{xs("nothing here"), xs("12")}, falseResult},
	{"has_number_gt", dmy, []types.XValue{xs("not great -12.51"), xs("-12.51")}, falseResult},
	{"has_number_gt", dmy, []types.XValue{xs("one"), xs("two"), xs("three")}, ERROR},
	{"has_number_gt", dmy, []types.XValue{}, ERROR},

	{"has_number_between", dmy, []types.XValue{xs("the number 10"), xs("8"), xs("12")}, result(xn("10"))},
	{"has_number_between", dmy, []types.XValue{xs("24ans"), xn("20"), xn("24")}, result(xn("24"))},
	{"has_number_between", dmy, []types.XValue{xs("another is -12.51"), xs("-12.51"), xs("-10")}, result(xn("-12.51"))},
	{"has_number_between", dmy, []types.XValue{xs("١٠"), xs("8"), xs("12")}, result(xn("10"))},
	{"has_number_between", dmy, []types.XValue{xs("nothing here"), xs("10"), xs("15")}, falseResult},
	{"has_number_between", dmy, []types.XValue{xs("one"), xs("two")}, ERROR},
	{"has_number_between", dmy, []types.XValue{xs("but foo"), nil, xs("10")}, ERROR},
	{"has_number_between", dmy, []types.XValue{nil, xs("but foo"), xs("10")}, ERROR},
	{"has_number_between", dmy, []types.XValue{xs("a string"), xs("10"), xs("not number")}, ERROR},
	{"has_number_between", dmy, []types.XValue{}, ERROR},

	{"has_date", dmy, []types.XValue{xs("31/02/2022")}, falseResult},
	{"has_date", dmy, []types.XValue{xs("33/11/2022")}, falseResult},
	{"has_date", dmy, []types.XValue{xs("31/11/2022")}, falseResult},
	{"has_date", dmy, []types.XValue{xs("30/02/2022")}, falseResult},
	{"has_date", dmy, []types.XValue{xs("last date was 1.10.2017")}, result(xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date", dmy, []types.XValue{xs("last date was 1.10.99")}, result(xd(time.Date(1999, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date", dmy, []types.XValue{xs("this isn't a valid date 33.2.99")}, falseResult},
	{"has_date", dmy, []types.XValue{xs("no date at all")}, falseResult},
	{"has_date", dmy, []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},
	{"has_date", dmy, []types.XValue{}, ERROR},

	{"has_date_lt", dmy, []types.XValue{xs("last date was 1.10.2017"), xs("3.10.2017")}, result(xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date_lt", dmy, []types.XValue{xs("last date was 1.10.99"), xs("3.10.98")}, falseResult},
	{"has_date_lt", dmy, []types.XValue{xs("no date at all"), xs("3.10.98")}, falseResult},
	{"has_date_lt", dmy, []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},
	{"has_date_lt", dmy, []types.XValue{xs("last date was 1.10.2017"), nil}, ERROR},
	{"has_date_lt", dmy, []types.XValue{nil, xs("but foo")}, ERROR},
	{"has_date_lt", dmy, []types.XValue{}, ERROR},

	{"has_date_eq", dmy, []types.XValue{xs("last date was 1.10.2017"), xs("1.10.2017")}, result(xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date_eq", dmy, []types.XValue{xs("last date was 1.10.99"), xs("3.10.98")}, falseResult},
	{"has_date_eq", dmy, []types.XValue{xs("2017-10-01T23:55:55.123456+02:00"), xs("1.10.2017")}, result(xd(time.Date(2017, 10, 1, 23, 55, 55, 123456000, kgl)))},
	{"has_date_eq", dmy, []types.XValue{xs("2017-10-01T23:55:55.123456+01:00"), xs("1.10.2017")}, falseResult}, // would have been 2017-10-02 in env timezone
	{"has_date_eq", dmy, []types.XValue{xs("no date at all"), xs("3.10.98")}, falseResult},
	{"has_date_eq", dmy, []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},
	{"has_date_eq", dmy, []types.XValue{}, ERROR},

	{"has_date_gt", dmy, []types.XValue{xs("last date was 1.10.2017"), xs("3.10.2016")}, result(xd(time.Date(2017, 10, 1, 15, 24, 30, 123456000, kgl)))},
	{"has_date_gt", dmy, []types.XValue{xs("last date was 1.10.99"), xs("3.10.01")}, falseResult},
	{"has_date_gt", dmy, []types.XValue{xs("no date at all"), xs("3.10.98")}, falseResult},
	{"has_date_gt", dmy, []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},
	{"has_date_gt", dmy, []types.XValue{}, ERROR},

	{"has_time", dmy, []types.XValue{xs("last time was 10:30")}, result(xt(dates.NewTimeOfDay(10, 30, 0, 0)))},
	{"has_time", dmy, []types.XValue{xs("this isn't a valid time 59:77")}, falseResult},
	{"has_time", dmy, []types.XValue{xs("no time at all")}, falseResult},
	{"has_time", dmy, []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},
	{"has_time", dmy, []types.XValue{}, ERROR},

	{"has_email", dmy, []types.XValue{xs("my email is foo@bar.com.")}, result(xs("foo@bar.com"))},
	{"has_email", dmy, []types.XValue{xs("my email is <foo~$1+spam@bar-2.com>")}, result(xs("foo~$1+spam@bar-2.com"))},
	{"has_email", dmy, []types.XValue{xs("FOO@bar.whatzit")}, result(xs("FOO@bar.whatzit"))},
	{"has_email", dmy, []types.XValue{xs("FOO@βήτα.whatzit")}, result(xs("FOO@βήτα.whatzit"))},
	{"has_email", dmy, []types.XValue{xs("email is foo @ bar . com")}, falseResult},
	{"has_email", dmy, []types.XValue{xs("email is foo@bar")}, falseResult},
	{"has_email", dmy, []types.XValue{nil}, falseResult},
	{"has_email", dmy, []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},
	{"has_email", dmy, []types.XValue{}, ERROR},

	// more has_phone tests in TestHasPhone below
	{"has_phone", dmy, []types.XValue{xs("my number is 0788123123"), xs("RW")}, result(xs("+250788123123"))},
	{"has_phone", dmy, []types.XValue{xs("my number is none of your business"), xs("US")}, falseResult},
	{"has_phone", dmy, []types.XValue{ERROR}, ERROR},
	{"has_phone", dmy, []types.XValue{xs("3245"), ERROR}, ERROR},
	{"has_phone", dmy, []types.XValue{xs("number"), nil}, falseResult},
	{"has_phone", dmy, []types.XValue{xs("too"), xs("many"), xs("args")}, ERROR},
	{"has_phone", dmy, []types.XValue{}, ERROR},

	{
		"has_group",
		dmy,
		[]types.XValue{
			xj(`[{"uuid": "group-uuid-1", "name": "Testers"}, {"uuid": "group-uuid-2", "name": "Customers"}]`),
			xs("group-uuid-2"),
		},
		result(xj(`{"uuid": "group-uuid-2", "name": "Customers"}`)),
	},
	{"has_group", dmy, []types.XValue{xa(ERROR), xs("group-uuid-2")}, ERROR},
	{"has_group", dmy, []types.XValue{xa(), xs("group-uuid-1")}, falseResult},
	{"has_group", dmy, []types.XValue{ERROR, xs("group-uuid-1")}, ERROR},
	{"has_group", dmy, []types.XValue{xa(), ERROR}, ERROR},
	{"has_group", dmy, []types.XValue{}, ERROR},

	{"has_state", dmy, []types.XValue{xs("Quebec")}, result(xs("Rwanda > Québec"))},
	{"has_state", dmy, []types.XValue{xs("Québec")}, result(xs("Rwanda > Québec"))},
	{"has_state", dmy, []types.XValue{xs("Je suis dans la province du Québec")}, result(xs("Rwanda > Québec"))},

	{"has_state", dmy, []types.XValue{xs("Je suis dans la province de l'Île-de-France")}, result(xs("Rwanda > Île-de-France"))},

	{"has_state", dmy, []types.XValue{xs("kigali city")}, result(xs("Rwanda > Kigali City"))},
	{"has_state", dmy, []types.XValue{xs("¡Kigali!")}, result(xs("Rwanda > Kigali City"))},

	{"has_state", dmy, []types.XValue{xs("kigari")}, result(xs("Rwanda > Kigali City"))},
	{"has_state", dmy, []types.XValue{xs("تروو")}, result(xs("Rwanda > Paktika"))},
	{"has_state", dmy, []types.XValue{xs("غم ځپلې هلمند")}, falseResult},
	{"has_state", dmy, []types.XValue{xs("\u063a\u0645 \u0681\u067e\u0644\u06d0 \u0647\u0644\u0645\u0646\u062f")}, falseResult},
	{"has_state", dmy, []types.XValue{xs("xyz")}, falseResult},
	{"has_state", dmy, []types.XValue{ERROR}, ERROR},

	{"has_district", dmy, []types.XValue{xs("Gasabo"), xs("kigali")}, result(xs("Rwanda > Kigali City > Gasabo"))},
	{"has_district", dmy, []types.XValue{xs("I live in gasabo"), xs("kigali")}, result(xs("Rwanda > Kigali City > Gasabo"))},
	{"has_district", dmy, []types.XValue{xs("Gasabo")}, result(xs("Rwanda > Kigali City > Gasabo"))},
	{"has_district", dmy, []types.XValue{xs("xyz"), xs("kigali")}, falseResult},
	{"has_district", dmy, []types.XValue{ERROR}, ERROR},

	{"has_ward", dmy, []types.XValue{xs("Gisozi"), xs("kigali"), xs("Gasabo")}, result(xs("Rwanda > Kigali City > Gasabo > Gisozi"))},
	{"has_ward", dmy, []types.XValue{xs("I live in gisozi"), xs("kigali"), xs("Gasabo")}, result(xs("Rwanda > Kigali City > Gasabo > Gisozi"))},
	{"has_ward", dmy, []types.XValue{xs("Gisozi")}, result(xs("Rwanda > Kigali City > Gasabo > Gisozi"))},
	{"has_ward", dmy, []types.XValue{xs("xyz"), xs("kigali"), xs("Gasabo")}, falseResult},
	{"has_ward", dmy, []types.XValue{ERROR}, ERROR},

	{
		"has_category",
		dmy,
		[]types.XValue{
			xj(`{
				"name": "Response 1",
				"value": "hi",
				"category": "Other",
				"input": "hello",
				"node_uuid": "0faca870-aca4-469d-89e2-a70df468ac68",
				"created_on": "2018-07-06T12:30:06.123456789Z"
			}`),
			xs("Chicken"),
			xs("Other"),
		},
		result(xs("Other")),
	},
	{
		"has_category",
		dmy,
		[]types.XValue{
			xj(`{
				"name": "Response 1",
				"value": "hi",
				"category": "All Responses",
				"input": "hello",
				"node_uuid": "0faca870-aca4-469d-89e2-a70df468ac68",
				"created_on": "2018-07-06T12:30:06.123456789Z"
			}`),
			xs("Chicken"),
		},
		falseResult,
	},
	{
		"has_category",
		dmy,
		[]types.XValue{
			xj(`{}`), // not a result
			xs("Chicken"),
		},
		ERROR,
	},

	{
		"has_intent",
		dmy,
		[]types.XValue{
			xj(`{
				"name": "Intention",
				"value": "0.92",
				"category": "success",
				"input": "book me a flight to Quito!",
				"node_uuid": "0faca870-aca4-469d-89e2-a70df468ac68",
				"created_on": "2018-07-06T12:30:06.123456789Z",
				"extra": {
					"intents": [
						{"name": "book_flight", "confidence": 0.7},
						{"name": "book_hotel", "confidence": 0.3}
					],
					"entities": {
						"location": [
							{"value": "Quito", "confidence": 1.0},
							{"value": "Cuenca", "confidence": 0.1} 
						],
						"date": [
							{"value": "May 21", "confidence": 0.6}
						]
					}
				}
			}`),
			xs("book_hotel"),
			xn("0.2"),
		},
		resultWithExtra(
			xs("book_hotel"),
			xj(`{"location": "Quito", "date": "May 21"}`).(*types.XObject),
		),
	},
	{
		"has_intent",
		dmy,
		[]types.XValue{
			xj(`{}`), // not a result
			xs("book_flight"),
			xn("0.5"),
		},
		ERROR,
	},
	{"has_intent", dmy, []types.XValue{}, ERROR},

	{
		"has_top_intent",
		dmy,
		[]types.XValue{
			xj(`{
				"name": "Intention",
				"value": "0.92",
				"category": "success",
				"input": "book me a flight to Quito!",
				"node_uuid": "0faca870-aca4-469d-89e2-a70df468ac68",
				"created_on": "2018-07-06T12:30:06.123456789Z",
				"extra": {
					"intents": [
						{"name": "book_flight", "confidence": 0.7},
						{"name": "book_hotel", "confidence": 0.3}
					],
					"entities": {
						"location": [
							{"value": "Quito", "confidence": 1.0},
							{"value": "Cuenca", "confidence": 0.1} 
						],
						"date": [
							{"value": "May 21", "confidence": 0.6}
						]
					}
				}
			}`),
			xs("book_flight"),
			xn("0.5"),
		},
		resultWithExtra(
			xs("book_flight"),
			xj(`{"location": "Quito", "date": "May 21"}`).(*types.XObject),
		),
	},
	{
		"has_top_intent",
		dmy,
		[]types.XValue{
			xj(`{
				"name": "Intention",
				"created_on": "2018-07-06T12:30:06.123456789Z",
				"extra": {
					"intents": [
						{"name": "book_flight", "confidence": 0.7},
						{"name": "book_hotel", "confidence": 0.3}
					]
				}
			}`),
			xs("book_flight"),
			xn("0.8"), // higher than the extracted confidence of book_flight
		},
		falseResult,
	},
	{
		"has_top_intent",
		dmy,
		[]types.XValue{
			xj(`{}`), // not a result
			xs("book_flight"),
			xn("0.5"),
		},
		ERROR,
	},
	{"has_top_intent", dmy, []types.XValue{}, ERROR},
}

func TestTests(t *testing.T) {
	dates.SetNowFunc(dates.NewFixedNow(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)))
	defer dates.SetNowFunc(time.Now)

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(dmy, source, nil)
	require.NoError(t, err)

	contact := flows.NewEmptyContact(sa, "", i18n.NilLanguage, nil)

	for _, tc := range testTests {
		trigger := triggers.NewBuilder(assets.NewFlowReference("76f0a02f-3b75-4b86-9064-e9195e1b3a02", "Test")).Manual().Build()
		eng := engine.NewBuilder().Build()

		session, _, err := eng.NewSession(context.Background(), sa, tc.env, contact, trigger, nil)
		require.NoError(t, err)

		env := session.MergedEnvironment()

		testID := fmt.Sprintf("%s(%#v)", tc.name, tc.args)

		testFunc, exists := cases.XTESTS[tc.name]
		require.True(t, exists, "no such registered function: %s", tc.name)

		result := testFunc.Call(env, tc.args)

		// don't check error equality - just check that we got an error if we expected one
		if tc.expected == ERROR {
			assert.True(t, types.IsXError(result), "expecting error, got %T{%s} for ", result, result, testID)
		} else {
			test.AssertXEqual(t, tc.expected, result, "result mismatch for %s", testID)
		}
	}
}

func TestEvaluateTemplate(t *testing.T) {
	ctx := types.NewXObject(map[string]types.XValue{
		"int1":   types.NewXNumberFromInt(1),
		"int2":   types.NewXNumberFromInt(2),
		"array1": types.NewXArray(xs("one"), xs("two"), xs("three")),
		"thing": types.NewXObject(map[string]types.XValue{
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
		{"@(has_error(array1[100]).match)", "index 100 out of range for 3 items", false}, // errors are like any other value
		{`@(has_error(round("foo", "bar")).match)`, "error calling round(...): unable to convert \"foo\" to a number", false},
		{`@(has_error(err).match)`, "an error", false},
		{"@(has_error(thing.foo).match)", "", false},
		{"@(has_error(thing.xxx).match)", "object has no property 'xxx'", false},
		{"@(has_error(1 / 0).match)", "division by zero", false},
	}

	env := envs.NewBuilder().Build()
	for _, test := range evalTests {
		eval := excellent.NewEvaluator()
		val, _, err := eval.Template(env, ctx, test.template, nil)

		if test.hasError {
			assert.Error(t, err, "expected error evaluating template '%s'", test.template)
		} else {
			assert.NoError(t, err, "unexpected error evaluating template '%s'", test.template)

			assert.Equal(t, test.expected, val, "actual '%s' does not match expected '%s' evaluating template: '%s'", eval, test.expected, test.template)
		}
	}
}

func TestHasPhone(t *testing.T) {
	tests := []struct {
		input    string
		country  string
		expected string
	}{
		{"+250788123123", "", "+250788123123"},
		{"u812111005611", "ID", "+62812111005611"}, // we try hard to find a number, but check it is valid, it is in this case for ID
		{"oioas812111", "US", ""},                  // in this case we also try hard, but the final result is not a valid US number
		{"+593979111111", "", "+593979111111"},
		{"0788123123", "", "+250788123123"}, // uses environment default
		{"0788123123", "RW", "+250788123123"},
		{"+250788123123", "RW", "+250788123123"},
		{"+12065551212", "RW", "+12065551212"}, // if num has country code, doesn't need to match test country
		{"12065551212", "US", "+12065551212"},
		{"206 555 1212", "US", "+12065551212"},
		{"5912705", "US", ""},                      // would be possible as a local number but not national
		{"+10001112222", "US", "+10001112222"},     // Invalid but possible US number
		{"0815 1053 7962", "ID", "+6281510537962"}, // Indonesian numbers with 12 digits
		{"0954 1053 7962", "ID", "+6295410537962"}, // Invalid but possible Indonesian number
		{"0811-1005-611", "ID", "+628111005611"},   // Valid with 11 digits
		{"10000", "US", ""},
		{"12067799294", "BW", "+12067799294"},
		{"oui", "CD", ""},
	}

	env := envs.NewBuilder().WithDefaultCountry(i18n.Country("RW")).Build()

	for _, tc := range tests {
		var actual, expected types.XValue
		if tc.country != "" {
			actual = cases.HasPhone(env, xs("my number is "+tc.input), xs(tc.country))
		} else {
			actual = cases.HasPhone(env, xs("my number is "+tc.input))
		}

		if tc.expected != "" {
			expected = cases.NewTrueResult(xs(tc.expected))
		} else {
			expected = falseResult
		}

		test.AssertXEqual(t, expected, actual, "has_phone mismatch for input=%s country=%s", tc.input, tc.country)
	}
}
