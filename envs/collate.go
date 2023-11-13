package envs

import (
	"strings"

	"github.com/nyaruka/gocommon/stringsx"
	"golang.org/x/text/unicode/norm"
)

type Collation string

const (
	CollationDefault        Collation = "default"
	CollationConfusables    Collation = "confusables"
	CollationArabicVariants Collation = "arabic_variants"
)

type collateTransformer func(string) string

// Based on https://en.wikipedia.org/wiki/Persian_alphabet#Deviations_from_the_Arabic_script
// and feedback from UNICEF Afghanistan
var arabicVariants = map[rune]rune{
	'٠': '۰', // U+0660 > U+06F0 (0 > ext arabic 0)
	'١': '۱', // U+0661 > U+06F1 (1 > ext arabic 1)
	'٢': '۲', // U+06F2 > U+0662 (2 > ext arabic 2)
	'٣': '۳', // U+06F3 > U+0663 (3 > ext arabic 3)
	'٤': '۴', // U+06F4 > U+0664 (4 > ext arabic 4)
	'٥': '۵', // U+06F5 > U+0665 (5 > ext arabic 5)
	'٦': '۶', // U+06F6 > U+0666 (6 > ext arabic 6)
	'٧': '۷', // U+06F7 > U+0667 (7 > ext arabic 7)
	'٨': '۸', // U+06F8 > U+0668 (8 > ext arabic 8)
	'٩': '۹', // U+06F9 > U+0669 (9 > ext arabic 9)
	'آ': 'ا', // U+0622 > U+0627 (alef with madda > alef)
	'ى': 'ی', // U+0649 > U+06CC (alef maksura > farsi yeh)
	'ي': 'ی', // U+064A > U+06CC (yeh > farsi yeh)
	'ې': 'ی', // U+06DO > U+06CC (eh > farsi yeh)
	'ۍ': 'ی', // U+06CD > U+06CC (yeh with tail > farsi yeh)
	'ئ': 'ی', // U+0626 > U+06CC (yeh with hamza > farsi yeh)
	'ك': 'ک', // U+0643 > U+06A9 (kāf > keheh)
	'ګ': 'ک', // U+06AB > U+06A9 (kāf with ring > keheh)
	'ټ': 'ت', // U+067C > U+062A (teh with ring > teh)
	'ډ': 'د', // U+0689 > U+062F (dal with ring > dal)
	'ړ': 'ر', // U+0693 > U+0631 (reh with ring > reh)
	'ڼ': 'ن', // U+06BC > U+0646 (noon with ring > noon)
	'ښ': 'ش', // U+069A > U+0634 (pashto seen > sheen)
}

var transformers = map[Collation]collateTransformer{
	CollationDefault: func(s string) string {
		return strings.ToLower(s)
	},
	CollationConfusables: func(s string) string {
		return strings.ToLower(stringsx.Skeleton(s))
	},
	CollationArabicVariants: func(s string) string {
		return strings.ToLower(replaceRunes(norm.NFKC.String(s), arabicVariants))
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
