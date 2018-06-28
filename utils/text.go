package utils

import (
	"regexp"
	"strings"
)

var snakedChars = regexp.MustCompile(`[^\p{L}\d_]+`)

// Snakify turns the passed in string into a context reference. We replace all whitespace
// characters with _ and replace any duplicate underscores
func Snakify(text string) string {
	return strings.Trim(strings.ToLower(snakedChars.ReplaceAllString(text, "_")), "_")
}

// see: https://en.wikipedia.org/wiki/Emoji for emoji ranges
var wordTokenRegex = regexp.MustCompile("((\\pL|\\pN|[\u20A0-\u20CF]|[\u2600-\u27BF])+|[\U0001F170-\U0001F9CF])")

// TokenizeString returns the words in the passed in string, split by non word characters including emojis
func TokenizeString(str string) []string {
	return wordTokenRegex.FindAllString(str, -1)
}

// TokenizeStringByChars returns the words in the passed in string, split by the chars in the given string
func TokenizeStringByChars(str string, chars string) []string {
	runes := []rune(chars)
	f := func(c rune) bool {
		for _, r := range runes {
			if c == r {
				return true
			}
		}
		return false
	}
	return strings.FieldsFunc(str, f)
}
