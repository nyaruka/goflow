package envs_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/stretchr/testify/assert"
)

func TestLocale(t *testing.T) {
	assert.Equal(t, envs.Locale(""), envs.NewLocale("", ""))
	assert.Equal(t, envs.Locale(""), envs.NewLocale("", "US"))     // invalid without language
	assert.Equal(t, envs.Locale("eng"), envs.NewLocale("eng", "")) // valid without country
	assert.Equal(t, envs.Locale("eng-US"), envs.NewLocale("eng", "US"))

	l, c := envs.Locale("eng-US").ToParts()
	assert.Equal(t, envs.Language("eng"), l)
	assert.Equal(t, envs.Country("US"), c)

	l, c = envs.NilLocale.ToParts()
	assert.Equal(t, envs.NilLanguage, l)
	assert.Equal(t, envs.NilCountry, c)

	v, err := envs.NewLocale("eng", "US").Value()
	assert.NoError(t, err)
	assert.Equal(t, "eng-US", v)

	v, err = envs.NilLanguage.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)

	var lc envs.Locale
	assert.NoError(t, lc.Scan("eng-US"))
	assert.Equal(t, envs.Locale("eng-US"), lc)

	assert.NoError(t, lc.Scan(nil))
	assert.Equal(t, envs.NilLocale, lc)
}

func TestToBCP47(t *testing.T) {
	tests := []struct {
		locale envs.Locale
		bcp47  string
	}{
		{``, ``},
		{`cat`, `ca`},
		{`deu`, `de`},
		{`eng`, `en`},
		{`fin`, `fi`},
		{`fra`, `fr`},
		{`jpn`, `ja`},
		{`kor`, `ko`},
		{`pol`, `pl`},
		{`por`, `pt`},
		{`rus`, `ru`},
		{`spa`, `es`},
		{`swe`, `sv`},
		{`zho`, `zh`},
		{`eng-US`, `en-US`},
		{`spa-EC`, `es-EC`},
		{`zho-CN`, `zh-CN`},

		{`yue`, ``}, // has no 2-letter represention
		{`und`, ``},
		{`mul`, ``},
		{`xyz`, ``}, // is not a language
	}

	for _, tc := range tests {
		assert.Equal(t, tc.bcp47, tc.locale.ToBCP47())
	}
}

func TesBCP47Matcher(t *testing.T) {
	tests := []struct {
		preferred []envs.Locale
		available []string
		best      string
	}{
		{preferred: []envs.Locale{"eng-US"}, available: []string{"es_EC", "en-US"}, best: "en-US"},
		{preferred: []envs.Locale{"eng-US"}, available: []string{"es", "en"}, best: "en"},
		{preferred: []envs.Locale{"eng"}, available: []string{"es-US", "en-UK"}, best: "en-UK"},
		{preferred: []envs.Locale{"eng", "fra"}, available: []string{"fr-CA", "en-RW"}, best: "en-RW"},
		{preferred: []envs.Locale{"eng", "fra"}, available: []string{"fra-CA", "eng-RW"}, best: "eng-RW"},
		{preferred: []envs.Locale{"fra", "eng"}, available: []string{"fra-CA", "eng-RW"}, best: "fra-CA"},
		{preferred: []envs.Locale{"spa"}, available: []string{"es-EC", "es-MX", "es-ES"}, best: "es-ES"},
		{preferred: []envs.Locale{}, available: []string{"es_EC", "en-US"}, best: "es_EC"},
	}

	for _, tc := range tests {
		m := envs.NewBCP47Matcher(tc.available...)
		best := m.ForLocales(tc.preferred...)

		assert.Equal(t, tc.best, best, "locale mismatch for preferred=%v available=%s", tc.preferred, tc.available)
	}
}
