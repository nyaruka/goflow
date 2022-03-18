package functions

import (
	"bytes"
	"fmt"
	"html"
	"math"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/random"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
)

var nanosPerSecond = decimal.RequireFromString("1000000000")
var nonPrintableRegex = regexp.MustCompile(`[\p{Cc}\p{C}]`)

func init() {
	builtin := map[string]types.XFunc{
		// type conversion
		"text":     OneArgFunction(Text),
		"boolean":  OneArgFunction(Boolean),
		"number":   OneArgFunction(Number),
		"date":     OneArgFunction(Date),
		"datetime": OneArgFunction(DateTime),
		"time":     OneArgFunction(Time),
		"array":    Array,
		"object":   Object,

		// text functions
		"char":              OneNumberFunction(Char),
		"code":              OneTextFunction(Code),
		"split":             TextAndOptionalTextFunction(Split, types.XTextEmpty),
		"trim":              TextAndOptionalTextFunction(Trim, types.XTextEmpty),
		"trim_left":         TextAndOptionalTextFunction(TrimLeft, types.XTextEmpty),
		"trim_right":        TextAndOptionalTextFunction(TrimRight, types.XTextEmpty),
		"title":             OneTextFunction(Title),
		"word":              InitialTextFunction(1, 2, Word),
		"remove_first_word": OneTextFunction(RemoveFirstWord),
		"word_count":        TextAndOptionalTextFunction(WordCount, types.XTextEmpty),
		"word_slice":        InitialTextFunction(1, 3, WordSlice),
		"field":             InitialTextFunction(2, 2, Field),
		"clean":             OneTextFunction(Clean),
		"text_slice":        InitialTextFunction(1, 3, TextSlice),
		"lower":             OneTextFunction(Lower),
		"regex_match":       InitialTextFunction(1, 2, RegexMatch),
		"text_length":       OneTextFunction(TextLength),
		"text_compare":      TwoTextFunction(TextCompare),
		"repeat":            TextAndIntegerFunction(Repeat),
		"replace":           MinAndMaxArgsCheck(3, 4, Replace),
		"upper":             OneTextFunction(Upper),
		"percent":           OneNumberFunction(Percent),
		"url_encode":        OneTextFunction(URLEncode),
		"html_decode":       OneTextFunction(HTMLDecode),

		// bool functions
		"and": MinArgsCheck(1, And),
		"if":  ThreeArgFunction(If),
		"or":  MinArgsCheck(1, Or),

		// number functions
		"round":        OneNumberAndOptionalIntegerFunction(Round, 0),
		"round_up":     OneNumberAndOptionalIntegerFunction(RoundUp, 0),
		"round_down":   OneNumberAndOptionalIntegerFunction(RoundDown, 0),
		"max":          MinArgsCheck(1, Max),
		"min":          MinArgsCheck(1, Min),
		"mean":         MinArgsCheck(1, Mean),
		"mod":          TwoNumberFunction(Mod),
		"rand":         NoArgFunction(Rand),
		"rand_between": TwoNumberFunction(RandBetween),
		"abs":          OneNumberFunction(Abs),

		// datetime functions
		"parse_datetime":      MinAndMaxArgsCheck(2, 3, ParseDateTime),
		"datetime_from_epoch": OneNumberFunction(DateTimeFromEpoch),
		"datetime_diff":       ThreeArgFunction(DateTimeDiff),
		"datetime_add":        DateTimeAdd,
		"replace_time":        TwoArgFunction(ReplaceTime),
		"tz":                  OneDateTimeFunction(TZ),
		"tz_offset":           OneDateTimeFunction(TZOffset),
		"now":                 NoArgFunction(Now),
		"epoch":               OneDateTimeFunction(Epoch),

		// date functions
		"date_from_parts": ThreeIntegerFunction(DateFromParts),
		"weekday":         OneDateFunction(Weekday),
		"week_number":     OneDateFunction(WeekNumber),
		"today":           NoArgFunction(Today),

		// time functions
		"parse_time":      TwoArgFunction(ParseTime),
		"time_from_parts": ThreeIntegerFunction(TimeFromParts),

		// array functions
		"join":    TwoArgFunction(Join),
		"reverse": OneArrayFunction(Reverse),
		"sort":    OneArrayFunction(Sort),
		"sum":     OneArrayFunction(Sum),
		"unique":  OneArrayFunction(Unique),
		"concat":  TwoArrayFunction(Concat),

		// encoded text functions
		"urn_parts":        OneTextFunction(URNParts),
		"attachment_parts": OneTextFunction(AttachmentParts),

		// json functions
		"json":       OneArgFunction(JSON),
		"parse_json": OneTextFunction(ParseJSON),

		// formatting functions
		"format":          OneArgFunction(Format),
		"format_date":     MinAndMaxArgsCheck(1, 2, FormatDate),
		"format_datetime": MinAndMaxArgsCheck(1, 3, FormatDateTime),
		"format_time":     MinAndMaxArgsCheck(1, 2, FormatTime),
		"format_location": OneTextFunction(FormatLocation),
		"format_number":   MinAndMaxArgsCheck(1, 3, FormatNumber),
		"format_urn":      OneTextFunction(FormatURN),

		// utility functions
		"is_error":       OneArgFunction(IsError),
		"count":          OneArgFunction(Count),
		"default":        TwoArgFunction(Default),
		"legacy_add":     TwoArgFunction(LegacyAdd),
		"read_chars":     OneTextFunction(ReadChars),
		"extract":        TwoArgFunction(Extract),
		"extract_object": MinArgsCheck(2, ExtractObject),
		"foreach":        MinArgsCheck(2, ForEach),
		"foreach_value":  MinArgsCheck(2, ForEachValue),
	}

	for name, fn := range builtin {
		RegisterXFunction(name, fn)
	}
}

//------------------------------------------------------------------------------------------
// Type Conversion Functions
//------------------------------------------------------------------------------------------

// Text tries to convert `value` to text.
//
// An error is returned if the value can't be converted.
//
//   @(text(3 = 3)) -> true
//   @(json(text(123.45))) -> "123.45"
//   @(text(1 / 0)) -> ERROR
//
// @function text(value)
func Text(env envs.Environment, value types.XValue) types.XValue {
	str, xerr := types.ToXText(env, value)
	if xerr != nil {
		return xerr
	}
	return str
}

// Boolean tries to convert `value` to a boolean.
//
// An error is returned if the value can't be converted.
//
//   @(boolean(array(1, 2))) -> true
//   @(boolean("FALSE")) -> false
//   @(boolean(1 / 0)) -> ERROR
//
// @function boolean(value)
func Boolean(env envs.Environment, value types.XValue) types.XValue {
	str, xerr := types.ToXBoolean(value)
	if xerr != nil {
		return xerr
	}
	return str
}

// Number tries to convert `value` to a number.
//
// An error is returned if the value can't be converted.
//
//   @(number(10)) -> 10
//   @(number("123.45000")) -> 123.45
//   @(number("what?")) -> ERROR
//
// @function number(value)
func Number(env envs.Environment, value types.XValue) types.XValue {
	num, xerr := types.ToXNumber(env, value)
	if xerr != nil {
		return xerr
	}
	return num
}

// Date tries to convert `value` to a date.
//
// If it is text then it will be parsed into a date using the default date format.
// An error is returned if the value can't be converted.
//
//   @(date("1979-07-18")) -> 1979-07-18
//   @(date("1979-07-18T10:30:45.123456Z")) -> 1979-07-18
//   @(date("10/05/2010")) -> 2010-05-10
//   @(date("NOT DATE")) -> ERROR
//
// @function date(value)
func Date(env envs.Environment, value types.XValue) types.XValue {
	d, err := types.ToXDate(env, value)
	if err != nil {
		return types.NewXError(err)
	}
	return d
}

// DateTime tries to convert `value` to a datetime.
//
// If it is text then it will be parsed into a datetime using the default date
// and time formats. An error is returned if the value can't be converted.
//
//   @(datetime("1979-07-18")) -> 1979-07-18T00:00:00.000000-05:00
//   @(datetime("1979-07-18T10:30:45.123456Z")) -> 1979-07-18T10:30:45.123456Z
//   @(datetime("10/05/2010")) -> 2010-05-10T00:00:00.000000-05:00
//   @(datetime("NOT DATE")) -> ERROR
//
// @function datetime(value)
func DateTime(env envs.Environment, value types.XValue) types.XValue {
	dt, err := types.ToXDateTime(env, value)
	if err != nil {
		return types.NewXError(err)
	}
	return dt
}

// Time tries to convert `value` to a time.
//
// If it is text then it will be parsed into a time using the default time format.
// An error is returned if the value can't be converted.
//
//   @(time("10:30")) -> 10:30:00.000000
//   @(time("10:30:45 PM")) -> 22:30:45.000000
//   @(time(datetime("1979-07-18T10:30:45.123456Z"))) -> 10:30:45.123456
//   @(time("what?")) -> ERROR
//
// @function time(value)
func Time(env envs.Environment, value types.XValue) types.XValue {
	t, xerr := types.ToXTime(env, value)
	if xerr != nil {
		return xerr
	}
	return t
}

// Array takes multiple `values` and returns them as an array.
//
//   @(array("a", "b", 356)[1]) -> b
//   @(join(array("a", "b", "c"), "|")) -> a|b|c
//   @(count(array())) -> 0
//   @(count(array("a", "b"))) -> 2
//
// @function array(values...)
func Array(env envs.Environment, values ...types.XValue) types.XValue {
	// check none of our args are errors
	for _, arg := range values {
		if types.IsXError(arg) {
			return arg
		}
	}

	return types.NewXArray(values...)
}

