package tests

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/nyaruka/phonenumbers"
	"github.com/shopspring/decimal"
)

//------------------------------------------------------------------------------------------
// Mapping
//------------------------------------------------------------------------------------------

func init() {
	// register our router tests as Excellent functions
	for name, testFunc := range XTESTS {
		functions.RegisterXFunction(name, testFunc)
	}
}

// RegisterXTest registers a new router test (and Excellent function)
func RegisterXTest(name string, function types.XFunction) {
	XTESTS[name] = function
	functions.RegisterXFunction(name, function)
}

// XTESTS is our mapping of the excellent test names to their actual functions
var XTESTS = map[string]types.XFunction{
	"has_error": functions.OneArgFunction(HasError),
	"has_value": functions.OneArgFunction(HasValue),

	"is_text_eq":      functions.TwoTextFunction(IsTextEQ),
	"has_phrase":      functions.TwoTextFunction(HasPhrase),
	"has_only_phrase": functions.TwoTextFunction(HasOnlyPhrase),
	"has_any_word":    functions.TwoTextFunction(HasAnyWord),
	"has_all_words":   functions.TwoTextFunction(HasAllWords),
	"has_beginning":   functions.TwoTextFunction(HasBeginning),
	"has_text":        functions.OneTextFunction(HasText),
	"has_pattern":     functions.TwoTextFunction(HasPattern),

	"has_number":         functions.OneTextFunction(HasNumber),
	"has_number_between": functions.ThreeArgFunction(HasNumberBetween),
	"has_number_lt":      functions.TextAndNumberFunction(HasNumberLT),
	"has_number_lte":     functions.TextAndNumberFunction(HasNumberLTE),
	"has_number_eq":      functions.TextAndNumberFunction(HasNumberEQ),
	"has_number_gte":     functions.TextAndNumberFunction(HasNumberGTE),
	"has_number_gt":      functions.TextAndNumberFunction(HasNumberGT),

	"has_date":    functions.OneTextFunction(HasDate),
	"has_date_lt": functions.TextAndDateFunction(HasDateLT),
	"has_date_eq": functions.TextAndDateFunction(HasDateEQ),
	"has_date_gt": functions.TextAndDateFunction(HasDateGT),

	"has_time":  functions.OneTextFunction(HasTime),
	"has_phone": functions.InitialTextFunction(0, 1, HasPhone),
	"has_email": functions.OneTextFunction(HasEmail),
	"has_group": functions.TwoArgFunction(HasGroup),

	"has_state":    functions.OneTextFunction(HasState),
	"has_district": HasDistrict,
	"has_ward":     HasWard,
}

//------------------------------------------------------------------------------------------
// Results
//------------------------------------------------------------------------------------------

// NewTrueResult creates a new true result with a match
func NewTrueResult(match types.XValue) *types.XDict {
	return types.NewXDict(map[string]types.XValue{"match": match})
}

// NewTrueResultWithExtra creates a new true result with a match and extra
func NewTrueResultWithExtra(match types.XValue, extra *types.XDict) *types.XDict {
	return types.NewXDict(map[string]types.XValue{"match": match, "extra": extra})
}

//------------------------------------------------------------------------------------------
// Tests
//------------------------------------------------------------------------------------------

// IsTextEQ returns whether two text values are equal (case sensitive). In the case that they
// are, it will return the text as the match.
//
//   @(is_text_eq("foo", "foo")) -> {match: foo}
//   @(is_text_eq("foo", "FOO")) ->
//   @(is_text_eq("foo", "bar")) ->
//   @(is_text_eq("foo", " foo ")) ->
//   @(is_text_eq(run.status, "completed")) -> {match: completed}
//   @(is_text_eq(results.webhook.category, "Success")) -> {match: Success}
//   @(is_text_eq(results.webhook.category, "Failure")) ->
//
// @test is_text_eq(text1, text2)
func IsTextEQ(env utils.Environment, text1 types.XText, text2 types.XText) types.XValue {
	if text1.Equals(text2) {
		return NewTrueResult(text1)
	}

	return nil
}

