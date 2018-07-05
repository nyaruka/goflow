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

// a time source which returns a fixed time
type fixedTimeSource struct {
	now time.Time
}

func (s fixedTimeSource) Now() time.Time {
	return s.now
}

// NewFixedTimeSource creates a new fixed time source
func NewFixedTimeSource(now time.Time) TimeSource {
	return &fixedTimeSource{now: now}
}

// a time source which returns a sequence of times 1 second after each other
type sequentialTimeSource struct {
	current time.Time
}

func (s sequentialTimeSource) Now() time.Time {
	now := s.current
	s.current.Add(time.Second * 1)
	return now
}

// NewSequentialTimeSource creates a new sequential time source
func NewSequentialTimeSource(start time.Time) TimeSource {
	return &sequentialTimeSource{current: start}
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
