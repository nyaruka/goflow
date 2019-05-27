package test_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestTimeSources(t *testing.T) {
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	d1 := time.Date(2018, 7, 5, 16, 29, 30, 123456, time.UTC)
	utils.SetTimeSource(test.NewFixedTimeSource(d1))

	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 30, 123456, time.UTC), utils.Now())
	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 30, 123456, time.UTC), utils.Now())

	utils.SetTimeSource(test.NewSequentialTimeSource(d1))

	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 30, 123456, time.UTC), utils.Now())
	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 31, 123456, time.UTC), utils.Now())
	assert.Equal(t, time.Date(2018, 7, 5, 16, 29, 32, 123456, time.UTC), utils.Now())
}
