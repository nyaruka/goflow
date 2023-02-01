package envs

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/nyaruka/null/v2"
	"golang.org/x/text/language"
)

// Locale is the combination of a language and optional country, e.g. US English, Brazilian Portuguese, encoded as the
// language code followed by the country code, e.g. eng-US, por-BR
type Locale string

// NewLocale creates a new locale
func NewLocale(l Language, c Country) Locale {
	if l == NilLanguage {
		return NilLocale
	}
	if c == NilCountry {
		return Locale(l) // e.g. "eng", "por"
	}
	return Locale(fmt.Sprintf("%s-%s", l, c)) // e.g. "eng-US", "por-BR"
}

// ToBCP47 returns the BCP47 code, e.g. en-US, pt, pt-BR
func (l Locale) ToBCP47() string {
	if l == NilLocale {
		return ""
	}

	lang, country := l.ToParts()

	base, err := language.ParseBase(string(lang))
	if err != nil {
		return ""
	}
	code := base.String()

	// not all languages have a 2-letter code
	if len(code) != 2 {
		return ""
	}

	if country != NilCountry {
		code += "-" + string(country)
	}
	return code
}

func (l Locale) ToParts() (Language, Country) {
	if l == NilLocale || len(l) < 3 {
		return NilLanguage, NilCountry
	}

	parts := strings.SplitN(string(l), "-", 2)
	lang := Language(parts[0])
	country := NilCountry
	if len(parts) > 1 {
		country = Country(parts[1])
	}

	return lang, country
}

var NilLocale = Locale("")

// Place nicely with NULLs if persisting to a database or JSON
func (l *Locale) Scan(value any) error         { return null.ScanString(value, l) }
func (l Locale) Value() (driver.Value, error)  { return null.StringValue(l) }
func (l Locale) MarshalJSON() ([]byte, error)  { return null.MarshalString(l) }
func (l *Locale) UnmarshalJSON(b []byte) error { return null.UnmarshalString(b, l) }
