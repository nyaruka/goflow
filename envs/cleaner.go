package envs

import (
	"encoding/json"
	"strings"

	"github.com/nyaruka/gocommon/stringsx"
	"github.com/pkg/errors"
)

// Cleaner is a named function that "cleans" a string.
type Cleaner struct {
	type_ string
	fn    func(string) string
}

func (c Cleaner) String() string {
	return c.type_
}

func (c Cleaner) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Cleaner) UnmarshalJSON(d []byte) error {
	var type_ string
	if err := json.Unmarshal(d, &type_); err != nil {
		return err
	}

	for _, cl := range cleaners {
		if type_ == cl.type_ {
			*c = cl
			return nil
		}
	}
	return errors.Errorf("%s is not a valid cleaner", type_)
}

var CleanConfusables = Cleaner{
	"confusables",
	func(s string) string { return stringsx.Skeleton(s) },
}

// https://en.wikipedia.org/wiki/Persian_alphabet#Deviations_from_the_Arabic_script
var farsiToArabic = map[rune]rune{
	'۰': '٠', // U+06F0 > U+0660 (0)
	'۱': '١', // U+06F1 > U+0661 (1)
	'۲': '٢', // U+06F2 > U+0662 (2)
	'۳': '٣', // U+06F3 > U+0663 (3)
	'۴': '٤', // U+06F4 > U+0664 (4)
	'۵': '٥', // U+06F5 > U+0665 (5)
	'۶': '٦', // U+06F6 > U+0666 (6)
	'۷': '۷', // U+06F7 > U+0667 (7)
	'۸': '٨', // U+06F8 > U+0668 (8)
	'۹': '٩', // U+06F9 > U+0669 (9)
	'ی': 'ي', // U+06CC > U+064A (yeh)
	'ک': 'ك', // U+06A9 > U+0643 (kāf)
}

var CleanFarsiToArabic = Cleaner{
	"farsi_to_arabic",
	func(s string) string { return replaceRunes(s, farsiToArabic) },
}

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

var CleanArabicToFarsi = Cleaner{
	"arabic_to_farsi",
	func(s string) string { return replaceRunes(s, arabicToFarsi) },
}

func CleanInput(env Environment, s string) string {
	for _, f := range env.InputCleaners() {
		s = f.fn(s)
	}
	return s
}

var cleaners = []Cleaner{CleanConfusables, CleanFarsiToArabic, CleanArabicToFarsi}

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