// HasError returns whether `value` is an error
//
//   @(has_error(datetime("foo"))) -> {match: error calling DATETIME: unable to convert "foo" to a datetime}
//   @(has_error(run.not.existing)) -> {match: dict has no property 'not'}
//   @(has_error(contact.fields.unset)) -> {match: dict has no property 'unset'}
//   @(has_error("hello")) ->
//
// @test has_error(value)
func HasError(env utils.Environment, value types.XValue) types.XValue {
	if types.IsXError(value) {
		return NewTrueResult(value)
	}

	return nil
}

// HasValue returns whether `value` is non-nil and not an error
//
// Note that `contact.fields` and `run.results` are considered dynamic, so it is not an error
// to try to retrieve a value from fields or results which don't exist, rather these return an empty
// value.
//
//   @(has_value(datetime("foo"))) ->
//   @(has_value(not.existing)) ->
//   @(has_value(contact.fields.unset)) ->
//   @(has_value("")) ->
//   @(has_value("hello")) -> {match: hello}
//
// @test has_value(value)
func HasValue(env utils.Environment, value types.XValue) types.XValue {
	if types.IsEmpty(value) || types.IsXError(value) {
		return nil
	}

	return NewTrueResult(value)
}

// HasGroup returns whether the `contact` is part of group with the passed in UUID
//
//   @(has_group(contact.groups, "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d")) -> {match: {name: Testers, uuid: b7cf0d83-f1c9-411c-96fd-c511a4cfa86d}}
//   @(has_group(array(), "97fe7029-3a15-4005-b0c7-277b884fc1d5")) ->
//
// @test has_group(contact, group_uuid)
func HasGroup(env utils.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
	// is the first argument an array
	array, xerr := types.ToXArray(env, arg1)
	if xerr != nil {
		return xerr
	}

	groupUUID, xerr := types.ToXText(env, arg2)
	if xerr != nil {
		return xerr
	}

	for i := 0; i < array.Length(); i++ {
		group, xerr := types.ToXDict(env, array.Get(i))
		if xerr != nil {
			return xerr
		}

		uuidValue, _ := group.Get("uuid")
		uuid, xerr := types.ToXText(env, uuidValue)
		if xerr != nil {
			return xerr
		}

		if uuid.Equals(groupUUID) {
			return NewTrueResult(group)
		}
	}

	return nil
}

// HasPhrase tests whether `phrase` is contained in `text`
//
// The words in the test phrase must appear in the same order with no other words
// in between.
//
//   @(has_phrase("the quick brown fox", "brown fox")) -> {match: brown fox}
//   @(has_phrase("the Quick Brown fox", "quick fox")) ->
//   @(has_phrase("the Quick Brown fox", "")) -> {match: }
//
// @test has_phrase(text, phrase)
func HasPhrase(env utils.Environment, text types.XText, test types.XText) types.XValue {
	return testStringTokens(env, text, test, hasPhraseTest)
}

// HasAllWords tests whether all the `words` are contained in `text`
//
// The words can be in any order and may appear more than once.
//
//   @(has_all_words("the quick brown FOX", "the fox")) -> {match: the FOX}
//   @(has_all_words("the quick brown fox", "red fox")) ->
//
// @test has_all_words(text, words)
func HasAllWords(env utils.Environment, text types.XText, test types.XText) types.XValue {
	return testStringTokens(env, text, test, hasAllWordsTest)
}

// HasAnyWord tests whether any of the `words` are contained in the `text`
//
// Only one of the words needs to match and it may appear more than once.
//
//   @(has_any_word("The Quick Brown Fox", "fox quick")) -> {match: Quick Fox}
//   @(has_any_word("The Quick Brown Fox", "red fox")) -> {match: Fox}
//
// @test has_any_word(text, words)
func HasAnyWord(env utils.Environment, text types.XText, test types.XText) types.XValue {
	return testStringTokens(env, text, test, hasAnyWordTest)
}

