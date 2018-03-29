package excellent

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"unicode/utf8"

	"math"

	"encoding/json"

	humanize "github.com/dustin/go-humanize"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
)

// XFunction defines the interface that Excellent functions must implement
type XFunction func(env utils.Environment, args ...interface{}) interface{}

// XFUNCTIONS is our map of functions available in Excellent which aren't tests
var XFUNCTIONS = map[string]XFunction{
	"and": And,
	"if":  If,
	"or":  Or,

	"array_length": ArrayLength,
	"default":      Default,

	"legacy_add": LegacyAdd,

	"round":      Round,
	"round_up":   RoundUp,
	"round_down": RoundDown,
	"int":        Int,
	"max":        Max,
	"min":        Min,
	"mean":       Mean,
	"mod":        Mod,
	"rand":       Rand,
	"abs":        Abs,

	"format_num": FormatNum,
	"read_code":  ReadCode,

	"to_json":    ToJSON,
	"from_json":  FromJSON,
	"url_encode": URLEncode,

	"char":              Char,
	"code":              Code,
	"split":             Split,
	"join":              Join,
	"title":             Title,
	"word":              Word,
	"remove_first_word": RemoveFirstWord,
	"word_count":        WordCount,
	"word_slice":        WordSlice,
	"field":             Field,
	"clean":             Clean,
	"left":              Left,
	"lower":             Lower,
	"length":            Length,
	"right":             Right,
	"string_length":     Length,
	"repeat":            Repeat,
	"replace":           Replace,
	"upper":             Upper,
	"percent":           Percent,

	"format_date":     FormatDate,
	"parse_date":      ParseDate,
	"date":            Date,
	"date_from_parts": DateFromParts,
	"date_diff":       DateDiff,
	"date_add":        DateAdd,
	"weekday":         Weekday,
	"tz":              TZ,
	"tz_offset":       TZOffset,
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
func LegacyAdd(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("LEGACY_ADD requires exactly two arguments, got %d", len(args))
	}

	// try to parse dates and decimals
	date1, date1Err := utils.ToDate(env, args[0])
	date2, date2Err := utils.ToDate(env, args[1])

	dec1, dec1Err := utils.ToDecimal(env, args[0])
	dec2, dec2Err := utils.ToDecimal(env, args[1])

	// if they are both dates, that's an error
	if date1Err == nil && date2Err == nil {
		return fmt.Errorf("LEGACY_ADD cannot operate on two dates")
	}

	// date and int, do a day addition
	if date1Err == nil && dec2Err == nil {
		if dec2.IntPart() < math.MinInt32 || dec2.IntPart() > math.MaxInt32 {
			return fmt.Errorf("LEGACY_ADD cannot operate on integers greater than 32 bit")
		}
		return date1.AddDate(0, 0, int(dec2.IntPart()))
	}

	// int and date, do a day addition
	if date2Err == nil && dec1Err == nil {
		if dec1.IntPart() < math.MinInt32 || dec1.IntPart() > math.MaxInt32 {
			return fmt.Errorf("LEGACY_ADD cannot operate on integers greater than 32 bit")
		}
		return date2.AddDate(0, 0, int(dec1.IntPart()))
	}

	// one of these doesn't look like a valid decimal either, bail
	if dec1Err != nil {
		return dec1Err
	}

	if dec2Err != nil {
		return dec2Err
	}

	// normal decimal addition
	return dec1.Add(dec2)
}

//------------------------------------------------------------------------------------------
// Utility Functions
//------------------------------------------------------------------------------------------

// ArrayLength returns the number of items in the passed in array
//
// array_length will return an error if it is passed an item which is not an array.
//
//    @(array_length(SPLIT("1 2 3", " "))) -> 3
//    @(array_length("123")) -> ERROR
//
// @function array_length(array)
func ArrayLength(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("ARRAY_LENGTH takes exactly one argument, got %d", len(args))
	}

	len, err := utils.SliceLength(args[0])
	if err != nil {
		return err
	}

	return len
}

// Default takes two arguments, returning `test` if not an error or nil, otherwise returning `default`
//
//   @(default(undeclared.var, "default_value")) -> default_value
//   @(default("10", "20")) -> 10
//   @(default(date("invalid-date"), "today")) -> today
//
// @function default(test, default)
func Default(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("DEFAULT takes exactly two arguments, got %d", len(args))
	}

	// first argument is nil, return arg2
	if args[0] == nil {
		return args[1]
	}

	// test whether arg1 is an error
	_, isErr := args[0].(error)
	if isErr {
		return args[1]
	}

	return args[0]
}

// FromJSON tries to parse `string` as JSON, returning a fragment you can index into
//
// If the passed in value is not JSON, then an error is returned
//
//   @(from_json("[1,2,3,4]").2) -> 3
//   @(from_json("invalid json")) -> ERROR
//
// @function from_json(string)
func FromJSON(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("FROM_JSON takes exactly one string argument, got %d", len(args))
	}

	arg, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	// unmarshal our string into a JSON fragment
	var fragment utils.JSONFragment
	err = json.Unmarshal([]byte(arg), &fragment)
	if err != nil {
		return err
	}
	return fragment
}

