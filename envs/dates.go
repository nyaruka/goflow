package envs

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.RegisterValidatorTag("date_format", validateDateFormat, func(validator.FieldError) string {
		return "is not a valid date format"
	})
	utils.RegisterValidatorTag("time_format", validateTimeFormat, func(validator.FieldError) string {
		return "is not a valid time format"
	})
}

func validateDateFormat(fl validator.FieldLevel) bool {
	// validate for parsing which has stricter requirements
	return dates.ValidateFormat(fl.Field().String(), dates.DateOnlyLayouts, dates.ParsingMode) == nil
}

func validateTimeFormat(fl validator.FieldLevel) bool {
	// validate for parsing which has stricter requirements
	return dates.ValidateFormat(fl.Field().String(), dates.TimeOnlyLayouts, dates.ParsingMode) == nil
}

// patterns for date and time formats supported for human-entered data
var patternDayMonthYear = regexp.MustCompile(`\b([0-9]{1,2})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{4}|[0-9]{2})\b`)
var patternMonthDayYear = regexp.MustCompile(`\b([0-9]{1,2})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{4}|[0-9]{2})\b`)
var patternYearMonthDay = regexp.MustCompile(`\b([0-9]{4}|[0-9]{2})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{1,2})\b`)

var patternTime = regexp.MustCompile(`\b(\d{1,2})(?:(?:\:)?(\d{2})(?:\:(\d{2})(?:\.(\d+))?)?)?\W*([aApP][mM])?\b`)

// DateFormat a date format string
type DateFormat string

// TimeFormat a time format string
type TimeFormat string

// standard date and time formats
const (
	DateFormatYearMonthDay DateFormat = "YYYY-MM-DD"
	DateFormatMonthDayYear DateFormat = "MM-DD-YYYY"
	DateFormatDayMonthYear DateFormat = "DD-MM-YYYY"

	TimeFormatHourMinute           TimeFormat = "tt:mm"
	TimeFormatHourMinuteAmPm       TimeFormat = "h:mm aa"
	TimeFormatHourMinuteSecond     TimeFormat = "tt:mm:ss"
	TimeFormatHourMinuteSecondAmPm TimeFormat = "h:mm:ss aa"
)

func (df DateFormat) String() string { return string(df) }
func (tf TimeFormat) String() string { return string(tf) }

// generic format for parsing any 8601 date
var iso8601Format = "2006-01-02T15:04:05Z07:00"
var iso8601NoSecondsFormat = "2006-01-02T15:04Z07:00"
var iso8601DateOnlyFormat = "2006-01-02"

var isoFormats = []string{iso8601Format, iso8601NoSecondsFormat}

// ZeroDateTime is our uninitialized datetime value
var ZeroDateTime = time.Time{}

func dateFromFormats(currentYear int, pattern *regexp.Regexp, d int, m int, y int, str string) (dates.Date, string, error) {

	matches := pattern.FindAllStringSubmatchIndex(str, -1)
	for _, match := range matches {
		groups := utils.StringSlices(str, match)

		// does our day look believable?
		day, _ := strconv.Atoi(groups[d])
		if day == 0 || day > 31 {
			continue
		}
		month, _ := strconv.Atoi(groups[m])
		if month == 0 || month > 12 {
			continue
		}

		year, _ := strconv.Atoi(groups[y])

		// convert to four digit year if necessary
		if len(groups[y]) == 2 {
			if year > currentYear%1000 {
				year += 1900
			} else {
				year += 2000
			}
		}

		remainder := str[match[1]:]

		// looks believable, go for it
		return dates.NewDate(year, month, day), remainder, nil
	}

	return dates.ZeroDate, str, errors.Errorf("string '%s' couldn't be parsed as a date", str)
}

