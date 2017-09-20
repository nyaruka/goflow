package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dd_mm_yyyy = regexp.MustCompile(`([0-9]{1,4})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{4,5})`)
var dd_mm_yy = regexp.MustCompile(`([0-9]{1,4})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{2,3})`)
var mm_dd_yyyy = regexp.MustCompile(`([0-9]{1,4})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{4,5})`)
var mm_dd_yy = regexp.MustCompile(`([0-9]{1,4})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{2,3})`)
var yyyy_mm_dd = regexp.MustCompile(`([0-9]{4})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{1,4})`)
var yy_mm_dd = regexp.MustCompile(`([0-9]{1,4})[-.\\/_ ]([0-9]{1,2})[-.\\/_ ]([0-9]{1,4})`)
var hh_mm_ss = regexp.MustCompile(`([0-9]{1,4}):([0-9]{2})(:([0-9]{2})(\.(\d+))?)?\W*([aApP][mM])?`)

type DateFormat string
type TimeFormat string

const (
	DateFormat_yyyy_MM_dd DateFormat = "yyyy-MM-dd"
	DateFormat_MM_dd_yyyy DateFormat = "MM-dd-yyyy"
	DateFormat_dd_MM_yyyy DateFormat = "dd-MM-yyyy"

	TimeFormat_HH_mm       TimeFormat = "hh:mm"
	TimeFormat_hh_mm_tt    TimeFormat = "hh:mm tt"
	TimeFormat_HH_mm_ss    TimeFormat = "HH:mm:ss"
	TimeFormat_hh_mm_ss_tt TimeFormat = "hh:mm:ss tt"
)

func (df DateFormat) String() string { return string(df) }
func (tf TimeFormat) String() string { return string(tf) }

// ZeroTime is our uninitialized time value
var ZeroTime = time.Time{}

func dateFromFormats(env Environment, currentYear int, fourDigit *regexp.Regexp, twoDigit *regexp.Regexp,
	d int, m int, y int, str string) (time.Time, error) {

	// four digit year comes first
	matches := fourDigit.FindAllStringSubmatch(str, -1)
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

		// looks believable, let's return it
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, env.Timezone()), nil
	}

	// then two digit
	matches = twoDigit.FindAllStringSubmatch(str, -1)
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

		// convert to four digit year
		year, _ := strconv.Atoi(match[y])
		if year > currentYear%1000 {
			year += 1900
		} else {
			year += 2000
		}

		// looks believable, go for it
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, env.Timezone()), nil
	}

	return ZeroTime, fmt.Errorf("No date found in string: %s", str)
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

var isoFormats = []string{iso8601Format, iso8601NoSecondsFormat}

// DateToISO converts the passed in time.Time to a string in ISO8601 format
func DateToISO(date time.Time) string {
	return date.Format(iso8601Default)
}

// DateToString converts the passed in time element to the right format based on the environment settings
func DateToString(env Environment, date time.Time) string {
	goFormat, _ := ToGoDateFormat(string(env.DateFormat()) + " " + string(env.TimeFormat()))
	return date.Format(goFormat)
}

// DateFromString returns a date constructed from the passed in string, or an error if we
// are unable to extract one
func DateFromString(env Environment, str string) (time.Time, error) {
	// first see if we can parse in any known iso formats, if so return that
	for _, format := range isoFormats {
		parsed, err := time.Parse(format, str)
		if err == nil {
			if env.Timezone() != nil {
				parsed = parsed.In(env.Timezone())
			}
			return parsed, nil
		}
	}

	// otherwise, try to parse according to their env settings
	parsed := ZeroTime
	currentYear := time.Now().Year()
	var err error
	switch env.DateFormat() {

	case DateFormat_yyyy_MM_dd:
		parsed, err = dateFromFormats(env, currentYear, yyyy_mm_dd, yy_mm_dd, 3, 2, 1, str)

	case DateFormat_dd_MM_yyyy:
		parsed, err = dateFromFormats(env, currentYear, dd_mm_yyyy, dd_mm_yy, 1, 2, 3, str)

	case DateFormat_MM_dd_yyyy:
		parsed, err = dateFromFormats(env, currentYear, mm_dd_yyyy, mm_dd_yy, 2, 1, 3, str)

	default:
		err = fmt.Errorf("unknown date format: %s", env.DateFormat())
	}

	// couldn't find a date? bail
	if err != nil {
		return parsed, err
	}

	// can we pull out a time?
	matches := hh_mm_ss.FindAllStringSubmatch(str, -1)
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

		seconds := 0
		if match[4] != "" {
			seconds, _ = strconv.Atoi(match[4])
			if seconds > 60 {
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

		parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), hour, minute, seconds, ns, env.Timezone())
	}

	// set our timezone if we have one
	if env.Timezone() != nil && parsed != ZeroTime {
		parsed = parsed.In(env.Timezone())
	}

	return parsed, nil
}

