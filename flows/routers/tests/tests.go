package tests

import (
	"regexp"
	"strings"
	"time"

	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/nyaruka/phonenumbers"
	"github.com/shopspring/decimal"
)

// TODO:
// InterruptTest
// TimeoutTest
// AirtimeStatusTest

//------------------------------------------------------------------------------------------
// Mapping
//------------------------------------------------------------------------------------------

func init() {
	// register our router tests as Excellent functions
	for name, testFunc := range XTESTS {
		functions.RegisterXFunction(name, testFunc)
	}
}

// XTESTS is our mapping of the excellent test names to their actual functions
var XTESTS = map[string]functions.XFunction{
	"is_error":           IsError,
	"has_value":          HasValue,
	"has_group":          HasGroup,
	"has_wait_timed_out": HasWaitTimedOut,

	"is_string_eq":    functions.TwoStringFunction("is_string_eq", IsStringEQ),
	"has_phrase":      functions.TwoStringFunction("has_phrase", HasPhrase),
	"has_only_phrase": functions.TwoStringFunction("has_only_phrase", HasOnlyPhrase),
	"has_any_word":    functions.TwoStringFunction("has_any_word", HasAnyWord),
	"has_all_words":   functions.TwoStringFunction("has_all_words", HasAllWords),
	"has_beginning":   functions.TwoStringFunction("has_beginning", HasBeginning),
	"has_text":        functions.OneStringFunction("has_text", HasText),
	"has_pattern":     functions.TwoStringFunction("has_pattern", HasPattern),

	"has_number":         functions.OneStringFunction("has_number", HasNumber),
	"has_number_between": HasNumberBetween,
	"has_number_lt":      functions.StringAndNumberFunction("has_number_lt", HasNumberLT),
	"has_number_lte":     functions.StringAndNumberFunction("has_number_lte", HasNumberLTE),
	"has_number_eq":      functions.StringAndNumberFunction("has_number_eq", HasNumberEQ),
	"has_number_gte":     functions.StringAndNumberFunction("has_number_gte", HasNumberGTE),
	"has_number_gt":      functions.StringAndNumberFunction("has_number_gt", HasNumberGT),

	"has_date":    functions.OneStringFunction("has_date", HasDate),
	"has_date_lt": functions.StringAndDateFunction("has_date_lt", HasDateLT),
	"has_date_eq": functions.StringAndDateFunction("has_date_eq", HasDateEQ),
	"has_date_gt": functions.StringAndDateFunction("has_date_gt", HasDateGT),

	"has_phone": functions.TwoStringFunction("has_phone", HasPhone),
	"has_email": functions.OneStringFunction("has_email", HasEmail),

	"has_state":    functions.OneStringFunction("has_state", HasState),
	"has_district": HasDistrict,
	"has_ward":     HasWard,
}

//------------------------------------------------------------------------------------------
// Tests
//------------------------------------------------------------------------------------------

// IsStringEQ returns whether two strings are equal (case sensitive). In the case that they
// are, it will return the string as the match.
//
//   @(is_string_eq("foo", "foo")) -> true
//   @(is_string_eq("foo", "FOO")) -> false
//   @(is_string_eq("foo", "bar")) -> false
//   @(is_string_eq("foo", " foo ")) -> false
//   @(is_string_eq(run.status, "completed")) -> true
//   @(is_string_eq(run.webhook.status, "success")) -> true
//   @(is_string_eq(run.webhook.status, "connection_error")) -> false
//
// @test is_string_eq(run)
func IsStringEQ(env utils.Environment, str1 types.XString, str2 types.XString) types.XValue {
	if str1.Native() == str2.Native() {
		return XTestResult{true, str1}
	}

	return XFalseResult
}

// IsError returns whether `value` is an error
//
// Note that `contact.fields` and `run.results` are considered dynamic, so it is not an error
// to try to retrieve a value from fields or results which don't exist, rather these return an empty
// value.
//
//   @(is_error(date("foo"))) -> true
//   @(is_error(run.not.existing)) -> true
//   @(is_error(contact.fields.unset)) -> true
//   @(is_error("hello")) -> false
//
// @test is_error(value)
func IsError(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("IS_ERROR takes exactly one argument, got %d", len(args))
	}

	if types.IsXError(args[0]) {
		return XTestResult{true, args[0]}
	}

	return XFalseResult
}