// HasOnlyPhrase tests whether the `text` contains only `phrase`
//
// The phrase must be the only text in the text to match
//
//   @(has_only_phrase("The Quick Brown Fox", "quick brown")) ->
//   @(has_only_phrase("Quick Brown", "quick brown")) -> {match: Quick Brown}
//   @(has_only_phrase("the Quick Brown fox", "")) ->
//   @(has_only_phrase("", "")) -> {match: }
//   @(has_only_phrase("The Quick Brown Fox", "red fox")) ->
//
// @test has_only_phrase(text, phrase)
func HasOnlyPhrase(env utils.Environment, text types.XText, test types.XText) types.XValue {
	return testStringTokens(env, text, test, hasOnlyPhraseTest)
}

// HasText tests whether there the text has any characters in it
//
//   @(has_text("quick brown")) -> {match: quick brown}
//   @(has_text("")) ->
//   @(has_text(" \n")) ->
//   @(has_text(123)) -> {match: 123}
//   @(has_text(contact.fields.not_set)) ->
//
// @test has_text(text)
func HasText(env utils.Environment, text types.XText) types.XValue {
	// trim any whitespace
	text = types.NewXText(strings.TrimSpace(text.Native()))

	// if there is anything left then we have text
	if text.Length() > 0 {
		return NewTrueResult(text)
	}

	return nil
}

// HasBeginning tests whether `text` starts with `beginning`
//
// Both text values are trimmed of surrounding whitespace, but otherwise matching is strict
// without any tokenization.
//
//   @(has_beginning("The Quick Brown", "the quick")) -> {match: The Quick}
//   @(has_beginning("The Quick Brown", "the   quick")) ->
//   @(has_beginning("The Quick Brown", "quick brown")) ->
//
// @test has_beginning(text, beginning)
func HasBeginning(env utils.Environment, text types.XText, beginning types.XText) types.XValue {
	// trim both
	hayStack := strings.TrimSpace(text.Native())
	pinCushion := strings.TrimSpace(beginning.Native())

	// either are empty, no match
	if hayStack == "" || pinCushion == "" {
		return nil
	}

	// haystack has to be at least length of needle
	if len(hayStack) < len(pinCushion) {
		return nil
	}

	segment := hayStack[:len(pinCushion)]
	if strings.ToLower(segment) == strings.ToLower(pinCushion) {
		return NewTrueResult(types.NewXText(segment))
	}

	return nil
}

// HasPattern tests whether `text` matches the regex `pattern`
//
// Both text values are trimmed of surrounding whitespace and matching is case-insensitive.
//
//   @(has_pattern("Sell cheese please", "buy (\w+)")) ->
//   @(has_pattern("Buy cheese please", "buy (\w+)")) -> {extra: {0: Buy cheese, 1: cheese}, match: Buy cheese}
//
// @test has_pattern(text, pattern)
func HasPattern(env utils.Environment, text types.XText, pattern types.XText) types.XValue {
	regex, err := regexp.Compile("(?mi)" + strings.TrimSpace(pattern.Native()))
	if err != nil {
		return types.NewXErrorf("must be called with a valid regular expression")
	}

	matches := regex.FindStringSubmatch(text.Native())
	if matches != nil {
		extra := make(map[string]types.XValue, len(matches))

		for g, group := range matches {
			extra[strconv.Itoa(g)] = types.NewXText(group)
		}
		return NewTrueResultWithExtra(types.NewXText(matches[0]), types.NewXDict(extra))
	}

	return nil
}

// HasNumber tests whether `text` contains a number
//
//   @(has_number("the number is 42")) -> {match: 42}
//   @(has_number("the number is forty two")) ->
//
// @test has_number(text)
func HasNumber(env utils.Environment, text types.XText) types.XValue {
	return testNumber(env, text, types.XNumberZero, types.XNumberZero, isNumberTest)
}

