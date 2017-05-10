package excellent

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
)

// TODO:
// SubflowTest
// RegexTest
// HasWardTest
// HasDistrictTest
// HasStateTest
// InterruptTest
// TimeoutTest
// AirtimeStatusTest

//------------------------------------------------------------------------------------------
// Mapping
//------------------------------------------------------------------------------------------

// XTESTS is our mapping of the excellent test names to their actual functions
var XTESTS = map[string]XFunction{
	"has_error":          HasError,
	"has_value":          HasValue,
	"has_group":          HasGroup,
	"has_run_status":     HasRunStatus,
	"has_webhook_status": HasWebhookStatus,

	"has_phrase":      HasPhrase,
	"has_only_phrase": HasOnlyPhrase,
	"has_any_word":    HasAnyWord,
	"has_all_words":   HasAllWords,
	"has_beginning":   HasBeginning,
	"has_text":        HasText,

	"has_number":         HasNumber,
	"has_number_between": HasNumberBetween,
	"has_number_lt":      HasNumberLT,
	"has_number_lte":     HasNumberLTE,
	"has_number_eq":      HasNumberEQ,
	"has_number_gte":     HasNumberGTE,
	"has_number_gt":      HasNumberGT,

	"has_date":    HasDate,
	"has_date_lt": HasDateLT,
	"has_date_eq": HasDateEQ,
	"has_date_gt": HasDateGT,

	"has_phone": HasPhone,
	"has_email": HasEmail,
}

//------------------------------------------------------------------------------------------
// Interfaces
//------------------------------------------------------------------------------------------

// XTestResult encapsulates not only if the test was true but what the match was
type XTestResult struct {
	matched bool
	match   interface{}
}

// Matched returns whether the test matched
func (t XTestResult) Matched() bool { return t.matched }

// Match returns the item which was matched
func (t XTestResult) Match() interface{} { return t.match }

// Default satisfies the utils.VariableResolver interface, we always default to whether we matched
func (t XTestResult) Default() interface{} {
	return t.Matched
}

// Resolve satisfies the utils.VariableResolver interface, users can look up the match or whether we matched
func (t XTestResult) Resolve(key string) interface{} {
	switch key {
	case "matched":
		return t.Matched

	case "match":
		return t.Match
	}
	return fmt.Errorf("No such key '%s' on test result", key)
}

// XFalseResult can be used as a singleton for false result values
var XFalseResult = XTestResult{}

// Enforce Variable Resolver interface
var _ utils.VariableResolver = XTestResult{}

//------------------------------------------------------------------------------------------
// Tests
//------------------------------------------------------------------------------------------

// HasError returns whether the passed in argument is an error
func HasError(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("HAS_ERROR takes exactly one argument, got %d", len(args))
	}

	// nil is not an error
	if args[0] == nil {
		return XFalseResult
	}

	err, isErr := args[0].(error)
	if isErr {
		return XTestResult{true, err}
	}

	return XFalseResult
}

// HasValue returns whether the passed in argument is non-nil
func HasValue(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("HAS_VALUE takes exactly one argument, got %d", len(args))
	}

	// nil is not an error
	if args[0] == nil {
		return XFalseResult
	}

	return XTestResult{true, args[0]}
}

// HasRunStatus returns whether the passed in run has the passed in status
func HasRunStatus(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("HAS_RUN_STATUS takes exactly two arguments, got %d", len(args))
	}

	// first parameter needs to be a flow run
	run, isRun := args[0].(flows.FlowRun)
	if !isRun {
		return fmt.Errorf("HAS_RUN_STATUS must be called with a run as first argument")
	}

	status, err := utils.ToString(env, args[1])
	if err != nil {
		return fmt.Errorf("HAS_RUN_STATUS must be called with a string as second argument")
	}

	if flows.RunStatus(strings.ToUpper(status)) == run.Status() {
		return XTestResult{true, run.Status()}
	}

	return XFalseResult
}

// HasWebhookStatus returns whether the passed in webhook response has the passed in status
func HasWebhookStatus(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("HAS_WEBHOOK_STATUS takes exactly two arguments, got %d", len(args))
	}

	// first parameter needs to be a request response
	rr, isRR := args[0].(utils.RequestResponse)
	if !isRR {
		return fmt.Errorf("HAS_WEBHOOK_STATUS must be called with webhook as first argument")
	}

	status, err := utils.ToString(env, args[1])
	if err != nil {
		return fmt.Errorf("HAS_WEBHOOK_STATUS must be called with a string as second argument")
	}

	if utils.RequestResponseStatus(strings.ToUpper(status)) == rr.Status() {
		return XTestResult{true, rr.Status()}
	}

	return XFalseResult
}