// HasValue returns whether `value` is non-nil and not an error
//
// Note that `contact.fields` and `run.results` are considered dynamic, so it is not an error
// to try to retrieve a value from fields or results which don't exist, rather these return an empty
// value.
//
//   @(has_value(date("foo"))) -> false
//   @(has_value(not.existing)) -> false
//   @(has_value(contact.fields.unset)) -> false
//   @(has_value("hello")) -> true
//
// @test has_value(value)
func HasValue(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("HAS_VALUE takes exactly one argument, got %d", len(args))
	}

	// nil is not a value
	if utils.IsNil(args[0]) {
		return XFalseResult
	}

	// error is not a value
	_, isErr := args[0].(error)
	if isErr {
		return XFalseResult
	}

	return XTestResult{true, args[0]}
}

// HasWaitTimedOut returns whether the last wait timed out.
//
//   @(has_wait_timed_out(run)) -> false
//
// @test has_wait_timed_out(run)
func HasWaitTimedOut(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 {
		return types.NewXErrorf("HAS_WAIT_TIMED_OUT takes exactly one argument, got %d", len(args))
	}

	// first parameter needs to be a flow run
	run, isRun := args[0].(flows.FlowRun)
	if !isRun {
		return types.NewXErrorf("HAS_WAIT_TIMED_OUT must be called with a run as first argument")
	}

	if run.Session().Wait() != nil && run.Session().Wait().HasTimedOut() {
		return XTestResult{true, nil}
	}

	return XFalseResult
}

// HasGroup returns whether the `contact` is part of group with the passed in UUID
//
//   @(has_group(contact, "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d")) -> true
//   @(has_group(contact, "97fe7029-3a15-4005-b0c7-277b884fc1d5")) -> false
//
// @test has_group(contact)
func HasGroup(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 2 {
		return types.NewXErrorf("HAS_GROUP takes exactly two arguments, got %d", len(args))
	}

	// is the first argument a contact?
	contact, isContact := args[0].(*flows.Contact)
	if !isContact {
		return types.NewXErrorf("HAS_GROUP must have a contact as its first argument")
	}

	groupUUID, xerr := types.ToXString(args[1])
	if xerr != nil {
		return xerr
	}

	// iterate through the groups looking for one with the same UUID as passed in
	group := contact.Groups().FindByUUID(flows.GroupUUID(groupUUID))
	if group != nil {
		return XTestResult{true, group}
	}

	return XFalseResult
}

// HasPhrase tests whether `phrase` is contained in `string`
//
// The words in the test phrase must appear in the same order with no other words
// in between.
//
//   @(has_phrase("the quick brown fox", "brown fox")) -> true
//   @(has_phrase("the Quick Brown fox", "quick fox")) -> false
//   @(has_phrase("the Quick Brown fox", "")) -> true
//   @(has_phrase("the.quick.brown.fox", "the quick").match) -> "the quick"
//
// @test has_phrase(string, phrase)
func HasPhrase(env utils.Environment, str types.XString, test types.XString) types.XValue {
	return testStringTokens(env, "HAS_PHRASE", str, test, hasPhraseTest)
}

// HasAllWords tests whether all the `words` are contained in `string`
//
// The words can be in any order and may appear more than once.
//
//   @(has_all_words("the quick brown FOX", "the fox")) -> true
//   @(has_all_words("the quick brown FOX", "the fox").match) -> "the FOX"
//   @(has_all_words("the quick brown fox", "red fox")) -> false
//
// @test has_all_words(string, words)
func HasAllWords(env utils.Environment, str types.XString, test types.XString) types.XValue {
	return testStringTokens(env, "HAS_ALL_WORDS", str, test, hasAllWordsTest)
}

