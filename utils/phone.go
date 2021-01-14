package utils

import (
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/urns"
)

var possiblePhone = regexp.MustCompile(`\+?[\d \.\-\(\)]{5,}`)
var onlyPhone = regexp.MustCompile(`^` + possiblePhone.String() + `$`)

// ParsePhoneNumber tries to parse the given string as a phone number. If successful, it returns it formatted as E164.
func ParsePhoneNumber(s, country string) string {
	s = strings.TrimSpace(s)

	if !onlyPhone.MatchString(s) {
		return ""
	}

	formatted, err := urns.ParseNumber(s, country)
	if err != nil {
		return ""
	}

	return formatted
}

// FindPhoneNumbers finds phone numbers anywhere in the given string
func FindPhoneNumbers(s, country string) []string {
	nums := make([]string, 0)
	for _, candidate := range possiblePhone.FindAllString(s, -1) {
		if num := ParsePhoneNumber(candidate, country); num != "" {
			nums = append(nums, num)
		}
	}
	return nums
}
