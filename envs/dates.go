package envs

import (
	"bytes"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/dates"

	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.Validator.RegisterValidation("date_format", func(fl validator.FieldLevel) bool {
		_, err := ToGoDateFormat(fl.Field().String(), DateOnlyFormatting)
		return err == nil
	})
	utils.Validator.RegisterValidation("time_format", func(fl validator.FieldLevel) bool {
		_, err := ToGoDateFormat(fl.Field().String(), TimeOnlyFormatting)
		return err == nil
	})
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
	asISO, err := time.ParseInLocation(dates.ISO8601Date, str[0:utils.MinInt(len(dates.ISO8601Date), len(str))], env.Timezone())
	if err == nil {
		return dates.ExtractDate(asISO), str[len(dates.ISO8601Date):], nil
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

// FormattingMode describe a mode of formatting dates, times, datetimes
type FormattingMode int

// supported formatting modes
const (
	DateOnlyFormatting FormattingMode = iota
	TimeOnlyFormatting
	DateTimeFormatting
)

var ignoredFormattingRunes = map[rune]bool{' ': true, ':': true, '/': true, '.': true, 'T': true, '-': true, '_': true}

// ToGoDateFormat converts the passed in format to a GoLang format string.
//
// If mode is DateOnlyFormatting or DateTimeFormatting, the following sequences are accepted:
//
//  `YY`        - last two digits of year 0-99
//  `YYYY`      - four digits of your 0000-9999
//  `M`         - month 1-12
//  `MM`        - month 01-12
//  `D`         - day of month, 1-31
//  `DD`        - day of month, zero padded 0-31
//
// If mode is TimeOnlyFormatting or DateTimeFormatting, the following sequences are accepted:
//
//  `h`         - hour of the day 1-12
//  `hh`        - hour of the day 01-12
//  `tt`        - twenty four hour of the day 00-23
//  `m`         - minute 0-59
//  `mm`        - minute 00-59
//  `s`         - second 0-59
//  `ss`        - second 00-59
//  `fff`       - milliseconds
//  `ffffff`    - microseconds
//  `fffffffff` - nanoseconds
//  `aa`        - am or pm
//  `AA`        - AM or PM
//  `Z`         - hour and minute offset from UTC, or Z for UTC
//  `ZZZ`       - hour and minute offset from UTC
//
// ignored chars: ' ', ':', ',', 'T', '-', '_', '/'
func ToGoDateFormat(format string, mode FormattingMode) (string, error) {
	runes := []rune(format)
	goFormat := bytes.Buffer{}

	repeatCount := func(runes []rune, offset int, test rune) int {
		count := 0
		for i := offset; i < len(runes); i++ {
			if runes[i] == test {
				count++
			} else {
				break
			}
		}
		return count
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		count := repeatCount(runes, i, r)

		if mode == DateOnlyFormatting || mode == DateTimeFormatting {
			switch r {
			case 'Y':
				if count == 2 {
					goFormat.WriteString("06")
					i++
				} else if count == 4 {
					goFormat.WriteString("2006")
					i += 3
				} else {
					return "", errors.Errorf("invalid date format, invalid count of 'Y' format: %d", count)
				}
				continue

			case 'M':
				if count == 1 {
					goFormat.WriteString("1")
				} else if count == 2 {
					goFormat.WriteString("01")
					i++
				} else {
					return "", errors.Errorf("invalid date format, invalid count of 'M' format: %d", count)
				}
				continue

			case 'D':
				if count == 1 {
					goFormat.WriteString("2")
				} else if count >= 2 {
					goFormat.WriteString("02")
					i++
				}
				continue
			}
		}

		if mode == TimeOnlyFormatting || mode == DateTimeFormatting {
			switch r {
			case 'f':
				if count == 9 {
					goFormat.WriteString("000000000")
					i += 8
				} else if count == 6 {
					goFormat.WriteString("000000")
					i += 5
				} else if count == 3 {
					goFormat.WriteString("000")
					i += 2
				} else {
					return "", errors.Errorf("invalid date format, invalid count of 'f' format: %d", count)
				}
				continue

			case 'h':
				if count == 1 {
					goFormat.WriteString("3")
				} else if count == 2 {
					goFormat.WriteString("03")
					i++
				}
				continue

			case 't':
				if count == 2 {
					goFormat.WriteString("15")
					i++
				} else {
					return "", errors.Errorf("invalid date format, invalid count of 't' format: %d", count)
				}
				continue

			case 'm':
				if count == 1 {
					goFormat.WriteString("4")
				} else if count == 2 {
					goFormat.WriteString("04")
					i++
				} else {
					return "", errors.Errorf("invalid date format, invalid count of 'm' format: %d", count)
				}
				continue

			case 's':
				if count == 1 {
					goFormat.WriteString("5")
				} else if count == 2 {
					goFormat.WriteString("05")
					i++
				} else {
					return "", errors.Errorf("invalid date format, invalid count of 's' format: %d", count)
				}
				continue

			case 'a':
				if count == 2 {
					goFormat.WriteString("pm")
					i++
				} else {
					return "", errors.Errorf("invalid date format, invalid count of 'a' format: %d", count)
				}
				continue

			case 'A':
				if count == 2 {
					goFormat.WriteString("PM")
					i++
				} else {
					return "", errors.Errorf("invalid date format, invalid count of 'A' format: %d", count)
				}
				continue
			}
		}

		if mode == DateTimeFormatting {
			switch r {
			case 'Z':
				if count == 1 {
					goFormat.WriteString("Z07:00")
				} else if count == 3 {
					goFormat.WriteString("-07:00")
					i += 2
				} else {
					return "", errors.Errorf("invalid date format, invalid count of 'Z' format: %d", count)
				}
				continue
			}
		}

		if ignoredFormattingRunes[r] {
			goFormat.WriteRune(r)
		} else {
			return "", errors.Errorf("invalid date format, unknown format char: %c", r)
		}
	}

	return goFormat.String(), nil
}
