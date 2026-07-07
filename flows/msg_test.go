package flows_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgIn(t *testing.T) {
	msg := core.NewMsgIn(
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference(assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), "My Android"),
		"Hi there",
		[]utils.Attachment{
			utils.Attachment("image/jpeg:https://example.com/test.jpg"),
			utils.Attachment("audio/mp3:https://example.com/test.mp3"),
		},
		"EX346436734",
	)

	// test marshaling our msg
	marshaled, err := jsonx.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"urn":"tel:+1234567890",
		"channel":{"uuid":"61f38f46-a856-4f90-899e-905691784159",
		"name":"My Android"},
		"text":"Hi there",
		"attachments":["image/jpeg:https://example.com/test.jpg",
		"audio/mp3:https://example.com/test.mp3"],
		"external_id":"EX346436734"
	}`), marshaled, "JSON mismatch")

	// test unmarshaling
	msg = &core.MsgIn{}
	err = utils.UnmarshalAndValidate(marshaled, msg)
	require.NoError(t, err)
	assert.Equal(t, urns.URN("tel:+1234567890"), msg.URN())
	assert.Equal(t, "Hi there", msg.Text())
	assert.Equal(t, assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), msg.Channel().UUID)
	assert.Equal(t, "My Android", msg.Channel().Name)
	assert.Equal(t, "EX346436734", msg.ExternalID())
}

func TestMsgOut(t *testing.T) {
	test.MockUniverse()

	msg := core.NewMsgOut(
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference(assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), "My Android"),
		&core.MsgContent{
			Text:        "Hi there",
			Attachments: []utils.Attachment{"image/jpeg:https://example.com/test.jpg", "audio/mp3:https://example.com/test.mp3"},
		},
		nil,
		"eng-US",
		"",
	)

	// test marshaling our msg
	marshaled, err := jsonx.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"urn": "tel:+1234567890",
		"channel": {"uuid":"61f38f46-a856-4f90-899e-905691784159", "name":"My Android"},
		"text": "Hi there",
		"attachments": ["image/jpeg:https://example.com/test.jpg", "audio/mp3:https://example.com/test.mp3"],
		"locale": "eng-US"
	}`), marshaled, "JSON mismatch")
}

