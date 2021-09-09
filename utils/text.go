package utils

import (
	"bytes"
	"regexp"
	"sort"
	"strings"

	"github.com/blevesearch/segment"
)

var snakedChars = regexp.MustCompile(`[^\p{L}\d_]+`)

// treats sequences of letters/numbers/_/' as tokens, and symbols as individual tokens
var wordTokenRegex = regexp.MustCompile(`[\pM\pL\pN_']+|\pS`)

// Snakify turns the passed in string into a context reference. We replace all whitespace
// characters with _ and replace any duplicate underscores
func Snakify(text string) string {
	return strings.ToLower(snakedChars.ReplaceAllString(strings.TrimSpace(text), "_"))
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

// TokenizeStringByUnicodeSeg tokenizes the given string using the Unicode Text Segmentation standard described at http://www.unicode.org/reports/tr29/
func TokenizeStringByUnicodeSeg(str string) []string {
	segmenter := segment.NewWordSegmenter(strings.NewReader(str))
	tokens := make([]string, 0)

	for segmenter.Segment() {
		token := string(segmenter.Bytes())
		ttype := segmenter.Type()
		if ttype != segment.None {
			tokens = append(tokens, token)
		}
	}

	return tokens
}

// PrefixOverlap returns the number of prefix characters which s1 and s2 have in common
func PrefixOverlap(s1, s2 string) int {
	r1 := []rune(s1)
	r2 := []rune(s2)
	i := 0
	for ; i < len(r1) && i < len(r2) && r1[i] == r2[i]; i++ {
	}
	return i
}

// StringSlices returns the slices of s defined by pairs of indexes in indices
func StringSlices(s string, indices []int) []string {
	slices := make([]string, 0, len(indices)/2)
	for i := 0; i < len(indices); i += 2 {
		slices = append(slices, s[indices[i]:indices[i+1]])
	}
	return slices
}

// StringSliceContains determines whether the given slice of strings contains the given string
func StringSliceContains(slice []string, str string, caseSensitive bool) bool {
	for _, s := range slice {
		if (caseSensitive && s == str) || (!caseSensitive && strings.EqualFold(s, str)) {
			return true
		}
	}
	return false
}

// StringSet converts a slice of strings to a set (a string > bool map)
func StringSet(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}

// StringSetKeys returns the keys of string set in lexical order
func StringSetKeys(m map[string]bool) []string {
	vals := make([]string, 0, len(m))
	for v := range m {
		vals = append(vals, v)
	}
	sort.Strings(vals)
	return vals
}

// Indent indents each non-empty line in the given string
func Indent(s string, prefix string) string {
	output := strings.Builder{}

	bol := true
	for _, c := range s {
		if bol && c != '\n' {
			output.WriteString(prefix)
		}
		output.WriteRune(c)
		bol = c == '\n'
	}
	return output.String()
}

// Truncate truncates the given string to ensure it's less than limit characters
func Truncate(s string, limit int) string {
	return truncate(s, limit, "")
}

// TruncateEllipsis truncates the given string and adds ellipsis where the input is cut
func TruncateEllipsis(s string, limit int) string {
	return truncate(s, limit, "...")
}

func truncate(s string, limit int, ending string) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	return string(runes[:limit-len(ending)]) + ending
}

// Redactor is a function which can redact the given string
type Redactor func(s string) string

// NewRedactor creates a new redaction function which replaces the given values
func NewRedactor(mask string, values ...string) Redactor {
	// convert list of redaction values to list of replacements with mask
	replacements := make([]string, len(values)*2)
	for i := range values {
		replacements[i*2] = values[i]
		replacements[i*2+1] = mask
	}
	return strings.NewReplacer(replacements...).Replace
}

// replaces any `\u0000` sequences with the given replacement sequence which may be empty.
// A sequence such as `\\u0000` is preserved as it is an escaped slash followed by the sequence `u0000`
func ReplaceEscapedNulls(data []byte, repl []byte) []byte {
	return nullEscapeRegex.ReplaceAllFunc(data, func(m []byte) []byte {
		slashes := bytes.Count(m, []byte(`\`))
		if slashes%2 == 0 {
			return m
		}

		return append(bytes.Repeat([]byte(`\`), slashes-1), repl...)
	})
}

var nullEscapeRegex = regexp.MustCompile(`\\+u0{4}`)
