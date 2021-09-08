package flows

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
)

func TestTemplateTranslation(t *testing.T) {
	tcs := []struct {
		Content   string
		Variables []string
		Expected  string
	}{
		{"Hi {{1}}, {{2}}", []string{"Chef"}, "Hi Chef, "},
		{"Good boy {{1}}! Who's the best {{1}}?", []string{"Chef"}, "Good boy Chef! Who's the best Chef?"},
		{"Orbit {{1}}! No, go around the {{2}}!", []string{"Chef", "sofa"}, "Orbit Chef! No, go around the sofa!"},
	}

	channel := assets.NewChannelReference("0bce5fd3-c215-45a0-bcb8-2386eb194175", "Test Channel")

	for i, tc := range tcs {
		tt := NewTemplateTranslation(static.NewTemplateTranslation(*channel, envs.Language("eng"), envs.Country("US"), tc.Content, len(tc.Variables), "a6a8863e_7879_4487_ad24_5e2ea429027c"))
		result := tt.Substitute(tc.Variables)
		assert.Equal(t, tc.Expected, result, "%d: unexpected template substitution", i)
	}
}

func TestTemplates(t *testing.T) {
	channel1 := assets.NewChannelReference("0bce5fd3-c215-45a0-bcb8-2386eb194175", "Test Channel")
	tt1 := static.NewTemplateTranslation(*channel1, envs.Language("eng"), envs.NilCountry, "Hello {{1}}", 1, "")
	tt2 := static.NewTemplateTranslation(*channel1, envs.Language("spa"), envs.Country("EC"), "Que tal {{1}}", 1, "")
	tt3 := static.NewTemplateTranslation(*channel1, envs.Language("spa"), envs.Country("ES"), "Hola {{1}}", 1, "")
	template := NewTemplate(static.NewTemplate("c520cbda-e118-440f-aaf6-c0485088384f", "greeting", []*static.TemplateTranslation{tt1, tt2, tt3}))

	tas := NewTemplateAssets([]assets.Template{template})

	tcs := []struct {
		UUID      assets.TemplateUUID
		Channel   *assets.ChannelReference
		Locales   []envs.Locale
		Variables []string
		Expected  string
	}{
		{
			"c520cbda-e118-440f-aaf6-c0485088384f",
			channel1,
			[]envs.Locale{{Language: "eng", Country: "US"}, {Language: "spa", Country: "CO"}},
			[]string{"Chef"},
			"Hello Chef",
		},
		{
			"c520cbda-e118-440f-aaf6-c0485088384f",
			channel1,
			[]envs.Locale{{Language: "eng", Country: ""}, {Language: "spa", Country: "CO"}},
			[]string{"Chef"},
			"Hello Chef",
		},
		{
			"c520cbda-e118-440f-aaf6-c0485088384f",
			channel1,
			[]envs.Locale{{Language: "deu", Country: "DE"}, {Language: "spa", Country: "ES"}},
			[]string{"Chef"},
			"Hola Chef",
		},
		{
			"c520cbda-e118-440f-aaf6-c0485088384f",
			nil,
			[]envs.Locale{{Language: "deu", Country: "DE"}, {Language: "spa", Country: "ES"}},
			[]string{"Chef"},
			"",
		},
		{
			"c520cbda-e118-440f-aaf6-c0485088384f",
			channel1,
			[]envs.Locale{{Language: "deu", Country: "DE"}},
			[]string{"Chef"},
			"",
		},
		{
			"8c5d4910-114a-4521-ba1d-bde8b024865a",
			channel1,
			[]envs.Locale{{Language: "eng", Country: "US"}, {Language: "spa", Country: "ES"}},
			[]string{"Chef"},
			"",
		},
	}

	for _, tc := range tcs {
		tr := tas.FindTranslation(tc.UUID, tc.Channel, tc.Locales)
		if tr == nil {
			assert.Equal(t, "", tc.Expected)
			continue
		}
		assert.NotNil(t, tr.Asset())

		assert.Equal(t, tc.Expected, tr.Substitute(tc.Variables))
	}

	template = tas.Get(assets.TemplateUUID("c520cbda-e118-440f-aaf6-c0485088384f"))
	assert.NotNil(t, template)
	assert.Equal(t, assets.NewTemplateReference("c520cbda-e118-440f-aaf6-c0485088384f", "greeting"), template.Reference())
	assert.Equal(t, (*assets.TemplateReference)(nil), (*Template)(nil).Reference())
}
