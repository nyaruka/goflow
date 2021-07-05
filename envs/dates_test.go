package envs_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateTimeFromString(t *testing.T) {

	testCases := []struct {
		DateFormat envs.DateFormat
		TimeFormat envs.TimeFormat
		Timezone   string
		FillTime   bool
		Value      string
		Expected   string
		Error      bool
	}{
		// valid cases, varying formats
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-02-2001", "01-02-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "date is 01.02.2001 yes", "01-02-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "date is 1-2-99 yes", "01-02-1999 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01/02/2001", "01-02-2001 00:00:00 +0000 UTC", false},

		// must be strict iso to match despite format
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001-01-02T10:34:56Z", "02-01-2001 10:34:56 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001-01-02T10:34:56+02:00", "02-01-2001 10:34:56 +0200 MST", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001-01-02", "02-01-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", true, "2001-01-02", "02-01-2001 13:36:30.123456789 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "America/Los_Angeles", true, "2001-01-02", "02-01-2001 06:36:30.123456789 -0800 PST", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, " 2001-01-02 ", "02-01-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "on 2001-01-02 ", "", true},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001_01_02", "", true},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2001-1-2", "", true},

		// month first
		{envs.DateFormatMonthDayYear, envs.TimeFormatHourMinute, "UTC", false, "01-02-2001", "02-01-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatMonthDayYear, envs.TimeFormatHourMinute, "UTC", false, "2001-01-02", "02-01-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatMonthDayYear, envs.TimeFormatHourMinute, "UTC", false, "2001-1-2", "", true},

		// year first
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01", "01-02-2001 00:00:00 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "99-02-01", "01-02-1999 00:00:00 +0000 UTC", false},

		// specific timezone
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "America/Los_Angeles", false, "01\\02\\2001", "01-02-2001 00:00:00 -0800 PST", false},

		// with time filling
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", true, "01-02-2001", "01-02-2001 13:36:30.123456789 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", true, "01-02-2001 04:23", "01-02-2001 04:23:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "America/Los_Angeles", true, "01-02-2001", "01-02-2001 06:36:30.123456789 -0800 PST", false},

		// illegal day
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "33-01-2001", "01-01-0001 00:00:00 +0000 UTC", true},

		// illegal month
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-13-2001", "01-01-0001 00:00:00 +0000 UTC", true},

		// valid two digit cases
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-01-99", "01-01-1999 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-01-16", "01-01-2016 00:00:00 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "01-01-16a", "", true},

		// iso dates
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2016-05-01T18:30:15-08:00", "01-05-2016 18:30:15 -0800 PST", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2016-05-01T18:30:15Z", "01-05-2016 18:30:15 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2016-05-01T18:30:15.250Z", "01-05-2016 18:30:15.250 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "1977-06-23T08:34:00.000-07:00", "23-06-1977 15:34:00.000 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "1977-06-23T08:34:00.000250-07:00", "23-06-1977 15:34:00.000250 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "1977-06-23T08:34:00.000250500-07:00", "23-06-1977 15:34:00.000250500 +0000 UTC", false},
		{envs.DateFormatDayMonthYear, envs.TimeFormatHourMinute, "UTC", false, "2017-06-10T17:34-06:00", "10-06-2017 23:34:00.000000 +0000 UTC", false},

		// with time
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15", "01-02-2001 03:15:00 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15pm", "01-02-2001 15:15:00 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15 AM", "01-02-2001 03:15:00 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15:34", "01-02-2001 03:15:34 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15:34.123", "01-02-2001 03:15:34.123 +0000 UTC", false},
		{envs.DateFormatYearMonthDay, envs.TimeFormatHourMinute, "UTC", false, "2001-02-01 03:15:34.123456", "01-02-2001 03:15:34.123456 +0000 UTC", false},
	}

	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 9, 13, 13, 36, 30, 123456789, time.UTC)))
	defer dates.SetNowSource(dates.DefaultNowSource)

	for _, tc := range testCases {
		timezone, err := time.LoadLocation(tc.Timezone)
		require.NoError(t, err)

		env := envs.NewBuilder().WithDateFormat(tc.DateFormat).WithTimeFormat(tc.TimeFormat).WithTimezone(timezone).Build()

		value, err := envs.DateTimeFromString(env, tc.Value, tc.FillTime)

		if tc.Error {
			assert.Error(t, err)
		} else {
			require.NoError(t, err, "error parsing date %s", tc.Value)

			expected, err := time.Parse("02-01-2006 15:04:05.999999999 -0700 MST", tc.Expected)
			require.NoError(t, err, "error parsing expected date %s", tc.Expected)

			if !expected.Equal(value) {
				assert.Fail(t, "", "mismatch for date input %s, expected %s, got %s", tc.Value, expected, value)
			}
		}
	}
}

