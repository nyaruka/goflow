package functions

import "github.com/nyaruka/goflow/utils"

func extractWords(text string, delimiters string) []string {
	if delimiters != "" {
		return utils.TokenizeStringByChars(text, delimiters)
	} else {
		return utils.TokenizeString(text)
	}
}
