package functions

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	humanize "github.com/dustin/go-humanize"
	"github.com/shopspring/decimal"
)

var randSource = rand.NewSource(time.Now().UnixNano())

// XFunction defines the interface that Excellent functions must implement
type XFunction func(env utils.Environment, args ...types.XValue) types.XValue

// RegisterXFunction registers a new function in Excellent
func RegisterXFunction(name string, function XFunction) {
	XFUNCTIONS[name] = function
}

// XFUNCTIONS is our map of functions available in Excellent which aren't tests
var XFUNCTIONS = map[string]XFunction{
	"and": And,
	"if":  If,
	"or":  Or,

	"length":  Length,
	"default": Default,
	"array":   Array,

	"legacy_add": LegacyAdd,

	"round":      Round,
	"round_up":   OneNumberFunction("round_up", RoundUp),
	"round_down": OneNumberFunction("round_down", RoundDown),
	"max":        Max,
	"min":        Min,
	"mean":       Mean,
	"mod":        TwoNumberFunction("mod", Mod),
	"rand":       Rand,
	"abs":        OneNumberFunction("abs", Abs),

	"format_num": FormatNum,
	"read_code":  OneStringFunction("read_code", ReadCode),

	"to_json":    ToJSON,
	"from_json":  FromJSON,
	"url_encode": OneStringFunction("url_encode", URLEncode),

	"char":              OneNumberFunction("char", Char),
	"code":              OneStringFunction("code", Code),
	"split":             TwoStringFunction("split", Split),
	"join":              Join,
	"title":             OneStringFunction("title", Title),
	"word":              Word,
	"remove_first_word": OneStringFunction("remove_first_word", RemoveFirstWord),
	"word_count":        OneStringFunction("word_count", WordCount),
	"word_slice":        WordSlice,
	"field":             Field,
	"clean":             OneStringFunction("clean", Clean),
	"left":              StringAndIntegerFunction("left", Left),
	"lower":             OneStringFunction("lower", Lower),
	"right":             StringAndIntegerFunction("right", Right),
	"string_cmp":        TwoStringFunction("string_cmp", StringCmp),
	"repeat":            StringAndIntegerFunction("repeat", Repeat),
	"replace":           Replace,
	"upper":             OneStringFunction("upper", Upper),
	"percent":           OneNumberFunction("percent", Percent),

	"format_date":     FormatDate,
	"parse_date":      ParseDate,
	"date":            Date,
	"date_from_parts": DateFromParts,
	"date_diff":       DateDiff,
	"date_add":        DateAdd,
	"weekday":         OneDateFunction("weekday", Weekday),
	"tz":              OneDateFunction("tz", TZ),
	"tz_offset":       OneDateFunction("tz_offset", TZOffset),
	"today":           Today,
	"now":             Now,
	"from_epoch":      FromEpoch,
	"to_epoch":        ToEpoch,

	"format_urn": FormatURN,
}

//------------------------------------------------------------------------------------------
// Legacy Functions
//------------------------------------------------------------------------------------------

// LegacyAdd simulates our old + operator, which operated differently based on whether
// one of the parameters was a date or not. If one is a date, then the other side is
// expected to be an integer with a number of days to add to the date, otherwise a normal
// decimal addition is attempted.
func LegacyAdd(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 2 {
		return types.NewXErrorf("LEGACY_ADD requires exactly two arguments, got %d", len(args))
	}

	// try to parse dates and decimals
	date1, date1Err := types.ToXDate(env, args[0])
	date2, date2Err := types.ToXDate(env, args[1])

	dec1, dec1Err := types.ToXNumber(args[0])
	dec2, dec2Err := types.ToXNumber(args[1])

	// if they are both dates, that's an error
	if date1Err == nil && date2Err == nil {
		return types.NewXErrorf("LEGACY_ADD cannot operate on two dates")
	}

	// date and int, do a day addition
	if date1Err == nil && dec2Err == nil {
		if dec2.Native().IntPart() < math.MinInt32 || dec2.Native().IntPart() > math.MaxInt32 {
			return types.NewXErrorf("LEGACY_ADD cannot operate on integers greater than 32 bit")
		}
		return types.NewXDate(date1.Native().AddDate(0, 0, int(dec2.Native().IntPart())))
	}

	// int and date, do a day addition
	if date2Err == nil && dec1Err == nil {
		if dec1.Native().IntPart() < math.MinInt32 || dec1.Native().IntPart() > math.MaxInt32 {
			return types.NewXErrorf("LEGACY_ADD cannot operate on integers greater than 32 bit")
		}
		return types.NewXDate(date2.Native().AddDate(0, 0, int(dec1.Native().IntPart())))
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

//------------------------------------------------------------------------------------------
// Utility Functions
//------------------------------------------------------------------------------------------

// Length returns the length of the passed in string or array.
//
// length will return an error if it is passed an item which doesn't have length.
//
//   @(length("Hello")) -> 5
//   @(length("😀😃😄😁")) -> 4
//   @(length(array())) -> 0
//   @(length(array("a", "b", "c"))) -> 3
//   @(length(1234)) -> ERROR
//
// @function length(object)
func Length(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("LENGTH takes exactly one argument, got %d", len(args))
	}

	// argument must be a value with length
	lengthable, isLengthable := args[0].(types.XLengthable)
	if isLengthable {
		return types.NewXNumberFromInt(lengthable.Length())
	}

	return types.NewXErrorf("LENGTH requires an object with length as its first argument, got %s", reflect.TypeOf(args[0]))
}

// Default takes two arguments, returning `test` if not an error or nil, otherwise returning `default`
//
//   @(default(undeclared.var, "default_value")) -> default_value
//   @(default("10", "20")) -> 10
//   @(default(date("invalid-date"), "today")) -> today
//
// @function default(test, default)
func Default(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 2 {
		return types.NewXErrorf("DEFAULT takes exactly two arguments, got %d", len(args))
	}

	// first argument is nil, return arg2
	if args[0] == nil {
		return args[1]
	}

	// test whether arg1 is an error
	_, isErr := args[0].(types.XError)
	if isErr {
		return args[1]
	}

	return args[0]
}

// Array takes a list of `values` and returns them as an array
//
//   @(array("a", "b", 356)[1]) -> b
//   @(join(array("a", "b", "c"), "|")) -> a|b|c
//   @(length(array())) -> 0
//   @(length(array("a", "b"))) -> 2
//
// @function array(values...)
func Array(env utils.Environment, args ...types.XValue) types.XValue {
	return types.NewXArray(args...)
}

// FromJSON tries to parse `string` as JSON, returning a fragment you can index into
//
// If the passed in value is not JSON, then an error is returned
//
//   @(from_json("[1,2,3,4]").2) -> 3
//   @(from_json("invalid json")) -> ERROR
//
// @function from_json(string)
func FromJSON(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("FROM_JSON takes exactly one string argument, got %d", len(args))
	}

	arg, err := types.ToXString(args[0])
	if err != nil {
		return err
	}

	return types.JSONToXValue([]byte(arg.Native()))
}

