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

// MaxInt returns the maximum of two integers
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MinInt returns the minimum of two integers
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// FindPhoneNumber attempts to parse a phone number from the given text. If it finds one, it returns it formatted as E164.
func FindPhoneNumber(s, country string) string {
	phone, err := phonenumbers.Parse(s, country)
	if err != nil {
		return ""
	}

	if !phonenumbers.IsValidNumber(phone) {
		return ""
	}

	return phonenumbers.Format(phone, phonenumbers.E164)
}

// DeriveCountryFromTel attempts to derive a country code (e.g. RW) from a phone number
func DeriveCountryFromTel(number string) string {
	parsed, err := phonenumbers.Parse(number, "")
	if err != nil {
		return ""
	}
	return phonenumbers.GetRegionCodeForNumber(parsed)
}

// Typed is an interface of objects that are marshalled as typed envelopes
type Typed interface {
	Type() string
}

// TypedEnvelope can be mixed into envelopes that have a type field
type TypedEnvelope struct {
	Type string `json:"type" validate:"required"`
}

// ReadTypeFromJSON reads a field called `type` from the given JSON
func ReadTypeFromJSON(data []byte) (string, error) {
	t := &TypedEnvelope{}
	if err := UnmarshalAndValidate(data, t); err != nil {
		return "", err
	}
	return t.Type, nil
}
