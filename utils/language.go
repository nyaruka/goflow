package utils

import (
	"fmt"

	"golang.org/x/text/language"
)

// Language is our internal representation of a language
type Language string

// NilLanguage represents our nil, or unknown language
var NilLanguage = Language("")

// ParseLanguage returns a new Language for the passed in language string, or an error if not found
func ParseLanguage(lang string) (Language, error) {
	if len(lang) != 3 {
		return NilLanguage, fmt.Errorf("iso-639-3 codes must be 3 characters, got: %s", lang)
	}

	base, err := language.ParseBase(lang)
	if err != nil {
		return NilLanguage, err
	}

	return Language(base.ISO3()), nil
}

type LanguageList []Language

func (ll LanguageList) RemoveDuplicates() LanguageList {
	result := LanguageList{}
	seen := map[Language]bool{}
	for _, val := range ll {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = true
		}
	}
	return result
}
