package envs

import "golang.org/x/text/language"

// Locale is the combination of a language and country, e.g. US English, Brazilian Portuguese
type Locale struct {
	Language Language
	Country  Country
}

// NewLocale creates a new locale
func NewLocale(language Language, country Country) Locale {
	return Locale{Language: language, Country: country}
}

// ToISO639_2 returns the ISO 639-2 code
func (l Locale) ToISO639_2() string {
	lang, err := language.ParseBase(string(l.Language))
	if err != nil {
		return ""
	}
	code := lang.String()

	// not all languages have a 2-letter code
	if len(code) != 2 {
		return ""
	}

	if l.Country != NilCountry {
		code += "-" + string(l.Country)
	}
	return code
}
