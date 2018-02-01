package excellent

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
)

var errorArg = fmt.Errorf("I am error")
var la, _ = time.LoadLocation("America/Los_Angeles")

var funcTests = []struct {
	name     string
	args     []interface{}
	expected interface{}
	error    bool
}{
	{"and", []interface{}{true}, true, false},
	{"and", []interface{}{true, false}, false, false},
	{"and", []interface{}{}, false, true},
	{"and", []interface{}{struct{}{}, true}, false, true},
	{"and", []interface{}{false, struct{}{}}, false, true},

	{"char", []interface{}{33}, "!", false},
	{"char", []interface{}{128513}, "😁", false},
	{"char", []interface{}{"not decimal"}, nil, true},
	{"char", []interface{}{}, false, true},

	{"or", []interface{}{true}, true, false},
	{"or", []interface{}{true, false}, true, false},
	{"or", []interface{}{}, false, true},
	{"or", []interface{}{struct{}{}, true}, false, true},
	{"or", []interface{}{true, struct{}{}}, false, true},

	{"if", []interface{}{true, "10", "20"}, "10", false},
	{"if", []interface{}{false, "10", "20"}, "20", false},
	{"if", []interface{}{true, errorArg, "20"}, errorArg, true},
	{"if", []interface{}{}, false, true},
	{"if", []interface{}{struct{}{}, "10", "20"}, false, true},

	{"round", []interface{}{"10.5", "0"}, decimal.NewFromFloat(11), false},
	{"round", []interface{}{"10.5", "1"}, decimal.NewFromFloat(10.5), false},
	{"round", []interface{}{"not_num", "1"}, nil, true},
	{"round", []interface{}{"10.5", "not_num"}, nil, true},
	{"round", []interface{}{"10.5", "-1"}, nil, true},
	{"round", []interface{}{"10.5"}, nil, true},

	{"round_up", []interface{}{"10.5"}, decimal.NewFromFloat(11), false},
	{"round_up", []interface{}{"10.2"}, decimal.NewFromFloat(11), false},
	{"round_up", []interface{}{"not_num"}, nil, true},
	{"round_up", []interface{}{}, nil, true},

	{"round_down", []interface{}{"10.5"}, decimal.NewFromFloat(10), false},
	{"round_down", []interface{}{"10.7"}, decimal.NewFromFloat(10), false},
	{"round_down", []interface{}{"not_num"}, nil, true},
	{"round_down", []interface{}{}, nil, true},

	{"int", []interface{}{"10.5"}, decimal.NewFromFloat(10), false},
	{"int", []interface{}{"10.7"}, decimal.NewFromFloat(10), false},
	{"int", []interface{}{"not_num"}, nil, true},
	{"int", []interface{}{}, nil, true},

	{"max", []interface{}{"10.5", "11"}, decimal.NewFromFloat(11), false},
	{"max", []interface{}{"10.2", "9"}, decimal.NewFromFloat(10.2), false},
	{"max", []interface{}{"not_num", "9"}, nil, true},
	{"max", []interface{}{"9", "not_num"}, nil, true},
	{"max", []interface{}{}, nil, true},
	{"min", []interface{}{"10.5", "11"}, decimal.NewFromFloat(10.5), false},
	{"min", []interface{}{"10.2", "9"}, decimal.NewFromFloat(9), false},
	{"min", []interface{}{"not_num", "9"}, nil, true},
	{"min", []interface{}{"9", "not_num"}, nil, true},
	{"min", []interface{}{}, nil, true},

	{"mean", []interface{}{"10", "11"}, decimal.NewFromFloat(10.5), false},
	{"mean", []interface{}{"10.2"}, decimal.NewFromFloat(10.2), false},
	{"mean", []interface{}{"not_num"}, nil, true},
	{"mean", []interface{}{"9", "not_num"}, nil, true},
	{"mean", []interface{}{}, nil, true},

	{"mod", []interface{}{"10", "3"}, decimal.NewFromFloat(1), false},
	{"mod", []interface{}{"10", "5"}, decimal.NewFromFloat(0), false},
	{"mod", []interface{}{"not_num", "3"}, nil, true},
	{"mod", []interface{}{"9", "not_num"}, nil, true},
	{"mod", []interface{}{}, nil, true},

	{"read_code", []interface{}{"123456"}, "1 2 3 , 4 5 6", false},
	{"read_code", []interface{}{"abcd"}, "a b c d", false},
	{"read_code", []interface{}{"12345678"}, "1 2 3 4 , 5 6 7 8", false},
	{"read_code", []interface{}{"12"}, "1 , 2", false},
	{"read_code", []interface{}{struct{}{}}, nil, true},
	{"read_code", []interface{}{}, nil, true},

	{"split", []interface{}{"1,2,3", ","}, []string{"1", "2", "3"}, false},
	{"split", []interface{}{"1,2,3", "."}, []string{"1,2,3"}, false},
	{"split", []interface{}{struct{}{}, "."}, nil, true},
	{"split", []interface{}{"1,2,3", struct{}{}}, nil, true},
	{"split", []interface{}{}, nil, true},

	{"join", []interface{}{[]interface{}{"1", "2", "3"}, ","}, "1,2,3", false},
	{"join", []interface{}{[]interface{}{}, ","}, "", false},
	{"join", []interface{}{[]interface{}{"1"}, ","}, "1", false},
	{"join", []interface{}{"1,2,3", struct{}{}}, nil, true},
	{"join", []interface{}{[]interface{}{"1,2,3"}, struct{}{}}, nil, true},
	{"join", []interface{}{[]interface{}{"1"}}, nil, true},

	{"title", []interface{}{"hello"}, "Hello", false},
	{"title", []interface{}{""}, "", false},
	{"title", []interface{}{struct{}{}}, nil, true},
	{"title", []interface{}{}, nil, true},

	{"word", []interface{}{"hello World", 2}, "World", false},
	{"word", []interface{}{"", 1}, "", true},
	{"word", []interface{}{"😁 hello World", 1}, "😁", false},
	{"word", []interface{}{" hello World", 3}, nil, true},
	{"word", []interface{}{"hello World", struct{}{}}, nil, true},
	{"word", []interface{}{struct{}{}, 3}, nil, true},
	{"word", []interface{}{struct{}{}}, nil, true},
	{"word", []interface{}{}, nil, true},

	{"remove_first_word", []interface{}{"hello World"}, "World", false},
	{"remove_first_word", []interface{}{"hello"}, "", false},
	{"remove_first_word", []interface{}{"😁hello"}, "hello", false},
	{"remove_first_word", []interface{}{""}, "", false},
	{"remove_first_word", []interface{}{struct{}{}}, nil, true},
	{"remove_first_word", []interface{}{}, nil, true},

	{"word_count", []interface{}{"hello World"}, decimal.NewFromFloat(2), false},
	{"word_count", []interface{}{"hello"}, decimal.NewFromFloat(1), false},
	{"word_count", []interface{}{""}, decimal.NewFromFloat(0), false},
	{"word_count", []interface{}{"😁😁"}, decimal.NewFromFloat(2), false},
	{"word_count", []interface{}{struct{}{}}, nil, true},
	{"word_count", []interface{}{}, nil, true},

	{"field", []interface{}{"hello,World", "1", ","}, "World", false},
	{"field", []interface{}{"hello,world", "2", ","}, "", false},
	{"field", []interface{}{"hello", "0", ","}, "hello", false},
	{"field", []interface{}{"hello,World", "-2", ","}, nil, true},
	{"field", []interface{}{"", "notnum", ","}, nil, true},
	{"field", []interface{}{struct{}{}, "0", ","}, nil, true},
	{"field", []interface{}{"hello", "0", struct{}{}}, nil, true},
	{"field", []interface{}{struct{}{}}, nil, true},

	{"clean", []interface{}{"hello"}, "hello", false},
	{"clean", []interface{}{"  hello  world\n\t"}, "hello  world", false},
	{"clean", []interface{}{""}, "", false},
	{"clean", []interface{}{struct{}{}}, nil, true},
	{"clean", []interface{}{}, nil, true},

	{"lower", []interface{}{"HEllo"}, "hello", false},
	{"lower", []interface{}{"  HELLO  WORLD"}, "  hello  world", false},
	{"lower", []interface{}{""}, "", false},
	{"lower", []interface{}{"😁"}, "😁", false},
	{"lower", []interface{}{struct{}{}}, nil, true},
	{"lower", []interface{}{}, nil, true},

	{"left", []interface{}{"hello", "2"}, "he", false},
	{"left", []interface{}{"  HELLO", "2"}, "  ", false},
	{"left", []interface{}{"hi", "4"}, "hi", false},
	{"left", []interface{}{"hi", "0"}, "", false},
	{"left", []interface{}{"😁hi", "2"}, "😁h", false},
	{"left", []interface{}{struct{}{}, "2"}, nil, true},
	{"left", []interface{}{"hello", struct{}{}}, nil, true},
	{"left", []interface{}{}, nil, true},

	{"right", []interface{}{"hello", "2"}, "lo", false},
	{"right", []interface{}{"  HELLO ", "2"}, "O ", false},
	{"right", []interface{}{"hi", "4"}, "hi", false},
	{"right", []interface{}{"hi", "0"}, "", false},
	{"right", []interface{}{"ho😁hi", "4"}, "o😁hi", false},
	{"right", []interface{}{struct{}{}, "2"}, nil, true},
	{"right", []interface{}{"hello", struct{}{}}, nil, true},
	{"right", []interface{}{}, nil, true},

	{"string_length", []interface{}{"hello"}, decimal.NewFromFloat(5), false},
	{"string_length", []interface{}{""}, decimal.NewFromFloat(0), false},
	{"string_length", []interface{}{"😁😁"}, decimal.NewFromFloat(2), false},
	{"string_length", []interface{}{struct{}{}}, nil, true},
	{"string_length", []interface{}{}, nil, true},
	// string_length doesn't work on arrays
	{"string_length", []interface{}{[]interface{}{"hello", "world"}}, decimal.NewFromFloat(2), true},
	{"string_length", []interface{}{[]interface{}{}}, decimal.NewFromFloat(0), true},

	{"array_length", []interface{}{[]string{"hello"}}, decimal.NewFromFloat(1), false},
	{"array_length", []interface{}{[]string{}}, decimal.NewFromFloat(0), false},
	{"array_length", []interface{}{struct{}{}}, nil, true},
	{"array_length", []interface{}{}, nil, true},

	{"default", []interface{}{"10", "20"}, "10", false},
	{"default", []interface{}{nil, "20"}, "20", false},
	{"default", []interface{}{fmt.Errorf("This is error"), "20"}, "20", false},
	{"default", []interface{}{struct{}{}}, nil, true},
	{"default", []interface{}{}, nil, true},

	{"repeat", []interface{}{"hi", "2"}, "hihi", false},
	{"repeat", []interface{}{"  ", "2"}, "    ", false},
	{"repeat", []interface{}{"", "4"}, "", false},
	{"repeat", []interface{}{"😁", "2"}, "😁😁", false},
	{"repeat", []interface{}{"hi", "0"}, "", false},
	{"repeat", []interface{}{"hi", "-1"}, "", true},
	{"repeat", []interface{}{struct{}{}, "2"}, nil, true},
	{"repeat", []interface{}{"hello", struct{}{}}, nil, true},
	{"repeat", []interface{}{}, nil, true},

	{"replace", []interface{}{"hi ho", "hi", "bye"}, "bye ho", false},
	{"replace", []interface{}{"foo bar ", " ", "."}, "foo.bar.", false},
	{"replace", []interface{}{"foo 😁 bar ", "😁", "😂"}, "foo 😂 bar ", false},
	{"replace", []interface{}{"foo bar", "zap", "zog"}, "foo bar", false},
	{"replace", []interface{}{struct{}{}, "foo bar", "foo"}, "", true},
	{"replace", []interface{}{"foo bar", struct{}{}, "foo"}, "", true},
	{"replace", []interface{}{"foo bar", "foo", struct{}{}}, "", true},
	{"replace", []interface{}{}, nil, true},

	{"upper", []interface{}{"HEllo"}, "HELLO", false},
	{"upper", []interface{}{"  HELLO  world"}, "  HELLO  WORLD", false},
	{"upper", []interface{}{"ß"}, "ß", false},
	{"upper", []interface{}{""}, "", false},
	{"upper", []interface{}{struct{}{}}, nil, true},
	{"upper", []interface{}{}, nil, true},

	{"percent", []interface{}{".54"}, "54%", false},
	{"percent", []interface{}{"1.246"}, "125%", false},
	{"percent", []interface{}{""}, nil, true},
	{"percent", []interface{}{struct{}{}}, nil, true},
	{"percent", []interface{}{}, nil, true},

	{"date", []interface{}{"01-12-2017"}, time.Date(2017, 12, 1, 0, 0, 0, 0, time.UTC), false},
	{"date", []interface{}{"01-12-2017 10:15pm"}, time.Date(2017, 12, 1, 22, 15, 0, 0, time.UTC), false},
	{"date", []interface{}{"01.15.2017"}, nil, true}, // month out of range
	{"date", []interface{}{"no date"}, nil, true},    // invalid date
	{"date", []interface{}{struct{}{}}, nil, true},
	{"date", []interface{}{}, nil, true},

	{"format_date", []interface{}{"1977-06-23T15:34:00.000000Z", "yyyy-MM-ddTHH:mm:ss.fffzzz", "America/Los_Angeles"}, "1977-06-23T08:34:00.000-07:00", false},
	{"format_date", []interface{}{"1977-06-23T15:34:00.000000Z", "yyyy-MM-ddTHH:mm:ss.fffK", "America/Los_Angeles"}, "1977-06-23T08:34:00.000-07:00", false},
	{"format_date", []interface{}{"1977-06-23T08:34:00.000-07:00", "yyyy-MM-ddTHH:mm:ss.fffK", "UTC"}, "1977-06-23T15:34:00.000Z", false},

	{"parse_date", []interface{}{"1977-06-23T15:34:00.000000Z", "yyyy-MM-ddTHH:mm:ss.ffffffK", "America/Los_Angeles"}, time.Date(1977, 06, 23, 8, 34, 0, 0, la), false},
	{"parse_date", []interface{}{"1977-06-23 15:34", "yyyy-MM-dd HH:mm", "America/Los_Angeles"}, time.Date(1977, 06, 23, 15, 34, 0, 0, la), false},
	{"parse_date", []interface{}{"1977-06-23 03:34 pm", "yyyy-MM-dd HH:mm tt", "America/Los_Angeles"}, time.Date(1977, 06, 23, 15, 34, 0, 0, la), false},
	{"parse_date", []interface{}{"1977-06-23 03:34 PM", "yyyy-MM-dd HH:mm TT", "America/Los_Angeles"}, time.Date(1977, 06, 23, 15, 34, 0, 0, la), false},

	{"date_diff", []interface{}{"03-12-2017", "01-12-2017", "d"}, 2, false},
	{"date_diff", []interface{}{"03-12-2017 10:15", "03-12-2017 18:15", "d"}, 0, false},
	{"date_diff", []interface{}{"03-12-2017", "01-12-2017", "w"}, 0, false},
	{"date_diff", []interface{}{"22-12-2017", "01-12-2017", "w"}, 3, false},
	{"date_diff", []interface{}{"03-12-2017", "03-12-2017", "M"}, 0, false},
	{"date_diff", []interface{}{"01-05-2018", "03-12-2017", "M"}, 5, false},
	{"date_diff", []interface{}{"01-12-2018", "03-12-2017", "y"}, 1, false},
	{"date_diff", []interface{}{"01-01-2017", "03-12-2017", "y"}, 0, false},
	{"date_diff", []interface{}{"04-12-2018 10:15", "03-12-2018 14:00", "h"}, 20, false},
	{"date_diff", []interface{}{"04-12-2018 10:15", "04-12-2018 14:00", "h"}, -3, false},
	{"date_diff", []interface{}{"04-12-2018 10:15", "04-12-2018 14:00", "m"}, -225, false},
	{"date_diff", []interface{}{"05-12-2018 10:15:15", "05-12-2018 10:15:35", "m"}, 0, false},
	{"date_diff", []interface{}{"05-12-2018 10:15:15", "05-12-2018 10:16:10", "m"}, 0, false},
	{"date_diff", []interface{}{"05-12-2018 10:15:15", "05-12-2018 10:15:35", "s"}, -20, false},
	{"date_diff", []interface{}{"05-12-2018 10:15:15", "05-12-2018 10:16:10", "s"}, -55, false},
	{"date_diff", []interface{}{"03-12-2017", "01-12-2017", "Z"}, nil, true},
	{"date_diff", []interface{}{struct{}{}, "01-12-2017", "y"}, nil, true},
	{"date_diff", []interface{}{"01-12-2017", struct{}{}, "y"}, nil, true},
	{"date_diff", []interface{}{"01-12-2017", "01-12-2017", struct{}{}}, nil, true},
	{"date_diff", []interface{}{struct{}{}}, nil, true},

	{"date_add", []interface{}{"03-12-2017 10:15pm", "2", "y"}, time.Date(2019, 12, 03, 22, 15, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "-2", "y"}, time.Date(2015, 12, 03, 22, 15, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "2", "M"}, time.Date(2018, 2, 03, 22, 15, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "-2", "M"}, time.Date(2017, 10, 3, 22, 15, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "2", "w"}, time.Date(2017, 12, 17, 22, 15, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "-2", "w"}, time.Date(2017, 11, 19, 22, 15, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017", "2", "d"}, time.Date(2017, 12, 5, 0, 0, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017", "-4", "d"}, time.Date(2017, 11, 29, 0, 0, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "2", "h"}, time.Date(2017, 12, 4, 0, 15, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "-2", "h"}, time.Date(2017, 12, 3, 20, 15, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "105", "m"}, time.Date(2017, 12, 4, 0, 0, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "-20", "m"}, time.Date(2017, 12, 3, 21, 55, 0, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "2", "s"}, time.Date(2017, 12, 3, 22, 15, 2, 0, time.UTC), false},
	{"date_add", []interface{}{"03-12-2017 10:15pm", "-2", "s"}, time.Date(2017, 12, 3, 22, 14, 58, 0, time.UTC), false},
	{"date_add", []interface{}{struct{}{}, "2", "d"}, nil, true},
	{"date_add", []interface{}{"03-12-2017 10:15", struct{}{}, "D"}, nil, true},
	{"date_add", []interface{}{"03-12-2017 10:15", "2", struct{}{}}, nil, true},
	{"date_add", []interface{}{"03-12-2017", "2", "Z"}, nil, true},
	{"date_add", []interface{}{"22-12-2017"}, nil, true},

	{"weekday", []interface{}{"01-12-2017"}, 5, false},
	{"weekday", []interface{}{"01-12-2017 10:15pm"}, 5, false},
	{"weekday", []interface{}{struct{}{}}, nil, true},
	{"weekday", []interface{}{}, nil, true},

	{"tz", []interface{}{"01-12-2017"}, "UTC", false},
	{"tz", []interface{}{"01-12-2017 10:15:33pm"}, "UTC", false},
	{"tz", []interface{}{struct{}{}}, nil, true},
	{"tz", []interface{}{}, nil, true},

	{"tz_offset", []interface{}{"01-12-2017"}, "+0000", false},
	{"tz_offset", []interface{}{"01-12-2017 10:15:33pm"}, "+0000", false},
	{"tz_offset", []interface{}{struct{}{}}, nil, true},
	{"tz_offset", []interface{}{}, nil, true},

	{"legacy_add", []interface{}{"01-12-2017", "2"}, time.Date(2017, 12, 3, 0, 0, 0, 0, time.UTC), false},
	{"legacy_add", []interface{}{"2", "01-12-2017 10:15:33pm"}, time.Date(2017, 12, 3, 22, 15, 33, 0, time.UTC), false},
	{"legacy_add", []interface{}{"2", "3.5"}, 5.5, false},
	{"legacy_add", []interface{}{"01-12-2017 10:15:33pm", "01-12-2017"}, nil, true},
	{"legacy_add", []interface{}{math.MaxInt32 + 1, "01-12-2017 10:15:33pm"}, nil, true},
	{"legacy_add", []interface{}{"01-12-2017 10:15:33pm", math.MaxInt32 + 1}, nil, true},
	{"legacy_add", []interface{}{struct{}{}, "10"}, nil, true},
	{"legacy_add", []interface{}{"10", struct{}{}}, nil, true},
	{"legacy_add", []interface{}{}, nil, true},

	{"format_urn", []interface{}{"tel:+250781234567"}, "0781 234 567", false},
	{"format_urn", []interface{}{[]string{"tel:+250781112222", "tel:+250781234567"}}, "0781 112 222", false},
	{"format_urn", []interface{}{"twitter:134252511151#billy_bob"}, "billy_bob", false},
	{"format_urn", []interface{}{"NOT URN"}, nil, true},
}