// Object takes property name value pairs and returns them as a new object.
//
//   @(object()) -> {}
//   @(object("a", 123, "b", "hello")) -> {a: 123, b: hello}
//   @(object("a")) -> ERROR
//
// @function object(pairs...)
func Object(env envs.Environment, pairs ...types.XValue) types.XValue {
	// check none of our args are errors
	for _, arg := range pairs {
		if types.IsXError(arg) {
			return arg
		}
	}

	if len(pairs)%2 != 0 {
		return types.NewXErrorf("requires an even number of arguments")
	}

	properties := make(map[string]types.XValue, len(pairs)/2)

	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i]
		value := pairs[i+1]

		keyAsText, xerr := types.ToXText(env, key)
		if xerr != nil {
			return xerr
		}

		properties[keyAsText.Native()] = value
	}

	return types.NewXObject(properties)
}

//------------------------------------------------------------------------------------------
// Bool Functions
//------------------------------------------------------------------------------------------

// And returns whether all the given `values` are truthy.
//
//   @(and(true)) -> true
//   @(and(true, false, true)) -> false
//
// @function and(values...)
func And(env envs.Environment, values ...types.XValue) types.XValue {
	for _, arg := range values {
		asBool, xerr := types.ToXBoolean(arg)
		if xerr != nil {
			return xerr
		}
		if !asBool.Native() {
			return types.XBooleanFalse
		}
	}
	return types.XBooleanTrue
}

// Or returns whether if any of the given `values` are truthy.
//
//   @(or(true)) -> true
//   @(or(true, false, true)) -> true
//
// @function or(values...)
func Or(env envs.Environment, values ...types.XValue) types.XValue {
	for _, arg := range values {
		asBool, xerr := types.ToXBoolean(arg)
		if xerr != nil {
			return xerr
		}
		if asBool.Native() {
			return types.XBooleanTrue
		}
	}
	return types.XBooleanFalse
}

// If returns `value1` if `test` is truthy or `value2` if not.
//
// If the first argument is an error that error is returned.
//
//   @(if(1 = 1, "foo", "bar")) -> foo
//   @(if("foo" > "bar", "foo", "bar")) -> ERROR
//
// @function if(test, value1, value2)
func If(env envs.Environment, test types.XValue, value1 types.XValue, value2 types.XValue) types.XValue {
	asBool, err := types.ToXBoolean(test)
	if err != nil {
		return err
	}

	if asBool.Native() {
		return value1
	}
	return value2
}

//------------------------------------------------------------------------------------------
// Text Functions
//------------------------------------------------------------------------------------------

// Code returns the UNICODE code for the first character of `text`.
//
// It is the inverse of [function:char].
//
//   @(code("a")) -> 97
//   @(code("abc")) -> 97
//   @(code("😀")) -> 128512
//   @(code("15")) -> 49
//   @(code(15)) -> 49
//   @(code("")) -> ERROR
//
// @function code(text)
func Code(env envs.Environment, text types.XText) types.XValue {
	if text.Length() == 0 {
		return types.NewXErrorf("requires a string of at least one character")
	}

	r, _ := utf8.DecodeRuneInString(text.Native())
	return types.NewXNumberFromInt(int(r))
}

// Split splits `text` into an array of separated words.
//
// Empty values are removed from the returned list. There is an optional final parameter `delimiters` which
// is string of characters used to split the text into words.
//
//   @(split("a b c")) -> [a, b, c]
//   @(split("a", " ")) -> [a]
//   @(split("abc..d", ".")) -> [abc, d]
//   @(split("a.b.c.", ".")) -> [a, b, c]
//   @(split("a|b,c  d", " .|,")) -> [a, b, c, d]
//
// @function split(text, [,delimiters])
func Split(env envs.Environment, text types.XText, delimiters types.XText) types.XValue {
	splits := extractWords(text.Native(), delimiters.Native())

	nonEmpty := make([]types.XValue, 0)
	for _, split := range splits {
		nonEmpty = append(nonEmpty, types.NewXText(split))
	}
	return types.NewXArray(nonEmpty...)
}

// Trim removes whitespace from either end of `text`.
//
// There is an optional final parameter `chars` which is string of characters to be removed instead of whitespace.
//
//   @(trim(" hello world    ")) -> hello world
//   @(trim("+123157568", "+")) -> 123157568
//
// @function trim(text, [,chars])
func Trim(env envs.Environment, text types.XText, chars types.XText) types.XValue {
	if chars != types.XTextEmpty {
		return types.NewXText(strings.Trim(text.Native(), chars.Native()))
	}

	return types.NewXText(strings.TrimSpace(text.Native()))
}

// TrimLeft removes whitespace from the start of `text`.
//
// There is an optional final parameter `chars` which is string of characters to be removed instead of whitespace.
//
//   @("*" & trim_left(" hello world   ") & "*") -> *hello world   *
//   @(trim_left("+12345+", "+")) -> 12345+
//
// @function trim_left(text, [,chars])
func TrimLeft(env envs.Environment, text types.XText, chars types.XText) types.XValue {
	if chars != types.XTextEmpty {
		return types.NewXText(strings.TrimLeft(text.Native(), chars.Native()))
	}

	return types.NewXText(strings.TrimLeftFunc(text.Native(), unicode.IsSpace))
}

// TrimRight removes whitespace from the end of `text`.
//
// There is an optional final parameter `chars` which is string of characters to be removed instead of whitespace.
//
//   @("*" & trim_right(" hello world   ") & "*") -> * hello world*
//   @(trim_right("+12345+", "+")) -> +12345
//
// @function trim_right(text, [,chars])
func TrimRight(env envs.Environment, text types.XText, chars types.XText) types.XValue {
	if chars != types.XTextEmpty {
		return types.NewXText(strings.TrimRight(text.Native(), chars.Native()))
	}

	return types.NewXText(strings.TrimRightFunc(text.Native(), unicode.IsSpace))
}

// Char returns the character for the given UNICODE `code`.
//
// It is the inverse of [function:code].
//
//   @(char(33)) -> !
//   @(char(128512)) -> 😀
//   @(char("foo")) -> ERROR
//
// @function char(code)
func Char(env envs.Environment, num types.XNumber) types.XValue {
	code, xerr := types.ToInteger(env, num)
	if xerr != nil {
		return xerr
	}

	return types.NewXText(string(rune(code)))
}

// Title capitalizes each word in `text`.
//
//   @(title("foo")) -> Foo
//   @(title("ryan lewis")) -> Ryan Lewis
//   @(title("RYAN LEWIS")) -> Ryan Lewis
//   @(title(123)) -> 123
//
// @function title(text)
func Title(env envs.Environment, text types.XText) types.XValue {
	return types.NewXText(strings.Title(strings.ToLower(text.Native())))
}

// Word returns the word at `index` in `text`.
//
// Indexes start at zero. There is an optional final parameter `delimiters` which
// is string of characters used to split the text into words.
//
//   @(word("bee cat dog", 0)) -> bee
//   @(word("bee.cat,dog", 0)) -> bee
//   @(word("bee.cat,dog", 1)) -> cat
//   @(word("bee.cat,dog", 2)) -> dog
//   @(word("bee.cat,dog", -1)) -> dog
//   @(word("bee.cat,dog", -2)) -> cat
//   @(word("bee.*cat,dog", 1, ".*=|")) -> cat,dog
//   @(word("O'Grady O'Flaggerty", 1, " ")) -> O'Flaggerty
//
// @function word(text, index [,delimiters])
func Word(env envs.Environment, text types.XText, args ...types.XValue) types.XValue {
	index, xerr := types.ToInteger(env, args[0])
	if xerr != nil {
		return xerr
	}

	delimiters := types.XTextEmpty
	if len(args) == 2 && args[1] != nil {
		delimiters, xerr = types.ToXText(env, args[1])
		if xerr != nil {
			return xerr
		}
	}

	words := extractWords(text.Native(), delimiters.Native())

	offset := index
	if offset < 0 {
		offset += len(words)
	}

	if !(offset >= 0 && offset < len(words)) {
		return types.NewXErrorf("index %d is out of range for the number of words %d", index, len(words))
	}

	return types.NewXText(words[offset])
}

// RemoveFirstWord removes the first word of `text`.
//
//   @(remove_first_word("foo bar")) -> bar
//   @(remove_first_word("Hi there. I'm a flow!")) -> there. I'm a flow!
//
// @function remove_first_word(text)
func RemoveFirstWord(env envs.Environment, text types.XText) types.XValue {
	s := text.Native()
	words := extractWords(s, "")
	if len(words) < 2 {
		return types.XTextEmpty
	}

	// find first word and remove
	w1Start := strings.Index(s, words[0])
	s = s[w1Start+len(words[0]):]

	// find where second word starts and discard everything up to that
	w2Start := strings.Index(s, words[1])
	s = s[w2Start:]

	return types.NewXText(s)
}

