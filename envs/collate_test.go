package envs_test

import (
	"strconv"
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/stretchr/testify/assert"
)

func TestCollation(t *testing.T) {

	tcs := []struct {
		collation envs.Collation
		input     string
		transform string
		equals    map[string]bool
	}{
		{envs.CollationDefault, "AbcD", "abcd", map[string]bool{
			"acde": false,
			"aBCd": true,
		}},
		{envs.CollationConfusables, "𝕟𝔂𝛼𝐫ᴜ𝞳𝕒", "nyaruka", map[string]bool{
			"trileet": false,
			"Nyaruka": true,
			"𝒩ɣaruka": true,
		}},
		{envs.CollationConfusables, "JOIN", "join", map[string]bool{ // capital I must not become l
			"join":  true,
			"Join":  true,
			"JOIN":  true,
			"jоin":  true, // cyrillic о
			"ЈОIN":  true, // cyrillic Ј and О
			"joln":  false,
			"jo1n":  false,
			"joint": false,
		}},
		{envs.CollationConfusables, "REGISTER", "register", map[string]bool{
			"register": true,
			"Register": true,
			"reglster": false,
		}},
		{envs.CollationConfusables, "𝐒𝐓𝐎𝐏", "stop", map[string]bool{ // math bold caps have no lowercase mapping
			"stop": true,
			"STOP": true,
			"Stop": true,
		}},
		{envs.CollationArabicVariants, "٠١٢٣٤٥٦۷٨٩", "۰۱۲۳۴۵۶۷۸۹", map[string]bool{
			"٤٥٦۷":       false,
			"٠١٢٣٤٥٦۷٨٩": true,
			"۰۱۲۳۴۵۶۷۸۹": true,
		}},
		{envs.CollationArabicVariants, "\u0628\u0644\u06CC", "\u0628\u0644\u06CC", map[string]bool{ // ends with farsi yeh (unchanged)
			"\u0628\u0644":       false,
			"\u0628\u0644\u0649": true, // ends with alef maksura
			"\u0628\u0644\u064A": true, // ends with arabic yeh
		}},
		{envs.CollationArabicVariants, "\u0628\u0644\u0649", "\u0628\u0644\u06CC", map[string]bool{ // ends with alef maksura
			"\u0628\u0644\u06CC": true, // ends with farsi yeh
			"\u0628\u0644\u064A": true, // ends with arabic yeh
		}},
		{envs.CollationArabicVariants, "\u0628\u0644\u064A", "\u0628\u0644\u06CC", map[string]bool{ // ends with arabic yeh
			"\u0628\u0644\u06CC": true, // ends with farsi yeh
			"\u0628\u0644\u0649": true, // ends with alef maksura
		}},
		{envs.CollationArabicVariants, "\u0643\u0627\u0641", "\u06A9\u0627\u0641", map[string]bool{ // starts with arabic kaf
			"\u0643\u0627\u0641": true, // starts with arabic kaf
			"\u06A9\u0627\u0641": true, // starts with farsi kaf
			"\uFEDB\u0627\u0641": true, // starts with explicit initial form kaf
		}},
		{envs.CollationArabicVariants, "\u0622", "\u0627", map[string]bool{}},
		{envs.CollationArabicVariants, "\uFE8F\uFEDD\uFBFC", "\u0628\u0644\u06CC", map[string]bool{}}, // Arabic Presentation forms
		{envs.CollationArabicVariants, "YES", "yes", map[string]bool{"yes": true, "no": false}},
	}

	for _, tc := range tcs {
		env := envs.NewBuilder().WithInputCollation(tc.collation).Build()

		assert.Equal(t, tc.transform, envs.CollateTransform(env, tc.input), "%s transform mismatch for input %s (%s)",
			tc.collation, strconv.QuoteToASCII(tc.input), strconv.QuoteToASCII(tc.transform))

		for eqStr, eqResult := range tc.equals {
			assert.Equal(t, eqResult, envs.CollateEquals(env, tc.input, eqStr))
		}
	}
}
