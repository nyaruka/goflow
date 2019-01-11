package utils

import (
	"reflect"

	"github.com/hashicorp/go-version"
	"github.com/nyaruka/phonenumbers"
)

// IsNil returns whether the given object is nil or an interface to a nil
func IsNil(v interface{}) bool {
	// if v doesn't have a type or value then v == nil
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)

	// if v is a typed nil pointer then v != nil but the value is nil
	if val.Kind() == reflect.Ptr {
		return val.IsNil()
	}

	return false
}

// MinInt returns the minimum of two integers
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// DeriveCountryFromTel attempts to derive a country code (e.g. RW) from a phone number
func DeriveCountryFromTel(number string) string {
	parsed, err := phonenumbers.Parse(number, "")
	if err != nil {
		return ""
	}
	return phonenumbers.GetRegionCodeForNumber(parsed)
}

// VersionCompare compares two version strings. Returns -1 if v1 is before v2, 0 if they are equal, 1 if v1 is after v2
func VersionCompare(v1, v2 string) (int, error) {
	p1, err := version.NewVersion(v1)
	if err != nil {
		return 0, err
	}
	p2, err := version.NewVersion(v2)
	if err != nil {
		return 0, err
	}

	return p1.Compare(p2), nil
}