// ToJSON tries to return a JSON representation of `value`. An error is returned if there is
// no JSON representation of that object.
//
//   @(to_json("string")) -> "string"
//   @(to_json(10)) -> 10
//   @(to_json(contact.uuid)) -> "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
//
// @function to_json(value)
func ToJSON(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("TO_JSON takes exactly one argument, got %d", len(args))
	}

	asJSON, err := types.ToXJSON(args[0])
	if err != nil {
		return err
	}
	return asJSON
}

// URLEncode URL encodes `string` for use in a URL parameter
//
//   @(url_encode("two words")) -> two+words
//   @(url_encode(10)) -> 10
//
// @function url_encode(string)
func URLEncode(env utils.Environment, str types.XString) types.XValue {
	return types.NewXString(url.QueryEscape(str.Native()))
}

//------------------------------------------------------------------------------------------
// Bool Functions
//------------------------------------------------------------------------------------------

// And returns whether all the passed in arguments are truthy
//
//   @(and(true)) -> true
//   @(and(true, false, true)) -> false
//
// @function and(tests...)
func And(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) == 0 {
		return types.NewXErrorf("AND requires at least one argument")
	}

	for _, arg := range args {
		asBool, err := types.ToXBool(arg)
		if err != nil {
			return err
		}
		if !asBool {
			return types.XBoolFalse
		}
	}
	return types.XBoolTrue
}

// Or returns whether if any of the passed in arguments are truthy
//
//   @(or(true)) -> true
//   @(or(true, false, true)) -> true
//
// @function or(tests...)
func Or(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) == 0 {
		return types.NewXErrorf("OR requires at least one argument")
	}

	for _, arg := range args {
		asBool, err := types.ToXBool(arg)
		if err != nil {
			return err
		}
		if asBool {
			return types.XBoolTrue
		}
	}
	return types.XBoolFalse
}

// If evaluates the `test` argument, and if truthy returns `true_value`, if not returning `false_value`
//
// If the first argument is an error that error is returned
//
//   @(if(1 = 1, "foo", "bar")) -> foo
//   @(if("foo" > "bar", "foo", "bar")) -> ERROR
//
// @function if(test, true_value, false_value)
func If(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 3 {
		return types.NewXErrorf("IF requires exactly 3 arguments, got %d", len(args))
	}

	asBool, err := types.ToXBool(args[0])
	if err != nil {
		return err
	}

	if asBool {
		return args[1]
	}
	return args[2]
}

//------------------------------------------------------------------------------------------
// Decimal Functions
//------------------------------------------------------------------------------------------

// Abs returns the absolute value of `num`
//
//   @(abs(-10)) -> 10
//   @(abs(10.5)) -> 10.5
//   @(abs("foo")) -> ERROR
//
// @function abs(num)
func Abs(env utils.Environment, num types.XNumber) types.XValue {
	return types.NewXNumber(num.Native().Abs())
}

// Round rounds `num` to the nearest value. You can optionally pass
// in the number of decimal places to round to as `places`.
//
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
//   @(round(12.141)) -> 12
//   @(round(12.6)) -> 13
//   @(round(12.141, 2)) -> 12.14
//   @(round(12.146, 2)) -> 12.15
//   @(round(12.146, -1)) -> 10
//   @(round("notnum", 2)) -> ERROR
//
// @function round(num [,places])
func Round(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) < 1 || len(args) > 2 {
		return types.NewXErrorf("ROUND takes either one or two arguments")
	}

	num, err := types.ToXNumber(args[0])
	if err != nil {
		return types.NewXErrorf("ROUND's first argument must be decimal")
	}

	round := 0
	if len(args) == 2 {
		round, err = types.ToInteger(args[1])
		if err != nil {
			return types.NewXErrorf("ROUND's decimal places argument must be integer")
		}
	}

	return types.NewXNumber(num.Native().Round(int32(round)))
}

// RoundUp rounds `num` up to the nearest integer value, also good at fighting weeds
//
//   @(round_up(12.141)) -> 13
//   @(round_up(12)) -> 12
//   @(round_up("foo")) -> ERROR
//
// @function round_up(num)
func RoundUp(env utils.Environment, num types.XNumber) types.XValue {
	return types.NewXNumber(num.Native().Ceil())
}

