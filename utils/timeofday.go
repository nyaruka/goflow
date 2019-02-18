package utils

import (
	"fmt"
	"time"
)

// TimeOfDay represents a local time of day value
type TimeOfDay struct {
	Hour   int
	Minute int
	Second int
	Nanos  int
}

// NewTimeOfDay creates a new time of day
func NewTimeOfDay(hour, minute, second, nanos int) TimeOfDay {
	return TimeOfDay{Hour: hour, Minute: minute, Second: second, Nanos: nanos}
}

// ExtractTimeOfDay extracts the time of day from the give datetime
func ExtractTimeOfDay(dt time.Time) TimeOfDay {
	return NewTimeOfDay(dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond())
}

// Equals determines equality for this type
func (t TimeOfDay) Equals(other TimeOfDay) bool {
	return t.Hour == other.Hour && t.Minute == other.Minute && t.Second == other.Second && t.Nanos == other.Nanos
}

// Compare compares this time of day to another
func (t TimeOfDay) Compare(other TimeOfDay) int {
	if t.Hour != other.Hour {
		return t.Hour - other.Hour
	}
	if t.Minute != other.Minute {
		return t.Minute - other.Minute
	}
	if t.Second != other.Second {
		return t.Second - other.Second
	}
	return t.Nanos - other.Nanos
}

// String returns the ISO8601 representation
func (t TimeOfDay) String() string {
	// TODO this matches our DateToISO accuracy.. but we're throwing away accuracy
	return fmt.Sprintf("%02d:%02d:%02d.%06d", t.Hour, t.Minute, t.Second, t.Nanos/1000)
}

// ZeroTimeOfDay is our uninitialized time of day value
var ZeroTimeOfDay = TimeOfDay{}