func TestDateFromString(t *testing.T) {
	testCases := []struct {
		dateFormat envs.DateFormat
		value      string
		expected   dates.Date
		hasError   bool
	}{
		{envs.DateFormatDayMonthYear, "it's 31-12-2018 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatDayMonthYear, "it's 31-12-18 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatMonthDayYear, "it's 12-31-2018 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatMonthDayYear, "it's 12-31-18 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatYearMonthDay, "it's 2018-12-31 ok", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatYearMonthDay, "it's 18-12-31 ok", dates.NewDate(2018, 12, 31), false},

		// valid ISO always accepted
		{envs.DateFormatDayMonthYear, "2018-12-31", dates.NewDate(2018, 12, 31), false},
		{envs.DateFormatDayMonthYear, "2018-12-31T18:30:15-08:00", dates.NewDate(2018, 12, 31), false},

		{envs.DateFormatDayMonthYear, "it's ok", dates.ZeroDate, true},
		{envs.DateFormatDayMonthYear, "it's 2018-13-01 ok", dates.ZeroDate, true},
		{envs.DateFormatDayMonthYear, "it's 2018-12-32 ok", dates.ZeroDate, true},
	}

	for _, tc := range testCases {
		env := envs.NewBuilder().WithDateFormat(tc.dateFormat).Build()
		parsed, err := envs.DateFromString(env, tc.value)

		if tc.hasError {
			assert.Error(t, err)
		} else {
			require.NoError(t, err, "error parsing time %s", tc.value)

			if !tc.expected.Equal(parsed) {
				assert.Fail(t, "", "mismatch for time input %s, expected %s, got %s", tc.value, tc.expected, parsed)
			}
		}
	}
}

func TestTimeFromString(t *testing.T) {
	testCases := []struct {
		value    string
		expected dates.TimeOfDay
		hasError bool
	}{
		{"it's 10 ok", dates.NewTimeOfDay(10, 0, 0, 0), false},
		{"it's 10 PM ok", dates.NewTimeOfDay(22, 0, 0, 0), false},
		{"it's 10:30 ok", dates.NewTimeOfDay(10, 30, 0, 0), false},
		{"it's 10:30pm ok", dates.NewTimeOfDay(22, 30, 0, 0), false},
		{"it's 10:30 pm ok", dates.NewTimeOfDay(22, 30, 0, 0), false},
		{"it's 1030 ok", dates.NewTimeOfDay(10, 30, 0, 0), false},
		{"it's 1030 PM ok", dates.NewTimeOfDay(22, 30, 0, 0), false},
		{"it's 10:30:45 ok", dates.NewTimeOfDay(10, 30, 45, 0), false},
		{"it's 10:30:45 pm ok", dates.NewTimeOfDay(22, 30, 45, 0), false},
		{"it's 10:30:45.123 ok", dates.NewTimeOfDay(10, 30, 45, 123000000), false},
		{"it's 10:30:45.123 pm ok", dates.NewTimeOfDay(22, 30, 45, 123000000), false},
		{"it's 10:30:45.123456 ok", dates.NewTimeOfDay(10, 30, 45, 123456000), false},
		{"it's 10:30:45.123456 pm ok", dates.NewTimeOfDay(22, 30, 45, 123456000), false},
		{"it's 10:30:45.123456789 ok", dates.NewTimeOfDay(10, 30, 45, 123456789), false},
		{"it's 10:30:45.123456789 pm ok", dates.NewTimeOfDay(22, 30, 45, 123456789), false},

		// 12 am and 12 pm
		{"it's 12:00 AM", dates.NewTimeOfDay(0, 0, 0, 0), false},
		{"it's 12pm", dates.NewTimeOfDay(12, 0, 0, 0), false},

		// fractional component can be any length
		{"it's 10:30:45.123456789123456789 ok", dates.NewTimeOfDay(10, 30, 45, 123456789), false},
		{"it's 10:30:45.1 ok", dates.NewTimeOfDay(10, 30, 45, 100000000), false},

		// 24 can be used to mean midnight
		{"it's 24:00 ok", dates.NewTimeOfDay(0, 0, 0, 0), false},
		{"it's 24:00:00 ok", dates.NewTimeOfDay(0, 0, 0, 0), false},

		{"it's ok", dates.ZeroTimeOfDay, true},
		{"it's 25:30", dates.ZeroTimeOfDay, true},
		{"it's 10:61", dates.ZeroTimeOfDay, true},
		{"it's 10:30:61", dates.ZeroTimeOfDay, true},
	}

	for _, tc := range testCases {
		parsed, err := envs.TimeFromString(tc.value)

		if tc.hasError {
			assert.Error(t, err)
		} else {
			require.NoError(t, err, "error parsing time %s", tc.value)

			if !tc.expected.Equal(parsed) {
				assert.Fail(t, "", "mismatch for time input %s, expected %s, got %s", tc.value, tc.expected, parsed)
			}
		}
	}
}
