package utils

import (
	"net/url"
	"regexp"
	"strings"
)

var snakedChars = regexp.MustCompile(`[^\p{L}\d_]+`)

// treats sequences of letters/numbers/_/' as tokens, and symbols as individual tokens
var wordTokenRegex = regexp.MustCompile(`[\pL\pN_']+|\pS`)

// Snakify turns the passed in string into a context reference. We replace all whitespace
// characters with _ and replace any duplicate underscores
func Snakify(text string) string {
	return strings.Trim(strings.ToLower(snakedChars.ReplaceAllString(text, "_")), "_")
}

// URLEscape escapes spaces as %20 matching urllib.quote(s, safe="") in Python
func URLEscape(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}

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