// WordSlice extracts a sub-sequence of words from `text`.
//
// The returned words are those from `start` up to but not-including `end`. Indexes start at zero and a negative
// end value means that all words after the start should be returned. There is an optional final parameter `delimiters`
// which is string of characters used to split the text into words.
//
//   @(word_slice("bee cat dog", 0, 1)) -> bee
//   @(word_slice("bee cat dog", 0, 2)) -> bee cat
//   @(word_slice("bee cat dog", 1, -1)) -> cat dog
//   @(word_slice("bee cat dog", 1)) -> cat dog
//   @(word_slice("bee cat dog", 2, 3)) -> dog
//   @(word_slice("bee cat dog", 3, 10)) ->
//   @(word_slice("bee.*cat,dog", 1, -1, ".*=|,")) -> cat dog
//   @(word_slice("O'Grady O'Flaggerty", 1, 2, " ")) -> O'Flaggerty
//
// @function word_slice(text, start, end [,delimiters])
func WordSlice(env envs.Environment, text types.XText, args ...types.XValue) types.XValue {
	start, xerr := types.ToInteger(env, args[0])
	if xerr != nil {
		return xerr
	}
	if start < 0 {
		return types.NewXErrorf("must start with a positive index")
	}

	end := -1
	if len(args) >= 2 {
		if end, xerr = types.ToInteger(env, args[1]); xerr != nil {
			return xerr
		}
	}
	if end > 0 && end <= start {
		return types.NewXErrorf("must have a end which is greater than the start")
	}

	delimiters := types.XTextEmpty
	if len(args) >= 3 && args[2] != nil {
		delimiters, xerr = types.ToXText(env, args[2])
		if xerr != nil {
			return xerr
		}
	}

	words := extractWords(text.Native(), delimiters.Native())

	if start >= len(words) {
		return types.XTextEmpty
	}
	if end >= len(words) {
		end = len(words)
	}

	if end > 0 {
		return types.NewXText(strings.Join(words[start:end], " "))
	}
	return types.NewXText(strings.Join(words[start:], " "))
}

// WordCount returns the number of words in `text`.
//
// There is an optional final parameter `delimiters` which is string of characters used
// to split the text into words.
//
//   @(word_count("foo bar")) -> 2
//   @(word_count(10)) -> 1
//   @(word_count("")) -> 0
//   @(word_count("😀😃😄😁")) -> 4
//   @(word_count("bee.*cat,dog", ".*=|")) -> 2
//   @(word_count("O'Grady O'Flaggerty", " ")) -> 2
//
// @function word_count(text [,delimiters])
func WordCount(env envs.Environment, text types.XText, delimiters types.XText) types.XValue {
	words := extractWords(text.Native(), delimiters.Native())

	return types.NewXNumberFromInt(len(words))
}

// Field splits `text` using the given `delimiter` and returns the field at `index`.
//
// The index starts at zero. When splitting with a space, the delimiter is considered to be all whitespace.
//
//   @(field("a,b,c", 1, ",")) -> b
//   @(field("a,,b,c", 1, ",")) ->
//   @(field("a   b c", 1, " ")) -> b
//   @(field("a		b	c	d", 1, "	")) ->
//   @(field("a\t\tb\tc\td", 1, " ")) ->
//   @(field("a,b,c", "foo", ",")) -> ERROR
//
// @function field(text, index, delimiter)
func Field(env envs.Environment, text types.XText, args ...types.XValue) types.XValue {
	field, xerr := types.ToInteger(env, args[0])
	if xerr != nil {
		return xerr
	}

	if field < 0 {
		return types.NewXErrorf("cannot use a negative index")
	}

	sep, xerr := types.ToXText(env, args[1])
	if xerr != nil {
		return xerr
	}

	fields := strings.Split(text.Native(), sep.Native())

	// when using a space as a delimiter, we consider it splitting on whitespace, so remove empty values
	if sep.Native() == " " {
		var newFields []string
		for _, f := range fields {
			if f != "" {
				newFields = append(newFields, f)
			}
		}
		fields = newFields
	}

	if field >= len(fields) {
		return types.XTextEmpty
	}

	return types.NewXText(strings.TrimSpace(fields[field]))
}

// Clean removes any non-printable characters from `text`.
//
//   @(clean("😃 Hello \nwo\tr\rld")) -> 😃 Hello world
//   @(clean(123)) -> 123
//
// @function clean(text)
func Clean(env envs.Environment, text types.XText) types.XValue {
	return types.NewXText(nonPrintableRegex.ReplaceAllString(text.Native(), ""))
}

// TextSlice returns the portion of `text` between `start` (inclusive) and `end` (exclusive).
//
// If `end` is not specified then the entire rest of `text` will be included. Negative values
// for `start` or `end` start at the end of `text`.
//
//   @(text_slice("hello", 2)) -> llo
//   @(text_slice("hello", 1, 3)) -> el
//   @(text_slice("hello😁", -3, -1)) -> lo
//   @(text_slice("hello", 7)) ->
//
// @function text_slice(text, start [, end])
func TextSlice(env envs.Environment, text types.XText, args ...types.XValue) types.XValue {
	length := utf8.RuneCountInString(text.Native())

	start, xerr := types.ToInteger(env, args[0])
	if xerr != nil {
		return xerr
	}
	if start < 0 {
		start = length + start
	}

	end := length
	if len(args) == 2 {
		if end, xerr = types.ToInteger(env, args[1]); xerr != nil {
			return xerr
		}
	}
	if end < 0 {
		end = length + end
	}

	var output bytes.Buffer
	i := 0
	for _, r := range text.Native() {
		if i >= start && i < end {
			output.WriteRune(r)
		}
		i++
	}

	return types.NewXText(output.String())
}

// Lower converts `text` to lowercase.
//
//   @(lower("HellO")) -> hello
//   @(lower("hello")) -> hello
//   @(lower("123")) -> 123
//   @(lower("😀")) -> 😀
//
// @function lower(text)
func Lower(env envs.Environment, text types.XText) types.XValue {
	return types.NewXText(strings.ToLower(text.Native()))
}

// RegexMatch returns the first match of the regular expression `pattern` in `text`.
//
// An optional third parameter `group` determines which matching group will be returned.
//
//   @(regex_match("sda34dfddg67", "\d+")) -> 34
//   @(regex_match("Bob Smith", "(\w+) (\w+)", 1)) -> Bob
//   @(regex_match("Bob Smith", "(\w+) (\w+)", 2)) -> Smith
//   @(regex_match("Bob Smith", "(\w+) (\w+)", 5)) -> ERROR
//   @(regex_match("abc", "[\.")) -> ERROR
//
// @function regex_match(text, pattern [,group])
func RegexMatch(env envs.Environment, text types.XText, args ...types.XValue) types.XValue {
	pattern, xerr := types.ToXText(env, args[0])
	if xerr != nil {
		return xerr
	}

	groupNum := 0
	if len(args) == 2 {
		groupNum, xerr = types.ToInteger(env, args[1])
		if xerr != nil {
			return xerr
		}
	}

	exp, err := regexp.Compile(`(?mi)` + pattern.Native())
	if err != nil {
		return types.NewXErrorf("invalid regular expression")
	}

	groups := exp.FindStringSubmatch(text.Native())

	if groupNum < 0 || groupNum >= len(groups) {
		return types.NewXErrorf("invalid regular expression group")
	}

	return types.NewXText(groups[groupNum])
}

// TextLength returns the length (number of characters) of `value` when converted to text.
//
//   @(text_length("abc")) -> 3
//   @(text_length(array(2, 3))) -> 6
//
// @function text_length(value)
func TextLength(env envs.Environment, value types.XText) types.XValue {
	return types.NewXNumberFromInt(value.Length())
}

// TextCompare returns the dictionary order of `text1` and `text2`.
//
// The return value will be -1 if `text1` comes before `text2`, 0 if they are equal
// and 1 if `text1` comes after `text2`.
//
//   @(text_compare("abc", "abc")) -> 0
//   @(text_compare("abc", "def")) -> -1
//   @(text_compare("zzz", "aaa")) -> 1
//
// @function text_compare(text1, text2)
func TextCompare(env envs.Environment, text1 types.XText, text2 types.XText) types.XValue {
	return types.NewXNumberFromInt(text1.Compare(text2))
}

// Repeat returns `text` repeated `count` number of times.
//
//   @(repeat("*", 8)) -> ********
//   @(repeat("*", "foo")) -> ERROR
//
// @function repeat(text, count)
func Repeat(env envs.Environment, text types.XText, count int) types.XValue {
	if count < 0 {
		return types.NewXErrorf("must be called with a positive integer, got %d", count)
	}

	var output bytes.Buffer
	for j := 0; j < count; j++ {
		output.WriteString(text.Native())
	}

	return types.NewXText(output.String())
}

// Replace replaces up to `count` occurrences of `needle` with `replacement` in `text`.
//
// If `count` is omitted or is less than 0 then all occurrences are replaced.
//
//   @(replace("foo bar foo", "foo", "zap")) -> zap bar zap
//   @(replace("foo bar foo", "foo", "zap", 1)) -> zap bar foo
//   @(replace("foo bar", "baz", "zap")) -> foo bar
//
// @function replace(text, needle, replacement [, count])
func Replace(env envs.Environment, args ...types.XValue) types.XValue {
	text, xerr := types.ToXText(env, args[0])
	if xerr != nil {
		return xerr
	}
	needle, xerr := types.ToXText(env, args[1])
	if xerr != nil {
		return xerr
	}
	replacement, xerr := types.ToXText(env, args[2])
	if xerr != nil {
		return xerr
	}

	count := -1
	if len(args) == 4 {
		count, xerr = types.ToInteger(env, args[3])
		if xerr != nil {
			return xerr
		}
	}

	return types.NewXText(strings.Replace(text.Native(), needle.Native(), replacement.Native(), count))
}