// HasNumberBetween tests whether `text` contains a number between `min` and `max` inclusive
//
//   @(has_number_between("the number is 42", 40, 44)) -> {match: 42}
//   @(has_number_between("the number is 42", 50, 60)) ->
//   @(has_number_between("the number is not there", 50, 60)) ->
//   @(has_number_between("the number is not there", "foo", 60)) -> ERROR
//
// @test has_number_between(text, min, max)
func HasNumberBetween(env utils.Environment, arg1 types.XValue, arg2 types.XValue, arg3 types.XValue) types.XValue {
	text, xerr := types.ToXText(env, arg1)
	if xerr != nil {
		return xerr
	}
	min, xerr := types.ToXNumber(env, arg2)
	if xerr != nil {
		return xerr
	}
	max, xerr := types.ToXNumber(env, arg3)
	if xerr != nil {
		return xerr
	}

	return testNumber(env, text, min, max, isNumberBetween)
}

// HasNumberLT tests whether `text` contains a number less than `max`
//
//   @(has_number_lt("the number is 42", 44)) -> {match: 42}
//   @(has_number_lt("the number is 42", 40)) ->
//   @(has_number_lt("the number is not there", 40)) ->
//   @(has_number_lt("the number is not there", "foo")) -> ERROR
//
// @test has_number_lt(text, max)
func HasNumberLT(env utils.Environment, text types.XText, num types.XNumber) types.XValue {
	return testNumber(env, text, num, types.XNumberZero, isNumberLT)
}

// HasNumberLTE tests whether `text` contains a number less than or equal to `max`
//
//   @(has_number_lte("the number is 42", 42)) -> {match: 42}
//   @(has_number_lte("the number is 42", 40)) ->
//   @(has_number_lte("the number is not there", 40)) ->
//   @(has_number_lte("the number is not there", "foo")) -> ERROR
//
// @test has_number_lte(text, max)
func HasNumberLTE(env utils.Environment, text types.XText, num types.XNumber) types.XValue {
	return testNumber(env, text, num, types.XNumberZero, isNumberLTE)
}

// HasNumberEQ tests whether `text` contains a number equal to the `value`
//
//   @(has_number_eq("the number is 42", 42)) -> {match: 42}
//   @(has_number_eq("the number is 42", 40)) ->
//   @(has_number_eq("the number is not there", 40)) ->
//   @(has_number_eq("the number is not there", "foo")) -> ERROR
//
// @test has_number_eq(text, value)
func HasNumberEQ(env utils.Environment, text types.XText, num types.XNumber) types.XValue {
	return testNumber(env, text, num, types.XNumberZero, isNumberEQ)
}

// HasNumberGTE tests whether `text` contains a number greater than or equal to `min`
//
//   @(has_number_gte("the number is 42", 42)) -> {match: 42}
//   @(has_number_gte("the number is 42", 45)) ->
//   @(has_number_gte("the number is not there", 40)) ->
//   @(has_number_gte("the number is not there", "foo")) -> ERROR
//
// @test has_number_gte(text, min)
func HasNumberGTE(env utils.Environment, text types.XText, num types.XNumber) types.XValue {
	return testNumber(env, text, num, types.XNumberZero, isNumberGTE)
}

// HasNumberGT tests whether `text` contains a number greater than `min`
//
//   @(has_number_gt("the number is 42", 40)) -> {match: 42}
//   @(has_number_gt("the number is 42", 42)) ->
//   @(has_number_gt("the number is not there", 40)) ->
//   @(has_number_gt("the number is not there", "foo")) -> ERROR
//
// @test has_number_gt(text, min)
func HasNumberGT(env utils.Environment, text types.XText, num types.XNumber) types.XValue {
	return testNumber(env, text, num, types.XNumberZero, isNumberGT)
}

// HasDate tests whether `text` contains a date formatted according to our environment
//
//   @(has_date("the date is 15/01/2017")) -> {match: 2017-01-15T13:24:30.123456-05:00}
//   @(has_date("there is no date here, just a year 2017")) ->
//
// @test has_date(text)
func HasDate(env utils.Environment, text types.XText) types.XValue {
	return testDate(env, text, types.XDateTimeZero, isDateTest)
}

