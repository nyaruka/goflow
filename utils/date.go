package utils

import (
	"time"
)

// Date is a local gregorian calendar date
type Date struct {
	Year  int
	Month int
	Day   int
}

// NewDate creates a new date
func NewDate(year, month, day int) Date {
	return Date{year, month, day}
}

// ExtractDate extracts the date from the give datetime
func ExtractDate(dt time.Time) Date {
	return NewDate(dt.Year(), int(dt.Month()), dt.Day())
}

// Equal determines equality for this type
func (d Date) Equal(other Date) bool {
	return d.Year == other.Year && d.Month == other.Month && d.Day == other.Day
}

// Compare compares this time of day to another
func (d Date) Compare(other Date) int {
	if d.Year != other.Year {
		return d.Year - other.Year
	}
	if d.Month != other.Month {
		return d.Month - other.Month
	}
	return d.Day - other.Day
}

// Format formats this date as a string
func (d Date) Format(layout string) string {
	// upgrade us to a date time so we can use standard time.Time formatting
	dt := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
	return dt.Format(layout)
}

// String returns the ISO8601 representation
func (d Date) String() string {
	return d.Format(iso8601Date)
}

// ZeroDate is our uninitialized date value
var ZeroDate = Date{}

// ParseDate parses the given string into a date
func ParseDate(layout string, value string) (Date, error) {
	dt, err := time.Parse(layout, value)
	if err != nil {
		return ZeroDate, err
	}

	return ExtractDate(dt), nil
}
