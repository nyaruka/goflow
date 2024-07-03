package flows_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
)

func TestFindTranslation(t *testing.T) {
	channel1 := test.NewChannel("WhatsApp 1", "+12345", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channel2 := test.NewChannel("WhatsApp 2", "+23456", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channel3 := test.NewChannel("WhatsApp 3", "+34567", []string{"whatsapp"}, []assets.ChannelRole{}, nil)
	channel1Ref := assets.NewChannelReference(channel1.UUID(), channel1.Name())
	channel2Ref := assets.NewChannelReference(channel2.UUID(), channel2.Name())

	tt1 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("eng"), "", []*static.TemplateComponent{}, []*static.TemplateVariable{})
	tt2 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("spa-EC"), "", []*static.TemplateComponent{}, []*static.TemplateVariable{})
	tt3 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("spa-ES"), "", []*static.TemplateComponent{}, []*static.TemplateVariable{})
	tt4 := static.NewTemplateTranslation(channel2Ref, i18n.Locale("kin"), "", []*static.TemplateComponent{}, []*static.TemplateVariable{})

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

func TestTemplateTranslationPreview(t *testing.T) {
	tcs := []struct {
		translation []byte
		variables   []*flows.TemplatingVariable
		expected    *flows.MsgContent
	}{
		{ // 0: empty translation
			translation: []byte(`{
				"channel": {"uuid": "79401ef2-8eb6-48f4-9f9d-0604530b1ac0", "name": "WhatsApp"}, 
				"locale": "eng",
				"components": [],
				"variables": []
			}`),
			variables: []*flows.TemplatingVariable{},
			expected:  &flows.MsgContent{},
		},
		{ // 1: body only
			translation: []byte(`{
				"channel": {"uuid": "79401ef2-8eb6-48f4-9f9d-0604530b1ac0", "name": "WhatsApp"}, 
				"locale": "eng",
				"components": [
					{
						"name": "body",
						"type": "body/text",
						"content": "Hi {{1}}, who's a good {{2}}?",
						"variables": {"1": 0, "2": 1}
					}
				],
				"variables": [
					{"type": "text"},
					{"type": "text"}
				]
			}`),
			variables: []*flows.TemplatingVariable{{Type: "text", Value: "Chef"}, {Type: "text", Value: "boy"}},
			expected:  &flows.MsgContent{Text: "Hi Chef, who's a good boy?"},
		},
		{ // 2: multiple text component types
			translation: []byte(`{
				"channel": {"uuid": "79401ef2-8eb6-48f4-9f9d-0604530b1ac0", "name": "WhatsApp"}, 
				"locale": "eng",
				"components": [
					{
						"name": "header",
						"type": "header/text",
						"content": "Header {{1}}",
						"variables": {"1": 0}
					},
					{
						"name": "body",
						"type": "body/text",
						"content": "Body {{1}}",
						"variables": {"1": 1}
					},
					{
						"name": "footer",
						"type": "footer/text",
						"content": "Footer {{1}}",
						"variables": {"1": 2}
					}
				],
				"variables": [
					{"type": "text"},
					{"type": "text"},
					{"type": "text"}
				]
			}`),
			variables: []*flows.TemplatingVariable{{Type: "text", Value: "A"}, {Type: "text", Value: "B"}, {Type: "text", Value: "C"}},
			expected:  &flows.MsgContent{Text: "Header A\n\nBody B\n\nFooter C"},
		},
		{ // 3: buttons become quick replies
			translation: []byte(`{
				"channel": {"uuid": "79401ef2-8eb6-48f4-9f9d-0604530b1ac0", "name": "WhatsApp"}, 
				"locale": "eng",
				"components": [
					{
						"name": "button.1",
						"type": "button/quick_reply",
						"content": "{{1}}",
						"variables": {"1": 0}
					},
					{
						"name": "button.2",
						"type": "button/quick_reply",
						"content": "{{1}}",
						"variables": {"1": 1}
					}
				],
				"variables": [
					{"type": "text"},
					{"type": "text"}
				]
			}`),
			variables: []*flows.TemplatingVariable{{Type: "text", Value: "Yes"}, {Type: "text", Value: "No"}},
			expected:  &flows.MsgContent{QuickReplies: []string{"Yes", "No"}},
		},
		{ // 4: header image becomes an attachment
			translation: []byte(`{
				"channel": {"uuid": "79401ef2-8eb6-48f4-9f9d-0604530b1ac0", "name": "WhatsApp"}, 
				"locale": "eng",
				"components": [
					{
						"name": "header",
						"type": "header/image",
						"content": "{{1}}",
						"variables": {"1": 0}
					}
				],
				"variables": [
					{"type": "image"}
				]
			}`),
			variables: []*flows.TemplatingVariable{{Type: "image", Value: "image/jpeg:http://example.com/test.jpg"}},
			expected:  &flows.MsgContent{Attachments: []utils.Attachment{"image/jpeg:http://example.com/test.jpg"}},
		},
	}

	for i, tc := range tcs {
		trans := &static.TemplateTranslation{}
		jsonx.MustUnmarshal(tc.translation, trans)

		actual := flows.NewTemplateTranslation(trans).Preview(tc.variables)
		assert.Equal(t, tc.expected, actual, "%d: preview mismatch", i)
	}
}