// RoundDown rounds `num` down to the nearest integer value
//
//   @(round_down(12.141)) -> 12
//   @(round_down(12.9)) -> 12
//   @(round_down("foo")) -> ERROR
//
// @function round_down(num)
func RoundDown(env utils.Environment, num types.XNumber) types.XValue {
	return types.NewXNumber(num.Native().Floor())
}

// Max takes a list of `values` and returns the greatest of them
//
//   @(max(1, 2)) -> 2
//   @(max(1, -1, 10)) -> 10
//   @(max(1, 10, "foo")) -> ERROR
//
// @function max(values...)
func Max(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) == 0 {
		return types.NewXErrorf("MAX takes at least one argument")
	}

	max, xerr := types.ToXNumber(args[0])
	if xerr != nil {
		return xerr
	}

	for _, v := range args[1:] {
		val, xerr := types.ToXNumber(v)
		if xerr != nil {
			return xerr
		}

		if val.Compare(max) > 0 {
			max = val
		}
	}
	return max
}

// Min takes a list of `values` and returns the smallest of them
//
//   @(min(1, 2)) -> 1
//   @(min(2, 2, -10)) -> -10
//   @(min(1, 2, "foo")) -> ERROR
//
// @function min(values)
func Min(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) == 0 {
		return types.NewXErrorf("MIN takes at least one argument")
	}

	max, xerr := types.ToXNumber(args[0])
	if xerr != nil {
		return xerr
	}

	for _, v := range args[1:] {
		val, xerr := types.ToXNumber(v)
		if xerr != nil {
			return xerr
		}

		if val.Compare(max) < 0 {
			max = val
		}
	}
	return max
}

// Mean takes a list of `values` and returns the arithmetic mean of them
//
//   @(mean(1, 2)) -> 1.5
//   @(mean(1, 2, 6)) -> 3
//   @(mean(1, "foo")) -> ERROR
//
// @function mean(values)
func Mean(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) == 0 {
		return types.NewXErrorf("MEAN requires at least one argument, got 0")
	}

	sum := decimal.Zero

	for _, val := range args {
		num, xerr := types.ToXNumber(val)
		if xerr != nil {
			return xerr
		}
		sum = sum.Add(num.Native())
	}

	return types.NewXNumber(sum.Div(decimal.New(int64(len(args)), 0)))
}

// Mod returns the remainder of the division of `divident` by `divisor`
//
//   @(mod(5, 2)) -> 1
//   @(mod(4, 2)) -> 0
//   @(mod(5, "foo")) -> ERROR
//
// @function mod(dividend, divisor)
func Mod(env utils.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXNumber(num1.Native().Mod(num2.Native()))
}

// Rand returns either a single random decimal between 0-1 or a random integer between `floor` and `ceiling` (inclusive)
//
//   @(rand() > 0) -> true
//   @(rand(1, 5) <= 5) -> true
//
// @function rand(floor, ceiling)
func Rand(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 0 && len(args) != 2 {
		return types.NewXErrorf("RAND takes either no arguments or two arguments, got %d", len(args))
	}

	if len(args) == 0 {
		return types.NewXNumber(decimal.NewFromFloat(rand.New(randSource).Float64()))
	}

	min, xerr := types.ToXNumber(args[0])
	if xerr != nil {
		return xerr
	}
	max, xerr := types.ToXNumber(args[1])
	if xerr != nil {
		return xerr
	}

	// turn to integers
	minDec := min.Native().Floor()
	maxDec := max.Native().Floor()

	spread := minDec.Sub(maxDec).Abs()

	// we add one here as the golang rand does is not inclusive, 2 will always return 1
	// since our contract is inclusive of both ends we need one more
	add := rand.New(randSource).Int63n(spread.IntPart() + 1)

	if minDec.Cmp(maxDec) <= 0 {
		return types.NewXNumber(minDec.Add(decimal.NewFromFloat(float64(add))))
	}
	return types.NewXNumber(maxDec.Add(decimal.NewFromFloat(float64(add))))
}

// FormatNum returns `num` formatted with the passed in number of decimal `places` and optional `commas` dividing thousands separators
//
//   @(format_num(31337, 2, true)) -> 31,337.00
//   @(format_num(31337, 0, false)) -> 31337
//   @(format_num("foo", 2, false)) -> ERROR
//
// @function format_num(num, places, commas)
func FormatNum(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 3 {
		return types.NewXErrorf("FORMAT_NUM takes exactly three arguments, got %d", len(args))
	}

	num, err := types.ToXNumber(args[0])
	if err != nil {
		return err
	}

	places, err := types.ToInteger(args[1])
	if err != nil {
		return err
	}
	if places < 0 || places > 9 {
		return types.NewXErrorf("FORMAT_NUM must take 0-9 number of places, got %d", args[1])
	}

	commas, err := types.ToXBool(args[2])
	if err != nil {
		return err
	}

	// build our format string
	formatStr := bytes.Buffer{}
	if commas {
		formatStr.WriteString("#,###.")
	} else {
		formatStr.WriteString("####.")
	}
	if places > 0 {
		for i := 0; i < places; i++ {
			formatStr.WriteString("#")
		}
	}
	f64, _ := num.Native().Float64()
	return types.NewXString(humanize.FormatFloat(formatStr.String(), f64))
}

//------------------------------------------------------------------------------------------
// IVR Functions
//------------------------------------------------------------------------------------------