// HasGroup returns whether the passed in contact is part of the passed in group
func HasGroup(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("HAS_GROUP takes exactly two arguments, got %d", len(args))
	}

	// is the first argument a contact?
	contact, isContact := args[0].(flows.Contact)
	if !isContact {
		return fmt.Errorf("HAS_GROUP must have a contact as its first argument")
	}

	groupUUID, err := utils.ToString(env, args[1])
	if err != nil {
		return err
	}

	// iterate through the groups looking for one with the same UUID as passed in
	group := contact.Groups().FindGroup(flows.GroupUUID(groupUUID))
	if group != nil {
		return XTestResult{true, group}
	}

	return XFalseResult
}

// HasPhrase tests whether the passed in phrase is contained in the text
func HasPhrase(env utils.Environment, args ...interface{}) interface{} {
	return testStringTokens(env, "HAS_PHRASE", hasPhraseTest, args)
}

// HasAllWords tests whether all the words are contained in the text
func HasAllWords(env utils.Environment, args ...interface{}) interface{} {
	return testStringTokens(env, "HAS_ALL_WORDS", hasAllWordsTest, args)
}

// HasAnyWord tests whether any of the words are contained in the text
func HasAnyWord(env utils.Environment, args ...interface{}) interface{} {
	return testStringTokens(env, "HAS_ANY_WORD", hasAnyWordTest, args)
}

// HasOnlyPhrase tests whether the text contains only the phrase
func HasOnlyPhrase(env utils.Environment, args ...interface{}) interface{} {
	return testStringTokens(env, "HAS_ONLY_PHRASE", hasOnlyPhraseTest, args)
}

// HasText tests whether there is any text
func HasText(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("HAS_TEXT takes exactly one arguments, got %d", len(args))
	}

	text, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	// trim any whitespace
	text = strings.TrimSpace(text)

	// if there is anything left then we have text
	if len(text) > 0 {
		return XTestResult{true, text}
	}

	return XFalseResult
}

// HasBeginning tests whether the text starts with the words
func HasBeginning(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("HAS_BEGINNING takes exactly two arguments, got %d", len(args))
	}

	hayStack, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	pinCushion, err := utils.ToString(env, args[1])
	if err != nil {
		return err
	}

	// trim both
	pinCushion = strings.TrimSpace(pinCushion)
	hayStack = strings.TrimSpace(hayStack)

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
		return XTestResult{true, segment}
	}

	return XFalseResult
}

// HasNumber tests whether the text contains a number
func HasNumber(env utils.Environment, args ...interface{}) interface{} {
	// only one argument for has number
	if len(args) != 1 {
		return fmt.Errorf("HAS_NUMBER takes exactly one arguments, got %d", len(args))
	}

	testArgs := make([]interface{}, 2)
	testArgs[0] = args[0]

	// set our second argument to a dummy, it isn't used but is need to satisfy our interface
	testArgs[1] = "0"

	return testDecimal(env, "HAS_NUMBER", isNumberTest, testArgs)
}

// HasNumberBetween tests whether the text contains a number between args[1] and args[2]
func HasNumberBetween(env utils.Environment, args ...interface{}) interface{} {
	// need three arguments, value being tested and min, max
	if len(args) != 3 {
		return fmt.Errorf("HAS_NUMBER_BETWEEN takes exactly three arguments, got %d", len(args))
	}

	values, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	min, err := utils.ToDecimal(env, args[1])
	if err != nil {
		return err
	}
	max, err := utils.ToDecimal(env, args[2])
	if err != nil {
		return err
	}

	// for each of our values, try to evaluate to a decimal
	for _, value := range strings.Fields(values) {
		decimalValue, err := utils.ToDecimal(env, value)
		if err == nil {
			if decimalValue.Cmp(min) >= 0 && decimalValue.Cmp(max) <= 0 {
				return XTestResult{true, decimalValue}
			}
		}
	}
	return XFalseResult
}

// HasNumberLT tests whether the text contains a number less than the value
func HasNumberLT(env utils.Environment, args ...interface{}) interface{} {
	return testDecimal(env, "HAS_NUMBER_LT", isNumberLT, args)
}

// HasNumberLTE tests whether the text contains a number less than or equal to the value
func HasNumberLTE(env utils.Environment, args ...interface{}) interface{} {
	return testDecimal(env, "HAS_NUMBER_LTE", isNumberLTE, args)
}

// HasNumberEQ tests whether the text contains a number equal to the value
func HasNumberEQ(env utils.Environment, args ...interface{}) interface{} {
	return testDecimal(env, "HAS_NUMBER_EQ", isNumberEQ, args)
}

// HasNumberGTE tests whether the text contains a number greater than or equal to the value
func HasNumberGTE(env utils.Environment, args ...interface{}) interface{} {
	return testDecimal(env, "HAS_NUMBER_GTE", isNumberGTE, args)
}

// HasNumberGT tests whether the text contains a number greater than the value
func HasNumberGT(env utils.Environment, args ...interface{}) interface{} {
	return testDecimal(env, "HAS_NUMBER_GT", isNumberGT, args)
}