func TestIVRMsgOut(t *testing.T) {
	test.MockUniverse()

	msg := core.NewIVRMsgOut(
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference(assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), "My Android"),
		"Hi there",
		"https://example.com/test.mp3",
		"eng-US",
	)

	// test marshaling our msg
	marshaled, err := jsonx.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"urn": "tel:+1234567890",
		"channel": {"uuid":"61f38f46-a856-4f90-899e-905691784159", "name":"My Android"},
		"text": "Hi there",
		"attachments": ["audio:https://example.com/test.mp3"],
		"locale": "eng-US"
	}`), marshaled, "JSON mismatch")
}

func TestMsgContent(t *testing.T) {
	assert.True(t, (&core.MsgContent{}).Empty())
	assert.False(t, (&core.MsgContent{Text: "hi"}).Empty())
	assert.False(t, (&core.MsgContent{Attachments: []utils.Attachment{"image:https://test.jpg"}}).Empty())
	assert.False(t, (&core.MsgContent{QuickReplies: []core.QuickReply{{Text: "Ok"}}}).Empty())

	// can unmarshal from object
	var c core.MsgContent
	err := json.Unmarshal([]byte(`{"text": "test1", "attachments": ["image:https://test.jpg"]}`), &c)
	assert.NoError(t, err)
	assert.Equal(t, "test1", c.Text)
	assert.Equal(t, []utils.Attachment{"image:https://test.jpg"}, c.Attachments)
}

func TestBroadcastTranslations(t *testing.T) {
	tcs := []struct {
		env             envs.Environment
		translations    core.BroadcastTranslations
		baseLanguage    i18n.Language
		contactLanguage i18n.Language
		expectedContent *core.MsgContent
		expectedLocale  i18n.Locale
	}{
		{ // 0: uses contact language
			env: envs.NewBuilder().WithAllowedLanguages("eng", "spa").WithDefaultCountry("US").Build(),
			translations: core.BroadcastTranslations{
				"eng": &core.MsgContent{Text: "Hello"},
				"spa": &core.MsgContent{Text: "Hola"},
			},
			baseLanguage:    "eng",
			contactLanguage: "spa",
			expectedContent: &core.MsgContent{Text: "Hola"},
			expectedLocale:  "spa-US",
		},
		{ // 1: ignores contact language because it's not in allowed languages, uses env default
			env: envs.NewBuilder().WithAllowedLanguages("eng", "spa").WithDefaultCountry("RW").Build(),
			translations: core.BroadcastTranslations{
				"eng": &core.MsgContent{Text: "Hello"},
				"kin": &core.MsgContent{Text: "Muraho"},
			},
			baseLanguage:    "eng",
			contactLanguage: "kin",
			expectedContent: &core.MsgContent{Text: "Hello"},
			expectedLocale:  "eng-RW",
		},
		{ // 2: ignores contact language because it's not translations, uses env default
			env: envs.NewBuilder().WithAllowedLanguages("spa", "fra", "eng").WithDefaultCountry("US").Build(),
			translations: core.BroadcastTranslations{
				"eng": &core.MsgContent{Text: "Hello"},
				"spa": &core.MsgContent{Text: "Hola"},
			},
			baseLanguage:    "eng",
			contactLanguage: "kin",
			expectedContent: &core.MsgContent{Text: "Hola"},
			expectedLocale:  "spa-US",
		},
		{ // 3: ignores contact language because it's not translations, uses base language
			env: envs.NewBuilder().WithAllowedLanguages("eng", "spa").WithDefaultCountry("US").Build(),
			translations: core.BroadcastTranslations{
				"fra": &core.MsgContent{Text: "Bonjour"},
			},
			baseLanguage:    "fra",
			contactLanguage: "eng",
			expectedContent: &core.MsgContent{Text: "Bonjour"},
			expectedLocale:  "fra-US",
		},
		{ // 4: merges content from different translations
			env: envs.NewBuilder().WithAllowedLanguages("eng", "spa").WithDefaultCountry("US").Build(),
			translations: core.BroadcastTranslations{
				"eng": &core.MsgContent{Text: "Hello", Attachments: []utils.Attachment{"image/jpeg:https://example.com/hello.jpg"}, QuickReplies: []core.QuickReply{{Text: "Yes"}, {Text: "No"}}},
				"spa": &core.MsgContent{Text: "Hola"},
			},
			baseLanguage:    "eng",
			contactLanguage: "spa",
			expectedContent: &core.MsgContent{Text: "Hola", Attachments: []utils.Attachment{"image/jpeg:https://example.com/hello.jpg"}, QuickReplies: []core.QuickReply{{Text: "Yes"}, {Text: "No"}}},
			expectedLocale:  "spa-US",
		},
		{ // 5: merges content from different translations
			env: envs.NewBuilder().WithAllowedLanguages("eng", "spa").WithDefaultCountry("US").Build(),
			translations: core.BroadcastTranslations{
				"eng": &core.MsgContent{QuickReplies: []core.QuickReply{{Text: "Yes"}, {Text: "No"}}},
				"spa": &core.MsgContent{Attachments: []utils.Attachment{"image/jpeg:https://example.com/hola.jpg"}},
				"kin": &core.MsgContent{Text: "Muraho"},
			},
			baseLanguage:    "kin",
			contactLanguage: "spa",
			expectedContent: &core.MsgContent{Text: "Muraho", Attachments: []utils.Attachment{"image/jpeg:https://example.com/hola.jpg"}, QuickReplies: []core.QuickReply{{Text: "Yes"}, {Text: "No"}}},
			expectedLocale:  "kin-US",
		},
	}

	for i, tc := range tcs {
		sa, err := engine.NewSessionAssets(tc.env, static.NewEmptySource(), nil)
		require.NoError(t, err)

		contact := flows.NewEmptyContact(sa, "Bob", tc.contactLanguage, nil)
		content, locale := flows.TranslationsForContact(tc.env, tc.translations, contact, tc.baseLanguage)

		assert.Equal(t, tc.expectedContent, content, "%d: content mismatch", i)
		assert.Equal(t, tc.expectedLocale, locale, "%d: locale mismatch", i)
	}
}

func TestMsgTemplating(t *testing.T) {
	templateRef := assets.NewTemplateReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Affirmation")

	msgTemplating := core.NewMsgTemplating(
		templateRef,
		[]*core.TemplatingComponent{{Type: "body/text", Name: "body", Variables: map[string]int{"1": 0, "2": 1}}},
		[]*core.TemplatingVariable{{Type: "text", Value: "Ryan Lewis"}, {Type: "text", Value: "boy"}},
	)

	assert.Equal(t, templateRef, msgTemplating.Template)
	assert.Equal(t, []*core.TemplatingComponent{{Type: "body/text", Name: "body", Variables: map[string]int{"1": 0, "2": 1}}}, msgTemplating.Components)
	assert.Equal(t, []*core.TemplatingVariable{{Type: "text", Value: "Ryan Lewis"}, {Type: "text", Value: "boy"}}, msgTemplating.Variables)

	// test marshaling our msg
	marshaled, err := jsonx.Marshal(msgTemplating)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"template":{
		  "name":"Affirmation",
		  "uuid":"61602f3e-f603-4c70-8a8f-c477505bf4bf"
		},
		"components":[
			{
				"name": "body",
				"type": "body/text",
				"variables": {"1": 0, "2": 1}
			}
		],
		"variables": [
			{
				"type": "text",
				"value": "Ryan Lewis"
			},
			{
				"type": "text",
				"value": "boy"
			}
		]
	  }`), marshaled, "JSON mismatch")
}