// HasDateLT tests whether `text` contains a date before the date `max`
//
//   @(has_date_lt("the date is 15/01/2017", "2017-06-01")) -> {match: 2017-01-15T13:24:30.123456-05:00}
//   @(has_date_lt("there is no date here, just a year 2017", "2017-06-01")) ->
//   @(has_date_lt("there is no date here, just a year 2017", "not date")) -> ERROR
//
// @test has_date_lt(text, max)
func HasDateLT(env utils.Environment, text types.XText, date types.XDateTime) types.XValue {
	return testDate(env, text, date, isDateLTTest)
}

// HasDateEQ tests whether `text` a date equal to `date`
//
//   @(has_date_eq("the date is 15/01/2017", "2017-01-15")) -> {match: 2017-01-15T13:24:30.123456-05:00}
//   @(has_date_eq("the date is 15/01/2017 15:00", "2017-01-15")) -> {match: 2017-01-15T15:00:00.000000-05:00}
//   @(has_date_eq("there is no date here, just a year 2017", "2017-06-01")) ->
//   @(has_date_eq("there is no date here, just a year 2017", "not date")) -> ERROR
//
// @test has_date_eq(text, date)
func HasDateEQ(env utils.Environment, text types.XText, date types.XDateTime) types.XValue {
	return testDate(env, text, date, isDateEQTest)
}

// HasDateGT tests whether `text` a date after the date `min`
//
//   @(has_date_gt("the date is 15/01/2017", "2017-01-01")) -> {match: 2017-01-15T13:24:30.123456-05:00}
//   @(has_date_gt("the date is 15/01/2017", "2017-03-15")) ->
//   @(has_date_gt("there is no date here, just a year 2017", "2017-06-01")) ->
//   @(has_date_gt("there is no date here, just a year 2017", "not date")) -> ERROR
//
// @test has_date_gt(text, min)
func HasDateGT(env utils.Environment, text types.XText, date types.XDateTime) types.XValue {
	return testDate(env, text, date, isDateGTTest)
}

// HasTime tests whether `text` contains a time.
//
//   @(has_time("the time is 10:30")) -> {match: 10:30:00.000000}
//   @(has_time("the time is 10 PM")) -> {match: 22:00:00.000000}
//   @(has_time("the time is 10:30:45")) -> {match: 10:30:45.000000}
//   @(has_time("there is no time here, just the number 25")) ->
//
// @test has_time(text)
func HasTime(env utils.Environment, text types.XText) types.XValue {
	t, xerr := types.ToXTime(env, text)
	if xerr == nil {
		return NewTrueResult(t)
	}

	return nil
}

var emailAddressRE = regexp.MustCompile(`([\pL\pN][-_.\pL\pN]*)@([\pL\pN][-_\pL\pN]*)(\.[\pL\pN][-_\pL\pN]*)+`)

// HasEmail tests whether an email is contained in `text`
//
//   @(has_email("my email is foo1@bar.com, please respond")) -> {match: foo1@bar.com}
//   @(has_email("my email is <foo@bar2.com>")) -> {match: foo@bar2.com}
//   @(has_email("i'm not sharing my email")) ->
//
// @test has_email(text)
func HasEmail(env utils.Environment, text types.XText) types.XValue {
	// split by whitespace
	email := emailAddressRE.FindString(text.Native())
	if email != "" {
		return NewTrueResult(types.NewXText(email))
	}

	return nil
}