// HasAnyWord tests whether any of the `words` are contained in the `string`
//
// Only one of the words needs to match and it may appear more than once.
//
//   @(has_any_word("The Quick Brown Fox", "fox quick")) -> true
//   @(has_any_word("The Quick Brown Fox", "red fox")) -> true
//   @(has_any_word("The Quick Brown Fox", "red fox").match) -> "Fox"
//
// @test has_any_word(string, words)
func HasAnyWord(env utils.Environment, str types.XString, test types.XString) types.XValue {
	return testStringTokens(env, "HAS_ANY_WORD", str, test, hasAnyWordTest)
}

// HasOnlyPhrase tests whether the `string` contains only `phrase`
//
// The phrase must be the only text in the string to match
//
//   @(has_only_phrase("The Quick Brown Fox", "quick brown")) -> false
//   @(has_only_phrase("Quick Brown", "quick brown")) -> true
//   @(has_only_phrase("the Quick Brown fox", "")) -> false
//   @(has_only_phrase("", "")) -> true
//   @(has_only_phrase("Quick Brown", "quick brown").match) -> "Quick Brown"
//   @(has_only_phrase("The Quick Brown Fox", "red fox")) -> false
//
// @test has_only_phrase(string, phrase)
func HasOnlyPhrase(env utils.Environment, str types.XString, test types.XString) types.XValue {
	return testStringTokens(env, "HAS_ONLY_PHRASE", str, test, hasOnlyPhraseTest)
}

// HasText tests whether there the string has any characters in it
//
//   @(has_text("quick brown")) -> true
//   @(has_text("quick brown").match) -> "quick brown"
//   @(has_text("")) -> false
//   @(has_text(" \n")) -> false
//   @(has_text(123)) -> true
//
// @test has_text(string)
func HasText(env utils.Environment, str types.XString) types.XValue {
	// trim any whitespace
	str = types.NewXString(strings.TrimSpace(str.Native()))

	// if there is anything left then we have text
	if len(str) > 0 {
		return XTestResult{true, str}
	}

	return XFalseResult
}

// HasBeginning tests whether `string` starts with `beginning`
//
// Both strings are trimmed of surrounding whitespace, but otherwise matching is strict
// without any tokenization.
//
//   @(has_beginning("The Quick Brown", "the quick")) -> true
//   @(has_beginning("The Quick Brown", "the quick").match) -> "The Quick"
//   @(has_beginning("The Quick Brown", "the   quick")) -> false
//   @(has_beginning("The Quick Brown", "quick brown")) -> false
//
// @test has_beginning(string, beginning)
func HasBeginning(env utils.Environment, str1 types.XString, str2 types.XString) types.XValue {
	// trim both
	hayStack := strings.TrimSpace(str1.Native())
	pinCushion := strings.TrimSpace(str2.Native())

	// either are empty, no match
	if hayStack == "" || pinCushion == "" {
		return XFalseResult
	}

	// haystack has to be at least length of needle
	if len(hayStack) < len(pinCushion) {
		return XFalseResult
	}

	segment := hayStack[:len(pinCushion)]
	if strings.ToLower(segment) == strings.ToLower(pinCushion) {
		return XTestResult{true, types.NewXString(segment)}
	}

	return XFalseResult
}

// Returned by the has_pattern test as its match value
type patternMatch struct {
	groups types.XArray
}

func newPatternMatch(matches []string) *patternMatch {
	groups := types.NewXArray()
	for _, match := range matches {
		groups.Append(types.NewXString(match))
	}
	return &patternMatch{groups: groups}
}

// Resolve resolves the given key when this match is referenced in an expression
func (m *patternMatch) Resolve(key string) types.XValue {
	switch key {
	case "groups":
		return m.groups
	}

	return types.NewXResolveError(m, key)
}

// Reduce is called when this object needs to be reduced to a primitive
func (m *patternMatch) Reduce() types.XPrimitive {
	return m.groups.Index(0).(types.XString)
}

func (m *patternMatch) ToJSON() types.XString { return types.NewXString("TODO") }

var _ types.XValue = (*patternMatch)(nil)
var _ types.XResolvable = (*patternMatch)(nil)

