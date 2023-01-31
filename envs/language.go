package envs

import (
	"database/sql/driver"

	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/null/v2"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.RegisterValidatorAlias("language", "len=3", func(validator.FieldError) string {
		return "is not a valid language code"
	})
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

// Place nicely with NULLs if persisting to a database or JSON
func (l *Language) Scan(value any) error         { return null.ScanString(value, l) }
func (l Language) Value() (driver.Value, error)  { return null.StringValue(l) }
func (l Language) MarshalJSON() ([]byte, error)  { return null.MarshalString(l) }
func (l *Language) UnmarshalJSON(b []byte) error { return null.UnmarshalString(b, l) }
