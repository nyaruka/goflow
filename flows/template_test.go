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

	tt1 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("eng"), []*static.TemplateComponent{}, []*static.TemplateVariable{})
	tt2 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("spa-EC"), []*static.TemplateComponent{}, []*static.TemplateVariable{})
	tt3 := static.NewTemplateTranslation(channel1Ref, i18n.Locale("spa-ES"), []*static.TemplateComponent{}, []*static.TemplateVariable{})
	tt4 := static.NewTemplateTranslation(channel2Ref, i18n.Locale("kin"), []*static.TemplateComponent{}, []*static.TemplateVariable{})

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

func TestTemplating(t *testing.T) {
	channel := flows.NewChannel(static.NewChannel("79401ef2-8eb6-48f4-9f9d-0604530b1ac0", "WhatsApp", "1234", []string{"whatsapp"}, nil, nil))

	tcs := []struct {
		template           []byte
		variables          []string
		expectedTemplating *flows.MsgTemplating
		expectedPreview    *flows.MsgContent
	}{
		{ // 0: empty translation
			template: []byte(`{
				"uuid": "4c01c732-e644-421c-af15-f5606c3e05f0",
				"name": "greeting",
				"translations": [
					{
						"channel": {"uuid": "79401ef2-8eb6-48f4-9f9d-0604530b1ac0", "name": "WhatsApp"}, 
						"locale": "eng",
						"components": [],
						"variables": []
					}
				]
			}`),
			variables: []string{},
			expectedTemplating: &flows.MsgTemplating{
				Template:   assets.NewTemplateReference("4c01c732-e644-421c-af15-f5606c3e05f0", "greeting"),
				Components: []*flows.TemplatingComponent{},
				Variables:  []*flows.TemplatingVariable{},
			},
			expectedPreview: &flows.MsgContent{},
		},
		{ // 1: body only
			template: []byte(`{
				"uuid": "4c01c732-e644-421c-af15-f5606c3e05f0",
				"name": "greeting",
				"translations": [
					{
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
					}
				]
			}`),
			variables: []string{"Chef", "boy"},
			expectedTemplating: &flows.MsgTemplating{
				Template: assets.NewTemplateReference("4c01c732-e644-421c-af15-f5606c3e05f0", "greeting"),
				Components: []*flows.TemplatingComponent{
					{
						Name:      "body",
						Type:      "body/text",
						Variables: map[string]int{"1": 0, "2": 1},
					},
				},
				Variables: []*flows.TemplatingVariable{
					{Type: "text", Value: "Chef"},
					{Type: "text", Value: "boy"},
				},
			},
			expectedPreview: &flows.MsgContent{Text: "Hi Chef, who's a good boy?"},
		},
		{ // 2: multiple text component types
			template: []byte(`{
				"uuid": "4c01c732-e644-421c-af15-f5606c3e05f0",
				"name": "greeting",
				"translations": [
					{
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
					}
				]
			}`),
			variables: []string{"A", "B", "C"},
			expectedTemplating: &flows.MsgTemplating{
				Template: assets.NewTemplateReference("4c01c732-e644-421c-af15-f5606c3e05f0", "greeting"),
				Components: []*flows.TemplatingComponent{
					{
						Name:      "header",
						Type:      "header/text",
						Variables: map[string]int{"1": 0},
					},
					{
						Name:      "body",
						Type:      "body/text",
						Variables: map[string]int{"1": 1},
					},
					{
						Name:      "footer",
						Type:      "footer/text",
						Variables: map[string]int{"1": 2},
					},
				},
				Variables: []*flows.TemplatingVariable{
					{Type: "text", Value: "A"},
					{Type: "text", Value: "B"},
					{Type: "text", Value: "C"},
				},
			},
			expectedPreview: &flows.MsgContent{Text: "Header A\n\nBody B\n\nFooter C"},
		},
		{ // 3: buttons become quick replies
			template: []byte(`{
				"uuid": "4c01c732-e644-421c-af15-f5606c3e05f0",
				"name": "greeting",
				"translations": [
					{
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
					}
				]
			}`),
			variables: []string{"Yes", "No"},
			expectedTemplating: &flows.MsgTemplating{
				Template: assets.NewTemplateReference("4c01c732-e644-421c-af15-f5606c3e05f0", "greeting"),
				Components: []*flows.TemplatingComponent{
					{
						Name:      "button.1",
						Type:      "button/quick_reply",
						Variables: map[string]int{"1": 0},
					},
					{
						Name:      "button.2",
						Type:      "button/quick_reply",
						Variables: map[string]int{"1": 1},
					},
				},
				Variables: []*flows.TemplatingVariable{
					{Type: "text", Value: "Yes"},
					{Type: "text", Value: "No"},
				},
			},
			expectedPreview: &flows.MsgContent{QuickReplies: []flows.QuickReply{{Text: "Yes"}, {Text: "No"}}},
		},
		{ // 4: header image becomes an attachment
			template: []byte(`{
				"uuid": "4c01c732-e644-421c-af15-f5606c3e05f0",
				"name": "greeting",
				"translations": [
					{
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
					}
				]
			}`),
			variables: []string{"image/jpeg:http://example.com/test.jpg"},
			expectedTemplating: &flows.MsgTemplating{
				Template: assets.NewTemplateReference("4c01c732-e644-421c-af15-f5606c3e05f0", "greeting"),
				Components: []*flows.TemplatingComponent{
					{
						Name:      "header",
						Type:      "header/image",
						Variables: map[string]int{"1": 0},
					},
				},
				Variables: []*flows.TemplatingVariable{
					{Type: "image", Value: "image/jpeg:http://example.com/test.jpg"},
				},
			},
			expectedPreview: &flows.MsgContent{Attachments: []utils.Attachment{"image/jpeg:http://example.com/test.jpg"}},
		},
		{ // 5: missing variables padded with empty
			template: []byte(`{
				"uuid": "4c01c732-e644-421c-af15-f5606c3e05f0",
				"name": "greeting",
				"translations": [
					{
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
					}
				]
			}`),
			variables: []string{"Chef"},
			expectedTemplating: &flows.MsgTemplating{
				Template: assets.NewTemplateReference("4c01c732-e644-421c-af15-f5606c3e05f0", "greeting"),
				Components: []*flows.TemplatingComponent{
					{
						Name:      "body",
						Type:      "body/text",
						Variables: map[string]int{"1": 0, "2": 1},
					},
				},
				Variables: []*flows.TemplatingVariable{
					{Type: "text", Value: "Chef"},
					{Type: "text", Value: ""},
				},
			},
			expectedPreview: &flows.MsgContent{Text: "Hi Chef, who's a good ?"},
		},
		{ // 6: excess variables ignored
			template: []byte(`{
				"uuid": "4c01c732-e644-421c-af15-f5606c3e05f0",
				"name": "greeting",
				"translations": [
					{
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
					}
				]
			}`),
			variables: []string{"Chef", "boy", "dog"},
			expectedTemplating: &flows.MsgTemplating{
				Template: assets.NewTemplateReference("4c01c732-e644-421c-af15-f5606c3e05f0", "greeting"),
				Components: []*flows.TemplatingComponent{
					{
						Name:      "body",
						Type:      "body/text",
						Variables: map[string]int{"1": 0, "2": 1},
					},
				},
				Variables: []*flows.TemplatingVariable{
					{Type: "text", Value: "Chef"},
					{Type: "text", Value: "boy"},
				},
			},
			expectedPreview: &flows.MsgContent{Text: "Hi Chef, who's a good boy?"},
		},
	}

	for i, tc := range tcs {
		tplAsset := &static.Template{}
		jsonx.MustUnmarshal(tc.template, tplAsset)

		tpl := flows.NewTemplate(tplAsset)
		trans := tpl.FindTranslation(channel, []i18n.Locale{"eng"})
		if assert.NotNil(t, trans, "%d: translation not found", i) {
			// check templating
			templating := tpl.Templating(trans, tc.variables)

			if assert.Equal(t, tc.expectedTemplating, templating, "%d: preview mismatch", i) {
				actualPreview := flows.NewTemplateTranslation(trans).Preview(templating.Variables)
				assert.Equal(t, tc.expectedPreview, actualPreview, "%d: preview mismatch", i)
			}
		}
	}
}
