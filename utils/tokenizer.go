package utils

import "regexp"

// see: https://en.wikipedia.org/wiki/Emoji for emoji ranges
var wordTokenRegex = regexp.MustCompile("((\\pL|\\pN|[\u20A0-\u20CF]|[\u2600-\u27BF])+|[\U0001F170-\U0001F9CF])")

// TokenizeString returns the words in the passed in string, split by non word characters including emojis
func TokenizeString(str string) []string {
	return wordTokenRegex.FindAllString(str, -1)
}
