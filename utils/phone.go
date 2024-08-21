package utils

import (
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/urns"
)

var possiblePhone = regexp.MustCompile(`\+?[\d \.\-\(\)]{5,}`)
var onlyPhone = regexp.MustCompile(`^` + possiblePhone.String() + `$`)

// ParsePhoneNumber tries to parse the given string as a phone number. If successful, it returns it formatted as E164.
func ParsePhoneNumber(s string, country i18n.Country) string {
	s = strings.TrimSpace(s)

	if !onlyPhone.MatchString(s) {
		return ""
	}

	formatted, err := urns.ParseNumber(s, country, false, false)
	if err != nil {
		return ""
	}

	// TODO there is a problem in urns package where a number +810000000977123456 is considered valid because it allows
	// unlimited zeros between the country code and the number... but our validation for phone URNs thinks it's too
	// long, so we validate explicitly here
	urn, err := urns.New(urns.Phone, formatted)
	if err != nil {
		return ""
	}

	return urn.Path()
}

// FindPhoneNumbers finds phone numbers anywhere in the given string
func FindPhoneNumbers(s string, country i18n.Country) []string {
	nums := make([]string, 0)
	for _, candidate := range possiblePhone.FindAllString(s, -1) {
		if num := ParsePhoneNumber(candidate, country); num != "" {
			nums = append(nums, num)
		}
	}
	return nums
}
