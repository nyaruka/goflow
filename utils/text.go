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