// Upper converts `text` to uppercase.
//
//   @(upper("Asdf")) -> ASDF
//   @(upper(123)) -> 123
//
// @function upper(text)
func Upper(env envs.Environment, text types.XText) types.XValue {
	return types.NewXText(strings.ToUpper(text.Native()))
}

// Percent formats `number` as a percentage.
//
//   @(percent(0.54234)) -> 54%
//   @(percent(1.2)) -> 120%
//   @(percent("foo")) -> ERROR
//
// @function percent(number)
func Percent(env envs.Environment, num types.XNumber) types.XValue {
	// multiply by 100 and floor
	percent := num.Native().Mul(decimal.NewFromFloat(100)).Round(0)

	// add on a %
	return types.NewXText(fmt.Sprintf("%d%%", percent.IntPart()))
}

// URLEncode encodes `text` for use as a URL parameter.
//
//   @(url_encode("two & words")) -> two%20%26%20words
//   @(url_encode(10)) -> 10
//
// @function url_encode(text)
func URLEncode(env envs.Environment, text types.XText) types.XValue {
	// escapes spaces as %20 matching urllib.quote(s, safe="") in Python
	encoded := strings.Replace(url.QueryEscape(text.Native()), "+", "%20", -1)
	return types.NewXText(encoded)
}

// HTMLDecode HTML decodes `text`
//
//   @(html_decode("Red &amp; Blue")) -> Red & Blue
//   @(html_decode("5 + 10")) -> 5 + 10
//
// @function html_decode(text)
func HTMLDecode(env envs.Environment, text types.XText) types.XValue {
	decoded := html.UnescapeString(text.Native())

	// the common nbsp; turns into a unicode non breaking space, convert to a normal space
	decoded = strings.ReplaceAll(decoded, "\U000000A0", " ")
	return types.NewXText(decoded)
}

//------------------------------------------------------------------------------------------
// Number Functions
//------------------------------------------------------------------------------------------

// Abs returns the absolute value of `number`.
//
//   @(abs(-10)) -> 10
//   @(abs(10.5)) -> 10.5
//   @(abs("foo")) -> ERROR
//
// @function abs(number)
func Abs(env envs.Environment, num types.XNumber) types.XValue {
	return types.NewXNumber(num.Native().Abs())
}

// Round rounds `number` to the nearest value.
//
// You can optionally pass in the number of decimal places to round to as `places`. If `places` < 0,
// it will round the integer part to the nearest 10^(-places).
//
//   @(round(12)) -> 12
//   @(round(12.141)) -> 12
//   @(round(12.6)) -> 13
//   @(round(12.141, 2)) -> 12.14
//   @(round(12.146, 2)) -> 12.15
//   @(round(12.146, -1)) -> 10
//   @(round("notnum", 2)) -> ERROR
//
// @function round(number [,places])
func Round(env envs.Environment, num types.XNumber, places int) types.XValue {
	return types.NewXNumber(num.Native().Round(int32(places)))
}

// RoundUp rounds `number` up to the nearest integer value.
//
// You can optionally pass in the number of decimal places to round to as `places`.
//
//   @(round_up(12)) -> 12
//   @(round_up(12.141)) -> 13
//   @(round_up(12.6)) -> 13
//   @(round_up(12.141, 2)) -> 12.15
//   @(round_up(12.146, 2)) -> 12.15
//   @(round_up("foo")) -> ERROR
//
// @function round_up(number [,places])
func RoundUp(env envs.Environment, num types.XNumber, places int) types.XValue {
	dec := num.Native()
	if dec.Round(int32(places)).Equal(dec) {
		return num
	}

	halfPrecision := decimal.New(5, -int32(places)-1)
	roundedDec := dec.Add(halfPrecision).Round(int32(places))

	return types.NewXNumber(roundedDec)
}

// RoundDown rounds `number` down to the nearest integer value.
//
// You can optionally pass in the number of decimal places to round to as `places`.
//
//   @(round_down(12)) -> 12
//   @(round_down(12.141)) -> 12
//   @(round_down(12.6)) -> 12
//   @(round_down(12.141, 2)) -> 12.14
//   @(round_down(12.146, 2)) -> 12.14
//   @(round_down("foo")) -> ERROR
//
// @function round_down(number [,places])
func RoundDown(env envs.Environment, num types.XNumber, places int) types.XValue {
	dec := num.Native()
	if dec.Round(int32(places)).Equal(dec) {
		return num
	}

	halfPrecision := decimal.New(5, -int32(places)-1)
	roundedDec := dec.Sub(halfPrecision).Round(int32(places))

	return types.NewXNumber(roundedDec)
}

// Max returns the maximum value in `numbers`.
//
//   @(max(1, 2)) -> 2
//   @(max(1, -1, 10)) -> 10
//   @(max(1, 10, "foo")) -> ERROR
//
// @function max(numbers...)
func Max(env envs.Environment, values ...types.XValue) types.XValue {
	max, xerr := types.ToXNumber(env, values[0])
	if xerr != nil {
		return xerr
	}

	for _, v := range values[1:] {
		val, xerr := types.ToXNumber(env, v)
		if xerr != nil {
			return xerr
		}

		if val.Compare(max) > 0 {
			max = val
		}
	}
	return max
}

// Min returns the minimum value in `numbers`.
//
//   @(min(1, 2)) -> 1
//   @(min(2, 2, -10)) -> -10
//   @(min(1, 2, "foo")) -> ERROR
//
// @function min(numbers...)
func Min(env envs.Environment, values ...types.XValue) types.XValue {
	max, xerr := types.ToXNumber(env, values[0])
	if xerr != nil {
		return xerr
	}

	for _, v := range values[1:] {
		val, xerr := types.ToXNumber(env, v)
		if xerr != nil {
			return xerr
		}

		if val.Compare(max) < 0 {
			max = val
		}
	}
	return max
}

// Mean returns the arithmetic mean of `numbers`.
//
//   @(mean(1, 2)) -> 1.5
//   @(mean(1, 2, 6)) -> 3
//   @(mean(1, "foo")) -> ERROR
//
// @function mean(numbers...)
func Mean(env envs.Environment, args ...types.XValue) types.XValue {
	sum := decimal.Zero

	for _, val := range args {
		num, xerr := types.ToXNumber(env, val)
		if xerr != nil {
			return xerr
		}
		sum = sum.Add(num.Native())
	}

	return types.NewXNumber(sum.Div(decimal.New(int64(len(args)), 0)))
}

// Mod returns the remainder of the division of `dividend` by `divisor`.
//
//   @(mod(5, 2)) -> 1
//   @(mod(4, 2)) -> 0
//   @(mod(5, "foo")) -> ERROR
//
// @function mod(dividend, divisor)
func Mod(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXNumber(num1.Native().Mod(num2.Native()))
}

// Rand returns a single random number between [0.0-1.0).
//
//   @(rand()) -> 0.6075520156746239
//   @(rand()) -> 0.48467757094734026
//
// @function rand()
func Rand(env envs.Environment) types.XValue {
	return types.NewXNumber(random.Decimal())
}

// RandBetween a single random integer in the given inclusive range.
//
//   @(rand_between(1, 10)) -> 10
//   @(rand_between(1, 10)) -> 2
//
// @function rand_between()
func RandBetween(env envs.Environment, min types.XNumber, max types.XNumber) types.XValue {
	span := (max.Native().Sub(min.Native())).Add(decimal.New(1, 0))

	val := random.Decimal().Mul(span).Add(min.Native()).Floor()

	return types.NewXNumber(val)
}

//------------------------------------------------------------------------------------------
// Date & Time Functions
//------------------------------------------------------------------------------------------

// ParseDateTime parses `text` into a date using the given `format`.
//
// The format string can consist of the following characters. The characters
// ' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.
//
// * `YY`        - last two digits of year 0-99
// * `YYYY`      - four digits of year 0000-9999
// * `M`         - month 1-12
// * `MM`        - month, zero padded 01-12
// * `D`         - day of month, 1-31
// * `DD`        - day of month, zero padded 01-31
// * `h`         - hour of the day 1-12
// * `hh`        - hour of the day 01-12
// * `t`         - twenty four hour of the day 1-23
// * `tt`        - twenty four hour of the day, zero padded 01-23
// * `m`         - minute 0-59
// * `mm`        - minute, zero padded 00-59
// * `s`         - second 0-59
// * `ss`        - second, zero padded 00-59
// * `fff`       - milliseconds
// * `ffffff`    - microseconds
// * `fffffffff` - nanoseconds
// * `aa`        - am or pm
// * `AA`        - AM or PM
// * `Z`         - hour and minute offset from UTC, or Z for UTC
// * `ZZZ`       - hour and minute offset from UTC
//
// Timezone should be a location name as specified in the IANA Time Zone database, such
// as "America/Guayaquil" or "America/Los_Angeles". If not specified, the current timezone
// will be used. An error will be returned if the timezone is not recognized.
//
// Note that fractional seconds will be parsed even without an explicit format identifier.
// You should only specify fractional seconds when you want to assert the number of places
// in the input format.
//
// parse_datetime will return an error if it is unable to convert the text to a datetime.
//
//   @(parse_datetime("1979-07-18", "YYYY-MM-DD")) -> 1979-07-18T00:00:00.000000-05:00
//   @(parse_datetime("2010 5 10", "YYYY M DD")) -> 2010-05-10T00:00:00.000000-05:00
//   @(parse_datetime("2010 5 10 12:50", "YYYY M DD tt:mm", "America/Los_Angeles")) -> 2010-05-10T12:50:00.000000-07:00
//   @(parse_datetime("NOT DATE", "YYYY-MM-DD")) -> ERROR
//
// @function parse_datetime(text, format [,timezone])
func ParseDateTime(env envs.Environment, args ...types.XValue) types.XValue {
	str, xerr := types.ToXText(env, args[0])
	if xerr != nil {
		return xerr
	}

	layout, xerr := types.ToXText(env, args[1])
	if xerr != nil {
		return xerr
	}

	// grab our location
	location := env.Timezone()
	if len(args) == 3 {
		tzStr, xerr := types.ToXText(env, args[2])
		if xerr != nil {
			return xerr
		}

		var err error
		location, err = time.LoadLocation(tzStr.Native())
		if err != nil {
			return types.NewXError(err)
		}
	}

	// finally try to parse the date
	parsed, err := dates.ParseDateTime(layout.Native(), str.Native(), location)
	if err != nil {
		return types.NewXError(err)
	}

	return types.NewXDateTime(parsed)
}

