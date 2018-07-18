package utils

import (
	"github.com/pariz/gountries"
)

// CountryCodeFromName translates a TranferTo country name to a ISO code
func CountryCodeFromName(name string) string {
	query := gountries.New()
	country, err := query.FindCountryByName(name)
	if err != nil {
		return ""
	}
	return country.Alpha2
}
