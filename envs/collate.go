package envs

import (
	"strings"

	"github.com/nyaruka/gocommon/stringsx"
	"golang.org/x/text/unicode/norm"
)

type Collation string

const (
	CollationDefault     Collation = "default"
	CollationConfusables Collation = "confusables"
	CollationArabicFarsi Collation = "arabic_farsi"
)

type collateTransformer func(string) string

// https://en.wikipedia.org/wiki/Persian_alphabet#Deviations_from_the_Arabic_script
var arabicToFarsi = map[rune]rune{
	'٠': '۰', // U+0660 > U+06F0 (0)
	'١': '۱', // U+0661 > U+06F1 (1)
	'٢': '۲', // U+06F2 > U+0662 (2)
	'٣': '۳', // U+06F3 > U+0663 (3)
	'٤': '۴', // U+06F4 > U+0664 (4)
	'٥': '۵', // U+06F5 > U+0665 (5)
	'٦': '۶', // U+06F6 > U+0666 (6)
	'٧': '۷', // U+06F7 > U+0667 (7)
	'٨': '۸', // U+06F8 > U+0668 (8)
	'٩': '۹', // U+06F9 > U+0669 (9)
	'ى': 'ی', // U+0649 > U+06CC (alef maksura)
	'ي': 'ی', // U+064A > U+06CC (yeh)
	'ك': 'ک', // U+0643 > U+06A9 (kāf)
}

var transformers = map[Collation]collateTransformer{
	CollationDefault: func(s string) string {
		return strings.ToLower(s)
	},
	CollationConfusables: func(s string) string {
		return strings.ToLower(stringsx.Skeleton(s))
	},
	CollationArabicFarsi: func(s string) string {
		return strings.ToLower(replaceRunes(norm.NFKD.String(s), arabicToFarsi))
	},
}

// CollateEquals returns true if the given strings are equal in the given environment's collation
func CollateEquals(env Environment, s, t string) bool {
	return CollateTransform(env, s) == CollateTransform(env, t)
}

// CollateTransform transforms the given string into it's form to be used for collation.
func CollateTransform(env Environment, s string) string {
	return transformers[env.InputCollation()](s)
}

func replaceRunes(s string, mapping map[rune]rune) string {
	var sb strings.Builder
	for _, r := range s {
		if repl, found := mapping[r]; found {
			sb.WriteRune(repl)
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
