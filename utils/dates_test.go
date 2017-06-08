package utils

import (
	"testing"
	"time"
)

var timeTests = []struct {
	DateFormat DateFormat
	TimeFormat TimeFormat
	Timezone   string
	Value      string
	Expected   string
	Error      bool
}{
	// valid cases, varying formats
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "01-02-2001", "01-02-2001 00:00:00 +0000 UTC", false},
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "date is 01.02.2001 yes", "01-02-2001 00:00:00 +0000 UTC", false},
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "date is 1-2-99 yes", "01-02-1999 00:00:00 +0000 UTC", false},
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "01/02/2001", "01-02-2001 00:00:00 +0000 UTC", false},

	// month first
	{DateFormat_MM_dd_yyyy, TimeFormat_HH_mm, "UTC", "01-02-2001", "02-01-2001 00:00:00 +0000 UTC", false},

	// year first
	{DateFormat_yyyy_MM_dd, TimeFormat_HH_mm, "UTC", "2001-02-01", "01-02-2001 00:00:00 +0000 UTC", false},

	// specific timezone
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "America/Los_Angeles", "01\\02\\2001", "01-02-2001 00:00:00 -0800 PST", false},

	// illegal day
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "33-01-2001", "01-01-0001 00:00:00 +0000 UTC", true},

	// illegal month
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "01-13-2001", "01-01-0001 00:00:00 +0000 UTC", true},

	// valid two digit cases
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "01-01-99", "01-01-1999 00:00:00 +0000 UTC", false},
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "01-01-16", "01-01-2016 00:00:00 +0000 UTC", false},

	// iso dates
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "2016-05-01T18:30:15-08:00", "01-05-2016 18:30:15 -0800 PST", false},
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "2016-05-01T18:30:15Z", "01-05-2016 18:30:15 -0000 UTC", false},
	{DateFormat_dd_MM_yyyy, TimeFormat_HH_mm, "UTC", "2016-05-01T18:30:15.250Z", "01-05-2016 18:30:15.250 -0000 UTC", false},

	// with time
	{DateFormat_yyyy_MM_dd, TimeFormat_HH_mm, "UTC", "2001-02-01 03:15", "01-02-2001 03:15:00 +0000 UTC", false},
	{DateFormat_yyyy_MM_dd, TimeFormat_HH_mm, "UTC", "2001-02-01 03:15pm", "01-02-2001 15:15:00 +0000 UTC", false},
	{DateFormat_yyyy_MM_dd, TimeFormat_HH_mm, "UTC", "2001-02-01 03:15 AM", "01-02-2001 03:15:00 +0000 UTC", false},
	{DateFormat_yyyy_MM_dd, TimeFormat_HH_mm, "UTC", "2001-02-01 03:15:34", "01-02-2001 03:15:34 +0000 UTC", false},
	{DateFormat_yyyy_MM_dd, TimeFormat_HH_mm, "UTC", "2001-02-01 03:15:34.123", "01-02-2001 03:15:34.123 +0000 UTC", false},
	{DateFormat_yyyy_MM_dd, TimeFormat_HH_mm, "UTC", "2001-02-01 03:15:34.123456", "01-02-2001 03:15:34.123456 +0000 UTC", false},
}

func TestDateFromString(t *testing.T) {
	env := NewEnvironment(DateFormat_dd_MM_yyyy, TimeFormat_HH_mm_ss, time.UTC)

	for _, test := range timeTests {
		env.SetDateFormat(test.DateFormat)
		env.SetTimeFormat(test.TimeFormat)
		timezone, err := time.LoadLocation(test.Timezone)
		env.SetTimezone(timezone)

		if err != nil {
			t.Errorf("Error parsing expected timezone: %s", err)
			continue
		}

		expected, err := time.Parse("02-01-2006 15:04:05 -0700 MST", test.Expected)
		if err != nil {
			t.Errorf("Error parsing expected date: %s", err)
			continue
		}

		value, err := DateFromString(env, test.Value)
		if err != nil && !test.Error {
			t.Errorf("Error parsing date: %s", err)
			continue
		}

		if !value.Equal(expected) {
			t.Errorf("Date '%s' not match expected date '%s' for input: '%s'", value, expected, test.Value)
		}
	}
}

var laTZ, _ = time.LoadLocation("America/Los_Angeles")

var daysBetweenTests = []struct {
	d1       time.Time
	d2       time.Time
	expected int
}{
	{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 3, 0, 30, 0, 0, time.UTC), -2},
	{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 3, 0, 30, 0, 0, laTZ), -2},
	{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2017, 12, 25, 0, 30, 0, 0, time.UTC), 7},
	{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), 0},
}

func TestDaysBetween(t *testing.T) {
	for _, test := range daysBetweenTests {
		actual := DaysBetween(test.d1, test.d2)
		if actual != test.expected {
			t.Errorf("Days between: %d did not match expected: %d for %s - %s", actual, test.expected, test.d1, test.d2)
		}
	}
}

var monthsBetweenTests = []struct {
	d1       time.Time
	d2       time.Time
	expected int
}{
	{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 3, 0, 30, 0, 0, time.UTC), 0},
	{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 3, 3, 0, 30, 0, 0, laTZ), -2},
	{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2017, 12, 25, 0, 30, 0, 0, time.UTC), 1},
	{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), 0},
}

func TestMonthsBetween(t *testing.T) {
	for _, test := range daysBetweenTests {
		actual := DaysBetween(test.d1, test.d2)
		if actual != test.expected {
			t.Errorf("Months between: %d did not match expected: %d for %s - %s", actual, test.expected, test.d1, test.d2)
		}
	}
}

func TestDateFormat(t *testing.T) {
	formatTests := []struct {
		input    string
		expected string
		hasErr   bool
	}{
		{"MM-dd-yyyy", "01-02-2006", false},
		{"M-d-yy", "1-2-06", false},
		{"h:m", "3:4", false},
		{"h:m:s tt", "3:4:5 PM", false},
		{"yyyy-MM-ddTHH:mm:sszzz", "2006-01-02T15:04:05-0700", false},
		{"yyyy-MM-ddTHH:mm:sszzz", "2006-01-02T15:04:05-0700", false},
		{"yyyy-MM-ddThh:mm:ss.fffzzz", "2006-01-02T03:04:05.000-0700", false},
		{"yyyy-MM-dd", "2006-01-02", false},
		{"2006-01-02", "", true},
	}

	for _, test := range formatTests {
		actual, err := ToGoDateFormat(test.input)
		if actual != test.expected {
			t.Errorf("Date format invalid for '%s'  Expected: '%s' Got: '%s'", test.input, test.expected, actual)
		}

		if err != nil && !test.hasErr {
			t.Errorf("Date format received error for '%s': %s", test.input, err)
		}

		if err == nil && test.hasErr {
			t.Errorf("Date format expected error for '%s'", test.input)
		}
	}
}