// ReadCode converts `code` into something that can be read by IVR systems
//
// ReadCode will split the numbers such as they are easier to understand. This includes
// splitting in 3s or 4s if appropriate.
//
//   @(read_code("1234")) -> 1 2 3 4
//   @(read_code("abc")) -> a b c
//   @(read_code("abcdef")) -> a b c , d e f
//
// @function read_code(code)
func ReadCode(env utils.Environment, val types.XString) types.XValue {
	var output bytes.Buffer

	// remove any leading +
	val = types.NewXString(strings.TrimLeft(val.Native(), "+"))

	length := len(val)

	// groups of three
	if length%3 == 0 {
		// groups of 3
		for i := 0; i < length; i += 3 {
			if i > 0 {
				output.WriteString(" , ")
			}
			output.WriteString(strings.Join(strings.Split(val.Native()[i:i+3], ""), " "))
		}
		return types.NewXString(output.String())
	}

	// groups of four
	if length%4 == 0 {
		for i := 0; i < length; i += 4 {
			if i > 0 {
				output.WriteString(" , ")
			}
			output.WriteString(strings.Join(strings.Split(val.Native()[i:i+4], ""), " "))
		}
		return types.NewXString(output.String())
	}

	// default, just do one at a time
	for i, c := range val {
		if i > 0 {
			output.WriteString(" , ")
		}
		output.WriteRune(c)
	}

	return types.NewXString(output.String())
}

//------------------------------------------------------------------------------------------
// String Functions
//------------------------------------------------------------------------------------------

// Code returns the numeric code for the first character in `string`, it is the inverse of char
//
//   @(code("a")) -> 97
//   @(code("abc")) -> 97
//   @(code("😀")) -> 128512
//   @(code("15")) -> 49
//   @(code(15)) -> 49
//   @(code("")) -> ERROR
//
// @function code(string)
func Code(env utils.Environment, str types.XString) types.XValue {
	if str.Length() == 0 {
		return types.NewXErrorf("CODE requires a string of at least one character")
	}

	r, _ := utf8.DecodeRuneInString(str.Native())
	return types.NewXNumberFromInt(int(r))
}

// Split splits `string` based on the passed in `delimeter`
//
// Empty values are removed from the returned list
//
//   @(split("a b c", " ")) -> ["a","b","c"]
//   @(split("a", " ")) -> ["a"]
//   @(split("abc..d", ".")) -> ["abc","d"]
//   @(split("a.b.c.", ".")) -> ["a","b","c"]
//   @(split("a && b && c", " && ")) -> ["a","b","c"]
//
// @function split(string, delimiter)
func Split(env utils.Environment, s types.XString, sep types.XString) types.XValue {
	splits := types.NewXArray()
	allSplits := strings.Split(s.Native(), sep.Native())
	for i := range allSplits {
		if allSplits[i] != "" {
			splits.Append(types.NewXString(allSplits[i]))
		}
	}
	return splits
}

// Join joins the passed in `array` of strings with the passed in `delimeter`
//
//   @(join(array("a", "b", "c"), "|")) -> a|b|c
//   @(join(split("a.b.c", "."), " ")) -> a b c
//
// @function join(array, delimeter)
func Join(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 2 {
		return types.NewXErrorf("JOIN takes exactly two arguments: the array to join and delimiter, got %d", len(args))
	}

	indexable, isIndexable := args[0].(types.XIndexable)
	if !isIndexable {
		return types.NewXErrorf("JOIN requires an indexable as its first argument, got %s", reflect.TypeOf(args[0]))
	}

	sep, err := types.ToXString(args[1])
	if err != nil {
		return err
	}

	var output bytes.Buffer
	for i := 0; i < indexable.Length(); i++ {
		if i > 0 {
			output.WriteString(sep.Native())
		}
		itemAsStr, err := types.ToXString(indexable.Index(i))
		if err != nil {
			return err
		}

		output.WriteString(itemAsStr.Native())
	}

	return types.NewXString(output.String())
}

// Char returns the rune for the passed in codepoint, `num`, which may be unicode, this is the reverse of code
//
//   @(char(33)) -> !
//   @(char(128512)) -> 😀
//   @(char("foo")) -> ERROR
//
// @function char(num)
func Char(env utils.Environment, num types.XNumber) types.XValue {
	code, xerr := types.ToInteger(num)
	if xerr != nil {
		return xerr
	}

	return types.NewXString(string(rune(code)))
}

// Title titlecases the passed in `string`, capitalizing each word
//
//   @(title("foo")) -> Foo
//   @(title("ryan lewis")) -> Ryan Lewis
//   @(title(123)) -> 123
//
// @function title(string)
func Title(env utils.Environment, str types.XString) types.XValue {
	return types.NewXString(strings.Title(str.Native()))
}

// Word returns the word at the passed in `offset` for the passed in `string`
//
//   @(word("foo bar", 0)) -> foo
//   @(word("foo.bar", 0)) -> foo
//   @(word("one two.three", 2)) -> three
//
// @function word(string, offset)
func Word(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 2 {
		return types.NewXErrorf("WORD takes exactly two arguments, got %d", len(args))
	}

	val, xerr := types.ToXString(args[0])
	if xerr != nil {
		return xerr
	}

	word, xerr := types.ToInteger(args[1])
	if xerr != nil {
		return xerr
	}

	words := utils.TokenizeString(val.Native())
	if word >= len(words) {
		return types.NewXErrorf("Word offset %d is greater than number of words %d", word, len(words))
	}

	return types.NewXString(words[word])
}

// RemoveFirstWord removes the 1st word of `string`
//
//   @(remove_first_word("foo bar")) -> bar
//
// @function remove_first_word(string)
func RemoveFirstWord(env utils.Environment, str types.XString) types.XValue {
	words := utils.TokenizeString(str.Native())
	if len(words) > 1 {
		return types.NewXString(strings.Join(words[1:], " "))
	}

	return types.XStringEmpty
}