func TestFunctions(t *testing.T) {
	env := utils.NewEnvironment(utils.DateFormat_dd_MM_yyyy, utils.TimeFormat_HH_mm_ss, time.UTC, utils.LanguageList{})

	for _, test := range funcTests {
		xFunc := XFUNCTIONS[test.name]
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic running function %s(%#v): %#v", test.name, test.args, r)
			}
		}()

		result := xFunc(env, test.args...)
		err, isErr := result.(error)

		// unexpected error
		if isErr != test.error {
			t.Errorf("Unexpected error value: %v running function %s(%#v): %s", isErr, test.name, test.args, err)
		}

		_, expectErr := test.expected.(error)

		// if this was an error and our expected isn't, move on, we have nothing to test against
		if isErr && !expectErr {
			continue
		}

		// and the match itself
		cmp, err := utils.Compare(env, result, test.expected)
		if err != nil {
			t.Errorf("Error while comparing expected: '%#v' with result: '%#v': %v for function %s(%#v)", test.expected, result, err, test.name, test.args)
		}

		if cmp != 0 {
			t.Errorf("Unexpected value, expected '%v', got '%v' for function %s(%#v)", test.expected, result, test.name, test.args)
		}
	}
}

var rangeTests = []struct {
	name        string
	args        []interface{}
	minExpected interface{}
	maxExpected interface{}
	error       bool
}{
	{"rand", []interface{}{}, 0, 1, false},
	{"rand", []interface{}{1, 10}, 1, 10, false},
	{"rand", []interface{}{10, -10}, -10, 10, false},
	{"rand", []interface{}{struct{}{}, 10}, nil, nil, true},
	{"rand", []interface{}{10, struct{}{}}, nil, nil, true},
	{"rand", []interface{}{struct{}{}}, nil, nil, true},

	{"now", []interface{}{}, time.Now().Add(time.Second * -5), time.Now().Add(time.Second * 5), false},
	{"now", []interface{}{struct{}{}}, nil, nil, true},

	{"today", []interface{}{}, time.Now().Add(time.Hour * -24), time.Now().Add(time.Second * 5), false},
	{"today", []interface{}{struct{}{}}, nil, nil, true},
}