// ToJSON tries to return a JSON representation of `value`. An error is returned if there is
// no JSON representation of that object.
//
//  @(to_json("string")) -> "string"
//  @(to_json(10)) -> 10
//  @(to_json(contact.uuid)) -> "ce2b5142-453b-4e43-868e-abdafafaa878"
//  @(to_json(now())) -> "2017-05-10T12:50:00.000000-07:00"
//
// @function to_json(value)
func ToJSON(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("TO_JSON takes exactly one argument, got %d", len(args))
	}

	json, err := utils.ToJSON(env, args[0])
	if err != nil {
		return err
	}

	return json
}

// URLEncode URL encodes `string` for use in a URL parameter
//
//  @(url_encode("two words")) -> two+words
//  @(url_encode(10)) -> 10
//
// @function url_encode(string)
func URLEncode(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("URL_ENCODE takes exactly one argument, got %d", len(args))
	}

	arg1, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	return url.QueryEscape(arg1)
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
func And(env utils.Environment, args ...interface{}) interface{} {
	if len(args) == 0 {
		return fmt.Errorf("AND requires at least one argument")
	}

	val, err := utils.ToBool(env, args[0])
	if err != nil {
		return err
	}
	for _, iArg := range args[1:] {
		iVal, err := utils.ToBool(env, iArg)
		if err != nil {
			return err
		}
		val = val && iVal
	}
	return val
}

// Or returns whether if any of the passed in arguments are truthy
//
//   @(or(true)) -> true
//   @(or(true, false, true)) -> true
//
// @function or(tests...)
func Or(env utils.Environment, args ...interface{}) interface{} {
	if len(args) == 0 {
		return fmt.Errorf("OR requires at least one argument")
	}

	val, err := utils.ToBool(env, args[0])
	if err != nil {
		return err
	}

	for _, iArg := range args[1:] {
		iVal, err := utils.ToBool(env, iArg)
		if err != nil {
			return err
		}
		val = val || iVal
	}
	return val
}