// WordSlice extracts a substring from `string` spanning from `start` up to but not-including `end`. (first word is 1)
//
//   @(word_slice("foo bar", 1, 1)) -> foo
//   @(word_slice("foo bar", 1, 3)) -> foo bar
//   @(word_slice("foo bar", 3, 4)) ->
//
// @function word_slice(string, start, end)
func WordSlice(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 3 {
		return types.NewXErrorf("WORD_SLICE takes exactly three arguments, got %d", len(args))
	}

	arg, xerr := types.ToXString(args[0])
	if xerr != nil {
		return xerr
	}

	start, xerr := types.ToInteger(args[1])
	if xerr != nil {
		return xerr
	}
	if start <= 0 {
		return types.NewXErrorf("WORD_SLICE must start with a positive index")
	}
	start--

	stop, xerr := types.ToInteger(args[2])
	if xerr != nil {
		return xerr
	}
	if start < 0 {
		return types.NewXErrorf("WORD_SLICE must have a stop of 0 or greater")
	}

	words := utils.TokenizeString(arg.Native())
	if start >= len(words) {
		return types.XStringEmpty
	}

	if stop >= len(words) {
		stop = len(words)
	}

	if stop > 0 {
		return types.NewXString(strings.Join(words[start:stop], " "))
	}
	return types.NewXString(strings.Join(words[start:], " "))
}

// WordCount returns the number of words in `string`
//
//   @(word_count("foo bar")) -> 2
//   @(word_count(10)) -> 1
//   @(word_count("")) -> 0
//   @(word_count("😀😃😄😁")) -> 4
//
// @function word_count(string)
func WordCount(env utils.Environment, str types.XString) types.XValue {
	words := utils.TokenizeString(str.Native())
	return types.NewXNumberFromInt(len(words))
}

// Field splits `string` based on the passed in `delimiter` and returns the field at `offset`.  When splitting
// with a space, the delimiter is considered to be all whitespace.  (first field is 0)
//
//   @(field("a,b,c", 1, ",")) -> b
//   @(field("a,,b,c", 1, ",")) ->
//   @(field("a   b c", 1, " ")) -> b
//   @(field("a		b	c	d", 1, "	")) ->
//   @(field("a\t\tb\tc\td", 1, " ")) ->
//   @(field("a,b,c", "foo", ",")) -> ERROR
//
// @function field(string, offset, delimeter)
func Field(env utils.Environment, args ...types.XValue) types.XValue {
	source, xerr := types.ToXString(args[0])
	if xerr != nil {
		return xerr
	}

	field, xerr := types.ToInteger(args[1])
	if xerr != nil {
		return xerr
	}

	if field < 0 {
		return types.NewXErrorf("cannot use a negative index to FIELD")
	}

	sep, xerr := types.ToXString(args[2])
	if xerr != nil {
		return xerr
	}

	fields := strings.Split(source.Native(), sep.Native())
	if field >= len(fields) {
		return types.XStringEmpty
	}

	// when using a space as a delimiter, we consider it splitting on whitespace, so remove empty values
	if sep.Native() == " " {
		var newFields []string
		for _, field := range fields {
			if field != "" {
				newFields = append(newFields, field)
			}
		}
		fields = newFields
	}

	return types.NewXString(strings.TrimSpace(fields[field]))
}

// Clean strips any leading or trailing whitespace from `string``
//
//   @(clean("\nfoo\t")) -> foo
//   @(clean(" bar")) -> bar
//   @(clean(123)) -> 123
//
// @function clean(string)
func Clean(env utils.Environment, str types.XString) types.XValue {
	return types.NewXString(strings.TrimSpace(str.Native()))
}

// Left returns the `count` most left characters of the passed in `string`
//
//   @(left("hello", 2)) -> he
//   @(left("hello", 7)) -> hello
//   @(left("😀😃😄😁", 2)) -> 😀😃
//   @(left("hello", -1)) -> ERROR
//
// @function left(string, count)
func Left(env utils.Environment, str types.XString, count int) types.XValue {
	if count < 0 {
		return types.NewXErrorf("LEFT can't take a negative count")
	}

	// this weird construct does the right thing for multi-byte unicode
	var output bytes.Buffer
	i := 0
	for _, r := range str.Native() {
		if i >= count {
			break
		}
		output.WriteRune(r)
		i++
	}

	return types.NewXString(output.String())
}

// Lower lowercases the passed in `string`
//
//   @(lower("HellO")) -> hello
//   @(lower("hello")) -> hello
//   @(lower("123")) -> 123
//   @(lower("😀")) -> 😀
//
// @function lower(string)
func Lower(env utils.Environment, str types.XString) types.XValue {
	return types.NewXString(strings.ToLower(str.Native()))
}

// Right returns the `count` most right characters of the passed in `string`
//
//   @(right("hello", 2)) -> lo
//   @(right("hello", 7)) -> hello
//   @(right("😀😃😄😁", 2)) -> 😄😁
//   @(right("hello", -1)) -> ERROR
//
// @function right(string, count)
func Right(env utils.Environment, str types.XString, count int) types.XValue {
	if count < 0 {
		return types.NewXErrorf("RIGHT can't take a negative count")
	}

	start := utf8.RuneCountInString(str.Native()) - count

	// this weird construct does the right thing for multi-byte unicode
	var output bytes.Buffer
	i := 0
	for _, r := range str.Native() {
		if i >= start {
			output.WriteRune(r)
		}
		i++
	}

	return types.NewXString(output.String())
}