// HasPattern tests whether `string` matches the regex `pattern`
//
// Both strings are trimmed of surrounding whitespace and matching is case-insensitive.
//
//   @(has_pattern("Sell cheese please", "buy (\w+)")) -> false
//   @(has_pattern("Buy cheese please", "buy (\w+)")) -> true
//   @(has_pattern("Buy cheese please", "buy (\w+)").match) -> "Buy cheese"
//   @(has_pattern("Buy cheese please", "buy (\w+)").match.groups[0]) -> "Buy cheese"
//   @(has_pattern("Buy cheese please", "buy (\w+)").match.groups[1]) -> "cheese"
//
// @test has_pattern(string, pattern)
func HasPattern(env utils.Environment, haystack types.XString, pattern types.XString) types.XValue {
	regex, err := regexp.Compile("(?i)" + strings.TrimSpace(pattern.Native()))
	if err != nil {
		return types.NewXErrorf("HAS_PATTERN must be called with a valid regular expression")
	}

	matches := regex.FindStringSubmatch(strings.TrimSpace(haystack.Native()))
	if matches != nil {
		return XTestResult{true, newPatternMatch(matches)}
	}

	return XFalseResult
}

// HasNumber tests whether `string` contains a number
//
//   @(has_number("the number is 42")) -> true
//   @(has_number("the number is 42").match) -> 42
//   @(has_number("the number is forty two")) -> false
//
// @test has_number(string)
func HasNumber(env utils.Environment, str types.XString) types.XValue {
	return testNumber(env, "HAS_NUMBER", str, types.XNumberZero, isNumberTest)
}

// HasNumberBetween tests whether `string` contains a number between `min` and `max` inclusive
//
//   @(has_number_between("the number is 42", 40, 44)) -> true
//   @(has_number_between("the number is 42", 40, 44).match) -> 42
//   @(has_number_between("the number is 42", 50, 60)) -> false
//   @(has_number_between("the number is not there", 50, 60)) -> false
//   @(has_number_between("the number is not there", "foo", 60)) -> ERROR
//
// @test has_number_between(string, min, max)
func HasNumberBetween(env utils.Environment, args ...types.XValue) types.XValue {
	// need three arguments, value being tested and min, max
	if len(args) != 3 {
		return types.NewXErrorf("HAS_NUMBER_BETWEEN takes exactly three arguments, got %d", len(args))
	}

	str, xerr := types.ToXString(args[0])
	if xerr != nil {
		return xerr
	}
	min, xerr := types.ToXNumber(args[1])
	if xerr != nil {
		return xerr
	}
	max, xerr := types.ToXNumber(args[2])
	if xerr != nil {
		return xerr
	}

	// for each of our values, try to evaluate to a decimal
	for _, value := range strings.Fields(str.Native()) {
		num, xerr := types.ToXNumber(types.NewXString(value))
		if xerr == nil {
			if num.Compare(min) >= 0 && num.Compare(max) <= 0 {
				return XTestResult{true, num}
			}
		}
	}
	return XFalseResult
}

// HasNumberLT tests whether `string` contains a number less than `max`
//
//   @(has_number_lt("the number is 42", 44)) -> true
//   @(has_number_lt("the number is 42", 44).match) -> 42
//   @(has_number_lt("the number is 42", 40)) -> false
//   @(has_number_lt("the number is not there", 40)) -> false
//   @(has_number_lt("the number is not there", "foo")) -> ERROR
//
// @test has_number_lt(string, max)
func HasNumberLT(env utils.Environment, str types.XString, num types.XNumber) types.XValue {
	return testNumber(env, "HAS_NUMBER_LT", str, num, isNumberLT)
}

// HasNumberLTE tests whether `value` contains a number less than or equal to `max`
//
//   @(has_number_lte("the number is 42", 42)) -> true
//   @(has_number_lte("the number is 42", 44).match) -> 42
//   @(has_number_lte("the number is 42", 40)) -> false
//   @(has_number_lte("the number is not there", 40)) -> false
//   @(has_number_lte("the number is not there", "foo")) -> ERROR
//
// @test has_number_lte(string, max)
func HasNumberLTE(env utils.Environment, str types.XString, num types.XNumber) types.XValue {
	return testNumber(env, "HAS_NUMBER_LTE", str, num, isNumberLTE)
}