// DateTimeFromEpoch converts the UNIX epoch time `seconds` into a new date.
//
//   @(datetime_from_epoch(1497286619)) -> 2017-06-12T11:56:59.000000-05:00
//   @(datetime_from_epoch(1497286619.123456)) -> 2017-06-12T11:56:59.123456-05:00
//
// @function datetime_from_epoch(seconds)
func DateTimeFromEpoch(env envs.Environment, num types.XNumber) types.XValue {
	nanos := num.Native().Mul(nanosPerSecond).IntPart()
	return types.NewXDateTime(time.Unix(0, nanos).In(env.Timezone()))
}

// DateTimeDiff returns the duration between `date1` and `date2` in the `unit` specified.
//
// Valid durations are "Y" for years, "M" for months, "W" for weeks, "D" for days, "h" for hour,
// "m" for minutes, "s" for seconds.
//
//   @(datetime_diff("2017-01-15", "2017-01-17", "D")) -> 2
//   @(datetime_diff("2017-01-15", "2017-05-15", "W")) -> 17
//   @(datetime_diff("2017-01-15", "2017-05-15", "M")) -> 4
//   @(datetime_diff("2017-01-17 10:50", "2017-01-17 12:30", "h")) -> 1
//   @(datetime_diff("2017-01-17", "2015-12-17", "Y")) -> -2
//
// @function datetime_diff(date1, date2, unit)
func DateTimeDiff(env envs.Environment, arg1 types.XValue, arg2 types.XValue, arg3 types.XValue) types.XValue {
	date1, xerr := types.ToXDateTime(env, arg1)
	if xerr != nil {
		return xerr
	}

	date2, xerr := types.ToXDateTime(env, arg2)
	if xerr != nil {
		return xerr
	}

	unit, xerr := types.ToXText(env, arg3)
	if xerr != nil {
		return xerr
	}

	// find the duration between our dates
	duration := date2.Native().Sub(date1.Native())

	// then convert based on our unit
	switch unit.Native() {
	case "s":
		return types.NewXNumberFromInt(int(duration / time.Second))
	case "m":
		return types.NewXNumberFromInt(int(duration / time.Minute))
	case "h":
		return types.NewXNumberFromInt(int(duration / time.Hour))
	case "D":
		return types.NewXNumberFromInt(dates.DaysBetween(date2.Native(), date1.Native()))
	case "W":
		return types.NewXNumberFromInt(int(dates.DaysBetween(date2.Native(), date1.Native()) / 7))
	case "M":
		return types.NewXNumberFromInt(dates.MonthsBetween(date2.Native(), date1.Native()))
	case "Y":
		return types.NewXNumberFromInt(date2.Native().Year() - date1.Native().Year())
	}

	return types.NewXErrorf("unknown unit: %s, must be one of s, m, h, D, W, M, Y", unit)
}

// DateTimeAdd calculates the date value arrived at by adding `offset` number of `unit` to the `datetime`
//
// Valid durations are "Y" for years, "M" for months, "W" for weeks, "D" for days, "h" for hour,
// "m" for minutes, "s" for seconds
//
//   @(datetime_add("2017-01-15", 5, "D")) -> 2017-01-20T00:00:00.000000-05:00
//   @(datetime_add("2017-01-15 10:45", 30, "m")) -> 2017-01-15T11:15:00.000000-05:00
//
// @function datetime_add(datetime, offset, unit)
func DateTimeAdd(env envs.Environment, args ...types.XValue) types.XValue {
	if len(args) != 3 {
		return types.NewXErrorf("takes exactly three arguments, received %d", len(args))
	}

	date, xerr := types.ToXDateTime(env, args[0])
	if xerr != nil {
		return xerr
	}

	duration, xerr := types.ToInteger(env, args[1])
	if xerr != nil {
		return xerr
	}

	unit, xerr := types.ToXText(env, args[2])
	if xerr != nil {
		return xerr
	}

	switch unit.Native() {
	case "s":
		return types.NewXDateTime(date.Native().Add(time.Duration(duration) * time.Second))
	case "m":
		return types.NewXDateTime(date.Native().Add(time.Duration(duration) * time.Minute))
	case "h":
		return types.NewXDateTime(date.Native().Add(time.Duration(duration) * time.Hour))
	case "D":
		return types.NewXDateTime(date.Native().AddDate(0, 0, duration))
	case "W":
		return types.NewXDateTime(date.Native().AddDate(0, 0, duration*7))
	case "M":
		return types.NewXDateTime(date.Native().AddDate(0, duration, 0))
	case "Y":
		return types.NewXDateTime(date.Native().AddDate(duration, 0, 0))
	}

	return types.NewXErrorf("unknown unit: %s, must be one of s, m, h, D, W, M, Y", unit)
}

// ReplaceTime returns a new datetime with the time part replaced by the `time`.
//
//   @(replace_time(now(), "10:30")) -> 2018-04-11T10:30:00.000000-05:00
//   @(replace_time("2017-01-15", "10:30")) -> 2017-01-15T10:30:00.000000-05:00
//   @(replace_time("foo", "10:30")) -> ERROR
//
// @function replace_time(datetime)
func ReplaceTime(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
	date, xerr := types.ToXDateTime(env, arg1)
	if xerr != nil {
		return xerr
	}
	t, xerr := types.ToXTime(env, arg2)
	if xerr != nil {
		return xerr
	}

	return date.ReplaceTime(t)
}

// TZ returns the name of the timezone of `date`.
//
// If no timezone information is present in the date, then the current timezone will be returned.
//
//   @(tz("2017-01-15T02:15:18.123456Z")) -> UTC
//   @(tz("2017-01-15 02:15:18PM")) -> America/Guayaquil
//   @(tz("2017-01-15")) -> America/Guayaquil
//   @(tz("foo")) -> ERROR
//
// @function tz(date)
func TZ(env envs.Environment, date types.XDateTime) types.XValue {
	return types.NewXText(date.Native().Location().String())
}

// TZOffset returns the offset of the timezone of `date`.
//
// The offset is returned in the format `[+/-]HH:MM`. If no timezone information is present in the date,
// then the current timezone offset will be returned.
//
//   @(tz_offset("2017-01-15T02:15:18.123456Z")) -> +0000
//   @(tz_offset("2017-01-15 02:15:18PM")) -> -0500
//   @(tz_offset("2017-01-15")) -> -0500
//   @(tz_offset("foo")) -> ERROR
//
// @function tz_offset(date)
func TZOffset(env envs.Environment, date types.XDateTime) types.XValue {
	// this looks like we are returning a set offset, but this is how go describes formats
	return types.NewXText(date.Native().Format("-0700"))
}

// Epoch converts `date` to a UNIX epoch time.
//
// The returned number can contain fractional seconds.
//
//   @(epoch("2017-06-12T16:56:59.000000Z")) -> 1497286619
//   @(epoch("2017-06-12T18:56:59.000000+02:00")) -> 1497286619
//   @(epoch("2017-06-12T16:56:59.123456Z")) -> 1497286619.123456
//   @(round_down(epoch("2017-06-12T16:56:59.123456Z"))) -> 1497286619
//
// @function epoch(date)
func Epoch(env envs.Environment, date types.XDateTime) types.XValue {
	nanos := decimal.New(date.Native().UnixNano(), 0)
	return types.NewXNumber(nanos.Div(nanosPerSecond))
}

// Now returns the current date and time in the current timezone.
//
//   @(now()) -> 2018-04-11T13:24:30.123456-05:00
//
// @function now()
func Now(env envs.Environment) types.XValue {
	return types.NewXDateTime(env.Now())
}

//------------------------------------------------------------------------------------------
// Date Functions
//------------------------------------------------------------------------------------------

// DateFromParts creates a date from `year`, `month` and `day`.
//
//   @(date_from_parts(2017, 1, 15)) -> 2017-01-15
//   @(date_from_parts(2017, 2, 31)) -> 2017-03-03
//   @(date_from_parts(2017, 13, 15)) -> ERROR
//
// @function date_from_parts(year, month, day)
func DateFromParts(env envs.Environment, year, month, day int) types.XValue {
	if month < 1 || month > 12 {
		return types.NewXErrorf("invalid value for month, must be 1-12")
	}

	return types.NewXDate(dates.NewDate(year, month, day))
}

// Weekday returns the day of the week for `date`.
//
// The week is considered to start on Sunday so a Sunday returns 0, a Monday returns 1 etc.
//
//   @(weekday("2017-01-15")) -> 0
//   @(weekday("foo")) -> ERROR
//
// @function weekday(date)
func Weekday(env envs.Environment, date types.XDate) types.XValue {
	return types.NewXNumberFromInt(int(date.Native().Weekday()))
}