// ToGoDateFormat converts the passed in format to a GoLang format string.
//
// Format strings we support:
//
// d       - day of month, 1-31
// dd      - day of month, zero padded 0-31
// fff     - thousandths of a second
// h       - hour of the day 1-12
// hh      - hour of the day 01-12
// H       - hour of the day 1-23
// HH      - hour of the day 01-23
// K       - hour and minute offset from UTC, or Z fo UTC
// m       - minute 0-59
// mm      - minute 00-59
// M       - month 1-12
// MM      - month 01-12
// s       - second 0-59
// ss      - second 00-59
// TT      - AM or PM
// tt      - am or pm
// yy      - last two digits of year 0-99
// yyyy    - four digits of your 0000-9999
// zzz     - hour and minute offset from UTC
// ignored chars: ' ', ':', ',', 'T', 'Z', '-', '_', '/'
func ToGoDateFormat(format string) (string, error) {
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

		switch r {
		case 'd':
			if count == 1 {
				goFormat.WriteString("2")
			} else if count >= 2 {
				goFormat.WriteString("02")
				i++
			}

		case 'f':
			if count >= 9 {
				goFormat.WriteString("000000000")
				i += 8
			} else if count >= 6 {
				goFormat.WriteString("000000")
				i += 5
			} else if count >= 3 {
				goFormat.WriteString("000")
				i += 2
			} else {
				return "", fmt.Errorf("invalid date format, invalid count of 'f' format: %d", count)
			}

		case 'h':
			if count == 1 {
				goFormat.WriteString("3")
			} else if count >= 2 {
				goFormat.WriteString("03")
				i++
			}

		case 'H':
			if count >= 2 {
				goFormat.WriteString("15")
				i++
			} else {
				return "", fmt.Errorf("invalid date format, invalid count of 'H' format: %d", count)
			}

		case 'K':
			goFormat.WriteString("Z07:00")

		case 'm':
			if count == 1 {
				goFormat.WriteString("4")
			} else if count >= 2 {
				goFormat.WriteString("04")
				i++
			}

		case 'M':
			if count == 1 {
				goFormat.WriteString("1")
			} else if count >= 2 {
				goFormat.WriteString("01")
				i++
			}

		case 's':
			if count == 1 {
				goFormat.WriteString("5")
			} else if count >= 2 {
				goFormat.WriteString("05")
				i++
			}

		case 't':
			if count >= 2 {
				goFormat.WriteString("pm")
				i++
			} else {
				return "", fmt.Errorf("invalid date format, invalid count of 't' format: %d", count)
			}

		case 'T':
			if count == 1 {
				goFormat.WriteString("T")
			} else if count >= 2 {
				goFormat.WriteString("PM")
				i++
			}

		case 'y':
			if count == 2 {
				goFormat.WriteString("06")
				i++
			} else if count >= 4 {
				goFormat.WriteString("2006")
				i += 3
			} else {
				return "", fmt.Errorf("invalid date format, invalid count of 'y' format: %d", count)
			}

		case 'z':
			if count == 3 {
				goFormat.WriteString("-07:00")
				i += 2
			} else {
				return "", fmt.Errorf("invalid date format, invalid count of 'z' format: %d", count)
			}

		case ' ', ':', '/', '.', 'Z', '-', '_':
			goFormat.WriteRune(r)

		default:
			return "", fmt.Errorf("invalid date format, unknown format char: %c", r)
		}
	}

	return goFormat.String(), nil
}

func DateToUTCRange(d time.Time, tz *time.Location) (time.Time, time.Time) {
	localMidnight := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	utcMidnight := localMidnight.In(tz)
	return utcMidnight, utcMidnight.Add(24 * time.Hour)
}
