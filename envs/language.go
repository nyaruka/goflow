package envs

import (
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func init() {
	utils.Validator.RegisterAlias("language", "eq=base|len=3")
}

// Language is our internal representation of a language
type Language string

// NilLanguage represents our nil, or unknown language
var NilLanguage = Language("")

// ParseLanguage returns a new Language for the passed in language string, or an error if not found
func ParseLanguage(lang string) (Language, error) {
	if len(lang) != 3 {
		return NilLanguage, errors.Errorf("iso-639-3 codes must be 3 characters, got: %s", lang)
	}

	base, err := language.ParseBase(lang)
	if err != nil {
		return NilLanguage, errors.Errorf("unrecognized language code: %s", lang)
	}

	return Language(base.ISO3()), nil
}
