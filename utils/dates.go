package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dd_mm_yyyy = regexp.MustCompile(`([0-9]{1,2})[-.\\/_]([0-9]{1,2})[-.\\/_]([0-9]{4})`)
var dd_mm_yy = regexp.MustCompile(`([0-9]{1,2})[-.\\/_]([0-9]{1,2})[-.\\/_]([0-9]{2})`)
var mm_dd_yyyy = regexp.MustCompile(`([0-9]{1,2})[-.\\/_]([0-9]{1,2})[-.\\/_]([0-9]{4})`)
var mm_dd_yy = regexp.MustCompile(`([0-9]{1,2})[-.\\/_]([0-9]{1,2})[-.\\/_]([0-9]{2})`)
var yyyy_mm_dd = regexp.MustCompile(`([0-9]{4})[-.\\/_]([0-9]{1,2})[-.\\/_]([0-9]{1,2})`)
var yy_mm_dd = regexp.MustCompile(`([0-9]{1,2})[-.\\/_]([0-9]{1,2})[-.\\/_]([0-9]{1,2})`)
var hh_mm_ss = regexp.MustCompile(`([0-9]{1,2}):([0-9]{2})(:([0-9]{2})(\.(\d+))?)?\W*([aApP][mM])?`)

type DateFormat string
type TimeFormat string

const (
	YYYY_MM_DD DateFormat = "YYYY-MM-DD"
	MM_DD_YYYY DateFormat = "MM-DD-YYYY"
	DD_MM_YYYY DateFormat = "DD-MM-YYYY"

	HH_MM       TimeFormat = "hh:mm"
	HH_MM_AP    TimeFormat = "hh:mm ap"
	HH_MM_SS    TimeFormat = "hh:mm:ss"
	HH_MM_SS_AP TimeFormat = "hh:mm:ss ap"
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
		if year >= currentYear%1000 {
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

// DateToString converts the passed in time element to the right format based on the environment settings
func DateToString(env Environment, date time.Time) string {
	var buf bytes.Buffer

	switch env.DateFormat() {
	case DD_MM_YYYY:
		buf.WriteString(fmt.Sprintf("%02d-%02d-%04d", date.Day(), date.Month(), date.Year()))

	case MM_DD_YYYY:
		buf.WriteString(fmt.Sprintf("%02d-%02d-%04d", date.Month(), date.Day(), date.Year()))

	case YYYY_MM_DD:
		buf.WriteString(fmt.Sprintf("%04d-%02d-%02d", date.Year(), date.Month(), date.Day()))
	}

	amPM := "AM"

	// write hour and minute
	switch env.TimeFormat() {
	case HH_MM, HH_MM_SS:
		buf.WriteString(fmt.Sprintf(" %02d:%02d", date.Hour(), date.Minute()))

	case HH_MM_AP, HH_MM_SS_AP:
		hour := date.Hour()
		if hour > 12 {
			hour -= 12
			amPM = "PM"
		}
		buf.WriteString(fmt.Sprintf(" %02d:%02d", hour, date.Minute()))
	}

	// write seconds if appropriate
	switch env.TimeFormat() {
	case HH_MM_SS, HH_MM_SS_AP:
		buf.WriteString(fmt.Sprintf(":%02d", date.Second()))
	}

	// write AM/PM if appropriate
	switch env.TimeFormat() {
	case HH_MM_AP, HH_MM_SS_AP:
		buf.WriteString(fmt.Sprintf(" %s", amPM))
	}

	return buf.String()
}

// DateFromString returns a date constructed from the passed in string, or an error if we
// are unable to extract one
func DateFromString(env Environment, str string) (time.Time, error) {
	currentYear := time.Now().Year()
	parsed := ZeroTime
	var err error
	switch env.DateFormat() {

	case DD_MM_YYYY:
		parsed, err = dateFromFormats(env, currentYear, dd_mm_yyyy, dd_mm_yy, 1, 2, 3, str)

	case MM_DD_YYYY:
		parsed, err = dateFromFormats(env, currentYear, mm_dd_yyyy, mm_dd_yy, 2, 1, 3, str)

	case YYYY_MM_DD:
		parsed, err = dateFromFormats(env, currentYear, yyyy_mm_dd, yy_mm_dd, 3, 2, 1, str)

	default:
		err = fmt.Errorf("Unknown date format: %s", env.DateFormat())
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

	return parsed, nil
}