// HasPhone tests whether `text` contains a phone number. The optional `country_code` argument specifies
// the country to use for parsing.
//
//   @(has_phone("my number is +12067799294")) -> {match: +12067799294}
//   @(has_phone("my number is 2067799294", "US")) -> {match: +12067799294}
//   @(has_phone("my number is 206 779 9294", "US")) -> {match: +12067799294}
//   @(has_phone("my number is none of your business", "US")) ->
//
// @test has_phone(text, country_code)
func HasPhone(env utils.Environment, text types.XText, args ...types.XValue) types.XValue {
	var country types.XText
	var xerr types.XError
	if len(args) == 1 {
		country, xerr = types.ToXText(env, args[0])
		if xerr != nil {
			return xerr
		}
	} else {
		country = types.NewXText(string(env.DefaultCountry()))
	}

	// try to find a phone number
	phone, err := phonenumbers.Parse(text.Native(), country.Native())
	if err != nil {
		return nil
	}

	if !phonenumbers.IsPossibleNumber(phone) {
		return nil
	}

	// format as E164 number
	formatted := phonenumbers.Format(phone, phonenumbers.E164)
	return NewTrueResult(types.NewXText(formatted))
}

// HasState tests whether a state name is contained in the `text`
//
//   @(has_state("Kigali")) -> {match: Rwanda > Kigali City}
//   @(has_state("¡Kigali!")) -> {match: Rwanda > Kigali City}
//   @(has_state("I live in Kigali")) -> {match: Rwanda > Kigali City}
//   @(has_state("Boston")) ->
//
// @test has_state(text)
func HasState(env utils.Environment, text types.XText) types.XValue {
	runEnv, _ := env.(flows.RunEnvironment)

	states, err := runEnv.FindLocationsFuzzy(text.Native(), flows.LocationLevelState, nil)
	if err != nil {
		return types.NewXError(err)
	}
	if len(states) > 0 {
		return NewTrueResult(types.NewXText(string(states[0].Path())))
	}
	return nil
}

// HasDistrict tests whether a district name is contained in the `text`. If `state` is also provided
// then the returned district must be within that state.
//
//   @(has_district("Gasabo", "Kigali")) -> {match: Rwanda > Kigali City > Gasabo}
//   @(has_district("I live in Gasabo", "Kigali")) -> {match: Rwanda > Kigali City > Gasabo}
//   @(has_district("Gasabo", "Boston")) ->
//   @(has_district("Gasabo")) -> {match: Rwanda > Kigali City > Gasabo}
//
// @test has_district(text, state)
func HasDistrict(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 && len(args) != 2 {
		return types.NewXErrorf("takes one or two arguments, got %d", len(args))
	}

	runEnv, _ := env.(flows.RunEnvironment)

	var text, stateText types.XText
	var xerr types.XError

	// grab the text we will search and the parent state name
	if text, xerr = types.ToXText(env, args[0]); xerr != nil {
		return xerr
	}
	if len(args) == 2 {
		if stateText, xerr = types.ToXText(env, args[1]); xerr != nil {
			return xerr
		}
	}

	states, err := runEnv.FindLocationsFuzzy(stateText.Native(), flows.LocationLevelState, nil)
	if err != nil {
		return types.NewXError(err)
	}
	if len(states) > 0 {
		districts, err := runEnv.FindLocationsFuzzy(text.Native(), flows.LocationLevelDistrict, states[0])
		if err != nil {
			return types.NewXError(err)
		}
		if len(districts) > 0 {
			return NewTrueResult(types.NewXText(string(districts[0].Path())))
		}
	}

	// try without a parent state - it's ok as long as we get a single match
	if stateText.Empty() {
		districts, err := runEnv.FindLocationsFuzzy(text.Native(), flows.LocationLevelDistrict, nil)
		if err != nil {
			return types.NewXError(err)
		}
		if len(districts) == 1 {
			return NewTrueResult(types.NewXText(string(districts[0].Path())))
		}
	}

	return nil
}