// WeekNumber returns the week number (1-54) of `date`.
//
// The week is considered to start on Sunday and week containing Jan 1st is week number 1.
//
//   @(week_number("2019-01-01")) -> 1
//   @(week_number("2019-07-23T16:56:59.000000Z")) -> 30
//   @(week_number("xx")) -> ERROR
//
// @function week_number(date)
func WeekNumber(env envs.Environment, date types.XDate) types.XValue {
	return types.NewXNumberFromInt(date.Native().WeekNum())
}

// Today returns the current date in the environment timezone.
//
//   @(today()) -> 2018-04-11
//
// @function today()
func Today(env envs.Environment) types.XValue {
	return types.NewXDate(dates.ExtractDate(env.Now()))
}

//------------------------------------------------------------------------------------------
// Time Functions
//------------------------------------------------------------------------------------------

// ParseTime parses `text` into a time using the given `format`.
//
// The format string can consist of the following characters. The characters
// ' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.
//
// * `h`         - hour of the day 1-12
// * `hh`        - hour of the day, zero padded 01-12
// * `t`         - twenty four hour of the day 1-23
// * `tt`        - twenty four hour of the day, zero padded 01-23
// * `m`         - minute 0-59
// * `mm`        - minute, zero padded 00-59
// * `s`         - second 0-59
// * `ss`        - second, zero padded 00-59
// * `fff`       - milliseconds
// * `ffffff`    - microseconds
// * `fffffffff` - nanoseconds
// * `aa`        - am or pm
// * `AA`        - AM or PM
//
// Note that fractional seconds will be parsed even without an explicit format identifier.
// You should only specify fractional seconds when you want to assert the number of places
// in the input format.
//
// parse_time will return an error if it is unable to convert the text to a time.
//
//   @(parse_time("15:28", "tt:mm")) -> 15:28:00.000000
//   @(parse_time("2:40 pm", "h:mm aa")) -> 14:40:00.000000
//   @(parse_time("NOT TIME", "tt:mm")) -> ERROR
//
// @function parse_time(text, format)
func ParseTime(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
	str, xerr := types.ToXText(env, arg1)
	if xerr != nil {
		return xerr
	}

	layout, xerr := types.ToXText(env, arg2)
	if xerr != nil {
		return xerr
	}

	// finally try to parse the time
	parsed, err := dates.ParseTimeOfDay(layout.Native(), str.Native())
	if err != nil {
		return types.NewXError(err)
	}

	return types.NewXTime(parsed)
}

// TimeFromParts creates a time from `hour`, `minute` and `second`
//
//   @(time_from_parts(14, 40, 15)) -> 14:40:15.000000
//   @(time_from_parts(8, 10, 0)) -> 08:10:00.000000
//   @(time_from_parts(25, 0, 0)) -> ERROR
//
// @function time_from_parts(hour, minute, second)
func TimeFromParts(env envs.Environment, hour, minute, second int) types.XValue {
	if hour < 0 || hour > 23 {
		return types.NewXErrorf("invalid value for hour, must be 0-23")
	}
	if minute < 0 || minute > 59 {
		return types.NewXErrorf("invalid value for minute, must be 0-59")
	}
	if second < 0 || second > 59 {
		return types.NewXErrorf("invalid value for second, must be 0-59")
	}

	return types.NewXTime(dates.NewTimeOfDay(hour, minute, second, 0))
}

//------------------------------------------------------------------------------------------
// Array Functions
//------------------------------------------------------------------------------------------

// Join joins the given `array` of strings with `separator` to make text.
//
//   @(join(array("a", "b", "c"), "|")) -> a|b|c
//   @(join(split("a.b.c", "."), " ")) -> a b c
//
// @function join(array, separator)
func Join(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
	array, xerr := types.ToXArray(env, arg1)
	if xerr != nil {
		return xerr
	}

	separator, xerr := types.ToXText(env, arg2)
	if xerr != nil {
		return xerr
	}

	var output bytes.Buffer
	for i := 0; i < array.Count(); i++ {
		if i > 0 {
			output.WriteString(separator.Native())
		}
		itemAsStr, xerr := types.ToXText(env, array.Get(i))
		if xerr != nil {
			return xerr
		}

		output.WriteString(itemAsStr.Native())
	}

	return types.NewXText(output.String())
}

// Reverse returns a new array with the values of `array` reversed.
//
//   @(reverse(array(3, 1, 2))) -> [2, 1, 3]
//   @(reverse(array("C", "A", "B"))) -> [B, A, C]
//
// @function reverse(array)
func Reverse(env envs.Environment, array *types.XArray) types.XValue {
	reversed := make([]types.XValue, array.Count())
	for i := 0; i < array.Count(); i++ {
		reversed[array.Count()-(i+1)] = array.Get(i)
	}
	return types.NewXArray(reversed...)
}

// Sort returns a new array with the values of `array` sorted.
//
//   @(sort(array(3, 1, 2))) -> [1, 2, 3]
//   @(sort(array("C", "A", "B"))) -> [A, B, C]
//
// @function sort(array)
func Sort(env envs.Environment, array *types.XArray) types.XValue {
	sorted := make([]types.XValue, array.Count())
	for i := 0; i < array.Count(); i++ {
		val := array.Get(i)

		_, isComparable := val.(types.XComparable)
		if !isComparable {
			return types.NewXErrorf("%s isn't a comparable type", types.Describe(val))
		}

		sorted[i] = val
	}

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].(types.XComparable).Compare(sorted[j]) < 0
	})

	return types.NewXArray(sorted...)
}

// Sum sums the items in the given `array`.
//
//   @(sum(array(1, 2, "3"))) -> 6
//
// @function sum(array)
func Sum(env envs.Environment, array *types.XArray) types.XValue {
	total := decimal.Zero
	for i := 0; i < array.Count(); i++ {
		itemAsNum, xerr := types.ToXNumber(env, array.Get(i))
		if xerr != nil {
			return xerr
		}

		total = total.Add(itemAsNum.Native())
	}

	return types.NewXNumber(total)
}

// Unique returns the unique values in `array`.
//
//   @(unique(array(1, 3, 2, 3))) -> [1, 3, 2]
//   @(unique(array("hi", "there", "hi"))) -> [hi, there]
//
// @function unique(array)
func Unique(env envs.Environment, array *types.XArray) types.XValue {
	unique := make([]types.XValue, 0, array.Count())
	for i := 0; i < array.Count(); i++ {
		val := array.Get(i)

		seen := false
		for j := 0; j < len(unique); j++ {
			if (val == nil && unique[j] == nil) || types.Equals(val, unique[j]) {
				seen = true
				break
			}
		}

		if !seen {
			unique = append(unique, val)
		}
	}

	return types.NewXArray(unique...)
}

// Concat returns the result of concatenating two arrays.
//
//   @(concat(array("a", "b"), array("c", "d"))) -> [a, b, c, d]
//   @(unique(concat(array(1, 2, 3), array(3, 4)))) -> [1, 2, 3, 4]
//
// @function concat(array1, array2)
func Concat(env envs.Environment, array1 *types.XArray, array2 *types.XArray) types.XValue {
	both := make([]types.XValue, 0, array1.Count()+array2.Count())

	for i := 0; i < array1.Count(); i++ {
		both = append(both, array1.Get(i))
	}
	for i := 0; i < array2.Count(); i++ {
		both = append(both, array2.Get(i))
	}

	return types.NewXArray(both...)
}

//------------------------------------------------------------------------------------------
// Encoded Text Functions
//------------------------------------------------------------------------------------------

// URNParts parses a URN into its different parts
//
//   @(urn_parts("tel:+593979012345")) -> {display: , path: +593979012345, scheme: tel}
//   @(urn_parts("twitterid:3263621177#bobby")) -> {display: bobby, path: 3263621177, scheme: twitterid}
//   @(urn_parts("not a urn")) -> ERROR
//
// @function urn_parts(urn)
func URNParts(env envs.Environment, urn types.XText) types.XValue {
	u, err := urns.Parse(urn.Native())
	if err != nil {
		return types.NewXErrorf("%s is not a valid URN: %s", urn.Native(), err)
	}

	scheme, path, _, display := u.ToParts()

	return types.NewXObject(map[string]types.XValue{
		"scheme":  types.NewXText(scheme),
		"path":    types.NewXText(path),
		"display": types.NewXText(display),
	})
}

// AttachmentParts parses an attachment into its different parts
//
//   @(attachment_parts("image/jpeg:https://example.com/test.jpg")) -> {content_type: image/jpeg, url: https://example.com/test.jpg}
//
// @function attachment_parts(attachment)
func AttachmentParts(env envs.Environment, attachment types.XText) types.XValue {
	a := utils.Attachment(attachment.Native())
	contentType, url := a.ToParts()

	return types.NewXObject(map[string]types.XValue{
		"content_type": types.NewXText(contentType),
		"url":          types.NewXText(url),
	})
}

//------------------------------------------------------------------------------------------
// JSON Functions
//------------------------------------------------------------------------------------------

// ParseJSON tries to parse `text` as JSON.
//
// If the given `text` is not valid JSON, then an error is returned
//
//   @(parse_json("{\"foo\": \"bar\"}").foo) -> bar
//   @(parse_json("[1,2,3,4]")[2]) -> 3
//   @(parse_json("invalid json")) -> ERROR
//
// @function parse_json(text)
func ParseJSON(env envs.Environment, text types.XText) types.XValue {
	return types.JSONToXValue([]byte(text.Native()))
}