// HasDate tests whether the text contains a date
func HasDate(env utils.Environment, args ...interface{}) interface{} {
	// only one argument for has date
	if len(args) != 1 {
		return fmt.Errorf("HAS_DATE takes exactly one arguments, got %d", len(args))
	}

	testArgs := make([]interface{}, 2)
	testArgs[0] = args[0]

	// set our second argument to a dummy, it isn't used but is need to satisfy our interface
	testArgs[1] = time.Now()

	return testDate(env, "HAS_DATE", isDateTest, testArgs)
}

// HasDateLT tests whether the text contains a date before the value
func HasDateLT(env utils.Environment, args ...interface{}) interface{} {
	return testDate(env, "HAS_DATE_LT", isDateLTTest, args)
}

// HasDateEQ tests whether the text contains a date equal to the value
func HasDateEQ(env utils.Environment, args ...interface{}) interface{} {
	return testDate(env, "HAS_DATE_EQ", isDateEQTest, args)
}

// HasDateGT tests whether the text contains a date after the value
func HasDateGT(env utils.Environment, args ...interface{}) interface{} {
	return testDate(env, "HAS_DATE_GT", isDateGTTest, args)
}

var emailAddressRE = regexp.MustCompile(`^([\pL][-_.\pL]*)@([\pL][-_\pL]*)(\.[\pL][-_\pL]*)+$`)

// HasEmail tests whether an email is contained in the text
func HasEmail(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("HAS_EMAIL takes exactly one argument, got %d", len(args))
	}

	// convert our arg to a string
	text, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	// split by whitespace
	for _, word := range strings.Fields(text) {
		email := emailAddressRE.FindString(word)
		if email != "" {
			return XTestResult{true, email}
		}
	}

	return XFalseResult
}

// TODO: plug in a real phone number parsing library
var phoneRE = regexp.MustCompile(`^\+?([0-9]{7,12})$`)

// HasPhone tests whether a phone number is contained in the text
func HasPhone(env utils.Environment, args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("HAS_PHONE takes exactly one argument, got %d", len(args))
	}

	// convert our arg to a string
	text, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	// split by whitespace
	for _, word := range strings.Fields(text) {
		phone := phoneRE.FindString(word)
		if phone != "" {
			return XTestResult{true, phone}
		}
	}

	return XFalseResult
}

//------------------------------------------------------------------------------------------
// String Test Functions
//------------------------------------------------------------------------------------------

type stringTokenTest func(origHayTokens []string, hayTokens []string, pinTokens []string) interface{}

func testStringTokens(env utils.Environment, name string, test stringTokenTest, args []interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("%s takes exactly two arguments, got %d", name, len(args))
	}

	hayStack, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	pinCushion, err := utils.ToString(env, args[1])
	if err != nil {
		return err
	}

	hayStack = strings.TrimSpace(hayStack)
	pinCushion = strings.TrimSpace(pinCushion)

	// either are empty, no match
	if hayStack == "" || pinCushion == "" {
		return XFalseResult
	}

	origHays := utils.TokenizeString(hayStack)
	hays := utils.TokenizeString(strings.ToLower(hayStack))
	pins := utils.TokenizeString(strings.ToLower(pinCushion))

	return test(origHays, hays, pins)
}

func hasPhraseTest(origHays []string, hays []string, pins []string) interface{} {
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
		return XTestResult{true, strings.Join(matches, " ")}
	}

	return XFalseResult
}

func hasAllWordsTest(origHays []string, hays []string, pins []string) interface{} {
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
		return XTestResult{true, strings.Join(matches, " ")}
	}

	return XFalseResult
}

func hasAnyWordTest(origHays []string, hays []string, pins []string) interface{} {
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
		return XTestResult{true, strings.Join(matches, " ")}
	}

	return XFalseResult
}

func hasOnlyPhraseTest(origHays []string, hays []string, pins []string) interface{} {
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

	return XTestResult{true, strings.Join(matches, " ")}
}

//------------------------------------------------------------------------------------------
// Decimal Test Functions
//------------------------------------------------------------------------------------------

type decimalTest func(value decimal.Decimal, test decimal.Decimal) bool

func testDecimal(env utils.Environment, name string, test decimalTest, args []interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("%s takes exactly two arguments, got %d", name, len(args))
	}

	values, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	decimalTest, err := utils.ToDecimal(env, args[1])
	if err != nil {
		return err
	}

	// for each of our values, try to evaluate to a decimal
	for _, value := range strings.Fields(values) {
		decimalValue, err := utils.ToDecimal(env, value)
		if err == nil {
			if test(decimalValue, decimalTest) {
				return XTestResult{true, decimalValue}
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

func testDate(env utils.Environment, name string, test dateTest, args []interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("%s takes exactly two arguments, got %d", name, len(args))
	}

	// if we can't convert this to a string, then that's an error
	_, err := utils.ToString(env, args[0])
	if err != nil {
		return err
	}

	// error is if we don't find a date on our test value, that's ok but no match
	value, err := utils.ToDate(env, args[0])
	if err != nil {
		return XFalseResult
	}

	dateTest, err := utils.ToDate(env, args[1])
	if err != nil {
		return err
	}

	if test(value, dateTest) {
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
