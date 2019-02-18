package utils

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	Validator.RegisterValidation("date_format", func(fl validator.FieldLevel) bool {
		_, err := ToGoDateFormat(fl.Field().String(), DateOnlyFormatting)
		return err == nil
	})
	Validator.RegisterValidation("time_format", func(fl validator.FieldLevel) bool {
		_, err := ToGoDateFormat(fl.Field().String(), TimeOnlyFormatting)
		return err == nil
	})
}

// patterns for date and time formats supported for human-entered data
var patternDayMonthYear = regexp.MustCompile(`\b([0-9]{1,2})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{4}|[0-9]{2})\b`)
var patternMonthDayYear = regexp.MustCompile(`\b([0-9]{1,2})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{4}|[0-9]{2})\b`)
var patternYearMonthDay = regexp.MustCompile(`\b([0-9]{4}|[0-9]{2})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{1,2})\b`)

var patternTime = regexp.MustCompile(`\b([0-9]{1,2}):([0-9]{2})(:([0-9]{2})(\.(\d+))?)?\W*([aApP][mM])?\b`)

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

// ZeroDateTime is our uninitialized datetime value
var ZeroDateTime = time.Time{}

func dateFromFormats(env Environment, currentYear int, pattern *regexp.Regexp, d int, m int, y int, str string) (time.Time, error) {

	matches := pattern.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		// does our day look believable?
		day, _ := strconv.Atoi(match[d])
		if day == 0 || day > 31 {
			continue
		}
		month, _ := strconv.Atoi(match[m])
		if month == 0 || month > 12 {
			continue
		}

		year, _ := strconv.Atoi(match[y])

		// convert to four digit year if necessary
		if len(match[y]) == 2 {
			if year > currentYear%1000 {
				year += 1900
			} else {
				year += 2000
			}
		}

		// looks believable, go for it
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, env.Timezone()), nil
	}

	return ZeroDateTime, errors.Errorf("string '%s' couldn't be parsed as a date", str)
}

// DaysBetween returns the number of calendar days (an int) between the two dates. Note
// that if these are in different timezones then the local calendar day is used for each
// and the difference is calculated from that.
func DaysBetween(date1 time.Time, date2 time.Time) int {
	d1 := time.Date(date1.Year(), date1.Month(), date1.Day(), 0, 0, 0, 0, time.UTC)
	d2 := time.Date(date2.Year(), date2.Month(), date2.Day(), 0, 0, 0, 0, time.UTC)

	return int(d1.Sub(d2) / (time.Hour * 24))
}

// MonthsBetween returns the number of calendar months (an int) between the two dates. Note
// that if these are in different timezones then the local calendar day is used for each
// and the difference is calculated from that.
func MonthsBetween(date1 time.Time, date2 time.Time) int {
	// difference in months
	months := int(date1.Month() - date2.Month())

	// difference in years
	months += (date1.Year() - date2.Year()) * 12

	return months
}

// format we use for output
var iso8601Default = "2006-01-02T15:04:05.000000Z07:00"

// generic format for parsing any 8601 date
var iso8601Format = "2006-01-02T15:04:05Z07:00"
var iso8601NoSecondsFormat = "2006-01-02T15:04Z07:00"
var iso8601Date = "2006-01-02"

var isoFormats = []string{iso8601Format, iso8601NoSecondsFormat, iso8601Date}

// DateToISO converts the passed in time.Time to a string in ISO8601 format
func DateToISO(date time.Time) string {
	return date.Format(iso8601Default)
}

// DateToString converts the passed in time element to the right format based on the environment settings
func DateToString(env Environment, date time.Time) string {
	goFormat, _ := ToGoDateFormat(string(env.DateFormat())+" "+string(env.TimeFormat()), DateTimeFormatting)
	return date.Format(goFormat)
}

// DateFromString returns a date constructed from the passed in string, or an error if we
// are unable to extract one
func DateFromString(env Environment, str string, fillTime bool) (time.Time, error) {
	// first see if we can parse in any known iso formats, if so return that
	for _, format := range isoFormats {
		parsed, err := time.ParseInLocation(format, strings.Trim(str, " \n\r\t"), env.Timezone())
		if err == nil {
			if format == iso8601Date && fillTime {
				parsed = replaceTime(parsed, Now().In(env.Timezone()))
			}
			return parsed, nil
		}
	}

	// otherwise, try to parse according to their env settings
	parsed := ZeroDateTime
	var err error
	currentYear := Now().Year()

	switch env.DateFormat() {
	case DateFormatYearMonthDay:
		parsed, err = dateFromFormats(env, currentYear, patternYearMonthDay, 3, 2, 1, str)
	case DateFormatDayMonthYear:
		parsed, err = dateFromFormats(env, currentYear, patternDayMonthYear, 1, 2, 3, str)
	case DateFormatMonthDayYear:
		parsed, err = dateFromFormats(env, currentYear, patternMonthDayYear, 2, 1, 3, str)
	default:
		err = errors.Errorf("unknown date format: %s", env.DateFormat())
	}

	// couldn't find a date? bail
	if err != nil {
		return parsed, err
	}

	// can we pull out a time?
	hasTime, timeOfDay := parseTime(str)
	if hasTime {
		parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), timeOfDay.Hour, timeOfDay.Minute, timeOfDay.Second, timeOfDay.Nanos, env.Timezone())
	} else if fillTime {
		parsed = replaceTime(parsed, Now().In(env.Timezone()))
	}

	// set our timezone if we have one
	if env.Timezone() != nil && parsed != ZeroDateTime {
		parsed = parsed.In(env.Timezone())
	}

	return parsed, nil
}

// TimeFromString returns a time of day constructed from the passed in string, or an error if we
// are unable to extract one
func TimeFromString(env Environment, str string) (TimeOfDay, error) {
	hasTime, timeOfDay := parseTime(str)
	if !hasTime {
		return ZeroTimeOfDay, errors.Errorf("string '%s' couldn't be parsed as a time", str)
	}
	return timeOfDay, nil
}

func replaceTime(d time.Time, t time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

func parseTime(str string) (bool, TimeOfDay) {
	matches := patternTime.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		hour, _ := strconv.Atoi(match[1])

		// do we have an AM/PM
		if match[7] != "" {
			if strings.ToLower(match[7]) == "pm" {
				hour += 12
			}
		}

		// is this a valid hour?
		if hour > 24 {
			continue
		}

		minute, _ := strconv.Atoi(match[2])
		if minute > 60 {
			continue
		}

		second := 0
		if match[4] != "" {
			second, _ = strconv.Atoi(match[4])
			if second > 60 {
				continue
			}
		}

		ns := 0
		if match[6] != "" {
			ns, _ = strconv.Atoi(match[6])

			if len(match[6]) == 3 {
				// these are milliseconds, multi by 1,000,000 for nano
				ns = ns * 1000000
			} else if len(match[6]) == 6 {
				// these are microseconds, times 1000 for nano
				ns = ns * 1000
			}
		}

		return true, NewTimeOfDay(hour, minute, second, ns)
	}

	return false, ZeroTimeOfDay
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
//  `tt`        - twenty four hour of the day 01-23
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

// DateToUTCRange returns the UTC time range of the given day
func DateToUTCRange(d time.Time, tz *time.Location) (time.Time, time.Time) {
	localMidnight := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	utcMidnight := localMidnight.In(tz)
	return utcMidnight, utcMidnight.Add(24 * time.Hour)
}
