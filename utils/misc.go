package utils

import (
	"reflect"

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

// StringSliceContains determines whether the given slice of strings contains the given string
func StringSliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
