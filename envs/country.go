package envs

import (
	"database/sql/driver"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/null/v2"
	"github.com/nyaruka/phonenumbers"
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

// Place nicely with NULLs if persisting to a database or JSON
func (c *Country) Scan(value any) error         { return null.ScanString(value, c) }
func (c Country) Value() (driver.Value, error)  { return null.StringValue(c) }
func (c Country) MarshalJSON() ([]byte, error)  { return null.MarshalString(c) }
func (c *Country) UnmarshalJSON(b []byte) error { return null.UnmarshalString(b, c) }
