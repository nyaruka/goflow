package functions_test

import (
	"math"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

var errorArg = types.NewXErrorf("I am error")
var la, _ = time.LoadLocation("America/Los_Angeles")

var xs = types.NewXString
var xn = types.RequireXNumberFromString
var xi = types.NewXNumberFromInt
var xd = types.NewXDate

var ERROR = types.NewXErrorf("any error")

var funcTests = []struct {
	name     string
	args     []types.XValue
	expected types.XValue
}{
	{"and", []types.XValue{types.XBoolTrue}, types.XBoolTrue},
	{"and", []types.XValue{types.XBoolFalse}, types.XBoolFalse},
	{"and", []types.XValue{types.XBoolTrue, types.XBoolFalse}, types.XBoolFalse},
	{"and", []types.XValue{}, ERROR},

	{"char", []types.XValue{xn("33")}, xs("!")},
	{"char", []types.XValue{xn("128513")}, xs("😁")},
	{"char", []types.XValue{xs("not a number")}, ERROR},
	{"char", []types.XValue{}, ERROR},

	{"or", []types.XValue{types.XBoolTrue}, types.XBoolTrue},
	{"or", []types.XValue{types.XBoolFalse}, types.XBoolFalse},
	{"or", []types.XValue{types.XBoolTrue, types.XBoolFalse}, types.XBoolTrue},
	{"or", []types.XValue{}, ERROR},

	{"if", []types.XValue{types.XBoolTrue, xs("10"), xs("20")}, xs("10")},
	{"if", []types.XValue{types.XBoolFalse, xs("10"), xs("20")}, xs("20")},
	{"if", []types.XValue{types.XBoolTrue, errorArg, xs("20")}, errorArg},
	{"if", []types.XValue{}, ERROR},
	{"if", []types.XValue{errorArg, xs("10"), xs("20")}, errorArg},

	{"rand", []types.XValue{}, xn("0.3849275689214193274523267973563633859157562255859375")},
	{"rand", []types.XValue{}, xn("0.607552015674623913099594574305228888988494873046875")},

	{"rand_between", []types.XValue{xn("1"), xn("10")}, xn("5")},
	{"rand_between", []types.XValue{xn("1"), xn("10")}, xn("10")},

	{"round", []types.XValue{xs("10.5"), xs("0")}, xi(11)},
	{"round", []types.XValue{xs("10.5"), xs("1")}, xn("10.5")},
	{"round", []types.XValue{xs("10.51"), xs("1")}, xn("10.5")},
	{"round", []types.XValue{xs("10.56"), xs("1")}, xn("10.6")},
	{"round", []types.XValue{xs("12.56"), xs("-1")}, xi(10)},
	{"round", []types.XValue{xs("10.5")}, xn("11")},
	{"round", []types.XValue{xs("not_num"), xs("1")}, ERROR},
	{"round", []types.XValue{xs("10.5"), xs("not_num")}, ERROR},
	{"round", []types.XValue{xs("10.5"), xs("1"), xs("30")}, ERROR},

	{"round_up", []types.XValue{xs("10.5")}, xi(11)},
	{"round_up", []types.XValue{xs("10.2")}, xi(11)},
	{"round_up", []types.XValue{xs("not_num")}, ERROR},
	{"round_up", []types.XValue{}, ERROR},

	{"round_down", []types.XValue{xs("10.5")}, xi(10)},
	{"round_down", []types.XValue{xs("10.7")}, xi(10)},
	{"round_down", []types.XValue{xs("not_num")}, ERROR},
	{"round_down", []types.XValue{}, ERROR},

	{"max", []types.XValue{xs("10.5"), xs("11")}, xi(11)},
	{"max", []types.XValue{xs("10.2"), xs("9")}, xn("10.2")},
	{"max", []types.XValue{xs("not_num"), xs("9")}, ERROR},
	{"max", []types.XValue{xs("9"), xs("not_num")}, ERROR},
	{"max", []types.XValue{}, ERROR},

	{"min", []types.XValue{xs("10.5"), xs("11")}, xn("10.5")},
	{"min", []types.XValue{xs("10.2"), xs("9")}, xi(9)},
	{"min", []types.XValue{xs("not_num"), xs("9")}, ERROR},
	{"min", []types.XValue{xs("9"), xs("not_num")}, ERROR},
	{"min", []types.XValue{}, ERROR},

	{"mean", []types.XValue{xs("10"), xs("11")}, xn("10.5")},
	{"mean", []types.XValue{xs("10.2")}, xn("10.2")},
	{"mean", []types.XValue{xs("not_num")}, ERROR},
	{"mean", []types.XValue{xs("9"), xs("not_num")}, ERROR},
	{"mean", []types.XValue{}, ERROR},

	{"mod", []types.XValue{xs("10"), xs("3")}, xi(1)},
	{"mod", []types.XValue{xs("10"), xs("5")}, xi(0)},
	{"mod", []types.XValue{xs("not_num"), xs("3")}, ERROR},
	{"mod", []types.XValue{xs("9"), xs("not_num")}, ERROR},
	{"mod", []types.XValue{}, ERROR},

	{"read_code", []types.XValue{xs("123456")}, xs("1 2 3 , 4 5 6")},
	{"read_code", []types.XValue{xs("abcd")}, xs("a b c d")},
	{"read_code", []types.XValue{xs("12345678")}, xs("1 2 3 4 , 5 6 7 8")},
	{"read_code", []types.XValue{xs("12")}, xs("1 , 2")},
	{"read_code", []types.XValue{}, ERROR},

	{"split", []types.XValue{xs("1,2,3"), xs(",")}, types.NewXArray(xs("1"), xs("2"), xs("3"))},
	{"split", []types.XValue{xs("1,2,3"), xs(".")}, types.NewXArray(xs("1,2,3"))},
	{"split", []types.XValue{xs("1,2,3"), nil}, types.NewXArray(xs("1"), xs(","), xs("2"), xs(","), xs("3"))},
	{"split", []types.XValue{}, ERROR},

	{"join", []types.XValue{types.NewXArray(xs("1"), xs("2"), xs("3")), xs(",")}, xs("1,2,3")},
	{"join", []types.XValue{types.NewXArray(), xs(",")}, xs("")},
	{"join", []types.XValue{types.NewXArray(xs("1")), xs(",")}, xs("1")},
	{"join", []types.XValue{xs("1,2,3"), nil}, ERROR},
	{"join", []types.XValue{types.NewXArray(xs("1,2,3")), nil}, xs("1,2,3")},
	{"join", []types.XValue{types.NewXArray(xs("1"))}, ERROR},

	{"title", []types.XValue{xs("hello")}, xs("Hello")},
	{"title", []types.XValue{xs("")}, xs("")},
	{"title", []types.XValue{nil}, xs("")},
	{"title", []types.XValue{}, ERROR},

	{"word", []types.XValue{xs("hello World"), xn("1.5")}, xs("World")},
	{"word", []types.XValue{xs(""), xi(0)}, ERROR},
	{"word", []types.XValue{xs("😁 hello World"), xi(0)}, xs("😁")},
	{"word", []types.XValue{xs(" hello World"), xi(2)}, ERROR},
	{"word", []types.XValue{xs("hello World"), nil}, xs("hello")},
	{"word", []types.XValue{}, ERROR},

	{"remove_first_word", []types.XValue{xs("hello World")}, xs("World")},
	{"remove_first_word", []types.XValue{xs("hello")}, xs("")},
	{"remove_first_word", []types.XValue{xs("😁hello")}, xs("hello")},
	{"remove_first_word", []types.XValue{xs("")}, xs("")},
	{"remove_first_word", []types.XValue{}, ERROR},

	{"word_count", []types.XValue{xs("hello World")}, xi(2)},
	{"word_count", []types.XValue{xs("hello")}, xi(1)},
	{"word_count", []types.XValue{xs("")}, xi(0)},
	{"word_count", []types.XValue{xs("😁😁")}, xi(2)},
	{"word_count", []types.XValue{}, ERROR},

	{"field", []types.XValue{xs("hello,World"), xs("1"), xs(",")}, xs("World")},
	{"field", []types.XValue{xs("hello,world"), xn("2.1"), xs(",")}, xs("")},
	{"field", []types.XValue{xs("hello"), xi(0), xs(",")}, xs("hello")},
	{"field", []types.XValue{xs("hello,World"), xn("-2"), xs(",")}, ERROR},
	{"field", []types.XValue{xs(""), xs("notnum"), xs(",")}, ERROR},
	{"field", []types.XValue{xs("hello"), xi(0), nil}, xs("h")},

	{"clean", []types.XValue{xs("hello")}, xs("hello")},
	{"clean", []types.XValue{xs("  hello  world\n\t")}, xs("hello  world")},
	{"clean", []types.XValue{xs("")}, xs("")},
	{"clean", []types.XValue{}, ERROR},

	{"lower", []types.XValue{xs("HEllo")}, xs("hello")},
	{"lower", []types.XValue{xs("  HELLO  WORLD")}, xs("  hello  world")},
	{"lower", []types.XValue{xs("")}, xs("")},
	{"lower", []types.XValue{xs("😁")}, xs("😁")},
	{"lower", []types.XValue{}, ERROR},

	{"left", []types.XValue{xs("hello"), xs("2")}, xs("he")},
	{"left", []types.XValue{xs("  HELLO"), xs("2")}, xs("  ")},
	{"left", []types.XValue{xs("hi"), xi(4)}, xs("hi")},
	{"left", []types.XValue{xs("hi"), xs("0")}, xs("")},
	{"left", []types.XValue{xs("😁hi"), xs("2")}, xs("😁h")},
	{"left", []types.XValue{xs("hello"), nil}, xs("")},
	{"left", []types.XValue{}, ERROR},

	{"right", []types.XValue{xs("hello"), xs("2")}, xs("lo")},
	{"right", []types.XValue{xs("  HELLO "), xs("2")}, xs("O ")},
	{"right", []types.XValue{xs("hi"), xi(4)}, xs("hi")},
	{"right", []types.XValue{xs("hi"), xs("0")}, xs("")},
	{"right", []types.XValue{xs("ho😁hi"), xs("4")}, xs("o😁hi")},
	{"right", []types.XValue{nil, xs("2")}, xs("")},
	{"right", []types.XValue{xs("hello"), nil}, xs("")},
	{"right", []types.XValue{}, ERROR},

	{"length", []types.XValue{xs("hello")}, xi(5)},
	{"length", []types.XValue{xs("")}, xi(0)},
	{"length", []types.XValue{xs("😁😁")}, xi(2)},
	{"length", []types.XValue{types.NewXArray(xs("hello"))}, xi(1)},
	{"length", []types.XValue{types.NewXArray()}, xi(0)},
	{"length", []types.XValue{}, ERROR},

	{"string_cmp", []types.XValue{xs("abc"), xs("abc")}, xi(0)},
	{"string_cmp", []types.XValue{xs("abc"), xs("def")}, xi(-1)},
	{"string_cmp", []types.XValue{xs("def"), xs("abc")}, xi(1)},
	{"string_cmp", []types.XValue{xs("abc"), types.NewXErrorf("error")}, ERROR},
	{"string_cmp", []types.XValue{}, ERROR},

	{"default", []types.XValue{xs("10"), xs("20")}, xs("10")},
	{"default", []types.XValue{nil, xs("20")}, xs("20")},
	{"default", []types.XValue{types.NewXErrorf("This is error"), xs("20")}, xs("20")},
	{"default", []types.XValue{}, ERROR},

	{"repeat", []types.XValue{xs("hi"), xs("2")}, xs("hihi")},
	{"repeat", []types.XValue{xs("  "), xs("2")}, xs("    ")},
	{"repeat", []types.XValue{xs(""), xi(4)}, xs("")},
	{"repeat", []types.XValue{xs("😁"), xs("2")}, xs("😁😁")},
	{"repeat", []types.XValue{xs("hi"), xs("0")}, xs("")},
	{"repeat", []types.XValue{xs("hi"), xs("-1")}, ERROR},
	{"repeat", []types.XValue{xs("hello"), nil}, xs("")},
	{"repeat", []types.XValue{}, ERROR},

	{"replace", []types.XValue{xs("hi ho"), xs("hi"), xs("bye")}, xs("bye ho")},
	{"replace", []types.XValue{xs("foo bar "), xs(" "), xs(".")}, xs("foo.bar.")},
	{"replace", []types.XValue{xs("foo 😁 bar "), xs("😁"), xs("😂")}, xs("foo 😂 bar ")},
	{"replace", []types.XValue{xs("foo bar"), xs("zap"), xs("zog")}, xs("foo bar")},
	{"replace", []types.XValue{nil, xs("foo bar"), xs("foo")}, xs("")},
	{"replace", []types.XValue{xs("foo bar"), nil, xs("foo")}, xs("fooffooofooofoo foobfooafoorfoo")},
	{"replace", []types.XValue{xs("foo bar"), xs("foo"), nil}, xs(" bar")},
	{"replace", []types.XValue{}, ERROR},

	{"upper", []types.XValue{xs("HEllo")}, xs("HELLO")},
	{"upper", []types.XValue{xs("  HELLO  world")}, xs("  HELLO  WORLD")},
	{"upper", []types.XValue{xs("ß")}, xs("ß")},
	{"upper", []types.XValue{xs("")}, xs("")},
	{"upper", []types.XValue{}, ERROR},

	{"percent", []types.XValue{xs(".54")}, xs("54%")},
	{"percent", []types.XValue{xs("1.246")}, xs("125%")},
	{"percent", []types.XValue{xs("")}, ERROR},
	{"percent", []types.XValue{}, ERROR},

	{"date", []types.XValue{xs("01-12-2017")}, xd(time.Date(2017, 12, 1, 0, 0, 0, 0, time.UTC))},
	{"date", []types.XValue{xs("01-12-2017 10:15pm")}, xd(time.Date(2017, 12, 1, 22, 15, 0, 0, time.UTC))},
	{"date", []types.XValue{xs("01.15.2017")}, ERROR}, // month out of range
	{"date", []types.XValue{xs("no date")}, ERROR},    // invalid date
	{"date", []types.XValue{}, ERROR},

	{"format_date", []types.XValue{xs("1977-06-23T15:34:00.000000Z")}, xs("23-06-1977 15:34:00")},
	{"format_date", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DDTtt:mm:ss.fffZZZ"), xs("America/Los_Angeles")}, xs("1977-06-23T08:34:00.000-07:00")},
	{"format_date", []types.XValue{xs("1977-06-23T15:34:00.123000Z"), xs("YYYY-MM-DDTtt:mm:ss.fffZ"), xs("America/Los_Angeles")}, xs("1977-06-23T08:34:00.123-07:00")},
	{"format_date", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DDTtt:mm:ss.ffffffZ"), xs("America/Los_Angeles")}, xs("1977-06-23T08:34:00.000000-07:00")},
	{"format_date", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YY-MM-DD h:mm:ss AA"), xs("America/Los_Angeles")}, xs("77-06-23 8:34:00 AM")},
	{"format_date", []types.XValue{xs("1977-06-23T08:34:00.000-07:00"), xs("YYYY-MM-DDTtt:mm:ss.fffZ"), xs("UTC")}, xs("1977-06-23T15:34:00.000Z")},

	{"parse_date", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DDTtt:mm:ss.ffffffZ"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 8, 34, 0, 0, la))},
	{"parse_date", []types.XValue{xs("1977-06-23T15:34:00.1234Z"), xs("YYYY-MM-DDTtt:mm:ssZ"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 8, 34, 0, 123400000, la))},
	{"parse_date", []types.XValue{xs("1977-06-23 15:34"), xs("YYYY-MM-DD tt:mm"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 15, 34, 0, 0, la))},
	{"parse_date", []types.XValue{xs("1977-06-23 03:34 pm"), xs("YYYY-MM-DD tt:mm aa"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 15, 34, 0, 0, la))},
	{"parse_date", []types.XValue{xs("1977-06-23 03:34 PM"), xs("YYYY-MM-DD tt:mm AA"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 15, 34, 0, 0, la))},

	{"date_diff", []types.XValue{xs("03-12-2017"), xs("01-12-2017"), xs("d")}, xi(2)},
	{"date_diff", []types.XValue{xs("03-12-2017 10:15"), xs("03-12-2017 18:15"), xs("d")}, xi(0)},
	{"date_diff", []types.XValue{xs("03-12-2017"), xs("01-12-2017"), xs("w")}, xi(0)},
	{"date_diff", []types.XValue{xs("22-12-2017"), xs("01-12-2017"), xs("w")}, xi(3)},
	{"date_diff", []types.XValue{xs("03-12-2017"), xs("03-12-2017"), xs("M")}, xi(0)},
	{"date_diff", []types.XValue{xs("01-05-2018"), xs("03-12-2017"), xs("M")}, xi(5)},
	{"date_diff", []types.XValue{xs("01-12-2018"), xs("03-12-2017"), xs("y")}, xi(1)},
	{"date_diff", []types.XValue{xs("01-01-2017"), xs("03-12-2017"), xs("y")}, xi(0)},
	{"date_diff", []types.XValue{xs("04-12-2018 10:15"), xs("03-12-2018 14:00"), xs("h")}, xi(20)},
	{"date_diff", []types.XValue{xs("04-12-2018 10:15"), xs("04-12-2018 14:00"), xs("h")}, xi(-3)},
	{"date_diff", []types.XValue{xs("04-12-2018 10:15"), xs("04-12-2018 14:00"), xs("m")}, xi(-225)},
	{"date_diff", []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:15:35"), xs("m")}, xi(0)},
	{"date_diff", []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:16:10"), xs("m")}, xi(0)},
	{"date_diff", []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:15:35"), xs("s")}, xi(-20)},
	{"date_diff", []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:16:10"), xs("s")}, xi(-55)},
	{"date_diff", []types.XValue{xs("03-12-2017"), xs("01-12-2017"), xs("Z")}, ERROR},
	{"date_diff", []types.XValue{xs("xxx"), xs("01-12-2017"), xs("y")}, ERROR},
	{"date_diff", []types.XValue{xs("01-12-2017"), xs("xxx"), xs("y")}, ERROR},
	{"date_diff", []types.XValue{xs("01-12-2017"), xs("01-12-2017"), xs("xxx")}, ERROR},
	{"date_diff", []types.XValue{}, ERROR},

	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("y")}, xd(time.Date(2019, 12, 03, 22, 15, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("y")}, xd(time.Date(2015, 12, 03, 22, 15, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("M")}, xd(time.Date(2018, 2, 03, 22, 15, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("M")}, xd(time.Date(2017, 10, 3, 22, 15, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("w")}, xd(time.Date(2017, 12, 17, 22, 15, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("w")}, xd(time.Date(2017, 11, 19, 22, 15, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017"), xs("2"), xs("d")}, xd(time.Date(2017, 12, 5, 0, 0, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017"), xs("-4"), xs("d")}, xd(time.Date(2017, 11, 29, 0, 0, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("h")}, xd(time.Date(2017, 12, 4, 0, 15, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("h")}, xd(time.Date(2017, 12, 3, 20, 15, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("105"), xs("m")}, xd(time.Date(2017, 12, 4, 0, 0, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-20"), xs("m")}, xd(time.Date(2017, 12, 3, 21, 55, 0, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("s")}, xd(time.Date(2017, 12, 3, 22, 15, 2, 0, time.UTC))},
	{"date_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("s")}, xd(time.Date(2017, 12, 3, 22, 14, 58, 0, time.UTC))},
	{"date_add", []types.XValue{xs("xxx"), xs("2"), xs("d")}, ERROR},
	{"date_add", []types.XValue{xs("03-12-2017 10:15"), xs("xxx"), xs("D")}, ERROR},
	{"date_add", []types.XValue{xs("03-12-2017 10:15"), xs("2"), xs("xxx")}, ERROR},
	{"date_add", []types.XValue{xs("03-12-2017"), xs("2"), xs("Z")}, ERROR},
	{"date_add", []types.XValue{xs("22-12-2017")}, ERROR},

	{"weekday", []types.XValue{xs("01-12-2017")}, xi(5)},
	{"weekday", []types.XValue{xs("01-12-2017 10:15pm")}, xi(5)},
	{"weekday", []types.XValue{xs("xxx")}, ERROR},
	{"weekday", []types.XValue{}, ERROR},

	{"tz", []types.XValue{xs("01-12-2017")}, xs("UTC")},
	{"tz", []types.XValue{xs("01-12-2017 10:15:33pm")}, xs("UTC")},
	{"tz", []types.XValue{xs("xxx")}, ERROR},
	{"tz", []types.XValue{}, ERROR},

	{"tz_offset", []types.XValue{xs("01-12-2017")}, xs("+0000")},
	{"tz_offset", []types.XValue{xs("01-12-2017 10:15:33pm")}, xs("+0000")},
	{"tz_offset", []types.XValue{xs("xxx")}, ERROR},
	{"tz_offset", []types.XValue{}, ERROR},

	{"legacy_add", []types.XValue{xs("01-12-2017"), xi(2)}, xd(time.Date(2017, 12, 3, 0, 0, 0, 0, time.UTC))},
	{"legacy_add", []types.XValue{xs("2"), xs("01-12-2017 10:15:33pm")}, xd(time.Date(2017, 12, 3, 22, 15, 33, 0, time.UTC))},
	{"legacy_add", []types.XValue{xs("2"), xs("3.5")}, xn("5.5")},
	{"legacy_add", []types.XValue{xs("01-12-2017 10:15:33pm"), xs("01-12-2017")}, ERROR},
	{"legacy_add", []types.XValue{types.NewXNumberFromInt64(int64(math.MaxInt32 + 1)), xs("01-12-2017 10:15:33pm")}, ERROR},
	{"legacy_add", []types.XValue{xs("01-12-2017 10:15:33pm"), types.NewXNumberFromInt64(int64(math.MaxInt32 + 1))}, ERROR},
	{"legacy_add", []types.XValue{xs("xxx"), xs("10")}, ERROR},
	{"legacy_add", []types.XValue{xs("10"), xs("xxx")}, ERROR},
	{"legacy_add", []types.XValue{}, ERROR},

	{"format_urn", []types.XValue{xs("tel:+250781234567")}, xs("0781 234 567")},
	{"format_urn", []types.XValue{types.NewXArray(xs("tel:+250781112222"), xs("tel:+250781234567"))}, xs("0781 112 222")},
	{"format_urn", []types.XValue{xs("twitter:134252511151#billy_bob")}, xs("billy_bob")},
	{"format_urn", []types.XValue{xs("NOT URN")}, ERROR},
}

func TestFunctions(t *testing.T) {
	env := utils.NewEnvironment(utils.DateFormatDayMonthYear, utils.TimeFormatHourMinuteSecond, time.UTC, utils.LanguageList{})

	utils.SetRand(utils.NewSeededRand(123456))
	defer utils.SetRand(utils.DefaultRand)

	for _, test := range funcTests {
		xFunc := functions.XFUNCTIONS[test.name]
		result := xFunc(env, test.args...)

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic running function %s(%#v): %#v", test.name, test.args, r)
			}
		}()

		// don't check error equality - just check that we got an error if we expected one
		errExpected, _ := test.expected.(types.XError)
		errReturned, _ := result.(types.XError)
		if errExpected != nil && errReturned != nil {
			continue
		}

		cmp, err := types.Compare(result, test.expected)
		if err != nil {
			assert.Fail(t, err.Error(), "error while comparing expected: '%s' with result: '%s': %v for function %s(%#v)", test.expected, result, err, test.name, test.args)
		}
		if cmp != 0 {
			assert.Fail(t, "", "unexpected value, expected '%v', got '%v' for function %s(%#v)", test.expected, result, test.name, test.args)
		}
	}
}