func TestQuickReplies(t *testing.T) {
	texts := []struct {
		text     string
		expected core.QuickReply
	}{
		{"", core.QuickReply{Type: "text", Text: ""}},
		{"Yes", core.QuickReply{Type: "text", Text: "Yes"}},
		{"Yes<extra>Really", core.QuickReply{Type: "text", Text: "Yes", Extra: "Really"}},
		{"<location>", core.QuickReply{Type: "location"}},
		{"<location>Click", core.QuickReply{Type: "location", Text: "Click"}},
	}
	for _, tc := range texts {
		qr := core.QuickReply{}
		err := qr.UnmarshalText([]byte(tc.text))
		require.NoError(t, err)
		assert.Equal(t, tc.expected, qr)

		marshaled, err := qr.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, tc.text, string(marshaled))
	}

	jsons := []struct {
		json     []byte
		expected core.QuickReply
	}{
		{[]byte(`"Yes"`), core.QuickReply{Type: "text", Text: "Yes"}},
		{[]byte(`"Yes<extra>Really"`), core.QuickReply{Type: "text", Text: "Yes", Extra: "Really"}},
		{[]byte(`{"text": "Yes"}`), core.QuickReply{Text: "Yes"}},
		{[]byte(`{"text": "Yes", "extra": "Really"}`), core.QuickReply{Text: "Yes", Extra: "Really"}},
		{[]byte(`{"type": "location"}`), core.QuickReply{Type: "location"}},
		{[]byte(`{"type": "location", "text": "Click"}`), core.QuickReply{Type: "location", Text: "Click"}},
	}
	for _, tc := range jsons {
		qr := core.QuickReply{}
		err := json.Unmarshal(tc.json, &qr)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, qr)
	}

	// marshaling is always as struct
	assert.Equal(t, []byte(`{"type":"text","text":"Yes"}`), jsonx.MustMarshal(core.QuickReply{Text: "Yes"}))
	assert.Equal(t, []byte(`{"type":"text","text":"Yes","extra":"Really"}`), jsonx.MustMarshal(core.QuickReply{Text: "Yes", Extra: "Really"}))
	assert.Equal(t, []byte(`{"type":"location"}`), jsonx.MustMarshal(core.QuickReply{Type: "location"}))
	assert.Equal(t, []byte(`[{"type":"text","text":"Yes"},{"type":"text","text":"No"}]`), jsonx.MustMarshal([]core.QuickReply{{Text: "Yes"}, {Text: "No"}}))
}