// DateTimeFromString returns a datetime constructed from the passed in string, or an error if we
// are unable to extract one
func DateTimeFromString(env Environment, str string, fillTime bool) (time.Time, error) {
	str = strings.Trim(str, " \n\r\t")

	// first see if we can parse in any known ISO formats, if so return that
	for _, format := range isoFormats {
		parsed, err := time.ParseInLocation(format, str, env.Timezone())
		if err == nil {
			return parsed, nil
		}
	}

	// otherwise, try to parse according to their env settings
	date, remainder, err := parseDate(env, str)

	// couldn't find a date? bail
	if err != nil {
		return ZeroDateTime, err
	}

	// can we pull out a time from the remainder of the string?
	hasTime, timeOfDay := parseTime(remainder)
	if !hasTime && fillTime {
		timeOfDay = dates.ExtractTimeOfDay(env.Now())
	}

	// combine our date and time
	return time.Date(date.Year, time.Month(date.Month), date.Day, timeOfDay.Hour, timeOfDay.Minute, timeOfDay.Second, timeOfDay.Nanos, env.Timezone()), nil
}

// DateFromString returns a date constructed from the passed in string, or an error if we
// are unable to extract one
func DateFromString(env Environment, str string) (dates.Date, error) {
	parsed, _, err := parseDate(env, str)
	return parsed, err
}

// TimeFromString returns a time of day constructed from the passed in string, or an error if we
// are unable to extract one
func TimeFromString(str string) (dates.TimeOfDay, error) {
	hasTime, timeOfDay := parseTime(str)
	if !hasTime {
		return dates.ZeroTimeOfDay, errors.Errorf("string '%s' couldn't be parsed as a time", str)
	}
	return timeOfDay, nil
}

func parseDate(env Environment, str string) (dates.Date, string, error) {
	str = strings.Trim(str, " \n\r\t")

	// try to parse as ISO date
	asISO, err := time.ParseInLocation(iso8601DateOnlyFormat, str[0:utils.MinInt(len(iso8601DateOnlyFormat), len(str))], env.Timezone())
	if err == nil {
		return dates.ExtractDate(asISO), str[len(iso8601DateOnlyFormat):], nil
	}

	// otherwise, try to parse according to their env settings
	currentYear := dates.Now().Year()

	switch env.DateFormat() {
	case DateFormatYearMonthDay:
		return dateFromFormats(currentYear, patternYearMonthDay, 3, 2, 1, str)
	case DateFormatDayMonthYear:
		return dateFromFormats(currentYear, patternDayMonthYear, 1, 2, 3, str)
	case DateFormatMonthDayYear:
		return dateFromFormats(currentYear, patternMonthDayYear, 2, 1, 3, str)
	}

	return dates.ZeroDate, "", errors.Errorf("unknown date format: %s", env.DateFormat())
}

func parseTime(str string) (bool, dates.TimeOfDay) {
	matches := patternTime.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		hour, _ := strconv.Atoi(match[1])
		minute, _ := strconv.Atoi(match[2])
		second, _ := strconv.Atoi(match[3])
		ampm := strings.ToLower(match[5])

		// do we have an AM/PM marker
		if hour < 12 && ampm == "pm" {
			hour += 12
		} else if hour == 12 && ampm == "am" {
			hour -= 12
		}

		nanosStr := match[4]
		nanos := 0
		if nanosStr != "" {
			// can only read nano second accuracy
			if len(nanosStr) > 9 {
				nanosStr = nanosStr[0:9]
			}
			nanos, _ = strconv.Atoi(nanosStr)
			nanos *= int(math.Pow(10, float64(9-len(nanosStr))))
		}

		// 24:00:00.000000 is a special case for midnight
		if hour == 24 && minute == 0 && second == 0 && nanos == 0 {
			hour = 0
		}

		// is our time valid?
		if hour > 24 {
			continue
		}
		if minute > 60 {
			continue
		}
		if second > 60 {
			continue
		}

		return true, dates.NewTimeOfDay(hour, minute, second, nanos)
	}

	return false, dates.ZeroTimeOfDay
}
