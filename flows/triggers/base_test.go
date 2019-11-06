package triggers_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var assetsJSON = `{
	"flows": [
		{
			"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc",
			"name": "Registration",
            "spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        }
	],
	"channels": [
		{
			"uuid": "8cd472c4-bb85-459a-8c9a-c04708af799e",
			"name": "Facebook",
			"address": "23532562626",
			"schemes": ["facebook"],
			"roles": ["send", "receive"]
		},
		{
            "uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7",
            "name": "Nexmo",
            "address": "+12345672222",
            "schemes": ["tel"],
            "roles": ["send", "receive"]
        }
	]
}`

func TestTriggerMarshaling(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 10, 20, 9, 49, 30, 1234567890, time.UTC)))

	uuids.SetGenerator(uuids.NewSeededGenerator(1234))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(source, nil)
	require.NoError(t, err)

	env := envs.NewBuilder().Build()
	flow := assets.NewFlowReference(assets.FlowUUID("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"), "Registration")
	channel := assets.NewChannelReference("3a05eaf5-cb1b-4246-bef1-f277419c83a7", "Nexmo")

	contact := flows.NewEmptyContact(sa, "Bob", envs.Language("eng"), nil)
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+12065551212"), nil))

	triggerTests := []struct {
		trigger   flows.Trigger
		marshaled string
	}{
		{
			triggers.NewCampaign(
				env,
				flow,
				contact,
				triggers.NewCampaignEvent("8d339613-f0be-48b7-92ee-155f4c7576f8", triggers.NewCampaignReference("8cd472c4-bb85-459a-8c9a-c04708af799e", "Reminders")),
			),
			`{
				"contact": {
					"created_on": "2018-10-20T09:49:31.23456789Z",
					"language": "eng",
					"name": "Bob",
					"urns": ["tel:+12065551212"],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"max_value_length": 640,
					"number_format": {
						"decimal_symbol": ".",
						"digit_grouping_symbol": ","
					},
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"event": {
					"campaign": {
						"name": "Reminders",
						"uuid": "8cd472c4-bb85-459a-8c9a-c04708af799e"
					},
					"uuid": "8d339613-f0be-48b7-92ee-155f4c7576f8"
				},
				"flow": {
					"name": "Registration",
					"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"
				},
				"triggered_on": "2018-10-20T09:49:31.23456789Z",
				"type": "campaign"
			}`,
		},
		{
			triggers.NewChannel(
				env,
				flow,
				contact,
				triggers.NewChannelEvent(triggers.ChannelEventTypeNewConversation, channel),
				nil,
			),
			`{
				"contact": {
					"created_on": "2018-10-20T09:49:31.23456789Z",
					"language": "eng",
					"name": "Bob",
					"urns": ["tel:+12065551212"],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"max_value_length": 640,
					"number_format": {
						"decimal_symbol": ".",
						"digit_grouping_symbol": ","
					},
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"event": {
					"channel": {
						"name": "Nexmo",
						"uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7"
					},
					"type": "new_conversation"
				},
				"flow": {
					"name": "Registration",
					"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"
				},
				"params": {},
				"triggered_on": "2018-10-20T09:49:31.23456789Z",
				"type": "channel"
			}`,
		},
		{
			triggers.NewFlowAction(
				env,
				flow,
				contact,
				json.RawMessage(`{"uuid": "084e4bed-667c-425e-82f7-bdb625e6ec9e"}`),
			),
			`{
				"contact": {
					"created_on": "2018-10-20T09:49:31.23456789Z",
					"language": "eng",
					"name": "Bob",
					"urns": ["tel:+12065551212"],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"max_value_length": 640,
					"number_format": {
						"decimal_symbol": ".",
						"digit_grouping_symbol": ","
					},
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"flow": {
					"name": "Registration",
					"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"
				},
				"run_summary": {
					"uuid": "084e4bed-667c-425e-82f7-bdb625e6ec9e"
				},
				"triggered_on": "2018-10-20T09:49:31.23456789Z",
				"type": "flow_action"
			}`,
		},
		{
			triggers.NewIncomingCall(
				env,
				flow,
				contact,
				urns.URN("tel:+12065551212"),
				channel,
			),
			`{
				"connection": {
					"channel": {
						"name": "Nexmo",
						"uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7"
					},
					"urn": "tel:+12065551212"
				},
				"contact": {
					"created_on": "2018-10-20T09:49:31.23456789Z",
					"language": "eng",
					"name": "Bob",
					"urns": ["tel:+12065551212"],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"max_value_length": 640,
					"number_format": {
						"decimal_symbol": ".",
						"digit_grouping_symbol": ","
					},
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"event": {
					"channel": {
						"name": "Nexmo",
						"uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7"
					},
					"type": "incoming_call"
				},
				"flow": {
					"name": "Registration",
					"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"
				},
				"params": {},
				"triggered_on": "2018-10-20T09:49:31.23456789Z",
				"type": "channel"
			}`,
		},
		{
			triggers.NewManual(
				env,
				flow,
				contact,
				types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")}),
			),
			`{
				"contact": {
					"created_on": "2018-10-20T09:49:31.23456789Z",
					"language": "eng",
					"name": "Bob",
					"urns": ["tel:+12065551212"],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"max_value_length": 640,
					"number_format": {
						"decimal_symbol": ".",
						"digit_grouping_symbol": ","
					},
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"flow": {
					"name": "Registration",
					"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"
				},
				"params": {
					"foo": "bar"
				},
				"triggered_on": "2018-10-20T09:49:31.23456789Z",
				"type": "manual"
			}`,
		},
		{
			triggers.NewManualVoice(
				env,
				flow,
				contact,
				flows.NewConnection(channel, "tel:+12065551212"),
				types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")}),
			),
			`{
				"connection": {
					"channel": {
						"name": "Nexmo",
						"uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7"
					},
					"urn": "tel:+12065551212"
				},
				"contact": {
					"created_on": "2018-10-20T09:49:31.23456789Z",
					"language": "eng",
					"name": "Bob",
					"urns": ["tel:+12065551212"],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"max_value_length": 640,
					"number_format": {
						"decimal_symbol": ".",
						"digit_grouping_symbol": ","
					},
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"flow": {
					"name": "Registration",
					"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"
				},
				"params": {
					"foo": "bar"
				},
				"triggered_on": "2018-10-20T09:49:31.23456789Z",
				"type": "manual"
			}`,
		},
		{
			triggers.NewMsg(
				env,
				flow,
				contact,
				flows.NewMsgIn(flows.MsgUUID("c8005ee3-4628-4d76-be66-906352cb1935"), urns.URN("tel:+1234567890"), channel, "Hi there", nil),
				triggers.NewKeywordMatch(triggers.KeywordMatchTypeFirstWord, "hi"),
			),
			`{
				"contact": {
					"created_on": "2018-10-20T09:49:31.23456789Z",
					"language": "eng",
					"name": "Bob",
					"urns": ["tel:+12065551212"],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"max_value_length": 640,
					"number_format": {
						"decimal_symbol": ".",
						"digit_grouping_symbol": ","
					},
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"flow": {
					"name": "Registration",
					"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"
				},
				"keyword_match": {
					"keyword": "hi",
					"type": "first_word"
				},
				"msg": {
					"channel": {
						"name": "Nexmo",
						"uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7"
					},
					"text": "Hi there",
					"urn": "tel:+1234567890",
					"uuid": "c8005ee3-4628-4d76-be66-906352cb1935"
				},
				"triggered_on": "2018-10-20T09:49:31.23456789Z",
				"type": "msg"
			}`,
		},
	}

	for _, tc := range triggerTests {
		triggerJSON, err := json.Marshal(tc.trigger)
		assert.NoError(t, err)

		test.AssertEqualJSON(t, []byte(tc.marshaled), triggerJSON, "trigger JSON mismatch")

		// then try to read from the JSON
		_, err = triggers.ReadTrigger(sa, triggerJSON, assets.PanicOnMissing)
		assert.NoError(t, err, "error reading trigger: %s", string(triggerJSON))
	}
}