// JSON returns the JSON representation of `value`.
//
//   @(json("string")) -> "string"
//   @(json(10)) -> 10
//   @(json(null)) -> null
//   @(json(contact.uuid)) -> "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
//
// @function json(value)
func JSON(env envs.Environment, value types.XValue) types.XValue {
	asJSON, xerr := types.ToXJSON(value)
	if xerr != nil {
		return xerr
	}
	return asJSON
}

//----------------------------------------------------------------------------------------
// Formatting Functions
//----------------------------------------------------------------------------------------

// Format formats `value` according to its type.
//
//   @(format(1234.5670)) -> 1,234.567
//   @(format(now())) -> 11-04-2018 13:24
//   @(format(today())) -> 11-04-2018
//
// @function format(value)
func Format(env envs.Environment, value types.XValue) types.XValue {
	if !utils.IsNil(value) {
		return types.NewXText(value.Format(env))
	}
	return types.XTextEmpty
}

// FormatDate formats `date` as text according to the given `format`.
//
// If `format` is not specified then the environment's default format is used. The format
// string can consist of the following characters. The characters ' ', ':', ',', 'T', '-'
// and '_' are ignored. Any other character is an error.
//
// * `YY`        - last two digits of year 0-99
// * `YYYY`      - four digits of year 0000-9999
// * `M`         - month 1-12
// * `MM`        - month, zero padded 01-12
// * `MMM`       - month Jan-Dec (localized)
// * `MMMM`      - month January-December (localized)
// * `D`         - day of month, 1-31
// * `DD`        - day of month, zero padded 01-31
// * `EEE`       - day of week Mon-Sun (localized)
// * `EEEE`      - day of week Monday-Sunday (localized)
//
//   @(format_date("1979-07-18T15:00:00.000000Z")) -> 18-07-1979
//   @(format_date("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) -> 1979-07-18
//   @(format_date("2010-05-10T19:50:00.000000Z", "YYYY M DD")) -> 2010 5 10
//   @(format_date("1979-07-18T15:00:00.000000Z", "YYYY")) -> 1979
//   @(format_date("1979-07-18T15:00:00.000000Z", "M")) -> 7
//   @(format_date("NOT DATE", "YYYY-MM-DD")) -> ERROR
//
// @function format_date(date, [,format])
func FormatDate(env envs.Environment, args ...types.XValue) types.XValue {
	date, xerr := types.ToXDate(env, args[0])
	if xerr != nil {
		return xerr
	}

	if len(args) >= 2 {
		layout, xerr := types.ToXText(env, args[1])
		if xerr != nil {
			return xerr
		}

		formatted, err := date.FormatCustom(env, layout.Native())
		if err != nil {
			return types.NewXError(err)
		}

		return types.NewXText(formatted)
	}

	return types.NewXText(date.Format(env))
}

// FormatDateTime formats `datetime` as text according to the given `format`.
//
// If `format` is not specified then the environment's default format is used. The format
// string can consist of the following characters. The characters ' ', ':', ',', 'T', '-'
// and '_' are ignored. Any other character is an error.
//
// * `YY`        - last two digits of year 0-99
// * `YYYY`      - four digits of year 0000-9999
// * `M`         - month 1-12
// * `MM`        - month, zero padded 01-12
// * `MMM`       - month Jan-Dec (localized)
// * `MMMM`      - month January-December (localized)
// * `D`         - day of month, 1-31
// * `DD`        - day of month, zero padded 01-31
// * `EEE`       - day of week Mon-Sun (localized)
// * `EEEE`      - day of week Monday-Sunday (localized)
// * `h`         - hour of the day 1-12
// * `hh`        - hour of the day, zero padded 01-12
// * `t`         - twenty four hour of the day 0-23
// * `tt`        - twenty four hour of the day, zero padded 00-23
// * `m`         - minute 0-59
// * `mm`        - minute, zero padded 00-59
// * `s`         - second 0-59
// * `ss`        - second, zero padded 00-59
// * `fff`       - milliseconds
// * `ffffff`    - microseconds
// * `fffffffff` - nanoseconds
// * `aa`        - am or pm (localized)
// * `AA`        - AM or PM (localized)
// * `Z`         - hour and minute offset from UTC, or Z for UTC
// * `ZZZ`       - hour and minute offset from UTC
//
// Timezone should be a location name as specified in the IANA Time Zone database, such
// as "America/Guayaquil" or "America/Los_Angeles". If not specified, the current timezone
// will be used. An error will be returned if the timezone is not recognized.
//
//   @(format_datetime("1979-07-18T15:00:00.000000Z")) -> 18-07-1979 10:00
//   @(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) -> 1979-07-18
//   @(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY M DD tt:mm")) -> 2010 5 10 14:50
//   @(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY-MM-DD hh:mm AA", "America/Los_Angeles")) -> 2010-05-10 12:50 PM
//   @(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY")) -> 1979
//   @(format_datetime("1979-07-18T15:00:00.000000Z", "M")) -> 7
//   @(format_datetime("NOT DATE", "YYYY-MM-DD")) -> ERROR
//
// @function format_datetime(datetime [,format [,timezone]])
func FormatDateTime(env envs.Environment, args ...types.XValue) types.XValue {
	date, xerr := types.ToXDateTime(env, args[0])
	if xerr != nil {
		return xerr
	}

	var format types.XText
	if len(args) >= 2 {
		format, xerr = types.ToXText(env, args[1])
		if xerr != nil {
			return xerr
		}
	} else {
		format = types.NewXText(fmt.Sprintf("%s %s", env.DateFormat().String(), env.TimeFormat().String()))
	}

	// grab our location
	var err error
	location := env.Timezone()
	if len(args) == 3 {
		arg3, xerr := types.ToXText(env, args[2])
		if xerr != nil {
			return xerr
		}

		location, err = time.LoadLocation(arg3.Native())
		if err != nil {
			return types.NewXError(err)
		}
	}

	formatted, err := date.FormatCustom(env, format.Native(), location)
	if err != nil {
		return types.NewXError(err)
	}

	return types.NewXText(formatted)
}

// FormatTime formats `time` as text according to the given `format`.
//
// If `format` is not specified then the environment's default format is used. The format
// string can consist of the following characters. The characters ' ', ':', ',', 'T', '-'
// and '_' are ignored. Any other character is an error.
//
// * `h`         - hour of the day 1-12
// * `hh`        - hour of the day, zero padded 01-12
// * `t`         - twenty four hour of the day 0-23
// * `tt`        - twenty four hour of the day, zero padded 00-23
// * `m`         - minute 0-59
// * `mm`        - minute, zero padded 00-59
// * `s`         - second 0-59
// * `ss`        - second, zero padded 00-59
// * `fff`       - milliseconds
// * `ffffff`    - microseconds
// * `fffffffff` - nanoseconds
// * `aa`        - am or pm (localized)
// * `AA`        - AM or PM (localized)
//
//   @(format_time("14:50:30.000000")) -> 14:50
//   @(format_time("14:50:30.000000", "h:mm aa")) -> 2:50 pm
//   @(format_time("15:00:27.000000", "s")) -> 27
//   @(format_time("NOT TIME", "hh:mm")) -> ERROR
//
// @function format_time(time [,format])
func FormatTime(env envs.Environment, args ...types.XValue) types.XValue {
	t, xerr := types.ToXTime(env, args[0])
	if xerr != nil {
		return xerr
	}

	if len(args) >= 2 {
		layout, xerr := types.ToXText(env, args[1])
		if xerr != nil {
			return xerr
		}

		formatted, err := t.FormatCustom(env, layout.Native())
		if err != nil {
			return types.NewXError(err)
		}

		return types.NewXText(formatted)
	}

	return types.NewXText(t.Format(env))
}

// FormatNumber formats `number` to the given number of decimal `places`.
//
// An optional third argument `humanize` can be false to disable the use of thousand separators.
//
//   @(format_number(1234)) -> 1,234
//   @(format_number(1234.5670)) -> 1,234.567
//   @(format_number(1234.5670, 2, true)) -> 1,234.57
//   @(format_number(1234.5678, 0, false)) -> 1235
//   @(format_number("foo", 2, false)) -> ERROR
//
// @function format_number(number, places [, humanize])
func FormatNumber(env envs.Environment, args ...types.XValue) types.XValue {
	num, err := types.ToXNumber(env, args[0])
	if err != nil {
		return err
	}

	places := -1
	if len(args) > 1 {
		if places, err = types.ToInteger(env, args[1]); err != nil {
			return err
		}
		if places < 0 || places > 9 {
			return types.NewXErrorf("must take 0-9 number of places, got %d", args[1])
		}
	}

	human := types.XBooleanTrue
	if len(args) > 2 {
		if human, err = types.ToXBoolean(args[2]); err != nil {
			return err
		}
	}

	return types.NewXText(num.FormatCustom(env.NumberFormat(), places, human.Native()))
}

// FormatLocation formats the given `location` as its name.
//
//   @(format_location("Rwanda")) -> Rwanda
//   @(format_location("Rwanda > Kigali")) -> Kigali
//
// @function format_location(location)
func FormatLocation(env envs.Environment, path types.XText) types.XValue {
	parts := strings.Split(path.Native(), ">")
	return types.NewXText(strings.TrimSpace(parts[len(parts)-1]))
}

