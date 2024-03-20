package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgIn(t *testing.T) {
	msg := flows.NewMsgIn(
		flows.MsgUUID("48c32bd4-ed68-4a21-b540-9da96217b022"),
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference(assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), "My Android"),
		"Hi there",
		[]utils.Attachment{
			utils.Attachment("image/jpeg:https://example.com/test.jpg"),
			utils.Attachment("audio/mp3:https://example.com/test.mp3"),
		},
	)
	msg.SetID(123)
	msg.SetExternalID("EX346436734")

	// test marshaling our msg
	marshaled, err := jsonx.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"uuid":"48c32bd4-ed68-4a21-b540-9da96217b022",
		"id":123,
		"urn":"tel:+1234567890",
		"channel":{"uuid":"61f38f46-a856-4f90-899e-905691784159",
		"name":"My Android"},
		"text":"Hi there",
		"attachments":["image/jpeg:https://example.com/test.jpg",
		"audio/mp3:https://example.com/test.mp3"],
		"external_id":"EX346436734"
	}`), marshaled, "JSON mismatch")

	// test unmarshaling
	msg = &flows.MsgIn{}
	err = utils.UnmarshalAndValidate(marshaled, msg)
	require.NoError(t, err)
	assert.Equal(t, flows.MsgUUID("48c32bd4-ed68-4a21-b540-9da96217b022"), msg.UUID())
	assert.Equal(t, flows.MsgID(123), msg.ID())
	assert.Equal(t, urns.URN("tel:+1234567890"), msg.URN())
	assert.Equal(t, "Hi there", msg.Text())
	assert.Equal(t, assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), msg.Channel().UUID)
	assert.Equal(t, "My Android", msg.Channel().Name)
	assert.Equal(t, "EX346436734", msg.ExternalID())
}

func TestMsgOut(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	msg := flows.NewMsgOut(
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference(assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), "My Android"),
		"Hi there",
		[]utils.Attachment{
			utils.Attachment("image/jpeg:https://example.com/test.jpg"),
			utils.Attachment("audio/mp3:https://example.com/test.mp3"),
		},
		nil,
		nil,
		flows.MsgTopicAgent,
		"eng-US",
		flows.NilUnsendableReason,
	)

	// test marshaling our msg
	marshaled, err := jsonx.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"uuid": "1ae96956-4b34-433e-8d1a-f05fe6923d6d",
		"urn": "tel:+1234567890",
		"channel": {"uuid":"61f38f46-a856-4f90-899e-905691784159", "name":"My Android"},
		"text": "Hi there",
		"attachments": ["image/jpeg:https://example.com/test.jpg", "audio/mp3:https://example.com/test.mp3"],
		"topic": "agent",
		"locale": "eng-US"
	}`), marshaled, "JSON mismatch")
}

