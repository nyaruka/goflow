package test

import (
	"time"

	"github.com/nyaruka/goflow/utils"
)

// TestEnvironment an extended environment that will let us override Now() so that it's constant
type TestEnvironment struct {
	utils.Environment

	now time.Time
}

// NewTestEnvironment creates a new test environment
func NewTestEnvironment(dateFormat utils.DateFormat, tz *time.Location, now *time.Time) utils.Environment {
	if now == nil {
		t := time.Date(2018, 4, 11, 13, 24, 30, 123456000, tz)
		now = &t
	}

	return &TestEnvironment{
		Environment: utils.NewEnvironment(dateFormat, utils.TimeFormatHourMinute, tz, utils.LanguageList{"eng", "spa"}),
		now:         *now,
	}
}

func (e *TestEnvironment) Now() time.Time {
	return e.now.In(e.Timezone())
}

func (e *TestEnvironment) SetNow(now time.Time) {
	e.now = now
}