// HasNumberEQ tests whether `strung` contains a number equal to the `value`
//
//   @(has_number_eq("the number is 42", 42)) -> true
//   @(has_number_eq("the number is 42", 42).match) -> 42
//   @(has_number_eq("the number is 42", 40)) -> false
//   @(has_number_eq("the number is not there", 40)) -> false
//   @(has_number_eq("the number is not there", "foo")) -> ERROR
//
// @test has_number_eq(string, value)
func HasNumberEQ(env utils.Environment, str types.XString, num types.XNumber) types.XValue {
	return testNumber(env, "HAS_NUMBER_EQ", str, num, isNumberEQ)
}

// HasNumberGTE tests whether `string` contains a number greater than or equal to `min`
//
//   @(has_number_gte("the number is 42", 42)) -> true
//   @(has_number_gte("the number is 42", 42).match) -> 42
//   @(has_number_gte("the number is 42", 45)) -> false
//   @(has_number_gte("the number is not there", 40)) -> false
//   @(has_number_gte("the number is not there", "foo")) -> ERROR
//
// @test has_number_gte(string, min)
func HasNumberGTE(env utils.Environment, str types.XString, num types.XNumber) types.XValue {
	return testNumber(env, "HAS_NUMBER_GTE", str, num, isNumberGTE)
}

// HasNumberGT tests whether `string` contains a number greater than `min`
//
//   @(has_number_gt("the number is 42", 40)) -> true
//   @(has_number_gt("the number is 42", 40).match) -> 42
//   @(has_number_gt("the number is 42", 42)) -> false
//   @(has_number_gt("the number is not there", 40)) -> false
//   @(has_number_gt("the number is not there", "foo")) -> ERROR
//
// @test has_number_gt(string, min)
func HasNumberGT(env utils.Environment, str types.XString, num types.XNumber) types.XValue {
	return testNumber(env, "HAS_NUMBER_GT", str, num, isNumberGT)
}

// HasDate tests whether `string` contains a date formatted according to our environment
//
//   @(has_date("the date is 2017-01-15")) -> true
//   @(has_date("the date is 2017-01-15").match) -> 2017-01-15T00:00:00.000000Z
//   @(has_date("there is no date here, just a year 2017")) -> false
//
// @test has_date(string)
func HasDate(env utils.Environment, str types.XString) types.XValue {
	return testDate(env, "HAS_DATE", str, types.XDateZero, isDateTest)
}

// HasDateLT tests whether `value` contains a date before the date `max`
//
//   @(has_date_lt("the date is 2017-01-15", "2017-06-01")) -> true
//   @(has_date_lt("the date is 2017-01-15", "2017-06-01").match) -> 2017-01-15T00:00:00.000000Z
//   @(has_date_lt("there is no date here, just a year 2017", "2017-06-01")) -> false
//   @(has_date_lt("there is no date here, just a year 2017", "not date")) -> ERROR
//
// @test has_date_lt(string, max)
func HasDateLT(env utils.Environment, str types.XString, date types.XDate) types.XValue {
	return testDate(env, "HAS_DATE_LT", str, date, isDateLTTest)
}

// HasDateEQ tests whether `string` a date equal to `date`
//
//   @(has_date_eq("the date is 2017-01-15", "2017-01-15")) -> true
//   @(has_date_eq("the date is 2017-01-15", "2017-01-15").match) -> 2017-01-15T00:00:00.000000Z
//   @(has_date_eq("the date is 2017-01-15 15:00", "2017-01-15")) -> false
//   @(has_date_eq("there is no date here, just a year 2017", "2017-06-01")) -> false
//   @(has_date_eq("there is no date here, just a year 2017", "not date")) -> ERROR
//
// @test has_date_eq(string, date)
func HasDateEQ(env utils.Environment, str types.XString, date types.XDate) types.XValue {
	return testDate(env, "HAS_DATE_EQ", str, date, isDateEQTest)
}

