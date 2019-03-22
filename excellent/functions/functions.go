package functions

import (
	"bytes"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	humanize "github.com/dustin/go-humanize"
	"github.com/shopspring/decimal"
)

var nanosPerSecond = decimal.RequireFromString("1000000000")
var nonPrintableRegex = regexp.MustCompile(`[\p{Cc}\p{C}]`)

// XFunction defines the interface that Excellent functions must implement
type XFunction func(env utils.Environment, args ...types.XValue) types.XValue

// RegisterXFunction registers a new function in Excellent
func RegisterXFunction(name string, function XFunction) {
	XFUNCTIONS[name] = function
}

// XFUNCTIONS is our map of functions available in Excellent which aren't tests
var XFUNCTIONS = map[string]XFunction{
	// type conversion
	"text":     OneArgFunction(Text),
	"boolean":  OneArgFunction(Boolean),
	"number":   OneArgFunction(Number),
	"date":     OneArgFunction(Date),
	"datetime": OneArgFunction(DateTime),
	"time":     OneArgFunction(Time),
	"array":    Array,

	// text functions
	"char":              OneNumberFunction(Char),
	"code":              OneTextFunction(Code),
	"split":             TwoTextFunction(Split),
	"join":              TwoArgFunction(Join),
	"title":             OneTextFunction(Title),
	"word":              InitialTextFunction(1, 2, Word),
	"remove_first_word": OneTextFunction(RemoveFirstWord),
	"word_count":        InitialTextFunction(0, 1, WordCount),
	"word_slice":        InitialTextFunction(1, 3, WordSlice),
	"field":             InitialTextFunction(2, 2, Field),
	"clean":             OneTextFunction(Clean),
	"left":              TextAndIntegerFunction(Left),
	"lower":             OneTextFunction(Lower),
	"right":             TextAndIntegerFunction(Right),
	"regex_match":       InitialTextFunction(1, 2, RegexMatch),
	"text_compare":      TwoTextFunction(TextCompare),
	"repeat":            TextAndIntegerFunction(Repeat),
	"replace":           ThreeTextFunction(Replace),
	"upper":             OneTextFunction(Upper),
	"percent":           OneNumberFunction(Percent),
	"url_encode":        OneTextFunction(URLEncode),

	// bool functions
	"and": ArgCountCheck(1, -1, And),
	"if":  ThreeArgFunction(If),
	"or":  ArgCountCheck(1, -1, Or),

	// number functions
	"round":        OneNumberAndOptionalIntegerFunction(Round, 0),
	"round_up":     OneNumberAndOptionalIntegerFunction(RoundUp, 0),
	"round_down":   OneNumberAndOptionalIntegerFunction(RoundDown, 0),
	"max":          ArgCountCheck(1, -1, Max),
	"min":          ArgCountCheck(1, -1, Min),
	"mean":         ArgCountCheck(1, -1, Mean),
	"mod":          TwoNumberFunction(Mod),
	"rand":         NoArgFunction(Rand),
	"rand_between": TwoNumberFunction(RandBetween),
	"abs":          OneNumberFunction(Abs),

	// datetime functions
	"parse_datetime":      ArgCountCheck(2, 3, ParseDateTime),
	"datetime_from_epoch": OneNumberFunction(DateTimeFromEpoch),
	"datetime_diff":       ThreeArgFunction(DateTimeDiff),
	"datetime_add":        DateTimeAdd,
	"replace_time":        ArgCountCheck(2, 2, ReplaceTime),
	"tz":                  OneDateTimeFunction(TZ),
	"tz_offset":           OneDateTimeFunction(TZOffset),
	"now":                 NoArgFunction(Now),
	"epoch":               OneDateTimeFunction(Epoch),

	// date functions
	"date_from_parts": ThreeIntegerFunction(DateFromParts),
	"weekday":         OneDateFunction(Weekday),
	"today":           NoArgFunction(Today),

	// time functions
	"parse_time":      ArgCountCheck(2, 2, ParseTime),
	"time_from_parts": ThreeIntegerFunction(TimeFromParts),

	// json functions
	"json":       OneArgFunction(JSON),
	"parse_json": OneTextFunction(ParseJSON),

	// formatting functions
	"format_date":     ArgCountCheck(1, 2, FormatDate),
	"format_datetime": ArgCountCheck(1, 3, FormatDateTime),
	"format_time":     ArgCountCheck(1, 2, FormatTime),
	"format_location": OneTextFunction(FormatLocation),
	"format_number":   FormatNumber,
	"format_urn":      OneTextFunction(FormatURN),

	// utility functions
	"length":     OneArgFunction(Length),
	"default":    TwoArgFunction(Default),
	"legacy_add": TwoArgFunction(LegacyAdd),
	"read_chars": OneTextFunction(ReadChars),
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
func Text(env utils.Environment, value types.XValue) types.XValue {
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
func Boolean(env utils.Environment, value types.XValue) types.XValue {
	str, xerr := types.ToXBoolean(env, value)
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
func Number(env utils.Environment, value types.XValue) types.XValue {
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
//   @(date("2010 05 10")) -> 2010-05-10
//   @(date("NOT DATE")) -> ERROR
//
// @function date(value)
func Date(env utils.Environment, value types.XValue) types.XValue {
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
//   @(datetime("2010 05 10")) -> 2010-05-10T00:00:00.000000-05:00
//   @(datetime("NOT DATE")) -> ERROR
//
// @function datetime(value)
func DateTime(env utils.Environment, value types.XValue) types.XValue {
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
func Time(env utils.Environment, value types.XValue) types.XValue {
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
//   @(length(array())) -> 0
//   @(length(array("a", "b"))) -> 2
//
// @function array(values...)
func Array(env utils.Environment, values ...types.XValue) types.XValue {
	// check none of our args are errors
	for _, arg := range values {
		if types.IsXError(arg) {
			return arg
		}
	}

	return types.NewXArray(values...)
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
func And(env utils.Environment, values ...types.XValue) types.XValue {
	for _, arg := range values {
		asBool, xerr := types.ToXBoolean(env, arg)
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
func Or(env utils.Environment, values ...types.XValue) types.XValue {
	for _, arg := range values {
		asBool, xerr := types.ToXBoolean(env, arg)
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
func If(env utils.Environment, test types.XValue, value1 types.XValue, value2 types.XValue) types.XValue {
	asBool, err := types.ToXBoolean(env, test)
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
func Code(env utils.Environment, text types.XText) types.XValue {
	if text.Length() == 0 {
		return types.NewXErrorf("requires a string of at least one character")
	}

	r, _ := utf8.DecodeRuneInString(text.Native())
	return types.NewXNumberFromInt(int(r))
}

// Split splits `text` based on the given characters in `delimiters`.
//
// Empty values are removed from the returned list.
//
//   @(split("a b c", " ")) -> [a, b, c]
//   @(split("a", " ")) -> [a]
//   @(split("abc..d", ".")) -> [abc, d]
//   @(split("a.b.c.", ".")) -> [a, b, c]
//   @(split("a|b,c  d", " .|,")) -> [a, b, c, d]
//
// @function split(text, delimiters)
func Split(env utils.Environment, text types.XText, delimiters types.XText) types.XValue {
	splits := types.NewXArray()
	allSplits := utils.TokenizeStringByChars(text.Native(), delimiters.Native())
	for i := range allSplits {
		splits.Append(types.NewXText(allSplits[i]))
	}
	return splits
}

// Join joins the given `array` of strings with `separator` to make text.
//
//   @(join(array("a", "b", "c"), "|")) -> a|b|c
//   @(join(split("a.b.c", "."), " ")) -> a b c
//
// @function join(array, separator)
func Join(env utils.Environment, array types.XValue, separator types.XValue) types.XValue {
	indexable, isIndexable := array.(types.XIndexable)
	if !isIndexable {
		return types.NewXErrorf("requires an indexable as its first argument")
	}

	sep, xerr := types.ToXText(env, separator)
	if xerr != nil {
		return xerr
	}

	var output bytes.Buffer
	for i := 0; i < indexable.Length(); i++ {
		if i > 0 {
			output.WriteString(sep.Native())
		}
		itemAsStr, xerr := types.ToXText(env, indexable.Index(i))
		if xerr != nil {
			return xerr
		}

		output.WriteString(itemAsStr.Native())
	}

	return types.NewXText(output.String())
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
func Char(env utils.Environment, num types.XNumber) types.XValue {
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
func Title(env utils.Environment, text types.XText) types.XValue {
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
func Word(env utils.Environment, text types.XText, args ...types.XValue) types.XValue {
	index, xerr := types.ToInteger(env, args[0])
	if xerr != nil {
		return xerr
	}

	var words []string
	if len(args) == 2 && args[1] != nil {
		delimiters, xerr := types.ToXText(env, args[1])
		if xerr != nil {
			return xerr
		}
		words = utils.TokenizeStringByChars(text.Native(), delimiters.Native())
	} else {
		words = utils.TokenizeString(text.Native())
	}

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
func RemoveFirstWord(env utils.Environment, text types.XText) types.XValue {
	firstWordVal := Word(env, text, types.XNumberZero)
	firstWord, isText := firstWordVal.(types.XText)
	if !isText || firstWord == types.XTextEmpty {
		return types.XTextEmpty
	}

	firstWordStart := strings.Index(text.Native(), firstWord.Native())
	firstWordEnd := firstWordStart + firstWord.Length()

	remainder := text.Slice(firstWordEnd, text.Length())

	// remove any white space left at start
	return types.NewXText(strings.TrimLeft(remainder.Native(), " "))
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
func WordSlice(env utils.Environment, text types.XText, args ...types.XValue) types.XValue {
	start, xerr := types.ToInteger(env, args[0])
	if xerr != nil {
		return xerr
	}
	if start < 0 {
		return types.NewXErrorf("must start with a positive index")
	}

	end := -1
	if len(args) == 2 {
		if end, xerr = types.ToInteger(env, args[1]); xerr != nil {
			return xerr
		}
	}
	if end > 0 && end <= start {
		return types.NewXErrorf("must have a end which is greater than the start")
	}

	var words []string
	if len(args) == 3 && args[2] != nil {
		delimiters, xerr := types.ToXText(env, args[2])
		if xerr != nil {
			return xerr
		}
		words = utils.TokenizeStringByChars(text.Native(), delimiters.Native())
	} else {
		words = utils.TokenizeString(text.Native())
	}

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
func WordCount(env utils.Environment, text types.XText, args ...types.XValue) types.XValue {
	var words []string
	if len(args) == 1 && args[0] != nil {
		delimiters, xerr := types.ToXText(env, args[0])
		if xerr != nil {
			return xerr
		}
		words = utils.TokenizeStringByChars(text.Native(), delimiters.Native())
	} else {
		words = utils.TokenizeString(text.Native())
	}

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
func Field(env utils.Environment, text types.XText, args ...types.XValue) types.XValue {
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

// Clean strips any non-printable characters from `text`.
//
//   @(clean("😃 Hello \nwo\tr\rld")) -> 😃 Hello world
//   @(clean(123)) -> 123
//
// @function clean(text)
func Clean(env utils.Environment, text types.XText) types.XValue {
	return types.NewXText(nonPrintableRegex.ReplaceAllString(text.Native(), ""))
}

// Left returns the `count` left-most characters in `text`
//
//   @(left("hello", 2)) -> he
//   @(left("hello", 7)) -> hello
//   @(left("😀😃😄😁", 2)) -> 😀😃
//   @(left("hello", -1)) -> ERROR
//
// @function left(text, count)
func Left(env utils.Environment, text types.XText, count int) types.XValue {
	if count < 0 {
		return types.NewXErrorf("can't take a negative count")
	}

	// this weird construct does the right thing for multi-byte unicode
	var output bytes.Buffer
	i := 0
	for _, r := range text.Native() {
		if i >= count {
			break
		}
		output.WriteRune(r)
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
func Lower(env utils.Environment, text types.XText) types.XValue {
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
func RegexMatch(env utils.Environment, text types.XText, args ...types.XValue) types.XValue {
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

// Right returns the `count` right-most characters in `text`
//
//   @(right("hello", 2)) -> lo
//   @(right("hello", 7)) -> hello
//   @(right("😀😃😄😁", 2)) -> 😄😁
//   @(right("hello", -1)) -> ERROR
//
// @function right(text, count)
func Right(env utils.Environment, text types.XText, count int) types.XValue {
	if count < 0 {
		return types.NewXErrorf("can't take a negative count")
	}

	start := utf8.RuneCountInString(text.Native()) - count

	// this weird construct does the right thing for multi-byte unicode
	var output bytes.Buffer
	i := 0
	for _, r := range text.Native() {
		if i >= start {
			output.WriteRune(r)
		}
		i++
	}

	return types.NewXText(output.String())
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
func TextCompare(env utils.Environment, text1 types.XText, text2 types.XText) types.XValue {
	return types.NewXNumberFromInt(text1.Compare(text2))
}

// Repeat returns `text` repeated `count` number of times.
//
//   @(repeat("*", 8)) -> ********
//   @(repeat("*", "foo")) -> ERROR
//
// @function repeat(text, count)
func Repeat(env utils.Environment, text types.XText, count int) types.XValue {
	if count < 0 {
		return types.NewXErrorf("must be called with a positive integer, got %d", count)
	}

	var output bytes.Buffer
	for j := 0; j < count; j++ {
		output.WriteString(text.Native())
	}

	return types.NewXText(output.String())
}

// Replace replaces all occurrences of `needle` with `replacement` in `text`.
//
//   @(replace("foo bar", "foo", "zap")) -> zap bar
//   @(replace("foo bar", "baz", "zap")) -> foo bar
//
// @function replace(text, needle, replacement)
func Replace(env utils.Environment, text types.XText, needle types.XText, replacement types.XText) types.XValue {
	return types.NewXText(strings.Replace(text.Native(), needle.Native(), replacement.Native(), -1))
}

// Upper converts `text` to lowercase.
//
//   @(upper("Asdf")) -> ASDF
//   @(upper(123)) -> 123
//
// @function upper(text)
func Upper(env utils.Environment, text types.XText) types.XValue {
	return types.NewXText(strings.ToUpper(text.Native()))
}

// Percent formats `num` as a percentage.
//
//   @(percent(0.54234)) -> 54%
//   @(percent(1.2)) -> 120%
//   @(percent("foo")) -> ERROR
//
// @function percent(num)
func Percent(env utils.Environment, num types.XNumber) types.XValue {
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
func URLEncode(env utils.Environment, text types.XText) types.XValue {
	// escapes spaces as %20 matching urllib.quote(s, safe="") in Python
	encoded := strings.Replace(url.QueryEscape(text.Native()), "+", "%20", -1)
	return types.NewXText(encoded)
}

//------------------------------------------------------------------------------------------
// Number Functions
//------------------------------------------------------------------------------------------

// Abs returns the absolute value of `num`.
//
//   @(abs(-10)) -> 10
//   @(abs(10.5)) -> 10.5
//   @(abs("foo")) -> ERROR
//
// @function abs(num)
func Abs(env utils.Environment, num types.XNumber) types.XValue {
	return types.NewXNumber(num.Native().Abs())
}

// Round rounds `num` to the nearest value.
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
// @function round(num [,places])
func Round(env utils.Environment, num types.XNumber, places int) types.XValue {
	return types.NewXNumber(num.Native().Round(int32(places)))
}

// RoundUp rounds `num` up to the nearest integer value.
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
// @function round_up(num [,places])
func RoundUp(env utils.Environment, num types.XNumber, places int) types.XValue {
	dec := num.Native()
	if dec.Round(int32(places)).Equal(dec) {
		return num
	}

	halfPrecision := decimal.New(5, -int32(places)-1)
	roundedDec := dec.Add(halfPrecision).Round(int32(places))

	return types.NewXNumber(roundedDec)
}

// RoundDown rounds `num` down to the nearest integer value.
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
// @function round_down(num [,places])
func RoundDown(env utils.Environment, num types.XNumber, places int) types.XValue {
	dec := num.Native()
	if dec.Round(int32(places)).Equal(dec) {
		return num
	}

	halfPrecision := decimal.New(5, -int32(places)-1)
	roundedDec := dec.Sub(halfPrecision).Round(int32(places))

	return types.NewXNumber(roundedDec)
}

// Max returns the maximum value in `values`.
//
//   @(max(1, 2)) -> 2
//   @(max(1, -1, 10)) -> 10
//   @(max(1, 10, "foo")) -> ERROR
//
// @function max(values...)
func Max(env utils.Environment, values ...types.XValue) types.XValue {
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

// Min returns the minimum value in `values`.
//
//   @(min(1, 2)) -> 1
//   @(min(2, 2, -10)) -> -10
//   @(min(1, 2, "foo")) -> ERROR
//
// @function min(values)
func Min(env utils.Environment, values ...types.XValue) types.XValue {
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

// Mean returns the arithmetic mean of the numbers in `values`.
//
//   @(mean(1, 2)) -> 1.5
//   @(mean(1, 2, 6)) -> 3
//   @(mean(1, "foo")) -> ERROR
//
// @function mean(values)
func Mean(env utils.Environment, args ...types.XValue) types.XValue {
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
func Mod(env utils.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXNumber(num1.Native().Mod(num2.Native()))
}

// Rand returns a single random number between [0.0-1.0).
//
//   @(rand()) -> 0.3849275689214193274523267973563633859157562255859375
//   @(rand()) -> 0.607552015674623913099594574305228888988494873046875
//
// @function rand()
func Rand(env utils.Environment) types.XValue {
	return types.NewXNumber(utils.RandDecimal())
}

// RandBetween a single random integer in the given inclusive range.
//
//   @(rand_between(1, 10)) -> 5
//   @(rand_between(1, 10)) -> 10
//
// @function rand_between()
func RandBetween(env utils.Environment, min types.XNumber, max types.XNumber) types.XValue {
	span := (max.Native().Sub(min.Native())).Add(decimal.New(1, 0))

	val := utils.RandDecimal().Mul(span).Add(min.Native()).Floor()

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
// * `MM`        - month 01-12
// * `D`         - day of month, 1-31
// * `DD`        - day of month, zero padded 0-31
// * `h`         - hour of the day 1-12
// * `hh`        - hour of the day 01-12
// * `tt`        - twenty four hour of the day 01-23
// * `m`         - minute 0-59
// * `mm`        - minute 00-59
// * `s`         - second 0-59
// * `ss`        - second 00-59
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
func ParseDateTime(env utils.Environment, args ...types.XValue) types.XValue {
	str, xerr := types.ToXText(env, args[0])
	if xerr != nil {
		return xerr
	}

	format, xerr := types.ToXText(env, args[1])
	if xerr != nil {
		return xerr
	}

	// try to turn it to a go format
	goFormat, err := utils.ToGoDateFormat(format.Native(), utils.DateTimeFormatting)
	if err != nil {
		return types.NewXError(err)
	}

	// grab our location
	location := env.Timezone()
	if len(args) == 3 {
		tzStr, xerr := types.ToXText(env, args[2])
		if xerr != nil {
			return xerr
		}

		location, err = time.LoadLocation(tzStr.Native())
		if err != nil {
			return types.NewXError(err)
		}
	}

	// finally try to parse the date
	parsed, err := time.ParseInLocation(goFormat, str.Native(), location)
	if err != nil {
		return types.NewXError(err)
	}

	return types.NewXDateTime(parsed.In(location))
}

// DateTimeFromEpoch converts the UNIX epoch time `seconds` into a new date.
//
//   @(datetime_from_epoch(1497286619)) -> 2017-06-12T11:56:59.000000-05:00
//   @(datetime_from_epoch(1497286619.123456)) -> 2017-06-12T11:56:59.123456-05:00
//
// @function datetime_from_epoch(seconds)
func DateTimeFromEpoch(env utils.Environment, num types.XNumber) types.XValue {
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
func DateTimeDiff(env utils.Environment, arg1 types.XValue, arg2 types.XValue, arg3 types.XValue) types.XValue {
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
		return types.NewXNumberFromInt(utils.DaysBetween(date2.Native(), date1.Native()))
	case "W":
		return types.NewXNumberFromInt(int(utils.DaysBetween(date2.Native(), date1.Native()) / 7))
	case "M":
		return types.NewXNumberFromInt(utils.MonthsBetween(date2.Native(), date1.Native()))
	case "Y":
		return types.NewXNumberFromInt(date2.Native().Year() - date1.Native().Year())
	}

	return types.NewXErrorf("unknown unit: %s, must be one of s, m, h, D, W, M, Y", unit)
}

// DateTimeAdd calculates the date value arrived at by adding `offset` number of `unit` to the `date`
//
// Valid durations are "Y" for years, "M" for months, "W" for weeks, "D" for days, "h" for hour,
// "m" for minutes, "s" for seconds
//
//   @(datetime_add("2017-01-15", 5, "D")) -> 2017-01-20T00:00:00.000000-05:00
//   @(datetime_add("2017-01-15 10:45", 30, "m")) -> 2017-01-15T11:15:00.000000-05:00
//
// @function datetime_add(date, offset, unit)
func DateTimeAdd(env utils.Environment, args ...types.XValue) types.XValue {
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

// ReplaceTime returns the a new date time with the time part replaced by the `time`.
//
//   @(replace_time(now(), "10:30")) -> 2018-04-11T10:30:00.000000-05:00
//   @(replace_time("2017-01-15", "10:30")) -> 2017-01-15T10:30:00.000000-05:00
//   @(replace_time("foo", "10:30")) -> ERROR
//
// @function replace_time(date)
func ReplaceTime(env utils.Environment, args ...types.XValue) types.XValue {
	date, xerr := types.ToXDateTime(env, args[0])
	if xerr != nil {
		return xerr
	}
	t, xerr := types.ToXTime(env, args[1])
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
func TZ(env utils.Environment, date types.XDateTime) types.XValue {
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
func TZOffset(env utils.Environment, date types.XDateTime) types.XValue {
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
func Epoch(env utils.Environment, date types.XDateTime) types.XValue {
	nanos := decimal.New(date.Native().UnixNano(), 0)
	return types.NewXNumber(nanos.Div(nanosPerSecond))
}

// Now returns the current date and time in the current timezone.
//
//   @(now()) -> 2018-04-11T13:24:30.123456-05:00
//
// @function now()
func Now(env utils.Environment) types.XValue {
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
func DateFromParts(env utils.Environment, year, month, day int) types.XValue {
	if month < 1 || month > 12 {
		return types.NewXErrorf("invalid value for month, must be 1-12")
	}

	return types.NewXDate(utils.NewDate(year, month, day))
}

// Weekday returns the day of the week for `date`.
//
// The week is considered to start on Sunday so a Sunday returns 0, a Monday returns 1 etc.
//
//   @(weekday("2017-01-15")) -> 0
//   @(weekday("foo")) -> ERROR
//
// @function weekday(date)
func Weekday(env utils.Environment, date types.XDate) types.XValue {
	return types.NewXNumberFromInt(int(date.Native().Weekday()))
}

// Today returns the current date in the environment timezone.
//
//   @(today()) -> 2018-04-11
//
// @function today()
func Today(env utils.Environment) types.XValue {
	return types.NewXDate(utils.ExtractDate(env.Now()))
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
// * `hh`        - hour of the day 01-12
// * `tt`        - twenty four hour of the day 01-23
// * `m`         - minute 0-59
// * `mm`        - minute 00-59
// * `s`         - second 0-59
// * `ss`        - second 00-59
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
func ParseTime(env utils.Environment, args ...types.XValue) types.XValue {
	str, xerr := types.ToXText(env, args[0])
	if xerr != nil {
		return xerr
	}

	format, xerr := types.ToXText(env, args[1])
	if xerr != nil {
		return xerr
	}

	// try to turn it to a go format
	goFormat, err := utils.ToGoDateFormat(format.Native(), utils.TimeOnlyFormatting)
	if err != nil {
		return types.NewXError(err)
	}

	// finally try to parse the date
	parsed, err := utils.ParseTimeOfDay(goFormat, str.Native())
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
func TimeFromParts(env utils.Environment, hour, minute, second int) types.XValue {
	if hour < 0 || hour > 23 {
		return types.NewXErrorf("invalid value for hour, must be 0-23")
	}
	if minute < 0 || minute > 59 {
		return types.NewXErrorf("invalid value for minute, must be 0-59")
	}
	if second < 0 || second > 59 {
		return types.NewXErrorf("invalid value for second, must be 0-59")
	}

	return types.NewXTime(utils.NewTimeOfDay(hour, minute, second, 0))
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
func ParseJSON(env utils.Environment, text types.XText) types.XValue {
	return types.JSONToXValue([]byte(text.Native()))
}

// JSON returns the JSON representation of `value`.
//
//   @(json("string")) -> "string"
//   @(json(10)) -> 10
//   @(json(contact.uuid)) -> "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
//
// @function json(value)
func JSON(env utils.Environment, value types.XValue) types.XValue {
	asJSON, xerr := types.ToXJSON(env, value)
	if xerr != nil {
		return xerr
	}
	return asJSON
}

//----------------------------------------------------------------------------------------
// Formatting Functions
//----------------------------------------------------------------------------------------

// FormatDate formats `date` as text according to the given `format`. If `format` is not
// specified then the environment's default format is used.
//
// The format string can consist of the following characters. The characters
// ' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.
//
// * `YY`        - last two digits of year 0-99
// * `YYYY`      - four digits of year 0000-9999
// * `M`         - month 1-12
// * `MM`        - month 01-12
// * `D`         - day of month, 1-31
// * `DD`        - day of month, zero padded 0-31
//
//   @(format_date("1979-07-18T15:00:00.000000Z")) -> 1979-07-18
//   @(format_date("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) -> 1979-07-18
//   @(format_date("2010-05-10T19:50:00.000000Z", "YYYY M DD")) -> 2010 5 10
//   @(format_date("1979-07-18T15:00:00.000000Z", "YYYY")) -> 1979
//   @(format_date("1979-07-18T15:00:00.000000Z", "M")) -> 7
//   @(format_date("NOT DATE", "YYYY-MM-DD")) -> ERROR
//
// @function format_date(date, [,format])
func FormatDate(env utils.Environment, args ...types.XValue) types.XValue {
	date, xerr := types.ToXDate(env, args[0])
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
		format = types.NewXText(env.DateFormat().String())
	}

	// try to turn it to a go format
	goFormat, err := utils.ToGoDateFormat(format.Native(), utils.DateOnlyFormatting)
	if err != nil {
		return types.NewXError(err)
	}

	return types.NewXText(date.Native().Format(goFormat))
}

// FormatDateTime formats `date` as text according to the given `format`. If `format` is not
// specified then the environment's default format is used.
//
// The format string can consist of the following characters. The characters
// ' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.
//
// * `YY`        - last two digits of year 0-99
// * `YYYY`      - four digits of year 0000-9999
// * `M`         - month 1-12
// * `MM`        - month 01-12
// * `D`         - day of month, 1-31
// * `DD`        - day of month, zero padded 0-31
// * `h`         - hour of the day 1-12
// * `hh`        - hour of the day 01-12
// * `tt`        - twenty four hour of the day 01-23
// * `m`         - minute 0-59
// * `mm`        - minute 00-59
// * `s`         - second 0-59
// * `ss`        - second 00-59
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
//   @(format_datetime("1979-07-18T15:00:00.000000Z")) -> 1979-07-18 10:00
//   @(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) -> 1979-07-18
//   @(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY M DD tt:mm")) -> 2010 5 10 14:50
//   @(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY-MM-DD tt:mm AA", "America/Los_Angeles")) -> 2010-05-10 12:50 PM
//   @(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY")) -> 1979
//   @(format_datetime("1979-07-18T15:00:00.000000Z", "M")) -> 7
//   @(format_datetime("NOT DATE", "YYYY-MM-DD")) -> ERROR
//
// @function format_datetime(date [,format [,timezone]])
func FormatDateTime(env utils.Environment, args ...types.XValue) types.XValue {
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

	// try to turn it to a go format
	goFormat, err := utils.ToGoDateFormat(format.Native(), utils.DateTimeFormatting)
	if err != nil {
		return types.NewXError(err)
	}

	// grab our location
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

	// convert to our timezone if we have one (otherwise we remain in the date's default)
	if location != nil {
		date = types.NewXDateTime(date.Native().In(location))
	}

	return types.NewXText(date.Native().Format(goFormat))
}

// FormatTime formats `time` as text according to the given `format`. If `format` is not
// specified then the environment's default format is used.
//
// The format string can consist of the following characters. The characters
// ' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.
//
// * `h`         - hour of the day 1-12
// * `hh`        - hour of the day 01-12
// * `tt`        - twenty four hour of the day 01-23
// * `m`         - minute 0-59
// * `mm`        - minute 00-59
// * `s`         - second 0-59
// * `ss`        - second 00-59
// * `fff`       - milliseconds
// * `ffffff`    - microseconds
// * `fffffffff` - nanoseconds
// * `aa`        - am or pm
// * `AA`        - AM or PM
//
//   @(format_time("14:50:30.000000")) -> 02:50
//   @(format_time("14:50:30.000000", "h:mm aa")) -> 2:50 pm
//   @(format_time("14:50:30.000000", "tt:mm")) -> 14:50
//   @(format_time("15:00:27.000000", "s")) -> 27
//   @(format_time("NOT TIME", "hh:mm")) -> ERROR
//
// @function format_time(time [,format])
func FormatTime(env utils.Environment, args ...types.XValue) types.XValue {
	t, xerr := types.ToXTime(env, args[0])
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
		format = types.NewXText(env.TimeFormat().String())
	}

	// try to turn it to a go format
	goFormat, err := utils.ToGoDateFormat(format.Native(), utils.TimeOnlyFormatting)
	if err != nil {
		return types.NewXError(err)
	}

	return types.NewXText(t.Native().Format(goFormat))
}

// FormatNumber formats `number` to the given number of decimal `places`.
//
// An optional third argument `humanize` can be false to disable the use of thousand separators.
//
//   @(format_number(31337)) -> 31,337.00
//   @(format_number(31337, 2)) -> 31,337.00
//   @(format_number(31337, 2, true)) -> 31,337.00
//   @(format_number(31337, 0, false)) -> 31337
//   @(format_number("foo", 2, false)) -> ERROR
//
// @function format_number(number, places [, humanize])
func FormatNumber(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) < 1 || len(args) > 3 {
		return types.NewXErrorf("takes 1 to 3 arguments, got %d", len(args))
	}

	num, err := types.ToXNumber(env, args[0])
	if err != nil {
		return err
	}

	places := 2
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
		if human, err = types.ToXBoolean(env, args[2]); err != nil {
			return err
		}
	}

	return types.NewXText(FormatDecimal(num.Native(), env.NumberFormat(), places, human.Native()))
}

// FormatDecimal formats the given decimal
func FormatDecimal(value decimal.Decimal, format *utils.NumberFormat, places int, groupDigits bool) string {
	// build our format string
	formatStr := strings.Builder{}
	if groupDigits {
		formatStr.WriteString(fmt.Sprintf("#%s###", format.DigitGroupingSymbol))
	} else {
		formatStr.WriteString("####")
	}
	formatStr.WriteString(format.DecimalSymbol)
	if places > 0 {
		for i := 0; i < places; i++ {
			formatStr.WriteString("#")
		}
	}
	f64, _ := value.Float64()
	return humanize.FormatFloat(formatStr.String(), f64)
}

// FormatLocation formats the given `location` as its name.
//
//   @(format_location("Rwanda")) -> Rwanda
//   @(format_location("Rwanda > Kigali")) -> Kigali
//
// @function format_location(location)
func FormatLocation(env utils.Environment, path types.XText) types.XValue {
	parts := strings.Split(path.Native(), ">")
	return types.NewXText(strings.TrimSpace(parts[len(parts)-1]))
}

// FormatURN formats `urn` into human friendly text.
//
//   @(format_urn("tel:+250781234567")) -> 0781 234 567
//   @(format_urn("twitter:134252511151#billy_bob")) -> billy_bob
//   @(format_urn(contact.urn)) -> (206) 555-1212
//   @(format_urn(contact.urns.telegram[0])) ->
//   @(format_urn(contact.urns[2])) -> foo@bar.com
//   @(format_urn(urns.mailto)) -> foo@bar.com
//   @(format_urn("NOT URN")) -> ERROR
//
// @function format_urn(urn)
func FormatURN(env utils.Environment, arg types.XText) types.XValue {
	urn := urns.URN(arg.Native())
	err := urn.Validate()
	if err != nil {
		return types.NewXErrorf("%s is not a valid URN: %s", arg.Native(), err)
	}

	return types.NewXText(urn.Format())
}

//------------------------------------------------------------------------------------------
// Utility Functions
//------------------------------------------------------------------------------------------

// Length returns the length of the passed in text or array.
//
// length will return an error if it is passed an item which doesn't have length.
//
//   @(length("Hello")) -> 5
//   @(length(contact.fields.gender)) -> 4
//   @(length("😀😃😄😁")) -> 4
//   @(length(array())) -> 0
//   @(length(array("a", "b", "c"))) -> 3
//   @(length(1234)) -> ERROR
//
// @function length(value)
func Length(env utils.Environment, value types.XValue) types.XValue {
	// a nil has length of zero
	if utils.IsNil(value) {
		return types.XNumberZero
	}

	// argument must be a value with length
	lengthable, isLengthable := value.(types.XLengthable)
	if isLengthable {
		return types.NewXNumberFromInt(lengthable.Length())
	}

	// or reduceable to something with length
	value = types.Reduce(env, value)
	lengthable, isLengthable = value.(types.XLengthable)
	if isLengthable {
		return types.NewXNumberFromInt(lengthable.Length())
	}

	return types.NewXErrorf("value doesn't have length")
}

// Default returns `value` if is not empty or an error, otherwise it returns `default`.
//
//   @(default(undeclared.var, "default_value")) -> default_value
//   @(default("10", "20")) -> 10
//   @(default("", "value")) -> value
//   @(default(array(1, 2), "value")) -> [1, 2]
//   @(default(array(), "value")) -> value
//   @(default(datetime("invalid-date"), "today")) -> today
//
// @function default(value, default)
func Default(env utils.Environment, value types.XValue, def types.XValue) types.XValue {
	if types.IsEmpty(value) || types.IsXError(value) {
		return def
	}

	return value
}

// LegacyAdd simulates our old + operator, which operated differently based on whether
// one of the parameters was a date or not. If one is a date, then the other side is
// expected to be an integer with a number of days to add to the date, otherwise a normal
// decimal addition is attempted.
func LegacyAdd(env utils.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {

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
func ReadChars(env utils.Environment, val types.XText) types.XValue {
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