// StringCmp returns the comparison between the strings `str1` and `str2`.
// The return value will be -1 if str1 is smaller than str2, 0 if they
// are equal and 1 if str1 is greater than str2
//
//   @(string_cmp("abc", "abc")) -> 0
//   @(string_cmp("abc", "def")) -> -1
//   @(string_cmp("zzz", "aaa")) -> 1
//
// @function string_cmp(str1, str2)
func StringCmp(env utils.Environment, str1 types.XString, str2 types.XString) types.XValue {
	return types.NewXNumberFromInt(str1.Compare(str2))
}

// Repeat return `string` repeated `count` number of times
//
//   @(repeat("*", 8)) -> ********
//   @(repeat("*", "foo")) -> ERROR
//
// @function repeat(string, count)
func Repeat(env utils.Environment, str types.XString, count int) types.XValue {
	if count < 0 {
		return types.NewXErrorf("REPEAT must be called with a positive integer, got %d", count)
	}

	var output bytes.Buffer
	for j := 0; j < count; j++ {
		output.WriteString(str.Native())
	}

	return types.NewXString(output.String())
}

// Replace replaces all occurrences of `needle` with `replacement` in `string`
//
//   @(replace("foo bar", "foo", "zap")) -> zap bar
//   @(replace("foo bar", "baz", "zap")) -> foo bar
//
// @function replace(string, needle, replacement)
func Replace(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 3 {
		return types.NewXErrorf("REPLACE takes exactly three arguments, got %d", len(args))
	}

	source, xerr := types.ToXString(args[0])
	if xerr != nil {
		return xerr
	}

	find, xerr := types.ToXString(args[1])
	if xerr != nil {
		return xerr
	}

	replace, xerr := types.ToXString(args[2])
	if xerr != nil {
		return xerr
	}

	return types.NewXString(strings.Replace(source.Native(), find.Native(), replace.Native(), -1))
}

// Upper uppercases all characters in the passed `string`
//
//   @(upper("Asdf")) -> ASDF
//   @(upper(123)) -> 123
//
// @function upper(string)
func Upper(env utils.Environment, str types.XString) types.XValue {
	return types.NewXString(strings.ToUpper(str.Native()))
}

// Percent converts `num` to a string represented as a percentage
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
	return types.NewXString(fmt.Sprintf("%d%%", percent.IntPart()))
}

//------------------------------------------------------------------------------------------
// Date & Time Functions
//------------------------------------------------------------------------------------------

// ParseDate turns `string` into a date according to the `format` and optional `timezone` specified
//
// The format string can consist of the following characters. The characters
// ' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.
//
// * `YY`        - last two digits of year 0-99
// * `YYYY`      - four digits of your 0000-9999
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
// as "America/Guayaquil" or "America/Los_Angeles". If not specified the timezone of your
// environment will be used. An error will be returned if the timezone is not recognized.
//
// Note that fractional seconds will be parsed even without an explicit format identifier.
// You should only specify fractional seconds when you want to assert the number of places
// in the input format.
//
// parse_date will return an error if it is unable to convert the string to a date.
//
//   @(parse_date("1979-07-18", "YYYY-MM-DD")) -> 1979-07-18T00:00:00.000000Z
//   @(parse_date("2010 5 10", "YYYY M DD")) -> 2010-05-10T00:00:00.000000Z
//   @(parse_date("2010 5 10 12:50", "YYYY M DD tt:mm", "America/Los_Angeles")) -> 2010-05-10T12:50:00.000000-07:00
//   @(parse_date("NOT DATE", "YYYY-MM-DD")) -> ERROR
//
// @function parse_date(string, format [,timezone])
func ParseDate(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) < 2 || len(args) > 3 {
		return types.NewXErrorf("PARSE_DATE requires at least two arguments, got %d", len(args))
	}

	arg1, xerr := types.ToXString(args[0])
	if xerr != nil {
		return xerr
	}

	format, xerr := types.ToXString(args[1])
	if xerr != nil {
		return xerr
	}

	// try to turn it to a go format
	goFormat, err := utils.ToGoDateFormat(format.Native())
	if err != nil {
		return types.NewXError(err)
	}

	// grab our location
	location := env.Timezone()
	if len(args) == 3 {
		arg3, xerr := types.ToXString(args[2])
		if xerr != nil {
			return xerr
		}

		location, err = time.LoadLocation(arg3.Native())
		if err != nil {
			return types.NewXError(err)
		}
	}

	// finally try to parse the date
	parsed, err := time.ParseInLocation(goFormat, arg1.Native(), location)
	if err != nil {
		return types.NewXError(err)
	}

	return types.NewXDate(parsed.In(location))
}

// FormatDate turns `date` into a string according to the `format` specified and in
// the optional `timezone`.
//
// The format string can consist of the following characters. The characters
// ' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.
//
// * `YY`        - last two digits of year 0-99
// * `YYYY`      - four digits of your 0000-9999
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
// as "America/Guayaquil" or "America/Los_Angeles". If not specified the timezone of your
// environment will be used. An error will be returned if the timezone is not recognized.
//
//   @(format_date("1979-07-18T15:00:00.000000Z")) -> 1979-07-18 15:00
//   @(format_date("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) -> 1979-07-18
//   @(format_date("2010-05-10T19:50:00.000000Z", "YYYY M DD tt:mm")) -> 2010 5 10 19:50
//   @(format_date("2010-05-10T19:50:00.000000Z", "YYYY-MM-DD tt:mm AA", "America/Los_Angeles")) -> 2010-05-10 12:50 PM
//   @(format_date("1979-07-18T15:00:00.000000Z", "YYYY")) -> 1979
//   @(format_date("1979-07-18T15:00:00.000000Z", "M")) -> 7
//   @(format_date("NOT DATE", "YYYY-MM-DD")) -> ERROR
//
// @function format_date(date, format [,timezone])
func FormatDate(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) < 1 || len(args) > 3 {
		return types.NewXErrorf("FORMAT_DATE takes one or two arguments, got %d", len(args))
	}
	date, xerr := types.ToXDate(env, args[0])
	if xerr != nil {
		return xerr
	}

	format := types.NewXString(fmt.Sprintf("%s %s", env.DateFormat().String(), env.TimeFormat().String()))
	if len(args) >= 2 {
		format, xerr = types.ToXString(args[1])
		if xerr != nil {
			return xerr
		}
	}

	// try to turn it to a go format
	goFormat, err := utils.ToGoDateFormat(format.Native())
	if err != nil {
		return types.NewXError(err)
	}

	// grab our location
	location := env.Timezone()
	if len(args) == 3 {
		arg3, xerr := types.ToXString(args[2])
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
		date = types.NewXDate(date.Native().In(location))
	}

	// return the formatted date
	return types.NewXString(date.Native().Format(goFormat))
}