// HasDateGT tests whether `string` a date after the date `min`
//
//   @(has_date_gt("the date is 2017-01-15", "2017-01-01")) -> true
//   @(has_date_gt("the date is 2017-01-15", "2017-01-01").match) -> 2017-01-15T00:00:00.000000Z
//   @(has_date_gt("the date is 2017-01-15", "2017-03-15")) -> false
//   @(has_date_gt("there is no date here, just a year 2017", "2017-06-01")) -> false
//   @(has_date_gt("there is no date here, just a year 2017", "not date")) -> ERROR
//
// @test has_date_gt(string, min)
func HasDateGT(env utils.Environment, str types.XString, date types.XDate) types.XValue {
	return testDate(env, "HAS_DATE_GT", str, date, isDateGTTest)
}

var emailAddressRE = regexp.MustCompile(`([\pL\pN][-_.\pL\pN]*)@([\pL\pN][-_\pL\pN]*)(\.[\pL\pN][-_\pL\pN]*)+`)

// HasEmail tests whether an email is contained in `string`
//
//   @(has_email("my email is foo1@bar.com, please respond")) -> true
//   @(has_email("my email is foo1@bar.com, please respond").match) -> "foo1@bar.com"
//   @(has_email("my email is <foo@bar2.com>")) -> true
//   @(has_email("i'm not sharing my email")) -> false
//
// @test has_email(string)
func HasEmail(env utils.Environment, str types.XString) types.XValue {
	// split by whitespace
	email := emailAddressRE.FindString(str.Native())
	if email != "" {
		return XTestResult{true, types.NewXString(email)}
	}

	return XFalseResult
}

// HasPhone tests whether a phone number (in the passed in `country_code`) is contained in the `string`
//
//   @(has_phone("my number is 2067799294", "US")) -> true
//   @(has_phone("my number is 206 779 9294", "US").match) -> "+12067799294"
//   @(has_phone("my number is none of your business", "US")) -> false
//
// @test has_phone(string, country_code)
func HasPhone(env utils.Environment, str types.XString, country types.XString) types.XValue {
	// try to find a phone number
	phone, err := phonenumbers.Parse(str.Native(), country.Native())
	if err != nil {
		return XFalseResult
	}

	// format as E164 number
	formatted := phonenumbers.Format(phone, phonenumbers.E164)
	return XTestResult{true, types.NewXString(formatted)}
}

// HasState tests whether a state name is contained in the `string`
//
//   @(has_state("Kigali")) -> true
//   @(has_state("Boston")) -> false
//   @(has_state("¡Kigali!")) -> true
//   @(has_state("I live in Kigali")) -> true
//
// @test has_state(string)
func HasState(env utils.Environment, str types.XString) types.XValue {
	runEnv, _ := env.(flows.RunEnvironment)

	states, err := runEnv.FindLocationsFuzzy(str.Native(), flows.LocationLevel(1), nil)
	if err != nil {
		return types.NewXError(err)
	}
	if len(states) > 0 {
		return XTestResult{true, states[0]}
	}
	return XFalseResult
}

// HasDistrict tests whether a district name is contained in the `string`. If `state` is also provided
// then the returned district must be within that state.
//
//   @(has_district("Gasabo", "Kigali")) -> true
//   @(has_district("I live in Gasabo", "Kigali")) -> true
//   @(has_district("Gasabo", "Boston")) -> false
//   @(has_district("Gasabo")) -> true
//
// @test has_district(string, state)
func HasDistrict(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 && len(args) != 2 {
		return types.NewXErrorf("HAS_DISTRICT takes one or two arguments, got %d", len(args))
	}

	runEnv, _ := env.(flows.RunEnvironment)

	var text, stateText types.XString
	var xerr types.XError

	// grab the text we will search and the parent state name
	if text, xerr = types.ToXString(args[0]); xerr != nil {
		return xerr
	}
	if len(args) == 2 {
		if stateText, xerr = types.ToXString(args[1]); xerr != nil {
			return xerr
		}
	}

	states, err := runEnv.FindLocationsFuzzy(stateText.Native(), flows.LocationLevel(1), nil)
	if err != nil {
		return types.NewXError(err)
	}
	if len(states) > 0 {
		districts, err := runEnv.FindLocationsFuzzy(text.Native(), flows.LocationLevel(2), states[0])
		if err != nil {
			return types.NewXError(err)
		}
		if len(districts) > 0 {
			return XTestResult{true, districts[0]}
		}
	}

	// try without a parent state - it's ok as long as we get a single match
	if stateText == "" {
		districts, err := runEnv.FindLocationsFuzzy(text.Native(), flows.LocationLevel(2), nil)
		if err != nil {
			return types.NewXError(err)
		}
		if len(districts) == 1 {
			return XTestResult{true, districts[0]}
		}
	}

	return XFalseResult
}

