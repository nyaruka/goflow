package dates

import (
	"time"
)

const (
	ISO8601Date = "2006-01-02"
)

// Date is a local gregorian calendar date
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// NewDate creates a new date
func NewDate(year, month, day int) Date {
	return Date{year, time.Month(month), day}
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
		return int(d.Month) - int(other.Month)
	}
	return d.Day - other.Day
}

// Combine combines this date and a time to make a datetime
func (d Date) Combine(timeOfDay TimeOfDay, tz *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, tz)
}

// Format formats this date as a string
func (d Date) Format(layout string) string {
	// upgrade us to a date time so we can use standard time.Time formatting
	return d.Combine(ZeroTimeOfDay, time.UTC).Format(layout)
}

// Weekday returns the day of the week
func (d Date) Weekday() time.Weekday {
	return d.Combine(ZeroTimeOfDay, time.UTC).Weekday()
}

// YearDay returns the day of the year (1-366)
func (d Date) YearDay() int {
	return d.Combine(ZeroTimeOfDay, time.UTC).YearDay()
}

// WeekNum returns the number of the week (1-54)
func (d Date) WeekNum() int {
	wday := int(d.Weekday()) // Sun = 0, Mon = 1
	yday := d.YearDay() - 1  // to 0-365

	week := (yday - wday + 7) / 7

	return week + 1
}

// String returns the ISO8601 representation
func (d Date) String() string {
	return d.Format(ISO8601Date)
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