// If evaluates the `test` argument, and if truthy returns `true_value`, if not returning `false_value`
//
// If the first argument is an error that error is returned
//
//   @(if(1 = 1, "foo", "bar")) -> "foo"
//   @(if("foo" > "bar", "foo", "bar")) -> ERROR
//
// @function if(test, true_value, false_value)
func If(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 3 {
		return fmt.Errorf("IF requires exactly 3 arguments, got %d", len(args))
	}

	truthy, err := utils.ToBool(env, args[0])
	if err != nil {
		return err
	}

	if truthy {
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
func Abs(env utils.Environment, args ...interface{}) interface{} {
	dec, err := checkOneDecimalArg(env, "ABS", args)
	if err != nil {
		return err
	}
	return dec.Abs()
}

// Round rounds `num` to the corresponding number of `places`
//
//   @(round(12.141, 2)) -> 12.14
//   @(round("notnum", 2)) -> ERROR
//
// @function round(num, places)
func Round(env utils.Environment, args ...interface{}) interface{} {
	dec, round, err := checkTwoDecimalArgs(env, "ROUND", args)
	if err != nil {
		return err
	}

	roundInt := round.IntPart()
	if roundInt < 0 {
		return fmt.Errorf("ROUND decimal places argument must be valid 32 bit integer")
	}

	return dec.Round(int32(roundInt))
}

// RoundUp rounds `num` up to the nearest integer value, also good at fighting weeds
//
//   @(round_up(12.141)) -> 13
//   @(round_up(12)) -> 12
//   @(round_up("foo")) -> ERROR
//
// @function round_up(num)
func RoundUp(env utils.Environment, args ...interface{}) interface{} {
	dec, err := checkOneDecimalArg(env, "ROUND_UP", args)
	if err != nil {
		return err
	}

	return dec.Ceil()
}

// RoundDown rounds `num` down to the nearest integer value
//
//   @(round_down(12.141)) -> 12
//   @(round_down(12.9)) -> 12
//   @(round_down("foo")) -> ERROR
//
// @function round_down(num)
func RoundDown(env utils.Environment, args ...interface{}) interface{} {
	dec, err := checkOneDecimalArg(env, "ROUND_DOWN", args)
	if err != nil {
		return err
	}

	return dec.Floor()
}

// Int takes `num` and returns the integer value (floored)
//
//   @(int(12.14)) -> 12
//   @(int(12.9)) -> 12
//   @(int("foo")) -> ERROR
//
// @function int(num)
func Int(env utils.Environment, args ...interface{}) interface{} {
	dec, err := checkOneDecimalArg(env, "INT", args)
	if err != nil {
		return err
	}

	return dec.Floor()
}

// Max takes a list of `values` and returns the greatest of them
//
//   @(max(1, 2)) -> 2
//   @(max(1, -1, 10)) -> 10
//   @(max(1, 10, "foo")) -> ERROR
//
// @function max(values...)
func Max(env utils.Environment, args ...interface{}) interface{} {
	if len(args) == 0 {
		return fmt.Errorf("MAX takes at least one argument")
	}

	max, err := utils.ToDecimal(env, args[0])
	if err != nil {
		return err
	}

	for _, v := range args[1:] {
		val, err := utils.ToDecimal(env, v)
		if err != nil {
			return err
		}

		if val.Cmp(max) > 0 {
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
func Min(env utils.Environment, args ...interface{}) interface{} {
	if len(args) == 0 {
		return fmt.Errorf("MIN takes at least one argument")
	}

	max, err := utils.ToDecimal(env, args[0])
	if err != nil {
		return err
	}

	for _, v := range args[1:] {
		val, err := utils.ToDecimal(env, v)
		if err != nil {
			return err
		}

		if val.Cmp(max) < 0 {
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
func Mean(env utils.Environment, args ...interface{}) interface{} {
	if len(args) == 0 {
		return fmt.Errorf("Mean requires at least one argument, got 0")
	}

	sum := decimal.Zero

	for _, val := range args {
		dec, err := utils.ToDecimal(env, val)
		if err != nil {
			return err
		}
		sum = sum.Add(dec)
	}

	return sum.Div(decimal.NewFromFloat(float64(len(args))))
}

// Mod returns the remainder of the division of `divident` by `divisor`
//
//   @(mod(5, 2)) -> 1
//   @(mod(4, 2)) -> 0
//   @(mod(5, "foo")) -> ERROR
//
// @function mod(dividend, divisor)
func Mod(env utils.Environment, args ...interface{}) interface{} {
	arg1, arg2, err := checkTwoDecimalArgs(env, "MOD", args)
	if err != nil {
		return err
	}

	return arg1.Mod(arg2)
}

var randSource = rand.NewSource(time.Now().UnixNano())

// Rand returns either a single random decimal between 0-1 or a random integer between `floor` and `ceiling` (inclusive)
//
//  @(rand()) == 0.5152
//  @(rand(1, 5)) == 3
//
// @function rand(floor, ceiling)
func Rand(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 0 && len(args) != 2 {
		return fmt.Errorf("RAND takes either no arguments or two arguments, got %d", len(args))
	}

	if len(args) == 0 {
		return decimal.NewFromFloat(rand.New(randSource).Float64())
	}

	min, err := utils.ToDecimal(env, args[0])
	if err != nil {
		return err
	}
	max, err := utils.ToDecimal(env, args[1])
	if err != nil {
		return err
	}

	// turn to integers
	min = min.Floor()
	max = max.Floor()

	spread := min.Sub(max).Abs()

	// we add one here as the golang rand does is not inclusive, 2 will always return 1
	// since our contract is inclusive of both ends we need one more
	add := rand.New(randSource).Int63n(spread.IntPart() + 1)

	if min.Cmp(max) <= 0 {
		return min.Add(decimal.NewFromFloat(float64(add)))
	}
	return max.Add(decimal.NewFromFloat(float64(add)))
}

// FormatNum returns `num` formatted with the passed in number of decimal `places` and optional `commas` dividing thousands separators
//
//   @(format_num(31337, 2, true)) -> "31,337.00"
//   @(format_num(31337, 0, false)) -> "31337"
//   @(format_num("foo", 2, false)) -> ERROR
//
// @function format_num(num, places, commas)
func FormatNum(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 3 {
		return fmt.Errorf("FORMAT_NUM takes exactly three arguments, got %d", len(args))
	}

	dec, err := utils.ToDecimal(env, args[0])
	if err != nil {
		return err
	}

	places, err := utils.ToInt(env, args[1])
	if err != nil {
		return err
	}
	if places < 0 || places > 9 {
		return fmt.Errorf("FORMAT_NUM must take 0-9 number of places, got %d", args[1])
	}

	commas, err := utils.ToBool(env, args[2])
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
	f64, _ := dec.Float64()
	return humanize.FormatFloat(formatStr.String(), f64)
}

//------------------------------------------------------------------------------------------
// IVR Functions
//------------------------------------------------------------------------------------------

// ReadCode converts `code` into something that can be read by IVR systems
//
// ReadCode will split the numbers such as they are easier to understand. This includes
// splitting in 3s or 4s if appropriate.
//
//   @(read_code("1234")) -> "1 2 3 4"
//   @(read_code("abc")) -> "a b c"
//   @(read_code("abcdef")) -> "a b c , d e f"
//
// @function read_code(code)
func ReadCode(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("READ_CODE takes exactly one argument, got %d", len(args))
	}

	// convert to a string
	val, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	var output bytes.Buffer

	// remove any leading +
	val = strings.TrimLeft(val, "+")

	length := len(val)

	// groups of three
	if length%3 == 0 {
		// groups of 3
		for i := 0; i < length; i += 3 {
			if i > 0 {
				output.WriteString(" , ")
			}
			output.WriteString(strings.Join(strings.Split(val[i:i+3], ""), " "))
		}
		return output.String()
	}

	// groups of four
	if length%4 == 0 {
		for i := 0; i < length; i += 4 {
			if i > 0 {
				output.WriteString(" , ")
			}
			output.WriteString(strings.Join(strings.Split(val[i:i+4], ""), " "))
		}
		return output.String()
	}

	// default, just do one at a time
	for i, c := range val {
		if i > 0 {
			output.WriteString(" , ")
		}
		output.WriteRune(c)
	}

	return output.String()
}

//------------------------------------------------------------------------------------------
// String Functions
//------------------------------------------------------------------------------------------

// Code returns the numeric code for the first character in `string`, it is the inverse of char
//
//   @(code("a")) -> "97"
//   @(code("abc")) -> "97"
//   @(code("😀")) -> "128512"
//   @(code("")) -> "ERROR"
//   @(code("15")) -> "49"
//   @(code(15)) -> "49"
//
// @function code(string)
func Code(env utils.Environment, args ...interface{}) interface{} {
	str, err := checkOneStringArg(env, "code", args)
	if err != nil {
		return err
	}

	if len(str) == 0 {
		return fmt.Errorf("CODE requires a string of at least one character")
	}

	r, _ := utf8.DecodeRuneInString(str)
	return int(r)
}

// Split splits `string` based on the passed in `delimeter`
//
// Empty values are removed from the returned list
//
//   @(split("a b c", " ")) -> "a, b, c"
//   @(split("a", " ")) -> "a"
//   @(split("abc..d", ".")) -> "abc, d"
//   @(split("a.b.c.", ".")) -> "a, b, c"
//   @(split("a && b && c", " && ")) -> "a, b, c"
//
// @function split(string, delimeter)
func Split(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("SPLIT takes exactly two arguments: string and delimiter, got %d", len(args))
	}

	s, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	sep, err := utils.ToString(env, args[1])
	if err != nil {
		return err
	}

	allSplits := strings.Split(s, sep)
	splits := make([]string, 0, len(allSplits))
	for i := range allSplits {
		if allSplits[i] != "" {
			splits = append(splits, allSplits[i])
		}
	}
	return splits
}

// Join joins the passed in `array` of strings with the passed in `delimeter`
//
//   @(join(split("a.b.c", "."), " ")) -> "a b c"
//
// @function join(array, delimeter)
func Join(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("JOIN takes exactly two arguments: the array to join and delimiter, got %d", len(args))
	}

	s, err := utils.ToStringArray(env, args[0])
	if err != nil {
		return err
	}

	sep, err := utils.ToString(env, args[1])
	if err != nil {
		return err
	}

	return strings.Join(s, sep)
}

// Char returns the rune for the passed in codepoint, `num`, which may be unicode, this is the reverse of code
//
//   @(char(33)) -> "!"
//   @(char(128512)) -> "😀"
//   @(char("foo")) -> ERROR
//
// @function char(num)
func Char(env utils.Environment, args ...interface{}) interface{} {
	arg, err := checkOneDecimalArg(env, "CHAR", args)
	if err != nil {
		return err
	}

	return string(rune(arg.IntPart()))
}

// Title titlecases the passed in `string`, capitalizing each word
//
//   @(title("foo")) -> "Foo"
//   @(title("ryan lewis")) -> "Ryan Lewis"
//   @(title(123)) -> "123"
//
// @function title(string)
func Title(env utils.Environment, args ...interface{}) interface{} {
	arg, err := checkOneStringArg(env, "TITLE", args)
	if err != nil {
		return err
	}

	return strings.Title(arg)
}

// Word returns the word at the passed in `offset` for the passed in `string`
//
//   @(word("foo bar", 0)) -> "foo"
//   @(word("foo.bar", 0)) -> "foo"
//   @(word("one two.three", 2)) -> "three"
//
// @function word(string, offset)
func Word(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("WORD takes exactly two arguments, got %d", len(args))
	}

	val, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	word, err := utils.ToInt(env, args[1])
	if err != nil {
		return err
	}

	words := utils.TokenizeString(val)
	if word >= len(words) {
		return fmt.Errorf("Word offset %d is greater than number of words %d", word, len(words))
	}

	return words[word]
}

// RemoveFirstWord removes the 1st word of `string`
//
//   @(remove_first_word("foo bar")) -> "bar"
//
// @function remove_first_word(string)
func RemoveFirstWord(env utils.Environment, args ...interface{}) interface{} {
	arg, err := checkOneStringArg(env, "REMOVE_FIRST_WORD", args)
	if err != nil {
		return err
	}

	words := utils.TokenizeString(arg)
	if len(words) > 1 {
		return strings.Join(words[1:], " ")
	}

	return ""
}

// WordSlice extracts a substring from `string` spanning from `start` up to but not-including `end`. (first word is 1)
//
//   @(word_slice("foo bar", 1, 1)) -> "foo"
//   @(word_slice("foo bar", 1, 3)) -> "foo bar"
//   @(word_slice("foo bar", 3, 4)) -> ""
//
// @function word_slice(string, start, end)
func WordSlice(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 3 {
		return fmt.Errorf("WORD_SLICE takes exactly three arguments, got %d", len(args))
	}

	arg, err := utils.ToString(env, args[0])
	if err != nil {
		return fmt.Errorf("WORD_SLICE requires a string as its first argument")
	}

	start, err := utils.ToInt(env, args[1])
	if err != nil || start <= 0 {
		return fmt.Errorf("WORD_SLICE must start with a positive index")
	}
	start--

	stop, err := utils.ToInt(env, args[2])
	if err != nil || start < 0 {
		return fmt.Errorf("WORD_SLICE must have a stop of 0 or greater")
	}

	words := utils.TokenizeString(arg)
	if start >= len(words) {
		return ""
	}

	if stop >= len(words) {
		stop = len(words)
	}

	if stop > 0 {
		return strings.Join(words[start:stop], " ")
	}
	return strings.Join(words[start:], " ")
}

// WordCount returns the number of words in `string`
//
//   @(word_count("foo bar")) -> 2
//   @(word_count(10)) -> 1
//   @(word_count("")) -> 0
//   @(word_count("😀😃😄😁")) -> 4
//
// @function word_count(string)
func WordCount(env utils.Environment, args ...interface{}) interface{} {
	arg, err := checkOneStringArg(env, "WORD_COUNT", args)
	if err != nil {
		return err
	}

	words := utils.TokenizeString(arg)
	return decimal.NewFromFloat(float64(len(words)))
}

// Field splits `string` based on the passed in `delimiter` and returns the field at `offset`.  When splitting
// with a space, the delimiter is considered to be all whitespace.  (first field is 0)
//
//   @(field("a,b,c", 1, ",")) -> "b"
//   @(field("a,,b,c", 1, ",")) -> ""
//   @(field("a   b c", 1, " ")) -> "b"
//   @(field("a		b	c	d", 1, "	")) -> ""
//   @(field("a\t\tb\tc\td", 1, " ")) -> ""
//   @(field("a,b,c", "foo", ",")) -> ERROR
//
// @function field(string, offset, delimeter)
func Field(env utils.Environment, args ...interface{}) interface{} {
	source, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	field, err := utils.ToInt(env, args[1])
	if err != nil {
		return err
	}

	if field < 0 {
		return fmt.Errorf("Cannot use a negative index to FIELD")
	}

	sep, err := utils.ToString(env, args[2])
	if err != nil {
		return err
	}

	fields := strings.Split(source, sep)
	if field >= len(fields) {
		return ""
	}

	// when using a space as a delimiter, we consider it splitting on whitespace, so remove empty values
	if sep == " " {
		var newFields []string
		for _, field := range fields {
			if field != "" {
				newFields = append(newFields, field)
			}
		}
		fields = newFields
	}

	return strings.TrimSpace(fields[field])
}

// Clean strips any leading or trailing whitespace from `string``
//
//   @(clean("\nfoo\t")) -> "foo"
//   @(clean(" bar")) -> "bar"
//   @(clean(123)) -> "123"
//
// @function clean(string)
func Clean(env utils.Environment, args ...interface{}) interface{} {
	arg, err := checkOneStringArg(env, "CLEAN", args)
	if err != nil {
		return err
	}

	return strings.TrimSpace(arg)
}

// Left returns the `len` most left characters of the passed in `string`
//
//   @(left("hello", 2)) -> "he"
//   @(left("hello", 7)) -> "hello"
//   @(left("😀😃😄😁", 2)) -> "😀😃"
//   @(left("hello", -1)) -> ERROR
//
// @function left(string, len)
func Left(env utils.Environment, args ...interface{}) interface{} {
	str, l, err := checkOneStringOneIntArg(env, "LEFT", args)
	if err != nil {
		return err
	}

	// this weird construct does the right thing for multi-byte unicode
	var output bytes.Buffer
	i := 0
	for _, r := range str {
		if i >= l {
			break
		}
		output.WriteRune(r)
		i++
	}

	return output.String()
}

// Lower lowercases the passed in `string`
//
//   @(lower("HellO")) -> "hello"
//   @(lower("hello")) -> "hello"
//   @(lower("123")) -> "123"
//   @(lower("😀")) -> "😀"
//
// @function lower(string)
func Lower(env utils.Environment, args ...interface{}) interface{} {
	arg, err := checkOneStringArg(env, "LOWER", args)
	if err != nil {
		return err
	}

	return strings.ToLower(arg)
}

// Right returns the `len` most right characters of the passed in `string`
//
//   @(right("hello", 2)) -> "lo"
//   @(right("hello", 7)) -> "hello"
//   @(right("😀😃😄😁", 2)) -> "😄😁"
//   @(right("hello", -1)) -> ERROR
//
// @function right(string, len)
func Right(env utils.Environment, args ...interface{}) interface{} {
	str, l, err := checkOneStringOneIntArg(env, "RIGHT", args)
	if err != nil {
		return err
	}

	start := utf8.RuneCountInString(str) - l

	// this weird construct does the right thing for multi-byte unicode
	var output bytes.Buffer
	i := 0
	for _, r := range str {
		if i >= start {
			output.WriteRune(r)
		}
		i++
	}

	return output.String()
}

// Length returns the number of unicode characters in `string`
//
//   @(length("Hello")) -> 5
//   @(length("😀😃😄😁")) -> 4
//   @(length(1234)) -> 4
//
// @function length(string)
func Length(env utils.Environment, args ...interface{}) interface{} {
	arg, err := checkOneStringArg(env, "LENGTH", args)
	if err != nil {
		return err
	}

	return utf8.RuneCountInString(arg)
}

// Repeat return `string` repeated `count` number of times
//
//   @(repeat("*", 8)) -> "********"
//   @(repeat("*", "foo")) -> ERROR
//
// @function repeat(string, count)
func Repeat(env utils.Environment, args ...interface{}) interface{} {
	str, i, err := checkOneStringOneIntArg(env, "REPEAT", args)
	if err != nil {
		return err
	}

	if i < 0 {
		return fmt.Errorf("REPEAT must be called with a positive integer, got %d", i)
	}

	var output bytes.Buffer
	for j := 0; j < i; j++ {
		output.WriteString(str)
	}

	return output.String()
}

// Replace replaces all occurrences of `needle` with `replacement` in `string`
//
//   @(replace("foo bar", "foo", "zap")) -> "zap bar"
//   @(replace("foo bar", "baz", "zap")) -> "foo bar"
//
// @function replace(string, needle, replacement)
func Replace(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 3 {
		return fmt.Errorf("REPLACE takes exactly three arguments, got %d", len(args))
	}

	source, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	find, err := utils.ToString(env, args[1])
	if err != nil {
		return err
	}

	replace, err := utils.ToString(env, args[2])
	if err != nil {
		return err
	}

	return strings.Replace(source, find, replace, -1)
}

// Upper uppercases all characters in the passed `string`
//
//   @(upper("Asdf")) -> "ASDF"
//   @(upper(123)) -> "123"
//
// @function upper(string)
func Upper(env utils.Environment, args ...interface{}) interface{} {
	str, err := checkOneStringArg(env, "UPPER", args)
	if err != nil {
		return err
	}
	return strings.ToUpper(str)
}

// Percent converts `num` to a string represented as a percentage
//
//   @(percent(0.54234)) -> "54%"
//   @(percent(1.2)) -> "120%"
//   @(percent("foo")) -> ERROR
//
// @function percent(num)
func Percent(env utils.Environment, args ...interface{}) interface{} {
	dec, err := checkOneDecimalArg(env, "PERCENT", args)
	if err != nil {
		return err
	}

	// multiply by 100 and floor
	percent := dec.Mul(decimal.NewFromFloat(100)).Round(0)

	// add on a %
	return fmt.Sprintf("%d%%", percent.IntPart())
}

//------------------------------------------------------------------------------------------
// Date & Time Functions
//------------------------------------------------------------------------------------------

// ParseDate turns `string` into a date according to the `format` and optional `timezone` specified
//
// The format string can consist of the following characters. The characters
// ' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.
//
// * `YY`    - last two digits of year 0-99
// * `YYYY`  - four digits of your 0000-9999
// * `M`     - month 1-12
// * `MM`    - month 01-12
// * `D`     - day of month, 1-31
// * `DD`    - day of month, zero padded 0-31
// * `h`     - hour of the day 1-12
// * `hh`    - hour of the day 01-12
// * `tt`    - twenty four hour of the day 01-23
// * `m`     - minute 0-59
// * `mm`    - minute 00-59
// * `s`     - second 0-59
// * `ss`    - second 00-59
// * `fff`   - thousandths of a second
// * `aa`    - am or pm
// * `AA`    - AM or PM
// * `Z`     - hour and minute offset from UTC, or Z for UTC
// * `ZZZ`   - hour and minute offset from UTC
//
// Timezone should be a location name as specified in the IANA Time Zone database, such
// as "America/Guayaquil" or "America/Los_Angeles". If not specified the timezone of your
// environment will be used. An error will be returned if the timezone is not recognized.
//
// parse_date will return an error if it is unable to convert the string to a date.
//
//   @(parse_date("1979-07-18", "YYYY-MM-DD")) -> 1979-07-18T00:00:00.000000Z
//   @(parse_date("2010 5 10", "YYYY M DD")) -> 2010-05-10T00:00:00.000000Z
//   @(parse_date("2010 5 10 12:50", "YYYY M DD tt:mm", "America/Los_Angeles")) -> 2010-05-10T12:50:00.000000-07:00
//   @(parse_date("NOT DATE", "YYYY-MM-DD")) -> ERROR
//
// @function parse_date(string, format [,timezone])
func ParseDate(env utils.Environment, args ...interface{}) interface{} {
	if len(args) < 2 || len(args) > 3 {
		return fmt.Errorf("PARSE_DATE requires at least two arguments, got %d", len(args))
	}
	arg1, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	format, err := utils.ToString(env, args[1])
	if err != nil {
		return err
	}

	// try to turn it to a go format
	goFormat, err := utils.ToGoDateFormat(format)
	if err != nil {
		return err
	}

	// grab our location
	location := env.Timezone()
	if len(args) == 3 {
		arg3, err := utils.ToString(env, args[2])
		if err != nil {
			return err
		}
		location, err = time.LoadLocation(arg3)
		if err != nil {
			return err
		}
	}

	// finally try to parse the date
	parsed, err := time.ParseInLocation(goFormat, arg1, location)
	if err != nil {
		return err
	}

	parsed = parsed.In(location)
	return parsed
}

// FormatDate turns `date` into a string according to the `format` specified and in
// the optional `timezone`.
//
// The format string can consist of the following characters. The characters
// ' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.
//
// * `YY`    - last two digits of year 0-99
// * `YYYY`  - four digits of your 0000-9999
// * `M`     - month 1-12
// * `MM`    - month 01-12
// * `D`     - day of month, 1-31
// * `DD`    - day of month, zero padded 0-31
// * `h`     - hour of the day 1-12
// * `hh`    - hour of the day 01-12
// * `tt`    - twenty four hour of the day 01-23
// * `m`     - minute 0-59
// * `mm`    - minute 00-59
// * `s`     - second 0-59
// * `ss`    - second 00-59
// * `fff`   - thousandths of a second
// * `aa`    - am or pm
// * `AA`    - AM or PM
// * `Z`     - hour and minute offset from UTC, or Z for UTC
// * `ZZZ`   - hour and minute offset from UTC
//
// Timezone should be a location name as specified in the IANA Time Zone database, such
// as "America/Guayaquil" or "America/Los_Angeles". If not specified the timezone of your
// environment will be used. An error will be returned if the timezone is not recognized.
//
//   @(format_date("1979-07-18T15:00:00.000000Z")) -> 1979-07-18 15:00
//   @(format_date("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) -> 1979-07-18
//   @(format_date("2010-05-10T19:50:00.000000Z", "YYYY M DD tt:mm")) -> 2010 5 10 19:50
//   @(format_date("2010-05-10T19:50:00.000000Z", "YYYY-MM-DD tt:mm AA", "America/Los_Angeles")) -> 2010-05-10 12:50 PM
//   @(format_date("NOT DATE", "YYYY-MM-DD")) -> ERROR
//
// @function format_date(date, format [,timezone])
func FormatDate(env utils.Environment, args ...interface{}) interface{} {
	if len(args) < 1 || len(args) > 3 {
		return fmt.Errorf("FORMAT_DATE takes one or two arguments, got %d", len(args))
	}
	date, err := utils.ToDate(env, args[0])
	if err != nil {
		return err
	}

	format := fmt.Sprintf("%s %s", env.DateFormat().String(), env.TimeFormat().String())
	if len(args) >= 2 {
		format, err = utils.ToString(env, args[1])
		if err != nil {
			return err
		}
	}

	// try to turn it to a go format
	goFormat, err := utils.ToGoDateFormat(format)
	if err != nil {
		return err
	}

	// grab our location
	location := env.Timezone()
	if len(args) == 3 {
		arg3, err := utils.ToString(env, args[2])
		if err != nil {
			return err
		}
		location, err = time.LoadLocation(arg3)
		if err != nil {
			return err
		}
	}

	// convert to our timezone if we have one (otherwise we remain in the date's default)
	if location != nil {
		date = date.In(location)
	}

	// return the formatted date
	return date.Format(goFormat)
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
func Date(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("DATE requires exactly one argument, got %d", len(args))
	}
	arg1, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	date, err := utils.DateFromString(env, arg1)
	if err != nil {
		return err
	}

	return date
}

// DateFromParts converts the passed in `year`, `month`` and `day`
//
//   @(date_from_parts(2017, 1, 15)) -> 2017-01-15T00:00:00.000000Z
//   @(date_from_parts(2017, 2, 31)) -> 2017-03-03T00:00:00.000000Z
//   @(date_from_parts(2017, 13, 15)) -> ERROR
//
// @function date_from_parts(year, month, day)
func DateFromParts(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 3 {
		return fmt.Errorf("DATE_FROM_PARTS requires three arguments, got %d", len(args))
	}
	year, err := utils.ToInt(env, args[0])
	if err != nil {
		return err
	}
	month, err := utils.ToInt(env, args[1])
	if err != nil {
		return err
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("Invalidate value for month, must be 1-12")
	}

	day, err := utils.ToInt(env, args[2])
	if err != nil {
		return err
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, env.Timezone())
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
func DateDiff(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 3 {
		return fmt.Errorf("DATE_DIFF takes exactly three arguments, received %d", len(args))
	}

	date1, err := utils.ToDate(env, args[0])
	if err != nil {
		return err
	}

	date2, err := utils.ToDate(env, args[1])
	if err != nil {
		return err
	}

	unit, err := utils.ToString(env, args[2])
	if err != nil {
		return err
	}

	// find the duration between our dates
	duration := date1.Sub(date2)

	// then convert based on our unit
	switch unit {

	case "s":
		return int(duration / time.Second)

	case "m":
		return int(duration / time.Minute)

	case "h":
		return int(duration / time.Hour)

	case "d":
		return utils.DaysBetween(date1, date2)

	case "w":
		return int(utils.DaysBetween(date1, date2) / 7)

	case "M":
		return utils.MonthsBetween(date1, date2)

	case "y":
		return date1.Year() - date2.Year()
	}

	return fmt.Errorf("Unknown unit: %s, must be one of s, m, h, D, W, M, Y", unit)
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
func DateAdd(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 3 {
		return fmt.Errorf("DATE_ADD takes exactly three arguments, received %d", len(args))
	}

	date, err := utils.ToDate(env, args[0])
	if err != nil {
		return err
	}

	duration, err := utils.ToInt(env, args[1])
	if err != nil {
		return err
	}

	unit, err := utils.ToString(env, args[2])
	if err != nil {
		return err
	}

	switch unit {

	case "s":
		return date.Add(time.Duration(duration) * time.Second)

	case "m":
		return date.Add(time.Duration(duration) * time.Minute)

	case "h":
		return date.Add(time.Duration(duration) * time.Hour)

	case "d":
		return date.AddDate(0, 0, duration)

	case "w":
		return date.AddDate(0, 0, duration*7)

	case "M":
		return date.AddDate(0, duration, 0)

	case "y":
		return date.AddDate(duration, 0, 0)
	}

	return fmt.Errorf("Unknown unit: %s, must be one of s, m, h, d, w, M, y", unit)
}

// Weekday returns the day of the week for `date`, 0 is sunday, 1 is monday..
//
//   @(weekday("2017-01-15")) -> 0
//   @(weekday("foo")) -> ERROR
//
// @function weekday(date)
func Weekday(env utils.Environment, args ...interface{}) interface{} {
	date, err := checkOneDateArg(env, "WEEKDAY", args)
	if err != nil {
		return err
	}

	return int(date.Weekday())
}

// TZ returns the timezone for `date``
//
// If not timezone information is present in the date, then the environment's
// timezone will be returned
//
//   @(tz("2017-01-15 02:15:18PM UTC")) -> "UTC"
//   @(tz("2017-01-15 02:15:18PM")) -> "UTC"
//   @(tz("2017-01-15")) -> "UTC"
//   @(tz("foo")) -> ERROR
//
// @function tz(date)
func TZ(env utils.Environment, args ...interface{}) interface{} {
	date, err := checkOneDateArg(env, "TZ", args)
	if err != nil {
		return err
	}

	return date.Location().String()
}

// TZOffset returns the offset for the timezone as a string +/- HHMM for `date`
//
// If no timezone information is present in the date, then the environment's
// timezone offset will be returned
//
//   @(tz_offset("2017-01-15 02:15:18PM UTC")) -> "+0000"
//   @(tz_offset("2017-01-15 02:15:18PM")) -> "+0000"
//   @(tz_offset("2017-01-15")) -> "+0000"
//   @(tz_offset("foo")) -> ERROR
//
// @function tz_offset(date)
func TZOffset(env utils.Environment, args ...interface{}) interface{} {
	date, err := checkOneDateArg(env, "TZ_OFFSET", args)
	if err != nil {
		return err
	}

	// this looks like we are returning a set offset, but this is how go describes formats
	return date.Format("-0700")

}

// Today returns the current date in the current timezone, time is set to midnight in the environment timezone
//
//  @(today()) -> 2017-01-20T00:00:00.000000Z
//
// @function today()
func Today(env utils.Environment, args ...interface{}) interface{} {
	if len(args) > 0 {
		return fmt.Errorf("TODAY takes no arguments, got %d", len(args))
	}

	nowTZ := time.Now().In(env.Timezone())
	return time.Date(nowTZ.Year(), nowTZ.Month(), nowTZ.Day(), 0, 0, 0, 0, env.Timezone())
}

// FromEpoch returns a new date created from `num` which represents number of nanoseconds since January 1st, 1970 GMT
//
//   @(from_epoch(1497286619000000000)) -> 2017-06-12T16:56:59.000000Z
//
// @function from_epoch(num)
func FromEpoch(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("FROM_EPOCH takes exactly one number argument, got %d", len(args))
	}

	offset, err := utils.ToDecimal(env, args[0])
	if err != nil {
		return err
	}

	return time.Unix(0, offset.IntPart()).In(env.Timezone())
}

// ToEpoch converts `date` to the number of nanoseconds since January 1st, 1970 GMT
//
//   @(to_epoch("2017-06-12T16:56:59.000000Z")) -> 1497286619000000000
//
// @function to_epoch(date)
func ToEpoch(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("TO_EPOCH takes exactly one date argument, got %d", len(args))
	}

	date, err := utils.ToDate(env, args[0])
	if err != nil {
		return err
	}

	return date.UnixNano()
}

// Now returns the current date and time in the environment timezone
//
//  @(now()) -> 2017-01-20T15:35:65.153654Z
//
// @function now()
func Now(env utils.Environment, args ...interface{}) interface{} {
	if len(args) > 0 {
		return fmt.Errorf("NOW takes no arguments, got %d", len(args))
	}

	return time.Now().In(env.Timezone())
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
//   @(format_urn(contact.urns.telegram)) -> ""
//   @(format_urn("NOT URN")) -> ERROR
//
// @function format_urn(urn)
func FormatURN(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("FORMAT_URN takes one argument, got %d", len(args))
	}

	// if we've been passed a slice like a URNList, use first item
	urnArg := args[0]
	if utils.IsSlice(urnArg) {
		sliceLen, _ := utils.SliceLength(urnArg)
		if sliceLen >= 1 {
			urnArg, _ = utils.LookupIndex(urnArg, 0)
		} else {
			return ""
		}
	}

	urnString, err := utils.ToString(env, urnArg)
	if err != nil {
		return err
	}

	urn := urns.URN(urnString)
	err = urn.Validate()
	if err != nil {
		return fmt.Errorf("%s is not a valid URN: %s", urnString, err)
	}

	return urn.Format()
}

//----------------------------------------------------------------------------------------
// Utility Functions
//----------------------------------------------------------------------------------------

func checkOneDecimalArg(env utils.Environment, funcName string, args []interface{}) (decimal.Decimal, error) {
	if len(args) != 1 {
		return decimal.Zero, fmt.Errorf("%s takes exactly one argument, got %d", funcName, len(args))
	}

	arg1, err := utils.ToDecimal(env, args[0])
	if err != nil {
		return decimal.Zero, err
	}

	return arg1, nil
}

func checkOneStringArg(env utils.Environment, funcName string, args []interface{}) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("%s takes exactly one argument, got %d", funcName, len(args))
	}

	arg1, err := utils.ToString(env, args[0])
	if err != nil {
		return "", err
	}

	return arg1, nil
}

func checkOneStringOneIntArg(env utils.Environment, funcName string, args []interface{}) (string, int, error) {
	if len(args) != 2 {
		return "", 0, fmt.Errorf("%s takes exactly two arguments, got %d", funcName, len(args))
	}

	arg1, err := utils.ToString(env, args[0])
	if err != nil {
		return "", 0, err
	}

	arg2, err := utils.ToInt(env, args[1])
	if err != nil {
		return "", 0, err
	}

	return arg1, arg2, err
}

func checkTwoDecimalArgs(env utils.Environment, funcName string, args []interface{}) (decimal.Decimal, decimal.Decimal, error) {
	if len(args) != 2 {
		return decimal.Zero, decimal.Zero, fmt.Errorf("%s takes exactly two arguments, got %d", funcName, len(args))
	}

	arg1, err := utils.ToDecimal(env, args[0])
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	arg2, err := utils.ToDecimal(env, args[1])
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	return arg1, arg2, nil
}

func checkOneDateArg(env utils.Environment, funcName string, args []interface{}) (time.Time, error) {
	if len(args) != 1 {
		return utils.ZeroTime, fmt.Errorf("%s takes exactly one argument, got %d", funcName, len(args))
	}

	arg1, err := utils.ToDate(env, args[0])
	if err != nil {
		return utils.ZeroTime, err
	}

	return arg1, err
}