// HasWard tests whether a ward name is contained in the `string`
//
//   @(has_ward("Gisozi", "Gasabo", "Kigali")) -> true
//   @(has_ward("I live in Gisozi", "Gasabo", "Kigali")) -> true
//   @(has_ward("Gisozi", "Gasabo", "Brooklyn")) -> false
//   @(has_ward("Gisozi", "Brooklyn", "Kigali")) -> false
//   @(has_ward("Brooklyn", "Gasabo", "Kigali")) -> false
//   @(has_ward("Gasabo")) -> false
//   @(has_ward("Gisozi")) -> true
//
// @test has_ward(string, district, state)
func HasWard(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 && len(args) != 3 {
		return types.NewXErrorf("HAS_WARD takes one or three arguments, got %d", len(args))
	}

	runEnv, _ := env.(flows.RunEnvironment)

	var text, districtText, stateText types.XString
	var xerr types.XError

	// grab the text we will search, as well as the parent district and state names
	if text, xerr = types.ToXString(args[0]); xerr != nil {
		return xerr
	}
	if len(args) == 3 {
		if districtText, xerr = types.ToXString(args[1]); xerr != nil {
			return xerr
		}
		if stateText, xerr = types.ToXString(args[2]); xerr != nil {
			return xerr
		}
	}

	states, err := runEnv.FindLocationsFuzzy(stateText.Native(), flows.LocationLevel(1), nil)
	if err != nil {
		return types.NewXError(err)
	}
	if len(states) > 0 {
		districts, err := runEnv.FindLocationsFuzzy(districtText.Native(), flows.LocationLevel(2), states[0])
		if err != nil {
			return types.NewXError(err)
		}
		if len(districts) > 0 {
			wards, err := runEnv.FindLocationsFuzzy(text.Native(), flows.LocationLevel(3), districts[0])
			if err != nil {
				return types.NewXError(err)
			}
			if len(wards) > 0 {
				return XTestResult{true, wards[0]}
			}
		}
	}

	// try without a parent district - it's ok as long as we get a single match
	if districtText == "" {
		wards, err := runEnv.FindLocationsFuzzy(text.Native(), flows.LocationLevel(3), nil)
		if err != nil {
			return types.NewXError(err)
		}
		if len(wards) == 1 {
			return XTestResult{true, wards[0]}
		}
	}

	return XFalseResult
}

//------------------------------------------------------------------------------------------
// String Test Functions
//------------------------------------------------------------------------------------------

type stringTokenTest func(origHayTokens []string, hayTokens []string, pinTokens []string) XTestResult

func testStringTokens(env utils.Environment, name string, str types.XString, testStr types.XString, testFunc stringTokenTest) types.XValue {
	hayStack := strings.TrimSpace(str.Native())
	needle := strings.TrimSpace(testStr.Native())

	origHays := utils.TokenizeString(hayStack)
	hays := utils.TokenizeString(strings.ToLower(hayStack))
	needles := utils.TokenizeString(strings.ToLower(needle))

	return testFunc(origHays, hays, needles)
}

func hasPhraseTest(origHays []string, hays []string, pins []string) XTestResult {
	if len(pins) == 0 {
		return XTestResult{true, types.XStringEmpty}
	}

	pinIdx := 0
	matches := make([]string, len(pins))
	for i, hay := range hays {
		if hay == pins[pinIdx] {
			matches[pinIdx] = origHays[i]
			pinIdx++
			if pinIdx == len(pins) {
				break
			}
		} else {
			pinIdx = 0
		}
	}

	if pinIdx == len(pins) {
		return XTestResult{true, types.NewXString(strings.Join(matches, " "))}
	}

	return XFalseResult
}

