package test

import (
	"time"

	"github.com/nyaruka/goflow/utils"
)

// a time source which returns a fixed time
type fixedTimeSource struct {
	now time.Time
}

func (s *fixedTimeSource) Now() time.Time {
	return s.now
}

// NewFixedTimeSource creates a new fixed time source
func NewFixedTimeSource(now time.Time) utils.TimeSource {
	return &fixedTimeSource{now: now}
}

// a time source which returns a sequence of times 1 second after each other
type sequentialTimeSource struct {
	current time.Time
}

func (s *sequentialTimeSource) Now() time.Time {
	now := s.current
	s.current = s.current.Add(time.Second * 1)
	return now
}

// NewSequentialTimeSource creates a new sequential time source
func NewSequentialTimeSource(start time.Time) utils.TimeSource {
	return &sequentialTimeSource{current: start}
}
