package functions

import "github.com/nyaruka/goflow/utils"

func extractWords(text string, delimiters string) []string {
	if delimiters != "" {
		return utils.TokenizeStringByChars(text, delimiters)
	} else {
		return utils.TokenizeString(text)
	}
}

// returns a slice of the array from start to end using same flexible rules as Python
func slice[T any](arr []T, start, end int) []T {
	length := len(arr)

	if start < 0 {
		start = max(length+start, 0)
	}
	if end < 0 {
		end = max(length+end, 0)
	}
	start = min(start, length)
	end = min(end, length)
	if start > end {
		return []T{}
	}

	return arr[start:end]
}