// Date turns `string` into a date according to the environment's settings
//
// date will return an error if it is unable to convert the string to a date.
//
//   @(date("1979-07-18")) -> 1979-07-18T00:00:00.000000Z
//   @(date("2010 05 10")) -> 2010-05-10T00:00:00.000000Z
//   @(date("NOT DATE")) -> ERROR
//
// @function date(string)
func Date(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("DATE requires exactly one argument, got %d", len(args))
	}

	arg1, xerr := types.ToXString(args[0])
	if xerr != nil {
		return xerr
	}

	date, err := utils.DateFromString(env, arg1.Native())
	if err != nil {
		return types.NewXError(err)
	}

	return types.NewXDate(date)
}

// DateFromParts converts the passed in `year`, `month`` and `day`
//
//   @(date_from_parts(2017, 1, 15)) -> 2017-01-15T00:00:00.000000Z
//   @(date_from_parts(2017, 2, 31)) -> 2017-03-03T00:00:00.000000Z
//   @(date_from_parts(2017, 13, 15)) -> ERROR
//
// @function date_from_parts(year, month, day)
func DateFromParts(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 3 {
		return types.NewXErrorf("DATE_FROM_PARTS requires three arguments, got %d", len(args))
	}
	year, xerr := types.ToInteger(args[0])
	if xerr != nil {
		return xerr
	}
	month, xerr := types.ToInteger(args[1])
	if xerr != nil {
		return xerr
	}
	if month < 1 || month > 12 {
		return types.NewXErrorf("invalid value for month, must be 1-12")
	}

	day, xerr := types.ToInteger(args[2])
	if xerr != nil {
		return xerr
	}

	return types.NewXDate(time.Date(year, time.Month(month), day, 0, 0, 0, 0, env.Timezone()))
}

// DateDiff returns the integer duration between `date1` and `date2` in the `unit` specified.
//
// Valid durations are "y" for years, "M" for months, "w" for weeks, "d" for days, h" for hour,
// "m" for minutes, "s" for seconds
//
//   @(date_diff("2017-01-17", "2017-01-15", "d")) -> 2
//   @(date_diff("2017-01-17 10:50", "2017-01-17 12:30", "h")) -> -1
//   @(date_diff("2017-01-17", "2015-12-17", "y")) -> 2
//
// @function date_diff(date1, date2, unit)
func DateDiff(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 3 {
		return types.NewXErrorf("DATE_DIFF takes exactly three arguments, received %d", len(args))
	}

	date1, xerr := types.ToXDate(env, args[0])
	if xerr != nil {
		return xerr
	}

	date2, xerr := types.ToXDate(env, args[1])
	if xerr != nil {
		return xerr
	}

	unit, xerr := types.ToXString(args[2])
	if xerr != nil {
		return xerr
	}

	// find the duration between our dates
	duration := date1.Native().Sub(date2.Native())

	// then convert based on our unit
	switch unit.Native() {
	case "s":
		return types.NewXNumberFromInt(int(duration / time.Second))
	case "m":
		return types.NewXNumberFromInt(int(duration / time.Minute))
	case "h":
		return types.NewXNumberFromInt(int(duration / time.Hour))
	case "d":
		return types.NewXNumberFromInt(utils.DaysBetween(date1.Native(), date2.Native()))
	case "w":
		return types.NewXNumberFromInt(int(utils.DaysBetween(date1.Native(), date2.Native()) / 7))
	case "M":
		return types.NewXNumberFromInt(utils.MonthsBetween(date1.Native(), date2.Native()))
	case "y":
		return types.NewXNumberFromInt(date1.Native().Year() - date2.Native().Year())
	}

	return types.NewXErrorf("Unknown unit: %s, must be one of s, m, h, D, W, M, Y", unit)
}

// DateAdd calculates the date value arrived at by adding `offset` number of `unit` to the `date`
//
// Valid durations are "y" for years, "M" for months, "w" for weeks, "d" for days, h" for hour,
// "m" for minutes, "s" for seconds
//
//   @(date_add("2017-01-15", 5, "d")) -> 2017-01-20T00:00:00.000000Z
//   @(date_add("2017-01-15 10:45", 30, "m")) -> 2017-01-15T11:15:00.000000Z
//
// @function date_add(date, offset, unit)
func DateAdd(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 3 {
		return types.NewXErrorf("DATE_ADD takes exactly three arguments, received %d", len(args))
	}

	date, xerr := types.ToXDate(env, args[0])
	if xerr != nil {
		return xerr
	}

	duration, xerr := types.ToInteger(args[1])
	if xerr != nil {
		return xerr
	}

	unit, xerr := types.ToXString(args[2])
	if xerr != nil {
		return xerr
	}

	switch unit.Native() {
	case "s":
		return types.NewXDate(date.Native().Add(time.Duration(duration) * time.Second))
	case "m":
		return types.NewXDate(date.Native().Add(time.Duration(duration) * time.Minute))
	case "h":
		return types.NewXDate(date.Native().Add(time.Duration(duration) * time.Hour))
	case "d":
		return types.NewXDate(date.Native().AddDate(0, 0, duration))
	case "w":
		return types.NewXDate(date.Native().AddDate(0, 0, duration*7))
	case "M":
		return types.NewXDate(date.Native().AddDate(0, duration, 0))
	case "y":
		return types.NewXDate(date.Native().AddDate(duration, 0, 0))
	}

	return types.NewXErrorf("Unknown unit: %s, must be one of s, m, h, d, w, M, y", unit)
}