// HasWard tests whether a ward name is contained in the `text`
//
//   @(has_ward("Gisozi", "Gasabo", "Kigali")) -> {match: Rwanda > Kigali City > Gasabo > Gisozi}
//   @(has_ward("I live in Gisozi", "Gasabo", "Kigali")) -> {match: Rwanda > Kigali City > Gasabo > Gisozi}
//   @(has_ward("Gisozi", "Gasabo", "Brooklyn")) ->
//   @(has_ward("Gisozi", "Brooklyn", "Kigali")) ->
//   @(has_ward("Brooklyn", "Gasabo", "Kigali")) ->
//   @(has_ward("Gasabo")) ->
//   @(has_ward("Gisozi")) -> {match: Rwanda > Kigali City > Gasabo > Gisozi}
//
// @test has_ward(text, district, state)
func HasWard(env utils.Environment, args ...types.XValue) types.XValue {
	if len(args) != 1 && len(args) != 3 {
		return types.NewXErrorf("takes one or three arguments, got %d", len(args))
	}

	runEnv, _ := env.(flows.RunEnvironment)

	var text, districtText, stateText types.XText
	var xerr types.XError

	// grab the text we will search, as well as the parent district and state names
	if text, xerr = types.ToXText(env, args[0]); xerr != nil {
		return xerr
	}
	if len(args) == 3 {
		if districtText, xerr = types.ToXText(env, args[1]); xerr != nil {
			return xerr
		}
		if stateText, xerr = types.ToXText(env, args[2]); xerr != nil {
			return xerr
		}
	}

	states, err := runEnv.FindLocationsFuzzy(stateText.Native(), flows.LocationLevelState, nil)
	if err != nil {
		return types.NewXError(err)
	}
	if len(states) > 0 {
		districts, err := runEnv.FindLocationsFuzzy(districtText.Native(), flows.LocationLevelDistrict, states[0])
		if err != nil {
			return types.NewXError(err)
		}
		if len(districts) > 0 {
			wards, err := runEnv.FindLocationsFuzzy(text.Native(), flows.LocationLevelWard, districts[0])
			if err != nil {
				return types.NewXError(err)
			}
			if len(wards) > 0 {
				return NewTrueResult(types.NewXText(string(wards[0].Path())))
			}
		}
	}

	// try without a parent district - it's ok as long as we get a single match
	if districtText.Empty() {
		wards, err := runEnv.FindLocationsFuzzy(text.Native(), flows.LocationLevelWard, nil)
		if err != nil {
			return types.NewXError(err)
		}
		if len(wards) == 1 {
			return NewTrueResult(types.NewXText(string(wards[0].Path())))
		}
	}

	return nil
}

//------------------------------------------------------------------------------------------
// Text Test Functions
//------------------------------------------------------------------------------------------

type stringTokenTest func(origHayTokens []string, hayTokens []string, pinTokens []string) types.XValue

func testStringTokens(env utils.Environment, str types.XText, testStr types.XText, testFunc stringTokenTest) types.XValue {
	hayStack := strings.TrimSpace(str.Native())
	needle := strings.TrimSpace(testStr.Native())

	origHays := utils.TokenizeString(hayStack)
	hays := utils.TokenizeString(strings.ToLower(hayStack))
	needles := utils.TokenizeString(strings.ToLower(needle))

	return testFunc(origHays, hays, needles)
}

func hasPhraseTest(origHays []string, hays []string, pins []string) types.XValue {
	if len(pins) == 0 {
		return NewTrueResult(types.XTextEmpty)
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
		return NewTrueResult(types.NewXText(strings.Join(matches, " ")))
	}

	return nil
}

func hasAllWordsTest(origHays []string, hays []string, pins []string) types.XValue {
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
		return NewTrueResult(types.NewXText(strings.Join(matches, " ")))
	}

	return nil
}

func hasAnyWordTest(origHays []string, hays []string, pins []string) types.XValue {
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
		return NewTrueResult(types.NewXText(strings.Join(matches, " ")))
	}

	return nil
}

func hasOnlyPhraseTest(origHays []string, hays []string, pins []string) types.XValue {
	// must be same length
	if len(hays) != len(pins) {
		return nil
	}

	// and every token must match
	matches := make([]string, 0, len(pins))
	for i := range hays {
		if hays[i] != pins[i] {
			return nil
		}
		matches = append(matches, origHays[i])
	}

	return NewTrueResult(types.NewXText(strings.Join(matches, " ")))
}

