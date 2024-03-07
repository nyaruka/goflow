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

func TestFindTranslation(t *testing.T) {
	channel1 := test.NewChannel("WhatsApp 1", "+12345", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channel2 := test.NewChannel("WhatsApp 2", "+23456", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channel3 := test.NewChannel("WhatsApp 3", "+34567", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channel1Ref := assets.NewChannelReference(channel1.UUID(), channel1.Name())
	channel2Ref := assets.NewChannelReference(channel2.UUID(), channel2.Name())

	tt1 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("eng"), "", map[string]*static.TemplateComponent{})
	tt2 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("spa-EC"), "", map[string]*static.TemplateComponent{})
	tt3 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("spa-ES"), "", map[string]*static.TemplateComponent{})
	tt4 := static.NewTemplateTranslation(channel2Ref, i18n.Locale("kin"), "", map[string]*static.TemplateComponent{})

	template := flows.NewTemplate(static.NewTemplate("c520cbda-e118-440f-aaf6-c0485088384f", "greeting", []*static.TemplateTranslation{tt1, tt2, tt3, tt4}))
	tas := flows.NewTemplateAssets([]assets.Template{template})

	tcs := []struct {
		channel  *flows.Channel
		locales  []i18n.Locale
		expected i18n.Locale
	}{
		{channel1, []i18n.Locale{"eng-US", "spa-CO"}, "eng"},
		{channel1, []i18n.Locale{"eng", "spa-CO"}, "eng"},
		{channel1, []i18n.Locale{"deu-DE", "spa-ES"}, "spa-ES"},
		{channel1, []i18n.Locale{"deu-DE"}, "eng"},
		{channel2, []i18n.Locale{"eng-US", "spa-ES"}, "kin"},
		{channel3, []i18n.Locale{"eng-US", "spa-ES"}, ""},
	}

	for _, tc := range tcs {
		testID := fmt.Sprintf("channel '%s' and locales %v", tc.channel.Name(), tc.locales)
		tr := template.FindTranslation(tc.channel, tc.locales)

		if tc.expected == "" {
			assert.Nil(t, tr, "unexpected translation found for %s", testID)
		} else {
			assert.Equal(t, tc.expected, tr.Locale(), "translation mismatch for %s", testID)
		}
	}

	template = tas.Get(assets.TemplateUUID("c520cbda-e118-440f-aaf6-c0485088384f"))
	assert.NotNil(t, template)
	assert.Equal(t, assets.NewTemplateReference("c520cbda-e118-440f-aaf6-c0485088384f", "greeting"), template.Reference())
	assert.Equal(t, (*assets.TemplateReference)(nil), (*flows.Template)(nil).Reference())
}

func TestTemplatePreview(t *testing.T) {
	channel := test.NewChannel("WhatsApp", "+12345", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channelRef := assets.NewChannelReference(channel.UUID(), channel.Name())

	tt := static.NewTemplateTranslation(channelRef, i18n.Locale("eng"), "", map[string]*static.TemplateComponent{
		"body": {
			Content_: "Hello {{1}}, {{2}}",
			Params_:  []*static.TemplateParam{{Type_: "text"}},
		},
		"buttons.0": {
			Content_: "Yes",
			Params_:  []*static.TemplateParam{},
		},
		"buttons.1": {
			Content_: "No {{1}}",
			Params_:  []*static.TemplateParam{{Type_: "text"}},
		},
	})

	template := flows.NewTemplate(static.NewTemplate("c520cbda-e118-440f-aaf6-c0485088384f", "greeting", []*static.TemplateTranslation{tt}))
	translation := template.FindTranslation(channel, []i18n.Locale{"eng"})

	tcs := []struct {
		templating map[string][]flows.TemplateParam
		expected   map[string]string
	}{
		{
			map[string][]flows.TemplateParam{}, // no params
			map[string]string{"body": "Hello , ", "buttons.0": "Yes", "buttons.1": "No "},
		},
		{
			map[string][]flows.TemplateParam{
				"body": {{Type: "text", Value: "Bob"}}, // missing 1 param for body
			},
			map[string]string{"body": "Hello Bob, ", "buttons.0": "Yes", "buttons.1": "No "},
		},
		{
			map[string][]flows.TemplateParam{
				"body": {{Type: "text", Value: "Bob"}, {Type: "text", Value: "how are you?"}, {Type: "text", Value: "xxx"}}, // 1 extra param
			},
			map[string]string{"body": "Hello Bob, how are you?", "buttons.0": "Yes", "buttons.1": "No "},
		},
		{
			map[string][]flows.TemplateParam{
				"body":   {{Type: "text", Value: "Bob"}},
				"header": {{Type: "text", Value: "Hi"}}, // extra component ignored
			},
			map[string]string{"body": "Hello Bob, ", "buttons.0": "Yes", "buttons.1": "No "},
		},
	}

	for _, tc := range tcs {
		templating := flows.NewMsgTemplating(template.Reference(), tc.templating, "")

		assert.Equal(t, tc.expected, translation.Preview(templating))
	}
}
