package envs

import (
	"github.com/nyaruka/goflow/utils"

	"github.com/nyaruka/phonenumbers"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.RegisterValidatorAlias("country", "len=2", func(validator.FieldError) string {
		return "is not a valid country code"
	})
}

// Country is a ISO 3166-1 alpha-2 country code
type Country string

// NilCountry represents our nil, or unknown country
var NilCountry = Country("")

// DeriveCountryFromTel attempts to derive a country code (e.g. RW) from a phone number
func DeriveCountryFromTel(number string) Country {
	parsed, err := phonenumbers.Parse(number, "")
	if err != nil {
		return ""
	}
	return Country(phonenumbers.GetRegionCodeForNumber(parsed))
}
