package dates_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/dates"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	d1 := dates.NewDate(2019, 2, 20)

	assert.Equal(t, d1.Year, 2019)
	assert.Equal(t, d1.Month, time.Month(2))
	assert.Equal(t, d1.Day, 20)
	assert.Equal(t, d1.Weekday(), time.Weekday(3))
	assert.Equal(t, "2019-02-20", d1.String())

	d2 := dates.NewDate(2020, 1, 1)

	assert.Equal(t, d2.Year, 2020)
	assert.Equal(t, d2.Month, time.Month(1))
	assert.Equal(t, d2.Day, 1)
	assert.Equal(t, "2020-01-01", d2.String())

	// differs from d1 by 1 day
	d3 := dates.NewDate(2019, 2, 19)

	assert.Equal(t, d3.Year, 2019)
	assert.Equal(t, d3.Month, time.Month(2))
	assert.Equal(t, d3.Day, 19)
	assert.Equal(t, "2019-02-19", d3.String())

	// should be same date value as d1
	d4 := dates.ExtractDate(time.Date(2019, 2, 20, 9, 38, 30, 123456789, time.UTC))

	assert.Equal(t, d4.Year, 2019)
	assert.Equal(t, d4.Month, time.Month(2))
	assert.Equal(t, d4.Day, 20)
	assert.Equal(t, "2019-02-20", d4.String())

	assert.False(t, d1.Equal(d2))
	assert.False(t, d2.Equal(d1))
	assert.False(t, d1.Equal(d3))
	assert.False(t, d3.Equal(d1))
	assert.True(t, d1.Equal(d4))
	assert.True(t, d4.Equal(d1))

	assert.True(t, d1.Compare(d2) < 0)
	assert.True(t, d2.Compare(d1) > 0)
	assert.True(t, d1.Compare(d3) > 0)
	assert.True(t, d3.Compare(d1) < 0)
	assert.True(t, d1.Compare(d4) == 0)
	assert.True(t, d4.Compare(d1) == 0)

	parsed, err := dates.ParseDate("2006-01-02", "2018-12-30")
	assert.NoError(t, err)
	assert.Equal(t, dates.NewDate(2018, 12, 30), parsed)
}