// FormatURN formats `urn` into human friendly text.
//
//   @(format_urn("tel:+250781234567")) -> 0781 234 567
//   @(format_urn("twitter:134252511151#billy_bob")) -> billy_bob
//   @(format_urn(contact.urn)) -> (202) 456-1111
//   @(format_urn(urns.tel)) -> (202) 456-1111
//   @(format_urn(urns.mailto)) -> foo@bar.com
//   @(format_urn("NOT URN")) -> ERROR
//
// @function format_urn(urn)
func FormatURN(env envs.Environment, arg types.XText) types.XValue {
	urn, err := urns.Parse(arg.Native())
	if err != nil {
		return types.NewXErrorf("%s is not a valid URN: %s", arg.Native(), err)
	}

	return types.NewXText(urn.Format())
}

//------------------------------------------------------------------------------------------
// Utility Functions
//------------------------------------------------------------------------------------------

// IsError returns whether `value` is an error
//
//   @(is_error(datetime("foo"))) -> true
//   @(is_error(run.not.existing)) -> true
//   @(is_error("hello")) -> false
//
// @function is_error(value)
func IsError(env envs.Environment, value types.XValue) types.XValue {
	return types.NewXBoolean(types.IsXError(value))
}

// Count returns the number of items in the given array or properties on an object.
//
// It will return an error if it is passed an item which isn't countable.
//
//   @(count(contact.fields)) -> 6
//   @(count(array())) -> 0
//   @(count(array("a", "b", "c"))) -> 3
//   @(count(1234)) -> ERROR
//
// @function count(value)
func Count(env envs.Environment, value types.XValue) types.XValue {
	// a nil has count of zero
	if utils.IsNil(value) {
		return types.XNumberZero
	}

	// argument must be a countable value
	countable, isCountable := value.(types.XCountable)
	if isCountable {
		return types.NewXNumberFromInt(countable.Count())
	}

	return types.NewXErrorf("value isn't countable")
}

// Default returns `value` if is not empty or an error, otherwise it returns `default`.
//
//   @(default(undeclared.var, "default_value")) -> default_value
//   @(default("10", "20")) -> 10
//   @(default("", "value")) -> value
//   @(default("  ", "value")) -> \x20\x20
//   @(default(datetime("invalid-date"), "today")) -> today
//   @(default(format_urn("invalid-urn"), "ok")) -> ok
//
// @function default(value, default)
func Default(env envs.Environment, value types.XValue, def types.XValue) types.XValue {
	asText, xerr := types.ToXText(env, value)
	if xerr != nil {
		return def
	}

	if len(asText.Native()) == 0 {
		return def
	}

	return value
}

// Extract takes an object and extracts the named property.
//
//   @(extract(contact, "name")) -> Ryan Lewis
//   @(extract(contact.groups[0], "name")) -> Testers
//
// @function extract(object, properties)
func Extract(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
	object, xerr := types.ToXObject(env, arg1)
	if xerr != nil {
		return xerr
	}

	property, xerr := types.ToXText(env, arg2)
	if xerr != nil {
		return xerr
	}

	value, _ := object.Get(property.Native())
	return value
}

// ExtractObject takes an object and returns a new object by extracting only the named properties.
//
//   @(extract_object(contact.groups[0], "name")) -> {name: Testers}
//
// @function extract_object(object, properties...)
func ExtractObject(env envs.Environment, args ...types.XValue) types.XValue {
	object, xerr := types.ToXObject(env, args[0])
	if xerr != nil {
		return xerr
	}

	properties := make([]string, 0, len(args)-1)
	for _, arg := range args[1:] {
		asText, xerr := types.ToXText(env, arg)
		if xerr != nil {
			return xerr
		}
		properties = append(properties, asText.Native())
	}

	result := make(map[string]types.XValue, len(properties))
	for _, prop := range properties {
		value, _ := object.Get(prop)
		result[prop] = value
	}

	return types.NewXObject(result)
}

// ForEach creates a new array by applying `func` to each value in `values`.
//
// If the given function takes more than one argument, you can pass additional arguments after the function.
//
//   @(foreach(array("a", "b", "c"), upper)) -> [A, B, C]
//   @(foreach(array("a", "b", "c"), (x) => x & "1")) -> [a1, b1, c1]
//   @(foreach(array("a", "b", "c"), (x) => object("v", x))) -> [{v: a}, {v: b}, {v: c}]
//   @(foreach(array("the man", "fox", "jumped up"), word, 0)) -> [the, fox, jumped]
//
// @function foreach(values, func, [args...])
func ForEach(env envs.Environment, args ...types.XValue) types.XValue {
	array, xerr := types.ToXArray(env, args[0])
	if xerr != nil {
		return xerr
	}

	function, isFunction := args[1].(*types.XFunction)
	if !isFunction {
		return types.NewXErrorf("requires an function as its second argument")
	}

	otherArgs := args[2:]

	result := make([]types.XValue, array.Count())

	for i := 0; i < array.Count(); i++ {
		oldItem := array.Get(i)
		funcArgs := append([]types.XValue{oldItem}, otherArgs...)

		newItem := function.Call(env, funcArgs)
		if types.IsXError(newItem) {
			return newItem
		}
		result[i] = newItem
	}

	return types.NewXArray(result...)
}

// ForEachValue creates a new object by applying `func` to each property value of `object`.
//
// If the given function takes more than one argument, you can pass additional arguments after the function.
//
//   @(foreach_value(object("a", "x", "b", "y"), upper)) -> {a: X, b: Y}
//   @(foreach_value(object("a", "hi there", "b", "good bye"), word, 1)) -> {a: there, b: bye}
//
// @function foreach_value(object, func, [args...])
func ForEachValue(env envs.Environment, args ...types.XValue) types.XValue {
	object, xerr := types.ToXObject(env, args[0])
	if xerr != nil {
		return xerr
	}

	function, isFunction := args[1].(*types.XFunction)
	if !isFunction {
		return types.NewXErrorf("requires an function as its second argument")
	}

	otherArgs := args[2:]

	props := object.Properties()
	result := make(map[string]types.XValue, len(props))

	for _, prop := range props {
		oldItem, _ := object.Get(prop)
		funcArgs := append([]types.XValue{oldItem}, otherArgs...)

		newItem := function.Call(env, funcArgs)
		if types.IsXError(newItem) {
			return newItem
		}
		result[prop] = newItem
	}

	return types.NewXObject(result)
}

// LegacyAdd simulates our old + operator, which operated differently based on whether
// one of the parameters was a date or not. If one is a date, then the other side is
// expected to be an integer with a number of days to add to the date, otherwise a normal
// decimal addition is attempted.
func LegacyAdd(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {

	// try to parse dates and decimals
	date1, date1Err := types.ToXDateTime(env, arg1)
	date2, date2Err := types.ToXDateTime(env, arg2)

	dec1, dec1Err := types.ToXNumber(env, arg1)
	dec2, dec2Err := types.ToXNumber(env, arg2)

	// if they are both dates, that's an error
	if date1Err == nil && date2Err == nil {
		return types.NewXErrorf("cannot operate on two dates")
	}

	// date and int, do a day addition
	if date1Err == nil && dec2Err == nil {
		if dec2.Native().IntPart() < math.MinInt32 || dec2.Native().IntPart() > math.MaxInt32 {
			return types.NewXErrorf("cannot operate on integers greater than 32 bit")
		}
		return types.NewXDateTime(date1.Native().AddDate(0, 0, int(dec2.Native().IntPart())))
	}

	// int and date, do a day addition
	if date2Err == nil && dec1Err == nil {
		if dec1.Native().IntPart() < math.MinInt32 || dec1.Native().IntPart() > math.MaxInt32 {
			return types.NewXErrorf("cannot operate on integers greater than 32 bit")
		}
		return types.NewXDateTime(date2.Native().AddDate(0, 0, int(dec1.Native().IntPart())))
	}

	// one of these doesn't look like a valid decimal either, bail
	if dec1Err != nil {
		return types.NewXError(dec1Err)
	}

	if dec2Err != nil {
		return types.NewXError(dec2Err)
	}

	// normal decimal addition
	return types.NewXNumber(dec1.Native().Add(dec2.Native()))
}

// ReadChars converts `text` into something that can be read by IVR systems.
//
// ReadChars will split the numbers such as they are easier to understand. This includes
// splitting in 3s or 4s if appropriate.
//
//   @(read_chars("1234")) -> 1 2 3 4
//   @(read_chars("abc")) -> a b c
//   @(read_chars("abcdef")) -> a b c , d e f
//
// @function read_chars(text)
func ReadChars(env envs.Environment, val types.XText) types.XValue {
	var output bytes.Buffer

	// remove any leading +
	val = types.NewXText(strings.TrimLeft(val.Native(), "+"))

	length := val.Length()

	// groups of three
	if length%3 == 0 {
		// groups of 3
		for i := 0; i < length; i += 3 {
			if i > 0 {
				output.WriteString(" , ")
			}
			output.WriteString(strings.Join(strings.Split(val.Native()[i:i+3], ""), " "))
		}
		return types.NewXText(output.String())
	}

	// groups of four
	if length%4 == 0 {
		for i := 0; i < length; i += 4 {
			if i > 0 {
				output.WriteString(" , ")
			}
			output.WriteString(strings.Join(strings.Split(val.Native()[i:i+4], ""), " "))
		}
		return types.NewXText(output.String())
	}

	// default, just do one at a time
	for i, c := range val.Native() {
		if i > 0 {
			output.WriteString(" , ")
		}
		output.WriteRune(c)
	}

	return types.NewXText(output.String())
}
