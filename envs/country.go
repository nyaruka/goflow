package envs

import (
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/phonenumbers"
)

func init() {
	utils.Validator.RegisterAlias("country", "len=2")
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
