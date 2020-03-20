package envs_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
)

func TestLanguage(t *testing.T) {
	lang, err := envs.ParseLanguage("ENG")
	assert.NoError(t, err)
	assert.Equal(t, envs.Language("eng"), lang)

	_, err = envs.ParseLanguage("base")
	assert.EqualError(t, err, "iso-639-3 codes must be 3 characters, got: base")

	_, err = envs.ParseLanguage("xzx")
	assert.EqualError(t, err, "unrecognized language code: xzx")
}

func TestToISO639_2(t *testing.T) {
	tests := []struct {
		lang     envs.Language
		country  envs.Country
		iso639_2 string
	}{
		{envs.NilLanguage, envs.NilCountry, ``},
		{envs.Language(`cat`), envs.NilCountry, `ca`},
		{envs.Language(`deu`), envs.NilCountry, `de`},
		{envs.Language(`eng`), envs.NilCountry, `en`},
		{envs.Language(`fin`), envs.NilCountry, `fi`},
		{envs.Language(`fra`), envs.NilCountry, `fr`},
		{envs.Language(`jpn`), envs.NilCountry, `ja`},
		{envs.Language(`kor`), envs.NilCountry, `ko`},
		{envs.Language(`pol`), envs.NilCountry, `pl`},
		{envs.Language(`por`), envs.NilCountry, `pt`},
		{envs.Language(`rus`), envs.NilCountry, `ru`},
		{envs.Language(`spa`), envs.NilCountry, `es`},
		{envs.Language(`swe`), envs.NilCountry, `sv`},
		{envs.Language(`zho`), envs.NilCountry, `zh`},
		{envs.Language(`eng`), envs.Country(`US`), `en-US`},
		{envs.Language(`spa`), envs.Country(`EC`), `es-EC`},
		{envs.Language(`zho`), envs.Country(`CN`), `zh-CN`},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.iso639_2, tc.lang.ToISO639_2(tc.country))
	}
}
