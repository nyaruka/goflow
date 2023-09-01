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

func (l Locale) ToTag() language.Tag {
	return language.MustParse(string(l))
}

// ToBCP47 returns the BCP47 code, e.g. en-US, pt, pt-BR
// deprecated
func (l Locale) ToBCP47() string {
	if l == NilLocale {
		return ""
	}

	lang, country := l.ToParts()
	lang2 := lang.ISO639_1()

	// not all languages have a 2-letter code
	if lang2 == "" {
		return ""
	}

	if country != NilCountry {
		lang2 += "-" + string(country)
	}
	return lang2
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

// BCP47Matcher helps find best matching locale from a set of available locales
type BCP47Matcher struct {
	codes   []string
	matcher language.Matcher
}

// NewBCP47Matcher creates a new BCP47 matcher from the set of available locales which must be valid BCP47 tags.
func NewBCP47Matcher(codes ...string) *BCP47Matcher {
	tags := make([]language.Tag, len(codes))
	for i := range codes {
		tags[i] = language.MustParse(codes[i])
	}
	return &BCP47Matcher{codes: codes, matcher: language.NewMatcher(tags)}
}

func (m *BCP47Matcher) ForLocales(preferred ...Locale) string {
	prefTags := make([]language.Tag, len(preferred))
	for i := range preferred {
		prefTags[i] = preferred[i].ToTag()
	}

	// see https://github.com/golang/go/issues/24211
	_, idx, _ := m.matcher.Match(prefTags...)
	return m.codes[idx]
}
