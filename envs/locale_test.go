package envs_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
)

func TestToBCP47(t *testing.T) {
	tests := []struct {
		lang    envs.Language
		country envs.Country
		bcp47   string
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

		{envs.Language(`yue`), envs.NilCountry, ``}, // has no 2-letter represention
		{envs.Language(`und`), envs.NilCountry, ``},
		{envs.Language(`mul`), envs.NilCountry, ``},
		{envs.Language(`xyz`), envs.NilCountry, ``}, // is not a language
	}

	for _, tc := range tests {
		assert.Equal(t, tc.bcp47, envs.NewLocale(tc.lang, tc.country).ToBCP47())
	}
}
