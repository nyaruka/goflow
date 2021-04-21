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

// ToBCP47 returns the BCP47 code, e.g. en-US, pt, pt-BR
func (l Locale) ToBCP47() string {
	if l == NilLocale {
		return ""
	}

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

var NilLocale = Locale{}
