package envs

import (
	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.Validator.RegisterAlias("country", "len=2")
}

// Country is a ISO 3166-1 alpha-2 country code
type Country string

// NilCountry represents our nil, or unknown country
var NilCountry = Country("")
