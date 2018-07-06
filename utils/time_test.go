package utils_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestTimeSources(t *testing.T) {
	d1 := time.Date(2018, 7, 5, 16, 29, 30, 123456, time.UTC)
	utils.SetTimeSource(utils.NewFixedTimeSource(d1))
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 30, 123456, time.UTC), utils.Now())
	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 30, 123456, time.UTC), utils.Now())

	utils.SetTimeSource(utils.NewSequentialTimeSource(d1))

	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 30, 123456, time.UTC), utils.Now())
	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 31, 123456, time.UTC), utils.Now())
	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 32, 123456, time.UTC), utils.Now())
}
