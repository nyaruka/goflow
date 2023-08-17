package envs_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/stretchr/testify/assert"
)

func TestCleaners(t *testing.T) {
	tcs := []struct {
		cleaner envs.Cleaner
		input   string
		cleaned string
	}{
		{envs.CleanConfusables, "", ""},
		{envs.CleanConfusables, "𝕟𝔂𝛼𝐫ᴜ𝞳𝕒", "nyaruka"},
		{envs.CleanFarsiToArabic, "۰۱۲۳۴۵۶۷۸۹", "٠١٢٣٤٥٦۷٨٩"},
		{envs.CleanFarsiToArabic, "بلی", "\u0628\u0644\u064A"}, // ends with farsi yeh
		{envs.CleanFarsiToArabic, "بلي", "\u0628\u0644\u064A"}, // ends with arabic yeh
		{envs.CleanArabicToFarsi, "٠١٢٣٤٥٦۷٨٩", "۰۱۲۳۴۵۶۷۸۹"},
		{envs.CleanArabicToFarsi, "بلى", "\u0628\u0644\u06CC"}, // ends with farsi yeh (unchanged)
		{envs.CleanArabicToFarsi, "بلى", "\u0628\u0644\u06CC"}, // ends with alef maksura
		{envs.CleanArabicToFarsi, "بلي", "\u0628\u0644\u06CC"}, // ends with arabic yeh
	}

	for _, tc := range tcs {
		env := envs.NewBuilder().WithInputCleaners(tc.cleaner).Build()

		assert.Equal(t, tc.cleaned, envs.CleanInput(env, tc.input), "%s mismatch for input %s (%s)",
			tc.cleaner, strconv.QuoteToASCII(tc.input), strconv.QuoteToASCII(tc.cleaned))
	}

	assert.Equal(t, `confusables`, envs.CleanConfusables.String())
	assert.Equal(t, []byte(`"confusables"`), jsonx.MustMarshal(envs.CleanConfusables))

	var cleaner envs.Cleaner
	jsonx.MustUnmarshal([]byte(`"arabic_to_farsi"`), &cleaner)
	assert.Equal(t, envs.CleanArabicToFarsi.String(), cleaner.String())

	err := json.Unmarshal([]byte(`"xxx"`), &cleaner)
	assert.EqualError(t, err, "xxx is not a valid cleaner")
}
