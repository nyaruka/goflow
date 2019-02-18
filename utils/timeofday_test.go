package utils_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestTimeOfDay(t *testing.T) {
	t1 := utils.NewTimeOfDay(9, 38, 30, 123456789)

	assert.Equal(t, t1.Hour, 9)
	assert.Equal(t, t1.Minute, 38)
	assert.Equal(t, t1.Second, 30)
	assert.Equal(t, t1.Nanos, 123456789)
	assert.Equal(t, "09:38:30.123456", t1.String())

	t2 := utils.NewTimeOfDay(14, 56, 15, 0)

	assert.Equal(t, t2.Hour, 14)
	assert.Equal(t, t2.Minute, 56)
	assert.Equal(t, t2.Second, 15)
	assert.Equal(t, t2.Nanos, 0)
	assert.Equal(t, "14:56:15.000000", t2.String())

	// differs from t1 by 1 nano
	t3 := utils.NewTimeOfDay(9, 38, 30, 123456781)

	assert.Equal(t, t3.Hour, 9)
	assert.Equal(t, t3.Minute, 38)
	assert.Equal(t, t3.Second, 30)
	assert.Equal(t, t3.Nanos, 123456781)
	assert.Equal(t, "09:38:30.123456", t3.String())

	// should be same time value as t1
	t4 := utils.ExtractTimeOfDay(time.Date(2019, 2, 18, 9, 38, 30, 123456789, time.UTC))

	assert.Equal(t, t4.Hour, 9)
	assert.Equal(t, t4.Minute, 38)
	assert.Equal(t, t4.Second, 30)
	assert.Equal(t, t4.Nanos, 123456789)
	assert.Equal(t, "09:38:30.123456", t4.String())

	assert.False(t, t1.Equals(t2))
	assert.False(t, t2.Equals(t1))
	assert.False(t, t1.Equals(t3))
	assert.False(t, t3.Equals(t1))
	assert.True(t, t1.Equals(t4))
	assert.True(t, t4.Equals(t1))

	assert.True(t, t1.Compare(t2) < 0)
	assert.True(t, t2.Compare(t1) > 0)
	assert.True(t, t1.Compare(t3) > 0)
	assert.True(t, t3.Compare(t1) < 0)
	assert.True(t, t1.Compare(t4) == 0)
	assert.True(t, t4.Compare(t1) == 0)

	parsed, err := utils.ParseTimeOfDay("15:04:05.000000", "11:02:30.123456")
	assert.NoError(t, err)
	assert.Equal(t, utils.NewTimeOfDay(11, 2, 30, 123456000), parsed)
}
