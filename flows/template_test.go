package flows_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
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
		tt := flows.NewTemplateTranslation(static.NewTemplateTranslation(channel, i18n.Locale("eng-US"), tc.Content, len(tc.Variables), "a6a8863e_7879_4487_ad24_5e2ea429027c", map[string][]assets.TemplateParam{}))
		result := tt.Substitute(tc.Variables)
		assert.Equal(t, tc.Expected, result, "%d: unexpected template substitution", i)
	}
}

func TestTemplate(t *testing.T) {
	channel1 := test.NewChannel("WhatsApp 1", "+12345", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channel2 := test.NewChannel("WhatsApp 2", "+23456", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channel3 := test.NewChannel("WhatsApp 3", "+34567", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channel1Ref := assets.NewChannelReference(channel1.UUID(), channel1.Name())
	channel2Ref := assets.NewChannelReference(channel2.UUID(), channel2.Name())

	tt1 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("eng"), "Hello {{1}}", 1, "", map[string][]assets.TemplateParam{})
	tt2 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("spa-EC"), "Que tal {{1}}", 1, "", map[string][]assets.TemplateParam{})
	tt3 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("spa-ES"), "Hola {{1}}", 1, "", map[string][]assets.TemplateParam{})
	tt4 := static.NewTemplateTranslation(channel2Ref, i18n.Locale("en"), "Hello {{1}}", 1, "", map[string][]assets.TemplateParam{})
	template := flows.NewTemplate(static.NewTemplate("c520cbda-e118-440f-aaf6-c0485088384f", "greeting", []*static.TemplateTranslation{tt1, tt2, tt3, tt4}))

	tas := flows.NewTemplateAssets([]assets.Template{template})

	tcs := []struct {
		channel   *flows.Channel
		locales   []i18n.Locale
		variables []string
		expected  string
	}{
		{
			channel1,
			[]i18n.Locale{"eng-US", "spa-CO"},
			[]string{"Chef"},
			"Hello Chef",
		},
		{
			channel1,
			[]i18n.Locale{"eng", "spa-CO"},
			[]string{"Chef"},
			"Hello Chef",
		},
		{
			channel1,
			[]i18n.Locale{"deu-DE", "spa-ES"},
			[]string{"Chef"},
			"Hola Chef",
		},
		{
			channel1,
			[]i18n.Locale{"deu-DE"},
			[]string{"Chef"},
			"Hello Chef",
		},
		{
			channel2,
			[]i18n.Locale{"eng-US", "spa-ES"},
			[]string{"Chef"},
			"Hello Chef",
		},
		{
			channel3,
			[]i18n.Locale{"eng-US", "spa-ES"},
			[]string{"Chef"},
			"",
		},
	}

	for _, tc := range tcs {
		testID := fmt.Sprintf("channel '%s' and locales %v", tc.channel.Name(), tc.locales)
		tr := template.FindTranslation(tc.channel, tc.locales)
		if tc.expected == "" {
			assert.Nil(t, tr, "unexpected translation found for %s", testID)
		} else {
			if assert.NotNil(t, tr, "expected translation to be found for %s", testID) {
				assert.Equal(t, tc.expected, tr.Substitute(tc.variables), "substition mismatch for %s", testID)
			}
		}
	}

	template = tas.Get(assets.TemplateUUID("c520cbda-e118-440f-aaf6-c0485088384f"))
	assert.NotNil(t, template)
	assert.Equal(t, assets.NewTemplateReference("c520cbda-e118-440f-aaf6-c0485088384f", "greeting"), template.Reference())
	assert.Equal(t, (*assets.TemplateReference)(nil), (*flows.Template)(nil).Reference())
}
