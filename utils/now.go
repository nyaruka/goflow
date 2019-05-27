package utils

import (
	"time"
)

// TimeSource is something that can provide a time
type TimeSource interface {
	Now() time.Time
}

// defaultTimeSource returns the current system time
type defaultTimeSource struct{}

func (s defaultTimeSource) Now() time.Time {
	return time.Now()
}

// DefaultTimeSource is the default time source
var DefaultTimeSource TimeSource = defaultTimeSource{}
var currentTimeSource = DefaultTimeSource

// Now returns the current time
func Now() time.Time {
	return currentTimeSource.Now()
}

// SetTimeSource sets the time source used by Now()
func SetTimeSource(source TimeSource) {
	currentTimeSource = source
}