func TestRangeFunctions(t *testing.T) {
	env := utils.NewEnvironment(utils.DateFormat_dd_MM_yyyy, utils.TimeFormat_HH_mm_ss, time.UTC, utils.LanguageList{})

	for _, test := range rangeTests {
		xFunc := XFUNCTIONS[test.name]
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic running function %s(%#v): %#v", test.name, test.args, r)
			}
		}()

		result := xFunc(env, test.args...)
		err, isErr := result.(error)

		// unexpected error
		if isErr != test.error {
			t.Errorf("Unexpected error value: %v running function %s(%#v): %s", isErr, test.name, test.args, err)
			continue
		}

		// expected error, continue, nothing to compare
		if isErr && test.error {
			continue
		}

		// and the match itself
		minCmp, err := utils.Compare(env, result, test.minExpected)
		if err != nil {
			t.Errorf("Error while comparing min expected: '%#v' with result: '%#v': %v for function %s(%#v)", test.minExpected, result, err, test.name, test.args)
		}

		maxCmp, err := utils.Compare(env, result, test.maxExpected)
		if err != nil {
			t.Errorf("Error while comparing max expected: '%#v' with result: '%#v': %v for function %s(%#v)", test.maxExpected, result, err, test.name, test.args)
		}

		if minCmp < 0 || maxCmp > 0 {
			t.Errorf("Unexpected value, expected '%v-%v', got '%v' for function %s(%#v)", test.minExpected, test.maxExpected, result, test.name, test.args)
		}
	}
}