// Weekday returns the day of the week for `date`, 0 is sunday, 1 is monday..
//
//   @(weekday("2017-01-15")) -> 0
//   @(weekday("foo")) -> ERROR
//
// @function weekday(date)
func Weekday(env utils.Environment, date types.XDate) types.XValue {
	return types.NewXNumberFromInt(int(date.Native().Weekday()))
}

// TZ returns the timezone for `date``
//
// If not timezone information is present in the date, then the environment's
// timezone will be returned
//
//   @(tz("2017-01-15 02:15:18PM UTC")) -> UTC
//   @(tz("2017-01-15 02:15:18PM")) -> UTC
//   @(tz("2017-01-15")) -> UTC
//   @(tz("foo")) -> ERROR
//
// @function tz(date)
func TZ(env utils.Environment, date types.XDate) types.XValue {
	return types.NewXString(date.Native().Location().String())
}

// TZOffset returns the offset for the timezone as a string +/- HHMM for `date`
//
// If no timezone information is present in the date, then the environment's
// timezone offset will be returned
//
//   @(tz_offset("2017-01-15 02:15:18PM UTC")) -> +0000
//   @(tz_offset("2017-01-15 02:15:18PM")) -> +0000
//   @(tz_offset("2017-01-15")) -> +0000
//   @(tz_offset("foo")) -> ERROR
//
// @function tz_offset(date)
func TZOffset(env utils.Environment, date types.XDate) types.XValue {
	// this looks like we are returning a set offset, but this is how go describes formats
	return types.NewXString(date.Native().Format("-0700"))

}

// Today returns the current date in the current timezone, time is set to midnight in the environment timezone
//
//   @(format_date(today(), "YYYY") >= 2018) -> true
//
// @function today()
func Today(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) > 0 {
		return types.NewXErrorf("TODAY takes no arguments, got %d", len(args))
	}

	nowTZ := time.Now().In(env.Timezone())
	return types.NewXDate(time.Date(nowTZ.Year(), nowTZ.Month(), nowTZ.Day(), 0, 0, 0, 0, env.Timezone()))
}

// FromEpoch returns a new date created from `num` which represents number of nanoseconds since January 1st, 1970 GMT
//
//   @(from_epoch(1497286619000000000)) -> 2017-06-12T16:56:59.000000Z
//
// @function from_epoch(num)
func FromEpoch(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("FROM_EPOCH takes exactly one number argument, got %d", len(args))
	}

	offset, xerr := types.ToXNumber(args[0])
	if xerr != nil {
		return xerr
	}

	return types.NewXDate(time.Unix(0, offset.Native().IntPart()).In(env.Timezone()))
}

// ToEpoch converts `date` to the number of nanoseconds since January 1st, 1970 GMT
//
//   @(to_epoch("2017-06-12T16:56:59.000000Z")) -> 1497286619000000000
//
// @function to_epoch(date)
func ToEpoch(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("TO_EPOCH takes exactly one date argument, got %d", len(args))
	}

	date, xerr := types.ToXDate(env, args[0])
	if xerr != nil {
		return xerr
	}

	return types.NewXNumberFromInt64(date.Native().UnixNano())
}

// Now returns the current date and time in the environment timezone
//
//   @(format_date(now(), "YYYY") >= 2018) -> true
//
// @function now()
func Now(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) > 0 {
		return types.NewXErrorf("NOW takes no arguments, got %d", len(args))
	}

	return types.NewXDate(time.Now().In(env.Timezone()))
}

//----------------------------------------------------------------------------------------
// URN Functions
//----------------------------------------------------------------------------------------

// FormatURN turns `urn` into a human friendly string
//
//   @(format_urn("tel:+250781234567")) -> 0781 234 567
//   @(format_urn("twitter:134252511151#billy_bob")) -> billy_bob
//   @(format_urn(contact.urns)) -> (206) 555-1212
//   @(format_urn(contact.urns.1)) -> foo@bar.com
//   @(format_urn(contact.urns.mailto)) -> foo@bar.com
//   @(format_urn(contact.urns.mailto.0)) -> foo@bar.com
//   @(format_urn(contact.urns.telegram)) ->
//   @(format_urn("NOT URN")) -> ERROR
//
// @function format_urn(urn)
func FormatURN(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("FORMAT_URN takes one argument, got %d", len(args))
	}

	// if we've been passed an indexable like a URNList, use first item
	urnArg := args[0]

	indexable, isIndexable := urnArg.(types.XIndexable)
	if isIndexable {
		if indexable.Length() >= 1 {
			urnArg = indexable.Index(0)
		} else {
			return types.XStringEmpty
		}
	}

	urnString, xerr := types.ToXString(urnArg)
	if xerr != nil {
		return xerr
	}

	urn := urns.URN(urnString)
	err := urn.Validate()
	if err != nil {
		return types.NewXErrorf("%s is not a valid URN: %s", urnString, err)
	}

	return types.NewXString(urn.Format())
}
