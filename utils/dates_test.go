package utils

import (
	"fmt"
	"testing"
	"time"
)

var timeTests = []struct {
	DateFormat string
	TimeFormat string
	Timezone   string
	Value      string
	Expected   string
	Error      bool
}{
	// valid cases, varying formats
	{"DD-MM-YYYY", "hh:mm", "UTC", "01-02-2001", "01-02-2001 00:00:00 +0000 UTC", false},
	{"DD-MM-YYYY", "hh:mm", "UTC", "date is 01.02.2001 yes", "01-02-2001 00:00:00 +0000 UTC", false},
	{"DD-MM-YYYY", "hh:mm", "UTC", "date is 1-2-99 yes", "01-02-1999 00:00:00 +0000 UTC", false},
	{"DD-MM-YYYY", "hh:mm", "UTC", "01/02/2001", "01-02-2001 00:00:00 +0000 UTC", false},

	// month first
	{"MM-DD-YYYY", "hh:mm", "UTC", "01-02-2001", "02-01-2001 00:00:00 +0000 UTC", false},

	// year first
	{"YYYY-MM-DD", "hh:mm", "UTC", "2001-02-01", "01-02-2001 00:00:00 +0000 UTC", false},

	// specific timezone
	{"DD-MM-YYYY", "hh:mm", "America/Los_Angeles", "01\\02\\2001", "01-02-2001 00:00:00 -0800 PST", false},

	// illegal day
	{"DD-MM-YYYY", "hh:mm", "UTC", "33-01-2001", "01-01-0001 00:00:00 +0000 UTC", true},

	// illegal month
	{"DD-MM-YYYY", "hh:mm", "UTC", "01-13-2001", "01-01-0001 00:00:00 +0000 UTC", true},

	// valid two digit cases
	{"DD-MM-YYYY", "hh:mm", "UTC", "01-01-99", "01-01-1999 00:00:00 +0000 UTC", false},
	{"DD-MM-YYYY", "hh:mm", "UTC", "01-01-16", "01-01-2016 00:00:00 +0000 UTC", false},

	// with time
	{"DD-MM-YYYY", "hh:mm", "UTC", "2001-02-01 03:15", "01-02-2001 03:15:00 +0000 UTC", false},
	{"DD-MM-YYYY", "hh:mm", "UTC", "2001-02-01 03:15pm", "01-02-2001 15:15:00 +0000 UTC", false},
	{"DD-MM-YYYY", "hh:mm", "UTC", "2001-02-01 03:15 AM", "01-02-2001 03:15:00 +0000 UTC", false},
	{"DD-MM-YYYY", "hh:mm", "UTC", "2001-02-01 03:15:34", "01-02-2001 03:15:34 +0000 UTC", false},
	{"DD-MM-YYYY", "hh:mm", "UTC", "2001-02-01 03:15:34.123", "01-02-2001 03:15:34.123 +0000 UTC", false},
	{"DD-MM-YYYY", "hh:mm", "UTC", "2001-02-01 03:15:34.123456", "01-02-2001 03:15:34.123456 +0000 UTC", false},
}

func TestDateFromString(t *testing.T) {
	env := NewEnvironment(DD_MM_YYYY, HH_MM_SS, time.UTC)

	for _, test := range timeTests {
		env.SetDateFormat(DateFormat(test.DateFormat))
		env.SetTimeFormat(TimeFormat(test.TimeFormat))
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
			fmt.Printf("value: %s  expected: %s\n", value.UTC(), expected.UTC())
			t.Errorf("Date '%s' not match expected date '%s'", value, expected)
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
