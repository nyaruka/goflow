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
		"eng": &flows.MsgContent{Text: "Hello"},
		"fra": &flows.MsgContent{Text: "Bonjour"},
		"spa": &flows.MsgContent{Text: "Hola"},
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
}

func TestMsgTemplating(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	templateRef := assets.NewTemplateReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Affirmation")

	msgTemplating := flows.NewMsgTemplating(
		templateRef,
		"0162a7f4_dfe4_4c96_be07_854d5dba3b2b",
		[]*flows.TemplatingComponent{{Type: "body/text", Name: "body", Variables: map[string]int{"1": 0, "2": 1}}},
		[]*flows.TemplatingVariable{{Type: "text", Value: "Ryan Lewis"}, {Type: "text", Value: "boy"}},
	)

	assert.Equal(t, templateRef, msgTemplating.Template)
	assert.Equal(t, "0162a7f4_dfe4_4c96_be07_854d5dba3b2b", msgTemplating.Namespace)
	assert.Equal(t, []*flows.TemplatingComponent{{Type: "body/text", Name: "body", Variables: map[string]int{"1": 0, "2": 1}}}, msgTemplating.Components)
	assert.Equal(t, []*flows.TemplatingVariable{{Type: "text", Value: "Ryan Lewis"}, {Type: "text", Value: "boy"}}, msgTemplating.Variables)

	// test marshaling our msg
	marshaled, err := jsonx.Marshal(msgTemplating)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"template":{
		  "name":"Affirmation",
		  "uuid":"61602f3e-f603-4c70-8a8f-c477505bf4bf"
		},
		"namespace":"0162a7f4_dfe4_4c96_be07_854d5dba3b2b",
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