//------------------------------------------------------------------------------------------
// Numerical Test Functions
//------------------------------------------------------------------------------------------

// ParseDecimalFuzzy parses a decimal from a string
func ParseDecimalFuzzy(val string, format *utils.NumberFormat) (decimal.Decimal, error) {
	// remove digit grouping symbol
	cleaned := strings.Replace(val, format.DigitGroupingSymbol, "", -1)

	// replace non-period decimal symbols
	cleaned = strings.Replace(cleaned, format.DecimalSymbol, ".", -1)

	return decimal.NewFromString(cleaned)
}

type decimalTest func(value decimal.Decimal, test1 decimal.Decimal, test2 decimal.Decimal) bool

func testNumber(env utils.Environment, str types.XText, testNum1 types.XNumber, testNum2 types.XNumber, testFunc decimalTest) types.XValue {
	// create a number finding regex based on current environment
	pattern := regexp.MustCompile(fmt.Sprintf(`[-+]?[\pNlO\%s]+(\%s[\pNlO]+)?`, env.NumberFormat().DigitGroupingSymbol, env.NumberFormat().DecimalSymbol))

	// look for number like things in the input and use the first one that we can actually parse
	for _, value := range pattern.FindAllString(str.Native(), -1) {
		num, err := ParseDecimalFuzzy(value, env.NumberFormat())
		if err == nil {
			if testFunc(num, testNum1.Native(), testNum2.Native()) {
				return NewTrueResult(types.NewXNumber(num))
			}
		}
	}

	return nil
}

func isNumberTest(value decimal.Decimal, _ decimal.Decimal, _ decimal.Decimal) bool {
	return true
}

func isNumberLT(value decimal.Decimal, test decimal.Decimal, _ decimal.Decimal) bool {
	return value.Cmp(test) < 0
}

func isNumberLTE(value decimal.Decimal, test decimal.Decimal, _ decimal.Decimal) bool {
	return value.Cmp(test) <= 0
}

func isNumberEQ(value decimal.Decimal, test decimal.Decimal, _ decimal.Decimal) bool {
	return value.Cmp(test) == 0
}

func isNumberGTE(value decimal.Decimal, test decimal.Decimal, _ decimal.Decimal) bool {
	return value.Cmp(test) >= 0
}

func isNumberGT(value decimal.Decimal, test decimal.Decimal, _ decimal.Decimal) bool {
	return value.Cmp(test) > 0
}

func isNumberBetween(value decimal.Decimal, test1 decimal.Decimal, test2 decimal.Decimal) bool {
	return value.Cmp(test1) >= 0 && value.Cmp(test2) <= 0
}

//------------------------------------------------------------------------------------------
// Date Test Functions
//------------------------------------------------------------------------------------------

type dateTest func(utils.Date, utils.Date) bool

func testDate(env utils.Environment, str types.XText, testDate types.XDateTime, testFunc dateTest) types.XValue {
	// first parse with time filling which will be the test result
	value, xerr := types.ToXDateTimeWithTimeFill(env, str)

	// but comparsion should be against only the date portions
	valueAsDate := utils.ExtractDate(value.In(env.Timezone()).Native())
	testAsDate := utils.ExtractDate(testDate.In(env.Timezone()).Native())

	if xerr != nil {
		return nil
	}

	if testFunc(valueAsDate, testAsDate) {
		return NewTrueResult(value)
	}

	return nil
}

func isDateTest(value utils.Date, test utils.Date) bool {
	return true
}

func isDateLTTest(value utils.Date, test utils.Date) bool {
	return value.Compare(test) < 0
}

func isDateEQTest(value utils.Date, test utils.Date) bool {
	return value.Compare(test) == 0
}

func isDateGTTest(value utils.Date, test utils.Date) bool {
	return value.Compare(test) > 0
}