func TestReadTrigger(t *testing.T) {
	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	sessionAssets, err := engine.NewSessionAssets(static.NewEmptySource(), nil)
	require.NoError(t, err)

	// error if no type field
	_, err = triggers.ReadTrigger(sessionAssets, []byte(`{"foo": "bar"}`), missing)
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = triggers.ReadTrigger(sessionAssets, []byte(`{"type": "do_the_foo", "foo": "bar"}`), missing)
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}

func TestTriggerSessionInitialization(t *testing.T) {
	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(source, nil)
	require.NoError(t, err)

	defaultEnv := envs.NewBuilder().Build()
	env := envs.NewBuilder().WithDateFormat(envs.DateFormatMonthDayYear).Build()

	flow := assets.NewFlowReference(assets.FlowUUID("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"), "Registration")

	contact := flows.NewEmptyContact(sa, "Bob", envs.Language("eng"), nil)
	contact.AddURN(flows.NewContactURN(urns.URN("tel:+12065551212"), nil))

	params := types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")})

	trigger := triggers.NewManual(
		env,
		flow,
		contact,
		params,
	)

	assert.Equal(t, triggers.TypeManual, trigger.Type())
	assert.Equal(t, env, trigger.Environment())
	assert.Equal(t, contact, trigger.Contact())
	assert.Nil(t, trigger.Connection())
	assert.Equal(t, params, trigger.Params())

	eng := engine.NewBuilder().Build()
	session, _, err := eng.NewSession(sa, trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.FlowTypeMessaging, session.Type())
	assert.Equal(t, contact, session.Contact())
	assert.Equal(t, env, session.Environment())
	assert.Equal(t, flow, session.Runs()[0].FlowReference())

	// contact, environment and params are optional
	trigger = triggers.NewManual(
		nil,
		flow,
		nil,
		nil,
	)

	assert.Equal(t, triggers.TypeManual, trigger.Type())
	assert.Nil(t, trigger.Environment())
	assert.Nil(t, trigger.Contact())
	assert.Equal(t, types.XObjectEmpty, trigger.Params())

	session, _, err = eng.NewSession(sa, trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.FlowTypeMessaging, session.Type())
	assert.Nil(t, session.Contact())
	assert.Equal(t, defaultEnv, session.Environment()) // uses defaults
}
