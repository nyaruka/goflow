package dates

import (
	"time"
)

// format we use for full ISO output
const iso8601Default = "2006-01-02T15:04:05.000000Z07:00"

// DateTimeToISO converts the passed in time.Time to a string in full ISO8601 format
func FormatISO(date time.Time) string {
	return date.Format(iso8601Default)
}

// DaysBetween returns the number of calendar days (an int) between the two dates. Note
// that if these are in different timezones then the local calendar day is used for each
// and the difference is calculated from that.
func DaysBetween(date1 time.Time, date2 time.Time) int {
	d1 := time.Date(date1.Year(), date1.Month(), date1.Day(), 0, 0, 0, 0, time.UTC)
	d2 := time.Date(date2.Year(), date2.Month(), date2.Day(), 0, 0, 0, 0, time.UTC)

	return int(d1.Sub(d2) / (time.Hour * 24))
}

// MonthsBetween returns the number of calendar months (an int) between the two dates. Note
// that if these are in different timezones then the local calendar day is used for each
// and the difference is calculated from that.
func MonthsBetween(date1 time.Time, date2 time.Time) int {
	// difference in months
	months := int(date1.Month() - date2.Month())

	// difference in years
	months += (date1.Year() - date2.Year()) * 12

	return months
}

// DayToUTCRange returns the UTC time range of the given day
func DayToUTCRange(d time.Time, tz *time.Location) (time.Time, time.Time) {
	localMidnight := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	utcMidnight := localMidnight.In(tz)
	return utcMidnight, utcMidnight.Add(24 * time.Hour)
}