func hasAllWordsTest(origHays []string, hays []string, pins []string) XTestResult {
	matches := make([]string, 0, len(pins))
	pinMatches := make([]int, len(pins))

	for i, hay := range hays {
		matched := false
		for j, pin := range pins {
			if hay == pin {
				matched = true
				pinMatches[j]++
			}
		}

		if matched {
			matches = append(matches, origHays[i])
		}
	}

	allMatch := true
	for _, matchCount := range pinMatches {
		if matchCount == 0 {
			allMatch = false
			break
		}

	}

	if allMatch {
		return XTestResult{true, types.NewXString(strings.Join(matches, " "))}
	}

	return XFalseResult
}

func hasAnyWordTest(origHays []string, hays []string, pins []string) XTestResult {
	matches := make([]string, 0, len(pins))
	for i, hay := range hays {
		matched := false
		for _, pin := range pins {
			if hay == pin {
				matched = true
				break
			}
		}
		if matched {
			matches = append(matches, origHays[i])
		}

	}

	if len(matches) > 0 {
		return XTestResult{true, types.NewXString(strings.Join(matches, " "))}
	}

	return XFalseResult
}

func hasOnlyPhraseTest(origHays []string, hays []string, pins []string) XTestResult {
	// must be same length
	if len(hays) != len(pins) {
		return XFalseResult
	}

	// and every token must match
	matches := make([]string, 0, len(pins))
	for i := range hays {
		if hays[i] != pins[i] {
			return XFalseResult
		}
		matches = append(matches, origHays[i])
	}

	return XTestResult{true, types.NewXString(strings.Join(matches, " "))}
}

//------------------------------------------------------------------------------------------
// Numerical Test Functions
//------------------------------------------------------------------------------------------

type decimalTest func(value decimal.Decimal, test decimal.Decimal) bool

func testNumber(env utils.Environment, name string, str types.XString, testNum types.XNumber, testFunc decimalTest) types.XValue {
	// for each of our values, try to evaluate to a decimal
	for _, value := range strings.Fields(str.Native()) {
		num, xerr := types.ToXNumber(types.NewXString(value))
		if xerr == nil {
			if testFunc(num.Native(), testNum.Native()) {
				return XTestResult{true, num}
			}
		}
	}

	return XFalseResult
}

func isNumberTest(value decimal.Decimal, test decimal.Decimal) bool {
	return true
}

func isNumberLT(value decimal.Decimal, test decimal.Decimal) bool {
	return value.Cmp(test) < 0
}

func isNumberLTE(value decimal.Decimal, test decimal.Decimal) bool {
	return value.Cmp(test) <= 0
}

func isNumberEQ(value decimal.Decimal, test decimal.Decimal) bool {
	return value.Cmp(test) == 0
}

func isNumberGTE(value decimal.Decimal, test decimal.Decimal) bool {
	return value.Cmp(test) >= 0
}

func isNumberGT(value decimal.Decimal, test decimal.Decimal) bool {
	return value.Cmp(test) > 0
}

//------------------------------------------------------------------------------------------
// Date Test Functions
//------------------------------------------------------------------------------------------

type dateTest func(value time.Time, test time.Time) bool

func testDate(env utils.Environment, name string, str types.XString, testDate types.XDate, testFunc dateTest) types.XValue {
	// error is if we don't find a date on our test value, that's ok but no match
	value, xerr := types.ToXDate(env, str)
	if xerr != nil {
		return XFalseResult
	}

	if testFunc(value.Native(), testDate.Native()) {
		return XTestResult{true, value}
	}

	return XFalseResult
}

func isDateTest(value time.Time, test time.Time) bool {
	return true
}

func isDateLTTest(value time.Time, test time.Time) bool {
	return value.Before(test)
}

func isDateEQTest(value time.Time, test time.Time) bool {
	return value.Equal(test)
}

func isDateGTTest(value time.Time, test time.Time) bool {
	return value.After(test)
}
