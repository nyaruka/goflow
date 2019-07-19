package dates_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils/dates"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatISO(t *testing.T) {
	la, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	assert.Equal(t, "2018-01-01T12:30:00.000000Z", dates.FormatISO(time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC)))
	assert.Equal(t, "2018-01-01T12:30:12.123456Z", dates.FormatISO(time.Date(2018, 1, 1, 12, 30, 12, 123456789, time.UTC)))
	assert.Equal(t, "2018-01-01T12:30:00.000000-08:00", dates.FormatISO(time.Date(2018, 1, 1, 12, 30, 0, 0, la)))
}

func TestDaysBetween(t *testing.T) {
	la, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	daysBetweenTests := []struct {
		d1       time.Time
		d2       time.Time
		expected int
	}{
		{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 3, 0, 30, 0, 0, time.UTC), -2},
		{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 3, 0, 30, 0, 0, la), -2},
		{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2017, 12, 25, 0, 30, 0, 0, time.UTC), 7},
		{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), 0},
	}

	for _, test := range daysBetweenTests {
		actual := dates.DaysBetween(test.d1, test.d2)
		assert.Equal(t, test.expected, actual, "mismatch for inputs %s - %s", test.d1, test.d2)
	}
}

func TestMonthsBetween(t *testing.T) {
	la, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)

	monthsBetweenTests := []struct {
		d1       time.Time
		d2       time.Time
		expected int
	}{
		{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 3, 0, 30, 0, 0, time.UTC), 0},
		{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 3, 3, 0, 30, 0, 0, la), -2},
		{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2017, 12, 25, 0, 30, 0, 0, time.UTC), 1},
		{time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC), 0},
	}

	for _, test := range monthsBetweenTests {
		actual := dates.MonthsBetween(test.d1, test.d2)
		assert.Equal(t, test.expected, actual, "mismatch for inputs %s - %s", test.d1, test.d2)
	}
}

func TestDayToUTCRange(t *testing.T) {
	ec, err := time.LoadLocation("America/Bogota")
	require.NoError(t, err)

	start, end := dates.DayToUTCRange(time.Date(2019, 7, 19, 9, 41, 0, 0, ec), time.UTC)
	assert.Equal(t, time.Date(2019, 7, 19, 5, 0, 0, 0, time.UTC), start)
	assert.Equal(t, time.Date(2019, 7, 20, 5, 0, 0, 0, time.UTC), end)
}