func TestIVRMsgOut(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	msg := flows.NewIVRMsgOut(
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
		"uuid": "1ae96956-4b34-433e-8d1a-f05fe6923d6d",
		"urn": "tel:+1234567890",
		"channel": {"uuid":"61f38f46-a856-4f90-899e-905691784159", "name":"My Android"},
		"text": "Hi there",
		"attachments": ["audio:https://example.com/test.mp3"],
		"locale": "eng-US"
	}`), marshaled, "JSON mismatch")
}

func TestBroadcastTranslations(t *testing.T) {
	bcastTrans := flows.BroadcastTranslations{
		"eng": &flows.BroadcastTranslation{Text: "Hello"},
		"fra": &flows.BroadcastTranslation{Text: "Bonjour"},
		"spa": &flows.BroadcastTranslation{Text: "Hola"},
	}
	baseLanguage := i18n.Language("eng")

	assertTranslation := func(contactLanguage i18n.Language, allowedLanguages []i18n.Language, expectedText string, expectedLang i18n.Language) {
		env := envs.NewBuilder().WithAllowedLanguages(allowedLanguages...).Build()
		sa, err := engine.NewSessionAssets(env, static.NewEmptySource(), nil)
		require.NoError(t, err)

		contact := flows.NewEmptyContact(sa, "Bob", contactLanguage, nil)
		trans, lang := bcastTrans.ForContact(env, contact, baseLanguage)

		assert.Equal(t, expectedText, trans.Text)
		assert.Equal(t, expectedLang, lang)
	}

	assertTranslation("eng", []i18n.Language{"eng"}, "Hello", "eng")          // uses contact language
	assertTranslation("fra", []i18n.Language{"eng", "fra"}, "Bonjour", "fra") // uses contact language
	assertTranslation("kin", []i18n.Language{"eng", "spa"}, "Hello", "eng")   // uses default flow language
	assertTranslation("kin", []i18n.Language{"spa", "eng"}, "Hola", "spa")    // uses default flow language
	assertTranslation("kin", []i18n.Language{"kin"}, "Hello", "eng")          // uses base language

	val, err := bcastTrans.Value()
	assert.NoError(t, err)
	assert.JSONEq(t, `{"eng": {"text": "Hello"}, "fra": {"text": "Bonjour"}, "spa": {"text": "Hola"}}`, string(val.([]byte)))

	var bt flows.BroadcastTranslations
	err = bt.Scan([]byte(`{"spa": {"text": "Adios"}}`))
	assert.NoError(t, err)
	assert.Equal(t, flows.BroadcastTranslations{"spa": {Text: "Adios"}}, bt)
}

func TestMsgTemplating(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	templateRef := assets.NewTemplateReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Affirmation")

	msgTemplating := flows.NewMsgTemplating(templateRef, map[string][]flows.TemplatingParam{"body": {{Type: "text", Value: "Ryan Lewis"}, {Type: "text", Value: "boy"}}}, []*flows.TemplatingComponent{{Type: "body", Params: []flows.TemplatingParam{{Type: "text", Value: "Ryan Lewis"}, {Type: "text", Value: "boy"}}}}, "0162a7f4_dfe4_4c96_be07_854d5dba3b2b")

	assert.Equal(t, templateRef, msgTemplating.Template())
	assert.Equal(t, "0162a7f4_dfe4_4c96_be07_854d5dba3b2b", msgTemplating.Namespace())
	assert.Equal(t, map[string][]flows.TemplatingParam{"body": {{Type: "text", Value: "Ryan Lewis"}, {Type: "text", Value: "boy"}}}, msgTemplating.Params())
	assert.Equal(t, []*flows.TemplatingComponent{{Type: "body", Params: []flows.TemplatingParam{{Type: "text", Value: "Ryan Lewis"}, {Type: "text", Value: "boy"}}}}, msgTemplating.Components())

	// test marshaling our msg
	marshaled, err := jsonx.Marshal(msgTemplating)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"namespace":"0162a7f4_dfe4_4c96_be07_854d5dba3b2b",
		"params": {
			"body": [
				{
					"type": "text",
					"value": "Ryan Lewis"
				},
				{
					"type": "text",
					"value": "boy"
				}
			]
		},
		"components":[{
		  "type": "body",
		  "params":[
			{
				"type": "text",
				"value": "Ryan Lewis"
			},
			{
				"type": "text",
				"value": "boy"
			}
		  ]
		}],
		"template":{
		  "name":"Affirmation",
		  "uuid":"61602f3e-f603-4c70-8a8f-c477505bf4bf"
		}
	  }`), marshaled, "JSON mismatch")
}

func TestTemplatingComponentPreview(t *testing.T) {
	tcs := []struct {
		templating      *flows.TemplatingComponent
		component       assets.TemplateComponent
		expectedContent string
		expectedDisplay string
	}{
		{ // 0: no params
			component:       static.NewTemplateComponent("body", "body", "Hello", "", []*static.TemplateParam{}),
			templating:      &flows.TemplatingComponent{Type: "body", Params: []flows.TemplatingParam{}},
			expectedContent: "Hello",
			expectedDisplay: "",
		},
		{ // 1: two params on component and two params in templating
			component:       static.NewTemplateComponent("body", "body", "Hello {{1}} {{2}}", "", []*static.TemplateParam{{Type_: "text"}, {Type_: "text"}}),
			templating:      &flows.TemplatingComponent{Type: "body", Params: []flows.TemplatingParam{{Type: "text", Value: "Dr"}, {Type: "text", Value: "Bob"}}},
			expectedContent: "Hello Dr Bob",
			expectedDisplay: "",
		},
		{ // 2: one less param in templating than on component
			component:       static.NewTemplateComponent("body", "body", "Hello {{1}} {{2}}", "", []*static.TemplateParam{{Type_: "text"}, {Type_: "text"}}),
			templating:      &flows.TemplatingComponent{Type: "body", Params: []flows.TemplatingParam{{Type: "text", Value: "Dr"}}},
			expectedContent: "Hello Dr ",
			expectedDisplay: "",
		},
		{ // 3
			component:       static.NewTemplateComponent("button/quick_reply", "button.0", "{{1}}", "", []*static.TemplateParam{{Type_: "text"}}),
			templating:      &flows.TemplatingComponent{Type: "button/quick_reply", Params: []flows.TemplatingParam{{Type: "text", Value: "Yes"}}},
			expectedContent: "Yes",
			expectedDisplay: "",
		},
		{ // 4: one param for content, one for display
			component:       static.NewTemplateComponent("button/url", "button.0", "example.com?p={{1}}", "{{1}}", []*static.TemplateParam{{Type_: "text"}}),
			templating:      &flows.TemplatingComponent{Type: "button/url", Params: []flows.TemplatingParam{{Type: "text", Value: "123"}, {Type: "text", Value: "Go"}}},
			expectedContent: "example.com?p=123",
			expectedDisplay: "Go",
		},
	}

	for i, tc := range tcs {
		actualContent, actualDisplay := tc.templating.Preview(tc.component)
		assert.Equal(t, tc.expectedContent, actualContent, "content mismatch in test %d", i)
		assert.Equal(t, tc.expectedDisplay, actualDisplay, "display mismatch in test %d", i)
	}
}
