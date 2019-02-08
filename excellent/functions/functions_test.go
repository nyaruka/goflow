package functions_test

import (
	"math"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errorArg = types.NewXErrorf("I am error")
var la, _ = time.LoadLocation("America/Los_Angeles")

var xs = types.NewXText
var xn = types.RequireXNumberFromString
var xi = types.NewXNumberFromInt
var xd = types.NewXDateTime

var ERROR = types.NewXErrorf("any error")

func TestFunctions(t *testing.T) {
	dmy := utils.NewEnvironment(utils.DateFormatDayMonthYear, utils.TimeFormatHourMinute, time.UTC, utils.NilLanguage, nil, utils.NilCountry, utils.DefaultNumberFormat, utils.RedactionPolicyNone, 640)
	mdy := utils.NewEnvironment(utils.DateFormatMonthDayYear, utils.TimeFormatHourMinute, la, utils.NilLanguage, nil, utils.NilCountry, utils.DefaultNumberFormat, utils.RedactionPolicyNone, 640)

	var funcTests = []struct {
		name     string
		env      utils.Environment
		args     []types.XValue
		expected types.XValue
	}{
		// tests for functions A-Z

		{"abs", dmy, []types.XValue{xi(33)}, xi(33)},
		{"abs", dmy, []types.XValue{xi(-33)}, xi(33)},
		{"abs", dmy, []types.XValue{xs("nan")}, ERROR},
		{"abs", dmy, []types.XValue{ERROR}, ERROR},
		{"abs", dmy, []types.XValue{}, ERROR},

		{"and", dmy, []types.XValue{types.XBooleanTrue}, types.XBooleanTrue},
		{"and", dmy, []types.XValue{types.XBooleanFalse}, types.XBooleanFalse},
		{"and", dmy, []types.XValue{types.XBooleanTrue, types.XBooleanFalse}, types.XBooleanFalse},
		{"and", dmy, []types.XValue{ERROR}, ERROR},
		{"and", dmy, []types.XValue{}, ERROR},

		{"array", dmy, []types.XValue{}, types.NewXArray()},
		{"array", dmy, []types.XValue{xi(123), xs("abc")}, types.NewXArray(xi(123), xs("abc"))},
		{"array", dmy, []types.XValue{xi(123), ERROR, xs("abc")}, ERROR},

		{"boolean", dmy, []types.XValue{xs("abc")}, types.XBooleanTrue},
		{"boolean", dmy, []types.XValue{xs("false")}, types.XBooleanFalse},
		{"boolean", dmy, []types.XValue{xs("FALSE")}, types.XBooleanFalse},
		{"boolean", dmy, []types.XValue{types.NewXArray()}, types.XBooleanFalse},
		{"boolean", dmy, []types.XValue{types.NewXArray(xi(1))}, types.XBooleanTrue},
		{"boolean", dmy, []types.XValue{ERROR}, ERROR},
		{"boolean", dmy, []types.XValue{}, ERROR},

		{"char", dmy, []types.XValue{xn("33")}, xs("!")},
		{"char", dmy, []types.XValue{xn("128513")}, xs("😁")},
		{"char", dmy, []types.XValue{xs("not a number")}, ERROR},
		{"char", dmy, []types.XValue{xn("12345678901234567890")}, ERROR},
		{"char", dmy, []types.XValue{}, ERROR},

		{"code", dmy, []types.XValue{xs(" ")}, xi(32)},
		{"code", dmy, []types.XValue{xs("😁")}, xi(128513)},
		{"code", dmy, []types.XValue{xs("abc")}, xi(97)},
		{"code", dmy, []types.XValue{xs("")}, ERROR},
		{"code", dmy, []types.XValue{ERROR}, ERROR},
		{"code", dmy, []types.XValue{}, ERROR},

		{"clean", dmy, []types.XValue{xs("hello")}, xs("hello")},
		{"clean", dmy, []types.XValue{xs("😃 Hello \nwo\tr\rld")}, xs("😃 Hello world")},
		{"clean", dmy, []types.XValue{xs("")}, xs("")},
		{"clean", dmy, []types.XValue{}, ERROR},

		{"datetime", dmy, []types.XValue{xs("01-12-2017")}, xd(time.Date(2017, 12, 1, 0, 0, 0, 0, time.UTC))},
		{"datetime", mdy, []types.XValue{xs("12-01-2017")}, xd(time.Date(2017, 12, 1, 0, 0, 0, 0, la))},
		{"datetime", dmy, []types.XValue{xs("01-12-2017 10:15pm")}, xd(time.Date(2017, 12, 1, 22, 15, 0, 0, time.UTC))},
		{"datetime", dmy, []types.XValue{xs("01.15.2017")}, ERROR}, // month out of range
		{"datetime", dmy, []types.XValue{xs("no date")}, ERROR},    // invalid date
		{"datetime", dmy, []types.XValue{}, ERROR},

		{"datetime_from_parts", dmy, []types.XValue{xi(2018), xi(11), xi(3)}, xd(time.Date(2018, 11, 3, 0, 0, 0, 0, time.UTC))},
		{"datetime_from_parts", mdy, []types.XValue{xi(2018), xi(11), xi(3)}, xd(time.Date(2018, 11, 3, 0, 0, 0, 0, la))},
		{"datetime_from_parts", dmy, []types.XValue{xi(2018), xi(15), xi(3)}, ERROR}, // month out of range
		{"datetime_from_parts", dmy, []types.XValue{ERROR, xi(11), xi(3)}, ERROR},
		{"datetime_from_parts", dmy, []types.XValue{xi(2018), ERROR, xi(3)}, ERROR},
		{"datetime_from_parts", dmy, []types.XValue{xi(2018), xi(11), ERROR}, ERROR},
		{"datetime_from_parts", dmy, []types.XValue{}, ERROR},

		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("Y")}, xd(time.Date(2019, 12, 03, 22, 15, 0, 0, time.UTC))},
		{"datetime_add", mdy, []types.XValue{xs("12-03-2017 10:15pm"), xs("2"), xs("Y")}, xd(time.Date(2019, 12, 03, 22, 15, 0, 0, la))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("Y")}, xd(time.Date(2015, 12, 03, 22, 15, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("M")}, xd(time.Date(2018, 2, 03, 22, 15, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("M")}, xd(time.Date(2017, 10, 3, 22, 15, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("W")}, xd(time.Date(2017, 12, 17, 22, 15, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("W")}, xd(time.Date(2017, 11, 19, 22, 15, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017"), xs("2"), xs("D")}, xd(time.Date(2017, 12, 5, 0, 0, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017"), xs("-4"), xs("D")}, xd(time.Date(2017, 11, 29, 0, 0, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("h")}, xd(time.Date(2017, 12, 4, 0, 15, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("h")}, xd(time.Date(2017, 12, 3, 20, 15, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("105"), xs("m")}, xd(time.Date(2017, 12, 4, 0, 0, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("-20"), xs("m")}, xd(time.Date(2017, 12, 3, 21, 55, 0, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("2"), xs("s")}, xd(time.Date(2017, 12, 3, 22, 15, 2, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15pm"), xs("-2"), xs("s")}, xd(time.Date(2017, 12, 3, 22, 14, 58, 0, time.UTC))},
		{"datetime_add", dmy, []types.XValue{xs("xxx"), xs("2"), xs("D")}, ERROR},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15"), xs("xxx"), xs("D")}, ERROR},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017 10:15"), xs("2"), xs("xxx")}, ERROR},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017"), xs("2"), xs("Z")}, ERROR},
		{"datetime_add", dmy, []types.XValue{xs("03-12-2017"), xs("2"), ERROR}, ERROR},
		{"datetime_add", dmy, []types.XValue{xs("22-12-2017")}, ERROR},

		{"datetime_diff", dmy, []types.XValue{xs("03-12-2017"), xs("01-12-2017"), xs("D")}, xi(-2)},
		{"datetime_diff", mdy, []types.XValue{xs("12-03-2017"), xs("12-01-2017"), xs("D")}, xi(-2)},
		{"datetime_diff", dmy, []types.XValue{xs("03-12-2017 10:15"), xs("03-12-2017 18:15"), xs("D")}, xi(0)},
		{"datetime_diff", dmy, []types.XValue{xs("03-12-2017"), xs("01-12-2017"), xs("W")}, xi(0)},
		{"datetime_diff", dmy, []types.XValue{xs("22-12-2017"), xs("01-12-2017"), xs("W")}, xi(-3)},
		{"datetime_diff", dmy, []types.XValue{xs("03-12-2017"), xs("03-12-2017"), xs("M")}, xi(0)},
		{"datetime_diff", dmy, []types.XValue{xs("01-05-2018"), xs("03-12-2017"), xs("M")}, xi(-5)},
		{"datetime_diff", dmy, []types.XValue{xs("01-12-2018"), xs("03-12-2017"), xs("Y")}, xi(-1)},
		{"datetime_diff", dmy, []types.XValue{xs("01-01-2017"), xs("03-12-2017"), xs("Y")}, xi(0)},
		{"datetime_diff", dmy, []types.XValue{xs("04-12-2018 10:15"), xs("03-12-2018 14:00"), xs("h")}, xi(-20)},
		{"datetime_diff", dmy, []types.XValue{xs("04-12-2018 10:15"), xs("04-12-2018 14:00"), xs("h")}, xi(3)},
		{"datetime_diff", dmy, []types.XValue{xs("04-12-2018 10:15"), xs("04-12-2018 14:00"), xs("m")}, xi(225)},
		{"datetime_diff", dmy, []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:15:35"), xs("m")}, xi(0)},
		{"datetime_diff", dmy, []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:16:10"), xs("m")}, xi(0)},
		{"datetime_diff", dmy, []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:15:35"), xs("s")}, xi(20)},
		{"datetime_diff", dmy, []types.XValue{xs("05-12-2018 10:15:15"), xs("05-12-2018 10:16:10"), xs("s")}, xi(55)},
		{"datetime_diff", dmy, []types.XValue{xs("03-12-2017"), xs("01-12-2017"), xs("Z")}, ERROR},
		{"datetime_diff", dmy, []types.XValue{xs("xxx"), xs("01-12-2017"), xs("Y")}, ERROR},
		{"datetime_diff", dmy, []types.XValue{xs("01-12-2017"), xs("xxx"), xs("Y")}, ERROR},
		{"datetime_diff", dmy, []types.XValue{xs("01-12-2017"), xs("01-12-2017"), xs("xxx")}, ERROR},
		{"datetime_diff", dmy, []types.XValue{xs("01-12-2017"), xs("01-12-2017"), ERROR}, ERROR},
		{"datetime_diff", dmy, []types.XValue{}, ERROR},

		{"default", dmy, []types.XValue{xs("10"), xs("20")}, xs("10")},
		{"default", dmy, []types.XValue{nil, xs("20")}, xs("20")},
		{"default", dmy, []types.XValue{types.NewXErrorf("This is error"), xs("20")}, xs("20")},
		{"default", dmy, []types.XValue{}, ERROR},

		{"epoch", dmy, []types.XValue{xd(time.Date(2017, 6, 12, 16, 56, 59, 0, time.UTC))}, xn("1497286619")},
		{"epoch", dmy, []types.XValue{ERROR}, ERROR},
		{"epoch", dmy, []types.XValue{}, ERROR},

		{"field", dmy, []types.XValue{xs("hello,World"), xs("1"), xs(",")}, xs("World")},
		{"field", dmy, []types.XValue{xs("hello,world"), xn("2.1"), xs(",")}, xs("")},
		{"field", dmy, []types.XValue{xs("hello world there now"), xn("2"), xs(" ")}, xs("there")},
		{"field", dmy, []types.XValue{xs("hello   world    there     now"), xn("1"), xs(" ")}, xs("world")},
		{"field", dmy, []types.XValue{xs("hello   world    there     now"), xn("5"), xs(" ")}, xs("")},
		{"field", dmy, []types.XValue{xs("hello"), xi(0), xs(",")}, xs("hello")},
		{"field", dmy, []types.XValue{xs("hello,World"), xn("-2"), xs(",")}, ERROR},
		{"field", dmy, []types.XValue{xs(""), xs("notnum"), xs(",")}, ERROR},
		{"field", dmy, []types.XValue{xs("hello"), xi(0), nil}, xs("h")},
		{"field", dmy, []types.XValue{ERROR, xs("1"), xs(",")}, ERROR},
		{"field", dmy, []types.XValue{xs("hello"), ERROR, xs(",")}, ERROR},
		{"field", dmy, []types.XValue{xs("hello"), xs("1"), ERROR}, ERROR},
		{"field", dmy, []types.XValue{}, ERROR},

		{"format_date", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z")}, xs("23-06-1977")},
		{"format_date", mdy, []types.XValue{xs("1977-06-23T15:34:00.000000Z")}, xs("06-23-1977")},
		{"format_date", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DD")}, xs("1977-06-23")},
		{"format_date", dmy, []types.XValue{xs("1977-06-23"), xs("YYYY/MM/DD")}, xs("1977/06/23")},
		{"format_date", dmy, []types.XValue{xs("NOT DATE")}, ERROR},
		{"format_date", dmy, []types.XValue{ERROR}, ERROR},
		{"format_date", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), ERROR}, ERROR},
		{"format_date", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYYYYY")}, ERROR},
		{"format_date", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY"), ERROR}, ERROR},
		{"format_date", dmy, []types.XValue{}, ERROR},

		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z")}, xs("23-06-1977 15:34")},
		{"format_datetime", mdy, []types.XValue{xs("1977-06-23T15:34:00.000000Z")}, xs("06-23-1977 08:34")},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DDTtt:mm:ss.fffZZZ"), xs("America/Los_Angeles")}, xs("1977-06-23T08:34:00.000-07:00")},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.123000Z"), xs("YYYY-MM-DDTtt:mm:ss.fffZ"), xs("America/Los_Angeles")}, xs("1977-06-23T08:34:00.123-07:00")},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DDTtt:mm:ss.ffffffZ"), xs("America/Los_Angeles")}, xs("1977-06-23T08:34:00.000000-07:00")},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YY-MM-DD h:mm:ss AA"), xs("America/Los_Angeles")}, xs("77-06-23 8:34:00 AM")},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T08:34:00.000-07:00"), xs("YYYY-MM-DDTtt:mm:ss.fffZ"), xs("UTC")}, xs("1977-06-23T15:34:00.000Z")},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T08:34:00.000-07:00"), xs("h"), xs("UTC")}, xs("3")},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T08:34:00.000-07:00"), xs("hh"), xs("UTC")}, xs("03")},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T08:34:00.000-07:00"), xs("tt"), xs("UTC")}, xs("15")},
		{"format_datetime", dmy, []types.XValue{xs("NOT DATE")}, ERROR},
		{"format_datetime", dmy, []types.XValue{ERROR}, ERROR},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), ERROR}, ERROR},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYYYYY")}, ERROR},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY"), ERROR}, ERROR},
		{"format_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY"), xs("Cuenca")}, ERROR},
		{"format_datetime", dmy, []types.XValue{}, ERROR},

		{"format_location", dmy, []types.XValue{xs("Rwanda")}, xs("Rwanda")},
		{"format_location", dmy, []types.XValue{xs("Rwanda > Kigali")}, xs("Kigali")},
		{"format_location", dmy, []types.XValue{ERROR}, ERROR},
		{"format_location", dmy, []types.XValue{}, ERROR},

		{"format_number", dmy, []types.XValue{xn("31337")}, xs("31,337.00")},
		{"format_number", dmy, []types.XValue{xn("31337"), xi(0), types.XBooleanFalse}, xs("31337")},
		{"format_number", dmy, []types.XValue{xn("31337"), xs("xxx")}, ERROR},
		{"format_number", dmy, []types.XValue{xn("31337"), xi(12345)}, ERROR},
		{"format_number", dmy, []types.XValue{xn("31337"), xi(2), ERROR}, ERROR},
		{"format_number", dmy, []types.XValue{ERROR}, ERROR},
		{"format_number", dmy, []types.XValue{}, ERROR},

		{"format_urn", dmy, []types.XValue{xs("tel:+250781234567")}, xs("0781 234 567")},
		{"format_urn", dmy, []types.XValue{xs("twitter:134252511151#billy_bob")}, xs("billy_bob")},
		{"format_urn", dmy, []types.XValue{xs("NOT URN")}, ERROR},
		{"format_urn", dmy, []types.XValue{ERROR}, ERROR},
		{"format_urn", dmy, []types.XValue{}, ERROR},

		{"from_epoch", dmy, []types.XValue{xn("1497286619.000000000")}, xd(time.Date(2017, 6, 12, 16, 56, 59, 0, time.UTC))},
		{"from_epoch", dmy, []types.XValue{ERROR}, ERROR},
		{"from_epoch", dmy, []types.XValue{}, ERROR},

		{"if", dmy, []types.XValue{types.XBooleanTrue, xs("10"), xs("20")}, xs("10")},
		{"if", dmy, []types.XValue{types.XBooleanFalse, xs("10"), xs("20")}, xs("20")},
		{"if", dmy, []types.XValue{types.XBooleanTrue, errorArg, xs("20")}, errorArg},
		{"if", dmy, []types.XValue{}, ERROR},
		{"if", dmy, []types.XValue{errorArg, xs("10"), xs("20")}, errorArg},

		{"join", dmy, []types.XValue{types.NewXArray(xs("1"), xs("2"), xs("3")), xs(",")}, xs("1,2,3")},
		{"join", dmy, []types.XValue{types.NewXArray(), xs(",")}, xs("")},
		{"join", dmy, []types.XValue{types.NewXArray(xs("1")), xs(",")}, xs("1")},
		{"join", dmy, []types.XValue{types.NewXArray(xs("1"), xs("2")), ERROR}, ERROR},
		{"join", dmy, []types.XValue{types.NewXArray(xs("1"), ERROR), xs(",")}, ERROR},
		{"join", dmy, []types.XValue{xs("1,2,3"), nil}, ERROR},
		{"join", dmy, []types.XValue{types.NewXArray(xs("1,2,3")), nil}, xs("1,2,3")},
		{"join", dmy, []types.XValue{types.NewXArray(xs("1"))}, ERROR},

		{"json", dmy, []types.XValue{xs("hello")}, xs(`"hello"`)},
		{"json", dmy, []types.XValue{ERROR}, ERROR},

		{"left", dmy, []types.XValue{xs("hello"), xs("2")}, xs("he")},
		{"left", dmy, []types.XValue{xs("  HELLO"), xs("2")}, xs("  ")},
		{"left", dmy, []types.XValue{xs("hi"), xi(4)}, xs("hi")},
		{"left", dmy, []types.XValue{xs("hi"), xs("0")}, xs("")},
		{"left", dmy, []types.XValue{xs("😁hi"), xs("2")}, xs("😁h")},
		{"left", dmy, []types.XValue{xs("hello"), nil}, ERROR},
		{"left", dmy, []types.XValue{xs("hello"), xi(-1)}, ERROR},
		{"left", dmy, []types.XValue{ERROR, xi(3)}, ERROR},
		{"left", dmy, []types.XValue{xs("hello"), ERROR}, ERROR},
		{"left", dmy, []types.XValue{}, ERROR},

		{"legacy_add", dmy, []types.XValue{xs("01-12-2017"), xi(2)}, xd(time.Date(2017, 12, 3, 0, 0, 0, 0, time.UTC))},
		{"legacy_add", dmy, []types.XValue{xs("2"), xs("01-12-2017 10:15:33pm")}, xd(time.Date(2017, 12, 3, 22, 15, 33, 0, time.UTC))},
		{"legacy_add", dmy, []types.XValue{xs("2"), xs("3.5")}, xn("5.5")},
		{"legacy_add", dmy, []types.XValue{xs("01-12-2017 10:15:33pm"), xs("01-12-2017")}, ERROR},
		{"legacy_add", dmy, []types.XValue{types.NewXNumberFromInt64(int64(math.MaxInt32 + 1)), xs("01-12-2017 10:15:33pm")}, ERROR},
		{"legacy_add", dmy, []types.XValue{xs("01-12-2017 10:15:33pm"), types.NewXNumberFromInt64(int64(math.MaxInt32 + 1))}, ERROR},
		{"legacy_add", dmy, []types.XValue{xs("xxx"), xs("10")}, ERROR},
		{"legacy_add", dmy, []types.XValue{xs("10"), xs("xxx")}, ERROR},
		{"legacy_add", dmy, []types.XValue{}, ERROR},

		{"length", dmy, []types.XValue{xs("hello")}, xi(5)},
		{"length", dmy, []types.XValue{xs("")}, xi(0)},
		{"length", dmy, []types.XValue{xs("😁😁")}, xi(2)},
		{"length", dmy, []types.XValue{types.NewXArray(xs("hello"))}, xi(1)},
		{"length", dmy, []types.XValue{types.NewXArray()}, xi(0)},
		{"length", dmy, []types.XValue{xi(1234)}, ERROR},
		{"length", dmy, []types.XValue{ERROR}, ERROR},
		{"length", dmy, []types.XValue{}, ERROR},

		{"lower", dmy, []types.XValue{xs("HEllo")}, xs("hello")},
		{"lower", dmy, []types.XValue{xs("  HELLO  WORLD")}, xs("  hello  world")},
		{"lower", dmy, []types.XValue{xs("")}, xs("")},
		{"lower", dmy, []types.XValue{xs("😁")}, xs("😁")},
		{"lower", dmy, []types.XValue{}, ERROR},

		{"max", dmy, []types.XValue{xs("10.5"), xs("11")}, xi(11)},
		{"max", dmy, []types.XValue{xs("10.2"), xs("9")}, xn("10.2")},
		{"max", dmy, []types.XValue{xs("not_num"), xs("9")}, ERROR},
		{"max", dmy, []types.XValue{xs("9"), xs("not_num")}, ERROR},
		{"max", dmy, []types.XValue{}, ERROR},

		{"min", dmy, []types.XValue{xs("10.5"), xs("11")}, xn("10.5")},
		{"min", dmy, []types.XValue{xs("10.2"), xs("9")}, xi(9)},
		{"min", dmy, []types.XValue{xs("not_num"), xs("9")}, ERROR},
		{"min", dmy, []types.XValue{xs("9"), xs("not_num")}, ERROR},
		{"min", dmy, []types.XValue{}, ERROR},

		{"mean", dmy, []types.XValue{xs("10"), xs("11")}, xn("10.5")},
		{"mean", dmy, []types.XValue{xs("10.2")}, xn("10.2")},
		{"mean", dmy, []types.XValue{xs("not_num")}, ERROR},
		{"mean", dmy, []types.XValue{xs("9"), xs("not_num")}, ERROR},
		{"mean", dmy, []types.XValue{}, ERROR},

		{"mod", dmy, []types.XValue{xs("10"), xs("3")}, xi(1)},
		{"mod", dmy, []types.XValue{xs("10"), xs("5")}, xi(0)},
		{"mod", dmy, []types.XValue{xs("not_num"), xs("3")}, ERROR},
		{"mod", dmy, []types.XValue{xs("9"), xs("not_num")}, ERROR},
		{"mod", dmy, []types.XValue{}, ERROR},

		{"now", dmy, []types.XValue{}, xd(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC))},
		{"now", dmy, []types.XValue{ERROR}, ERROR},

		{"number", dmy, []types.XValue{xn("10")}, xn("10")},
		{"number", dmy, []types.XValue{xs("123.45000")}, xn("123.45")},
		{"number", dmy, []types.XValue{xs("what?")}, ERROR},

		{"or", dmy, []types.XValue{types.XBooleanTrue}, types.XBooleanTrue},
		{"or", dmy, []types.XValue{types.XBooleanFalse}, types.XBooleanFalse},
		{"or", dmy, []types.XValue{types.XBooleanTrue, types.XBooleanFalse}, types.XBooleanTrue},
		{"or", dmy, []types.XValue{ERROR}, ERROR},
		{"or", dmy, []types.XValue{}, ERROR},

		{"parse_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.000000Z"), xs("YYYY-MM-DDTtt:mm:ss.ffffffZ"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 8, 34, 0, 0, la))},
		{"parse_datetime", dmy, []types.XValue{xs("1977-06-23T15:34:00.1234Z"), xs("YYYY-MM-DDTtt:mm:ssZ"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 8, 34, 0, 123400000, la))},
		{"parse_datetime", dmy, []types.XValue{xs("1977-06-23 15:34"), xs("YYYY-MM-DD tt:mm"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 15, 34, 0, 0, la))},
		{"parse_datetime", dmy, []types.XValue{xs("1977-06-23 03:34 pm"), xs("YYYY-MM-DD tt:mm aa"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 15, 34, 0, 0, la))},
		{"parse_datetime", dmy, []types.XValue{xs("1977-06-23 03:34 PM"), xs("YYYY-MM-DD tt:mm AA"), xs("America/Los_Angeles")}, xd(time.Date(1977, 06, 23, 15, 34, 0, 0, la))},
		{"parse_datetime", dmy, []types.XValue{xs("1977-06-23 15:34"), xs("ttttttttt")}, ERROR},                // invalid format
		{"parse_datetime", dmy, []types.XValue{xs("1977-06-23 15:34"), xs("YYYY-MM-DD"), xs("Cuenca")}, ERROR}, // invalid timezone
		{"parse_datetime", dmy, []types.XValue{xs("abcd"), xs("YYYY-MM-DD")}, ERROR},                           // unparseable date
		{"parse_datetime", dmy, []types.XValue{ERROR, xs("YYYY-MM-DD")}, ERROR},
		{"parse_datetime", dmy, []types.XValue{xs("1977-06-23 15:34"), ERROR}, ERROR},
		{"parse_datetime", dmy, []types.XValue{xs("1977-06-23 15:34"), xs("YYYY-MM-DD"), ERROR}, ERROR},
		{"parse_datetime", dmy, []types.XValue{}, ERROR},

		{"parse_json", dmy, []types.XValue{xs(`"hello"`)}, xs(`hello`)},
		{"parse_json", dmy, []types.XValue{ERROR}, ERROR},

		{"percent", dmy, []types.XValue{xs(".54")}, xs("54%")},
		{"percent", dmy, []types.XValue{xs("1.246")}, xs("125%")},
		{"percent", dmy, []types.XValue{xs("")}, ERROR},
		{"percent", dmy, []types.XValue{}, ERROR},

		{"rand", dmy, []types.XValue{}, xn("0.3849275689214193274523267973563633859157562255859375")},
		{"rand", dmy, []types.XValue{}, xn("0.607552015674623913099594574305228888988494873046875")},

		{"rand_between", dmy, []types.XValue{xn("1"), xn("10")}, xn("5")},
		{"rand_between", dmy, []types.XValue{xn("1"), xn("10")}, xn("10")},

		{"read_chars", dmy, []types.XValue{xs("123456")}, xs("1 2 3 , 4 5 6")},
		{"read_chars", dmy, []types.XValue{xs("abcd")}, xs("a b c d")},
		{"read_chars", dmy, []types.XValue{xs("12345678")}, xs("1 2 3 4 , 5 6 7 8")},
		{"read_chars", dmy, []types.XValue{xs("12")}, xs("1 , 2")},
		{"read_chars", dmy, []types.XValue{}, ERROR},

		{"regex_match", dmy, []types.XValue{xs("zAbc"), xs(`a\w`)}, xs(`Ab`)},
		{"regex_match", dmy, []types.XValue{xs("<html>"), xs(`<(\w+)>`), xn("1")}, xs(`html`)},
		{"regex_match", dmy, []types.XValue{xs("<html>"), xs(`<(\w+)>`), xn("2")}, ERROR},

		{"remove_first_word", dmy, []types.XValue{xs("hello World")}, xs("World")},
		{"remove_first_word", dmy, []types.XValue{xs("hello")}, xs("")},
		{"remove_first_word", dmy, []types.XValue{xs("😁hello")}, xs("hello")},
		{"remove_first_word", dmy, []types.XValue{xs("")}, xs("")},
		{"remove_first_word", dmy, []types.XValue{}, ERROR},

		{"repeat", dmy, []types.XValue{xs("hi"), xs("2")}, xs("hihi")},
		{"repeat", dmy, []types.XValue{xs("  "), xs("2")}, xs("    ")},
		{"repeat", dmy, []types.XValue{xs(""), xi(4)}, xs("")},
		{"repeat", dmy, []types.XValue{xs("😁"), xs("2")}, xs("😁😁")},
		{"repeat", dmy, []types.XValue{xs("hi"), xs("0")}, xs("")},
		{"repeat", dmy, []types.XValue{xs("hi"), xs("-1")}, ERROR},
		{"repeat", dmy, []types.XValue{xs("hello"), nil}, ERROR},
		{"repeat", dmy, []types.XValue{}, ERROR},

		{"replace", dmy, []types.XValue{xs("hi ho"), xs("hi"), xs("bye")}, xs("bye ho")},
		{"replace", dmy, []types.XValue{xs("foo bar "), xs(" "), xs(".")}, xs("foo.bar.")},
		{"replace", dmy, []types.XValue{xs("foo 😁 bar "), xs("😁"), xs("😂")}, xs("foo 😂 bar ")},
		{"replace", dmy, []types.XValue{xs("foo bar"), xs("zap"), xs("zog")}, xs("foo bar")},
		{"replace", dmy, []types.XValue{nil, xs("foo bar"), xs("foo")}, xs("")},
		{"replace", dmy, []types.XValue{xs("foo bar"), nil, xs("foo")}, xs("fooffooofooofoo foobfooafoorfoo")},
		{"replace", dmy, []types.XValue{xs("foo bar"), xs("foo"), nil}, xs(" bar")},
		{"replace", dmy, []types.XValue{ERROR, xs("hi"), xs("bye")}, ERROR},
		{"replace", dmy, []types.XValue{xs("hi ho"), ERROR, xs("bye")}, ERROR},
		{"replace", dmy, []types.XValue{xs("hi ho"), xs("bye"), ERROR}, ERROR},
		{"replace", dmy, []types.XValue{}, ERROR},

		{"right", dmy, []types.XValue{xs("hello"), xs("2")}, xs("lo")},
		{"right", dmy, []types.XValue{xs("  HELLO "), xs("2")}, xs("O ")},
		{"right", dmy, []types.XValue{xs("hi"), xi(4)}, xs("hi")},
		{"right", dmy, []types.XValue{xs("hi"), xs("0")}, xs("")},
		{"right", dmy, []types.XValue{xs("ho😁hi"), xs("4")}, xs("o😁hi")},
		{"right", dmy, []types.XValue{nil, xs("2")}, xs("")},
		{"right", dmy, []types.XValue{xs("hello"), nil}, ERROR},
		{"right", dmy, []types.XValue{xs("hello"), xi(-1)}, ERROR},
		{"right", dmy, []types.XValue{ERROR, xi(3)}, ERROR},
		{"right", dmy, []types.XValue{xs("hello"), ERROR}, ERROR},
		{"right", dmy, []types.XValue{}, ERROR},

		{"round", dmy, []types.XValue{xs("10.5"), xs("0")}, xi(11)},
		{"round", dmy, []types.XValue{xs("10.5"), xs("1")}, xn("10.5")},
		{"round", dmy, []types.XValue{xs("10.51"), xs("1")}, xn("10.5")},
		{"round", dmy, []types.XValue{xs("10.56"), xs("1")}, xn("10.6")},
		{"round", dmy, []types.XValue{xs("12.56"), xs("-1")}, xi(10)},
		{"round", dmy, []types.XValue{xs("10.5")}, xn("11")},
		{"round", dmy, []types.XValue{xs("not_num"), xs("1")}, ERROR},
		{"round", dmy, []types.XValue{xs("10.5"), xs("not_num")}, ERROR},
		{"round", dmy, []types.XValue{xs("10.5"), xs("1"), xs("30")}, ERROR},

		{"round_down", dmy, []types.XValue{xs("10")}, xi(10)},
		{"round_down", dmy, []types.XValue{xs("10.5")}, xi(10)},
		{"round_down", dmy, []types.XValue{xs("10.7")}, xi(10)},
		{"round_down", dmy, []types.XValue{xs("not_num")}, ERROR},
		{"round_down", dmy, []types.XValue{}, ERROR},

		{"round_up", dmy, []types.XValue{xs("10")}, xi(10)},
		{"round_up", dmy, []types.XValue{xs("10.5")}, xi(11)},
		{"round_up", dmy, []types.XValue{xs("10.2")}, xi(11)},
		{"round_up", dmy, []types.XValue{xs("not_num")}, ERROR},
		{"round_up", dmy, []types.XValue{}, ERROR},

		{"split", dmy, []types.XValue{xs("1,2,3"), xs(",")}, types.NewXArray(xs("1"), xs("2"), xs("3"))},
		{"split", dmy, []types.XValue{xs("1,2,3"), xs(".")}, types.NewXArray(xs("1,2,3"))},
		{"split", dmy, []types.XValue{xs("1,2,3"), nil}, types.NewXArray(xs("1,2,3"))},
		{"split", dmy, []types.XValue{ERROR, xs(",")}, ERROR},
		{"split", dmy, []types.XValue{xs("1,2,3"), ERROR}, ERROR},
		{"split", dmy, []types.XValue{}, ERROR},

		{"text", dmy, []types.XValue{xs("abc")}, xs("abc")},
		{"text", dmy, []types.XValue{xi(123)}, xs("123")},
		{"text", dmy, []types.XValue{ERROR}, ERROR},
		{"text", dmy, []types.XValue{}, ERROR},

		{"text_compare", dmy, []types.XValue{xs("abc"), xs("abc")}, xi(0)},
		{"text_compare", dmy, []types.XValue{xs("abc"), xs("def")}, xi(-1)},
		{"text_compare", dmy, []types.XValue{xs("def"), xs("abc")}, xi(1)},
		{"text_compare", dmy, []types.XValue{xs("abc"), types.NewXErrorf("error")}, ERROR},
		{"text_compare", dmy, []types.XValue{}, ERROR},

		{"title", dmy, []types.XValue{xs("hello")}, xs("Hello")},
		{"title", dmy, []types.XValue{xs("")}, xs("")},
		{"title", dmy, []types.XValue{nil}, xs("")},
		{"title", dmy, []types.XValue{}, ERROR},

		{"today", dmy, []types.XValue{}, xd(time.Date(2018, 4, 11, 0, 0, 0, 0, time.UTC))},
		{"today", mdy, []types.XValue{}, xd(time.Date(2018, 4, 11, 0, 0, 0, 0, la))},
		{"today", dmy, []types.XValue{ERROR}, ERROR},

		{"tz", dmy, []types.XValue{xs("01-12-2017")}, xs("UTC")},
		{"tz", mdy, []types.XValue{xs("01-12-2017")}, xs("America/Los_Angeles")},
		{"tz", dmy, []types.XValue{xs("01-12-2017 10:15:33pm")}, xs("UTC")},
		{"tz", dmy, []types.XValue{xs("xxx")}, ERROR},
		{"tz", dmy, []types.XValue{}, ERROR},

		{"tz_offset", dmy, []types.XValue{xs("01-12-2017")}, xs("+0000")},
		{"tz_offset", mdy, []types.XValue{xs("01-12-2017")}, xs("-0800")},
		{"tz_offset", dmy, []types.XValue{xs("01-12-2017 10:15:33pm")}, xs("+0000")},
		{"tz_offset", dmy, []types.XValue{xs("xxx")}, ERROR},
		{"tz_offset", dmy, []types.XValue{}, ERROR},

		{"upper", dmy, []types.XValue{xs("HEllo")}, xs("HELLO")},
		{"upper", dmy, []types.XValue{xs("  HELLO  world")}, xs("  HELLO  WORLD")},
		{"upper", dmy, []types.XValue{xs("ß")}, xs("ß")},
		{"upper", dmy, []types.XValue{xs("")}, xs("")},
		{"upper", dmy, []types.XValue{}, ERROR},

		{"word", dmy, []types.XValue{xs("hello World"), xn("1.5")}, xs("World")},
		{"word", dmy, []types.XValue{xs(""), xi(0)}, ERROR},
		{"word", dmy, []types.XValue{xs("cat dog bee"), xi(-1)}, xs("bee")},
		{"word", dmy, []types.XValue{xs("😁 hello World"), xi(0)}, xs("😁")},
		{"word", dmy, []types.XValue{xs(" hello World"), xi(2)}, ERROR},
		{"word", dmy, []types.XValue{xs("hello World"), nil}, ERROR},
		{"word", dmy, []types.XValue{}, ERROR},

		{"word_slice", dmy, []types.XValue{xs("hello-world from mars"), xi(0), xi(2)}, xs("hello world")},
		{"word_slice", dmy, []types.XValue{xs("hello-world from mars"), xi(2)}, xs("from mars")},
		{"word_slice", dmy, []types.XValue{xs("hello-world from mars"), xi(10)}, xs("")},
		{"word_slice", dmy, []types.XValue{xs("hello-world from mars"), xi(3), xi(10)}, xs("mars")},
		{"word_slice", dmy, []types.XValue{xs("hello-world from mars"), xi(-1), xi(2)}, ERROR},
		{"word_slice", dmy, []types.XValue{xs("hello-world from mars"), xi(3), xi(1)}, ERROR},
		{"word_slice", dmy, []types.XValue{xs("hello-world from mars"), xs("x"), xi(3)}, ERROR},
		{"word_slice", dmy, []types.XValue{xs("hello-world from mars"), xi(3), xs("x")}, ERROR},
		{"word_slice", dmy, []types.XValue{xs("hello-world from mars"), ERROR, xi(2)}, ERROR},
		{"word_slice", dmy, []types.XValue{ERROR, xi(0), xi(2)}, ERROR},
		{"word_slice", dmy, []types.XValue{ERROR}, ERROR},

		{"word_count", dmy, []types.XValue{xs("hello World")}, xi(2)},
		{"word_count", dmy, []types.XValue{xs("hello")}, xi(1)},
		{"word_count", dmy, []types.XValue{xs("")}, xi(0)},
		{"word_count", dmy, []types.XValue{xs("😁😁")}, xi(2)},
		{"word_count", dmy, []types.XValue{}, ERROR},

		{"weekday", dmy, []types.XValue{xs("01-12-2017")}, xi(5)},
		{"weekday", mdy, []types.XValue{xs("12-01-2017")}, xi(5)},
		{"weekday", dmy, []types.XValue{xs("01-12-2017 10:15pm")}, xi(5)},
		{"weekday", dmy, []types.XValue{xs("xxx")}, ERROR},
		{"weekday", dmy, []types.XValue{}, ERROR},

		{"url_encode", dmy, []types.XValue{xs(`hi-% ?/`)}, xs(`hi-%25%20%3F%2F`)},
		{"url_encode", dmy, []types.XValue{ERROR}, ERROR},
		{"url_encode", dmy, []types.XValue{}, ERROR},
	}

	defer utils.SetRand(utils.DefaultRand)
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	utils.SetRand(utils.NewSeededRand(123456))
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)))

	for _, test := range funcTests {
		xFunc, exists := functions.XFUNCTIONS[test.name]
		require.True(t, exists, "no such registered function: %s", test.name)

		result := xFunc(test.env, test.args...)

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic running function %s(%#v): %#v", test.name, test.args, r)
			}
		}()

		// don't check error equality - just check that we got an error if we expected one
		if test.expected == ERROR {
			assert.True(t, types.IsXError(result), "expecting error, got %T{%s} for function %s(%T{%s})", result, result, test.name, test.args, test.args)
		} else {
			if !types.Equals(test.env, result, test.expected) {
				assert.Fail(t, "", "unexpected value, expected %T{%s}, got %T{%s} for function %s(%T{%s})", test.expected, test.expected, result, result, test.name, test.args, test.args)
			}
		}
	}
}

func TestFormatDecimal(t *testing.T) {
	fmtTests := []struct {
		input       decimal.Decimal
		format      *utils.NumberFormat
		places      int
		groupDigits bool
		expected    string
	}{
		{decimal.RequireFromString("1234"), utils.DefaultNumberFormat, 2, true, "1,234.00"},
		{decimal.RequireFromString("1234"), utils.DefaultNumberFormat, 0, false, "1234"},
		{decimal.RequireFromString("1234.567"), utils.DefaultNumberFormat, 2, true, "1,234.57"},
		{decimal.RequireFromString("1234.567"), utils.DefaultNumberFormat, 2, false, "1234.57"},
		{decimal.RequireFromString("1234.567"), &utils.NumberFormat{DecimalSymbol: ",", DigitGroupingSymbol: "."}, 2, true, "1.234,57"},
	}

	for _, test := range fmtTests {
		val := functions.FormatDecimal(test.input, test.format, test.places, test.groupDigits)

		assert.Equal(t, test.expected, val, "format decimal failed for input '%s'", test.input)
	}
}
