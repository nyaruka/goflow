package functions_test

import (
	"github.com/stretchr/testify/require"
	"math"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

var errorArg = types.NewXErrorf("I am error")
var la, _ = time.LoadLocation("America/Los_Angeles")

var xs = types.NewXText
var xn = types.RequireXNumberFromString
var xi = types.NewXNumberFromInt
var xd = types.NewXDateTime

var ERROR = types.NewXErrorf("any error")

var funcTests = []struct {
	name     string
	args     []types.XValue
	expected types.XValue
}{
	// tests for functions A-Z

	{"abs", []types.XValue{xi(33)}, xi(33)},
	{"abs", []types.XValue{xi(-33)}, xi(33)},
	{"abs", []types.XValue{xs("nan")}, ERROR},
	{"abs", []types.XValue{ERROR}, ERROR},
	{"abs", []types.XValue{}, ERROR},

	{"and", []types.XValue{types.XBooleanTrue}, types.XBooleanTrue},
	{"and", []types.XValue{types.XBooleanFalse}, types.XBooleanFalse},
	{"and", []types.XValue{types.XBooleanTrue, types.XBooleanFalse}, types.XBooleanFalse},
	{"and", []types.XValue{ERROR}, ERROR},
	{"and", []types.XValue{}, ERROR},

	{"array", []types.XValue{}, types.NewXArray()},
	{"array", []types.XValue{xi(123), xs("abc")}, types.NewXArray(xi(123), xs("abc"))},
	{"array", []types.XValue{xi(123), ERROR, xs("abc")}, ERROR},

	{"boolean", []types.XValue{xs("abc")}, types.XBooleanTrue},
	{"boolean", []types.XValue{xs("false")}, types.XBooleanFalse},
	{"boolean", []types.XValue{xs("FALSE")}, types.XBooleanFalse},
	{"boolean", []types.XValue{types.NewXArray()}, types.XBooleanFalse},
	{"boolean", []types.XValue{types.NewXArray(xi(1))}, types.XBooleanTrue},
	{"boolean", []types.XValue{ERROR}, ERROR},
	{"boolean", []types.XValue{}, ERROR},

	{"char", []types.XValue{xn("33")}, xs("!")},
	{"char", []types.XValue{xn("128513")}, xs("😁")},
	{"char", []types.XValue{xs("not a number")}, ERROR},
	{"char", []types.XValue{xn("12345678901234567890")}, ERROR},
	{"char", []types.XValue{}, ERROR},

	{"code", []types.XValue{xs(" ")}, xi(32)},
	{"code", []types.XValue{xs("😁")}, xi(128513)},
	{"code", []types.XValue{xs("abc")}, xi(97)},
	{"code", []types.XValue{xs("")}, ERROR},
	{"code", []types.XValue{ERROR}, ERROR},
	{"code", []types.XValue{}, ERROR},

	{"clean", []types.XValue{xs("hello")}, xs("hello")},
	{"clean", []types.XValue{xs("😃 Hello \nwo\tr\rld")}, xs("😃 Hello world")},
	{"clean", []types.XValue{xs("")}, xs("")},
	{"clean", []types.XValue{}, ERROR},

	{"datetime", []types.XValue{xs("01-12-2017")}, xd(time.Date(2017, 12, 1, 0, 0, 0, 0, time.UTC))},
	{"datetime", []types.XValue{xs("01-12-2017 10:15pm")}, xd(time.Date(2017, 12, 1, 22, 15, 0, 0, time.UTC))},
	{"datetime", []types.XValue{xs("01.15.2017")}, ERROR}, // month out of range
	{"datetime", []types.XValue{xs("no date")}, ERROR},    // invalid date
	{"datetime", []types.XValue{}, ERROR},

	{"datetime_from_parts", []types.XValue{xi(2018), xi(11), xi(3)}, xd(time.Date(2018, 11, 3, 0, 0, 0, 0, time.UTC))},
	{"datetime_from_parts", []types.XValue{xi(2018), xi(15), xi(3)}, ERROR}, // month out of range
	{"datetime_from_parts", []types.XValue{ERROR, xi(11), xi(3)}, ERROR},
	{"datetime_from_parts", []types.XValue{xi(2018), ERROR, xi(3)}, ERROR},
	{"datetime_from_parts", []types.XValue{xi(2018), xi(11), ERROR}, ERROR},
	{"datetime_from_parts", []types.XValue{}, ERROR},

	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("Y")}, xd(time.Date(2019, 12, 03, 22, 15, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("Y")}, xd(time.Date(2015, 12, 03, 22, 15, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("M")}, xd(time.Date(2018, 2, 03, 22, 15, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("M")}, xd(time.Date(2017, 10, 3, 22, 15, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("W")}, xd(time.Date(2017, 12, 17, 22, 15, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("W")}, xd(time.Date(2017, 11, 19, 22, 15, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017"), xs("2"), xs("D")}, xd(time.Date(2017, 12, 5, 0, 0, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017"), xs("-4"), xs("D")}, xd(time.Date(2017, 11, 29, 0, 0, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("h")}, xd(time.Date(2017, 12, 4, 0, 15, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("h")}, xd(time.Date(2017, 12, 3, 20, 15, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("105"), xs("m")}, xd(time.Date(2017, 12, 4, 0, 0, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-20"), xs("m")}, xd(time.Date(2017, 12, 3, 21, 55, 0, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("s")}, xd(time.Date(2017, 12, 3, 22, 15, 2, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("s")}, xd(time.Date(2017, 12, 3, 22, 14, 58, 0, time.UTC))},
	{"datetime_add", []types.XValue{xs("xxx"), xs("2"), xs("D")}, ERROR},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15"), xs("xxx"), xs("D")}, ERROR},
	{"datetime_add", []types.XValue{xs("03-12-2017 10:15"), xs("2"), xs("xxx")}, ERROR},
	{"datetime_add", []types.XValue{xs("03-12-2017"), xs("2"), xs("Z")}, ERROR},
	{"datetime_add", []types.XValue{xs("03-12-2017"), xs("2"), ERROR}, ERROR},
	{"datetime_add", []types.XValue{xs("22-12-2017")}, ERROR},

	{"datetime_diff", []types.XValue{xs("03-12-2017"), xs("01-12-2017"), xs("D")}, xi(2)},
	{"datetime_diff", []types.XValue{xs("03-12-2017 10:15"), xs("03-12-2017 18:15"), xs("D")}, xi(0)},
	{"datetime_diff", []types.XValue{xs("03-12-2017"), xs("01-12-2017"), xs("W")}, xi(0)},
	{"datetime_diff", []types.XValue{xs("22-12-2017"), xs("01-12-2017"), xs("W")}, xi(3)},
	{"datetime_diff", []types.XValue{xs("03-12-2017"), xs("03-12-2017"), xs("M")}, xi(0)},
	{"datetime_diff", []types.XValue{xs("01-05-2018"), xs("03-12-2017"), xs("M")}, xi(5)},
	{"datetime_diff", []types.XValue{xs("01-12-2018"), xs("03-12-2017"), xs("Y")}, xi(1)},
	{"datetime_diff", []types.XValue{xs("01-01-2017"), xs("03-12-2017"), xs("Y")}, xi(0)},
	{"datetime_diff", []types.XValue{xs("04-12-2018 10:15"), xs("03-12-2018 14:00"), xs("h")}, xi(20)},
	{"datetime_diff", []types.XValue{xs("04-12-2018 10:15"), xs("04-12-2018 14:00"), xs("h")}, xi(-3)},
	{"datetime_diff", []types.XValue{xs("04-12-2018 10:15"), xs("04-12-2018 14:00"), xs("m")}, xi(-225)},
	{"datetime_diff", []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:15:35"), xs("m")}, xi(0)},
	{"datetime_diff", []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:16:10"), xs("m")}, xi(0)},
	{"datetime_diff", []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:15:35"), xs("s")}, xi(-20)},
	{"datetime_diff", []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:16:10"), xs("s")}, xi(-55)},
	{"datetime_diff", []types.XValue{xs("03-12-2017"), xs("01-12-2017"), xs("Z")}, ERROR},
	{"datetime_diff", []types.XValue{xs("xxx"), xs("01-12-2017"), xs("Y")}, ERROR},
	{"datetime_diff", []types.XValue{xs("01-12-2017"), xs("xxx"), xs("Y")}, ERROR},
	{"datetime_diff", []types.XValue{xs("01-12-2017"), xs("01-12-2017"), xs("xxx")}, ERROR},
	{"datetime_diff", []types.XValue{xs("01-12-2017"), xs("01-12-2017"), ERROR}, ERROR},
	{"datetime_diff", []types.XValue{}, ERROR},

	{"default", []types.XValue{xs("10"), xs("20")}, xs("10")},
	{"default", []types.XValue{nil, xs("20")}, xs("20")},
	{"default", []types.XValue{types.NewXErrorf("This is error"), xs("20")}, xs("20")},
	{"default", []types.XValue{}, ERROR},

	{"field", []types.XValue{xs("hello,World"), xs("1"), xs(",")}, xs("World")},
	{"field", []types.XValue{xs("hello,world"), xn("2.1"), xs(",")}, xs("")},
	{"field", []types.XValue{xs("hello world there now"), xn("2"), xs(" ")}, xs("there")},
	{"field", []types.XValue{xs("hello"), xi(0), xs(",")}, xs("hello")},
	{"field", []types.XValue{xs("hello,World"), xn("-2"), xs(",")}, ERROR},
	{"field", []types.XValue{xs(""), xs("notnum"), xs(",")}, ERROR},
	{"field", []types.XValue{xs("hello"), xi(0), nil}, xs("h")},
	{"field", []types.XValue{ERROR, xs("1"), xs(",")}, ERROR},
	{"field", []types.XValue{xs("hello"), ERROR, xs(",")}, ERROR},
	{"field", []types.XValue{xs("hello"), xs("1"), ERROR}, ERROR},
	{"field", []types.XValue{}, ERROR},

	{"format_datetime", []types.XValue{xs("1977-06-23T15:34:00.000000Z")}, xs("23-06-1977 15:34")},
	{"format_datetime", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DDTtt:mm:ss.fffZZZ"), xs("America/Los_Angeles")}, xs("1977-06-23T08:34:00.000-07:00")},
	{"format_datetime", []types.XValue{xs("1977-06-23T15:34:00.123000Z"), xs("YYYY-MM-DDTtt:mm:ss.fffZ"), xs("America/Los_Angeles")}, xs("1977-06-23T08:34:00.123-07:00")},
	{"format_datetime", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DDTtt:mm:ss.ffffffZ"), xs("America/Los_Angeles")}, xs("1977-06-23T08:34:00.000000-07:00")},
	{"format_datetime", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YY-MM-DD h:mm:ss AA"), xs("America/Los_Angeles")}, xs("77-06-23 8:34:00 AM")},
	{"format_datetime", []types.XValue{xs("1977-06-23T08:34:00.000-07:00"), xs("YYYY-MM-DDTtt:mm:ss.fffZ"), xs("UTC")}, xs("1977-06-23T15:34:00.000Z")},
	{"format_datetime", []types.XValue{xs("NOT DATE")}, ERROR},
	{"format_datetime", []types.XValue{ERROR}, ERROR},
	{"format_datetime", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), ERROR}, ERROR},
	{"format_datetime", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYYYYY")}, ERROR},
	{"format_datetime", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY"), ERROR}, ERROR},
	{"format_datetime", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY"), xs("Cuenca")}, ERROR},
	{"format_datetime", []types.XValue{}, ERROR},

	{"format_location", []types.XValue{xs("Rwanda")}, xs("Rwanda")},
	{"format_location", []types.XValue{xs("Rwanda > Kigali")}, xs("Kigali")},
	{"format_location", []types.XValue{ERROR}, ERROR},
	{"format_location", []types.XValue{}, ERROR},

	{"format_number", []types.XValue{xn("31337")}, xs("31,337.00")},
	{"format_number", []types.XValue{xn("31337"), xi(0), types.XBooleanFalse}, xs("31337")},
	{"format_number", []types.XValue{xn("31337"), xs("xxx")}, ERROR},
	{"format_number", []types.XValue{xn("31337"), xi(12345)}, ERROR},
	{"format_number", []types.XValue{xn("31337"), xi(2), ERROR}, ERROR},
	{"format_number", []types.XValue{ERROR}, ERROR},
	{"format_number", []types.XValue{}, ERROR},

	{"format_urn", []types.XValue{xs("tel:+250781234567")}, xs("0781 234 567")},
	{"format_urn", []types.XValue{types.NewXArray(xs("tel:+250781112222"), xs("tel:+250781234567"))}, xs("0781 112 222")},
	{"format_urn", []types.XValue{xs("twitter:134252511151#billy_bob")}, xs("billy_bob")},
	{"format_urn", []types.XValue{types.NewXArray()}, xs("")},
	{"format_urn", []types.XValue{xs("NOT URN")}, ERROR},
	{"format_urn", []types.XValue{ERROR}, ERROR},
	{"format_urn", []types.XValue{}, ERROR},

	{"from_epoch", []types.XValue{xn("1497286619000000000")}, xd(time.Date(2017, 6, 12, 16, 56, 59, 0, time.UTC))},
	{"from_epoch", []types.XValue{ERROR}, ERROR},
	{"from_epoch", []types.XValue{}, ERROR},

	{"if", []types.XValue{types.XBooleanTrue, xs("10"), xs("20")}, xs("10")},
	{"if", []types.XValue{types.XBooleanFalse, xs("10"), xs("20")}, xs("20")},
	{"if", []types.XValue{types.XBooleanTrue, errorArg, xs("20")}, errorArg},
	{"if", []types.XValue{}, ERROR},
	{"if", []types.XValue{errorArg, xs("10"), xs("20")}, errorArg},

	{"join", []types.XValue{types.NewXArray(xs("1"), xs("2"), xs("3")), xs(",")}, xs("1,2,3")},
	{"join", []types.XValue{types.NewXArray(), xs(",")}, xs("")},
	{"join", []types.XValue{types.NewXArray(xs("1")), xs(",")}, xs("1")},
	{"join", []types.XValue{types.NewXArray(xs("1"), xs("2")), ERROR}, ERROR},
	{"join", []types.XValue{types.NewXArray(xs("1"), ERROR), xs(",")}, ERROR},
	{"join", []types.XValue{xs("1,2,3"), nil}, ERROR},
	{"join", []types.XValue{types.NewXArray(xs("1,2,3")), nil}, xs("1,2,3")},
	{"join", []types.XValue{types.NewXArray(xs("1"))}, ERROR},

	{"json", []types.XValue{xs("hello")}, xs(`"hello"`)},
	{"json", []types.XValue{ERROR}, ERROR},

	{"left", []types.XValue{xs("hello"), xs("2")}, xs("he")},
	{"left", []types.XValue{xs("  HELLO"), xs("2")}, xs("  ")},
	{"left", []types.XValue{xs("hi"), xi(4)}, xs("hi")},
	{"left", []types.XValue{xs("hi"), xs("0")}, xs("")},
	{"left", []types.XValue{xs("😁hi"), xs("2")}, xs("😁h")},
	{"left", []types.XValue{xs("hello"), nil}, ERROR},
	{"left", []types.XValue{xs("hello"), xi(-1)}, ERROR},
	{"left", []types.XValue{ERROR, xi(3)}, ERROR},
	{"left", []types.XValue{xs("hello"), ERROR}, ERROR},
	{"left", []types.XValue{}, ERROR},

	{"legacy_add", []types.XValue{xs("01-12-2017"), xi(2)}, xd(time.Date(2017, 12, 3, 0, 0, 0, 0, time.UTC))},
	{"legacy_add", []types.XValue{xs("2"), xs("01-12-2017 10:15:33pm")}, xd(time.Date(2017, 12, 3, 22, 15, 33, 0, time.UTC))},
	{"legacy_add", []types.XValue{xs("2"), xs("3.5")}, xn("5.5")},
	{"legacy_add", []types.XValue{xs("01-12-2017 10:15:33pm"), xs("01-12-2017")}, ERROR},
	{"legacy_add", []types.XValue{types.NewXNumberFromInt64(int64(math.MaxInt32 + 1)), xs("01-12-2017 10:15:33pm")}, ERROR},
	{"legacy_add", []types.XValue{xs("01-12-2017 10:15:33pm"), types.NewXNumberFromInt64(int64(math.MaxInt32 + 1))}, ERROR},
	{"legacy_add", []types.XValue{xs("xxx"), xs("10")}, ERROR},
	{"legacy_add", []types.XValue{xs("10"), xs("xxx")}, ERROR},
	{"legacy_add", []types.XValue{}, ERROR},

	{"length", []types.XValue{xs("hello")}, xi(5)},
	{"length", []types.XValue{xs("")}, xi(0)},
	{"length", []types.XValue{xs("😁😁")}, xi(2)},
	{"length", []types.XValue{types.NewXArray(xs("hello"))}, xi(1)},
	{"length", []types.XValue{types.NewXArray()}, xi(0)},
	{"length", []types.XValue{xi(1234)}, ERROR},
	{"length", []types.XValue{ERROR}, ERROR},
	{"length", []types.XValue{}, ERROR},

	{"lower", []types.XValue{xs("HEllo")}, xs("hello")},
	{"lower", []types.XValue{xs("  HELLO  WORLD")}, xs("  hello  world")},
	{"lower", []types.XValue{xs("")}, xs("")},
	{"lower", []types.XValue{xs("😁")}, xs("😁")},
	{"lower", []types.XValue{}, ERROR},

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

	{"now", []types.XValue{}, xd(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC))},
	{"now", []types.XValue{ERROR}, ERROR},

	{"number", []types.XValue{xn("10")}, xn("10")},
	{"number", []types.XValue{xs("123.45000")}, xn("123.45")},
	{"number", []types.XValue{xs("what?")}, ERROR},

	{"or", []types.XValue{types.XBooleanTrue}, types.XBooleanTrue},
	{"or", []types.XValue{types.XBooleanFalse}, types.XBooleanFalse},
	{"or", []types.XValue{types.XBooleanTrue, types.XBooleanFalse}, types.XBooleanTrue},
	{"or", []types.XValue{ERROR}, ERROR},
	{"or", []types.XValue{}, ERROR},

	{"parse_datetime", []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DDTtt:mm:ss.ffffffZ"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 8, 34, 0, 0, la))},
	{"parse_datetime", []types.XValue{xs("1977-06-23T15:34:00.1234Z"), xs("YYYY-MM-DDTtt:mm:ssZ"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 8, 34, 0, 123400000, la))},
	{"parse_datetime", []types.XValue{xs("1977-06-23 15:34"), xs("YYYY-MM-DD tt:mm"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 15, 34, 0, 0, la))},
	{"parse_datetime", []types.XValue{xs("1977-06-23 03:34 pm"), xs("YYYY-MM-DD tt:mm aa"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 15, 34, 0, 0, la))},
	{"parse_datetime", []types.XValue{xs("1977-06-23 03:34 PM"), xs("YYYY-MM-DD tt:mm AA"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 15, 34, 0, 0, la))},
	{"parse_datetime", []types.XValue{xs("1977-06-23 15:34"), xs("ttttttttt")}, ERROR},                // invalid format
	{"parse_datetime", []types.XValue{xs("1977-06-23 15:34"), xs("YYYY-MM-DD"), xs("Cuenca")}, ERROR}, // invalid timezone
	{"parse_datetime", []types.XValue{xs("abcd"), xs("YYYY-MM-DD")}, ERROR},                           // unparseable date
	{"parse_datetime", []types.XValue{ERROR, xs("YYYY-MM-DD")}, ERROR},
	{"parse_datetime", []types.XValue{xs("1977-06-23 15:34"), ERROR}, ERROR},
	{"parse_datetime", []types.XValue{xs("1977-06-23 15:34"), xs("YYYY-MM-DD"), ERROR}, ERROR},
	{"parse_datetime", []types.XValue{}, ERROR},

	{"parse_json", []types.XValue{xs(`"hello"`)}, xs(`hello`)},
	{"parse_json", []types.XValue{ERROR}, ERROR},

	{"percent", []types.XValue{xs(".54")}, xs("54%")},
	{"percent", []types.XValue{xs("1.246")}, xs("125%")},
	{"percent", []types.XValue{xs("")}, ERROR},
	{"percent", []types.XValue{}, ERROR},

	{"rand", []types.XValue{}, xn("0.3849275689214193274523267973563633859157562255859375")},
	{"rand", []types.XValue{}, xn("0.607552015674623913099594574305228888988494873046875")},

	{"rand_between", []types.XValue{xn("1"), xn("10")}, xn("5")},
	{"rand_between", []types.XValue{xn("1"), xn("10")}, xn("10")},

	{"read_chars", []types.XValue{xs("123456")}, xs("1 2 3 , 4 5 6")},
	{"read_chars", []types.XValue{xs("abcd")}, xs("a b c d")},
	{"read_chars", []types.XValue{xs("12345678")}, xs("1 2 3 4 , 5 6 7 8")},
	{"read_chars", []types.XValue{xs("12")}, xs("1 , 2")},
	{"read_chars", []types.XValue{}, ERROR},

	{"remove_first_word", []types.XValue{xs("hello World")}, xs("World")},
	{"remove_first_word", []types.XValue{xs("hello")}, xs("")},
	{"remove_first_word", []types.XValue{xs("😁hello")}, xs("hello")},
	{"remove_first_word", []types.XValue{xs("")}, xs("")},
	{"remove_first_word", []types.XValue{}, ERROR},

	{"repeat", []types.XValue{xs("hi"), xs("2")}, xs("hihi")},
	{"repeat", []types.XValue{xs("  "), xs("2")}, xs("    ")},
	{"repeat", []types.XValue{xs(""), xi(4)}, xs("")},
	{"repeat", []types.XValue{xs("😁"), xs("2")}, xs("😁😁")},
	{"repeat", []types.XValue{xs("hi"), xs("0")}, xs("")},
	{"repeat", []types.XValue{xs("hi"), xs("-1")}, ERROR},
	{"repeat", []types.XValue{xs("hello"), nil}, ERROR},
	{"repeat", []types.XValue{}, ERROR},

	{"replace", []types.XValue{xs("hi ho"), xs("hi"), xs("bye")}, xs("bye ho")},
	{"replace", []types.XValue{xs("foo bar "), xs(" "), xs(".")}, xs("foo.bar.")},
	{"replace", []types.XValue{xs("foo 😁 bar "), xs("😁"), xs("😂")}, xs("foo 😂 bar ")},
	{"replace", []types.XValue{xs("foo bar"), xs("zap"), xs("zog")}, xs("foo bar")},
	{"replace", []types.XValue{nil, xs("foo bar"), xs("foo")}, xs("")},
	{"replace", []types.XValue{xs("foo bar"), nil, xs("foo")}, xs("fooffooofooofoo foobfooafoorfoo")},
	{"replace", []types.XValue{xs("foo bar"), xs("foo"), nil}, xs(" bar")},
	{"replace", []types.XValue{ERROR, xs("hi"), xs("bye")}, ERROR},
	{"replace", []types.XValue{xs("hi ho"), ERROR, xs("bye")}, ERROR},
	{"replace", []types.XValue{xs("hi ho"), xs("bye"), ERROR}, ERROR},
	{"replace", []types.XValue{}, ERROR},

	{"right", []types.XValue{xs("hello"), xs("2")}, xs("lo")},
	{"right", []types.XValue{xs("  HELLO "), xs("2")}, xs("O ")},
	{"right", []types.XValue{xs("hi"), xi(4)}, xs("hi")},
	{"right", []types.XValue{xs("hi"), xs("0")}, xs("")},
	{"right", []types.XValue{xs("ho😁hi"), xs("4")}, xs("o😁hi")},
	{"right", []types.XValue{nil, xs("2")}, xs("")},
	{"right", []types.XValue{xs("hello"), nil}, ERROR},
	{"right", []types.XValue{xs("hello"), xi(-1)}, ERROR},
	{"right", []types.XValue{ERROR, xi(3)}, ERROR},
	{"right", []types.XValue{xs("hello"), ERROR}, ERROR},
	{"right", []types.XValue{}, ERROR},

	{"round", []types.XValue{xs("10.5"), xs("0")}, xi(11)},
	{"round", []types.XValue{xs("10.5"), xs("1")}, xn("10.5")},
	{"round", []types.XValue{xs("10.51"), xs("1")}, xn("10.5")},
	{"round", []types.XValue{xs("10.56"), xs("1")}, xn("10.6")},
	{"round", []types.XValue{xs("12.56"), xs("-1")}, xi(10)},
	{"round", []types.XValue{xs("10.5")}, xn("11")},
	{"round", []types.XValue{xs("not_num"), xs("1")}, ERROR},
	{"round", []types.XValue{xs("10.5"), xs("not_num")}, ERROR},
	{"round", []types.XValue{xs("10.5"), xs("1"), xs("30")}, ERROR},

	{"round_down", []types.XValue{xs("10")}, xi(10)},
	{"round_down", []types.XValue{xs("10.5")}, xi(10)},
	{"round_down", []types.XValue{xs("10.7")}, xi(10)},
	{"round_down", []types.XValue{xs("not_num")}, ERROR},
	{"round_down", []types.XValue{}, ERROR},

	{"round_up", []types.XValue{xs("10")}, xi(10)},
	{"round_up", []types.XValue{xs("10.5")}, xi(11)},
	{"round_up", []types.XValue{xs("10.2")}, xi(11)},
	{"round_up", []types.XValue{xs("not_num")}, ERROR},
	{"round_up", []types.XValue{}, ERROR},

	{"split", []types.XValue{xs("1,2,3"), xs(",")}, types.NewXArray(xs("1"), xs("2"), xs("3"))},
	{"split", []types.XValue{xs("1,2,3"), xs(".")}, types.NewXArray(xs("1,2,3"))},
	{"split", []types.XValue{xs("1,2,3"), nil}, types.NewXArray(xs("1"), xs(","), xs("2"), xs(","), xs("3"))},
	{"split", []types.XValue{ERROR, xs(",")}, ERROR},
	{"split", []types.XValue{xs("1,2,3"), ERROR}, ERROR},
	{"split", []types.XValue{}, ERROR},

	{"text", []types.XValue{xs("abc")}, xs("abc")},
	{"text", []types.XValue{xi(123)}, xs("123")},
	{"text", []types.XValue{ERROR}, ERROR},
	{"text", []types.XValue{}, ERROR},

	{"text_compare", []types.XValue{xs("abc"), xs("abc")}, xi(0)},
	{"text_compare", []types.XValue{xs("abc"), xs("def")}, xi(-1)},
	{"text_compare", []types.XValue{xs("def"), xs("abc")}, xi(1)},
	{"text_compare", []types.XValue{xs("abc"), types.NewXErrorf("error")}, ERROR},
	{"text_compare", []types.XValue{}, ERROR},

	{"title", []types.XValue{xs("hello")}, xs("Hello")},
	{"title", []types.XValue{xs("")}, xs("")},
	{"title", []types.XValue{nil}, xs("")},
	{"title", []types.XValue{}, ERROR},

	{"to_epoch", []types.XValue{xd(time.Date(2017, 6, 12, 16, 56, 59, 0, time.UTC))}, xn("1497286619000000000")},
	{"to_epoch", []types.XValue{ERROR}, ERROR},
	{"to_epoch", []types.XValue{}, ERROR},

	{"today", []types.XValue{}, xd(time.Date(2018, 4, 11, 0, 0, 0, 0, time.UTC))},
	{"today", []types.XValue{ERROR}, ERROR},

	{"tz", []types.XValue{xs("01-12-2017")}, xs("UTC")},
	{"tz", []types.XValue{xs("01-12-2017 10:15:33pm")}, xs("UTC")},
	{"tz", []types.XValue{xs("xxx")}, ERROR},
	{"tz", []types.XValue{}, ERROR},

	{"tz_offset", []types.XValue{xs("01-12-2017")}, xs("+0000")},
	{"tz_offset", []types.XValue{xs("01-12-2017 10:15:33pm")}, xs("+0000")},
	{"tz_offset", []types.XValue{xs("xxx")}, ERROR},
	{"tz_offset", []types.XValue{}, ERROR},

	{"upper", []types.XValue{xs("HEllo")}, xs("HELLO")},
	{"upper", []types.XValue{xs("  HELLO  world")}, xs("  HELLO  WORLD")},
	{"upper", []types.XValue{xs("ß")}, xs("ß")},
	{"upper", []types.XValue{xs("")}, xs("")},
	{"upper", []types.XValue{}, ERROR},

	{"word", []types.XValue{xs("hello World"), xn("1.5")}, xs("World")},
	{"word", []types.XValue{xs(""), xi(0)}, ERROR},
	{"word", []types.XValue{xs("cat dog bee"), xi(-1)}, xs("bee")},
	{"word", []types.XValue{xs("😁 hello World"), xi(0)}, xs("😁")},
	{"word", []types.XValue{xs(" hello World"), xi(2)}, ERROR},
	{"word", []types.XValue{xs("hello World"), nil}, ERROR},
	{"word", []types.XValue{}, ERROR},

	{"word_slice", []types.XValue{xs("hello-world from mars"), xi(0), xi(2)}, xs("hello world")},
	{"word_slice", []types.XValue{xs("hello-world from mars"), xi(2)}, xs("from mars")},
	{"word_slice", []types.XValue{xs("hello-world from mars"), xi(10)}, xs("")},
	{"word_slice", []types.XValue{xs("hello-world from mars"), xi(3), xi(10)}, xs("mars")},
	{"word_slice", []types.XValue{xs("hello-world from mars"), xi(-1), xi(2)}, ERROR},
	{"word_slice", []types.XValue{xs("hello-world from mars"), xi(3), xi(1)}, ERROR},
	{"word_slice", []types.XValue{xs("hello-world from mars"), xs("x"), xi(3)}, ERROR},
	{"word_slice", []types.XValue{xs("hello-world from mars"), xi(3), xs("x")}, ERROR},
	{"word_slice", []types.XValue{xs("hello-world from mars"), ERROR, xi(2)}, ERROR},
	{"word_slice", []types.XValue{ERROR, xi(0), xi(2)}, ERROR},
	{"word_slice", []types.XValue{ERROR}, ERROR},

	{"word_count", []types.XValue{xs("hello World")}, xi(2)},
	{"word_count", []types.XValue{xs("hello")}, xi(1)},
	{"word_count", []types.XValue{xs("")}, xi(0)},
	{"word_count", []types.XValue{xs("😁😁")}, xi(2)},
	{"word_count", []types.XValue{}, ERROR},

	{"weekday", []types.XValue{xs("01-12-2017")}, xi(5)},
	{"weekday", []types.XValue{xs("01-12-2017 10:15pm")}, xi(5)},
	{"weekday", []types.XValue{xs("xxx")}, ERROR},
	{"weekday", []types.XValue{}, ERROR},

	{"url_encode", []types.XValue{xs(`hi-% ?/`)}, xs(`hi-%25+%3F%2F`)},
	{"url_encode", []types.XValue{ERROR}, ERROR},
	{"url_encode", []types.XValue{}, ERROR},
}

func TestFunctions(t *testing.T) {
	env := test.NewTestEnvironment(utils.DateFormatDayMonthYear, time.UTC, nil)

	utils.SetRand(utils.NewSeededRand(123456))
	defer utils.SetRand(utils.DefaultRand)

	for _, test := range funcTests {
		xFunc, exists := functions.XFUNCTIONS[test.name]
		require.True(t, exists, "no such registered function: %s", test.name)

		result := xFunc(env, test.args...)

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic running function %s(%#v): %#v", test.name, test.args, r)
			}
		}()

		// don't check error equality - just check that we got an error if we expected one
		if test.expected == ERROR {
			assert.True(t, types.IsXError(result), "expecting error, got %T{%s} for function %s(%T{%s})", result, result, test.name, test.args, test.args)
		} else {
			if !types.Equals(env, result, test.expected) {
				assert.Fail(t, "", "unexpected value, expected %T{%s}, got %T{%s} for function %s(%T{%s})", test.expected, test.expected, result, result, test.name, test.args, test.args)
			}
		}
	}
}
