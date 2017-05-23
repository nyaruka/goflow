package utils

import "regexp"

// see: https://en.wikipedia.org/wiki/Emoji for emoji ranges
var work_token_regex = regexp.MustCompile("((\\pL|\\pN|[\u20A0-\u20CF]|[\u2600-\u27BF])+|[\U0001F170-\U0001F9CF])")

func TokenizeString(str string) []string {
	return work_token_regex.FindAllString(str, -1)
}
