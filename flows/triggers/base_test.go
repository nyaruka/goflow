package triggers_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var assetsJSON = `{
	"flows": [
		{
			"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc",
			"name": "Registration",
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
		}
	]
}`

func TestTriggerMarshaling(t *testing.T) {
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(1234))
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	source, err := static.NewStaticSource([]byte(assetsJSON))
	require.NoError(t, err)

	sessionAssets, err := engine.NewSessionAssets(source)
	require.NoError(t, err)

	env := utils.NewDefaultEnvironment()
	contact := flows.NewEmptyContact("Bob", utils.Language("eng"), nil)
	flow := assets.NewFlowReference(assets.FlowUUID("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"), "Registration")
	channel := assets.NewChannelReference("8cd472c4-bb85-459a-8c9a-c04708af799e", "Facebook")
	triggeredOn := time.Date(2018, 10, 20, 9, 49, 30, 1234567890, time.UTC)

	triggerTests := []struct {
		trigger   flows.Trigger
		marshaled string
	}{
		{
			triggers.NewCampaignTrigger(
				env,
				flow,
				contact,
				triggers.NewCampaignEvent("8d339613-f0be-48b7-92ee-155f4c7576f8", triggers.NewCampaignReference("8cd472c4-bb85-459a-8c9a-c04708af799e", "Reminders")),
				triggeredOn,
			),
			`{
				"contact": {
					"created_on": "2018-10-18T14:20:30.000123456Z",
					"id": 0,
					"language": "eng",
					"name": "Bob",
					"urns": [],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
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
			triggers.NewChannelTrigger(
				env,
				flow,
				contact,
				triggers.NewChannelEvent("new-conversation", channel),
				types.NewEmptyXMap(),
				triggeredOn,
			),
			`{
				"contact": {
					"created_on": "2018-10-18T14:20:30.000123456Z",
					"id": 0,
					"language": "eng",
					"name": "Bob",
					"urns": [],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"event": {
					"channel": {
						"name": "Facebook",
						"uuid": "8cd472c4-bb85-459a-8c9a-c04708af799e"
					},
					"type": "new-conversation"
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
			triggers.NewFlowActionTrigger(
				env,
				flow,
				contact,
				json.RawMessage(`{"uuid": "084e4bed-667c-425e-82f7-bdb625e6ec9e"}`),
				triggeredOn,
			),
			`{
				"contact": {
					"created_on": "2018-10-18T14:20:30.000123456Z",
					"id": 0,
					"language": "eng",
					"name": "Bob",
					"urns": [],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
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
			triggers.NewManualTrigger(
				env,
				flow,
				contact,
				types.NewXArray(types.NewXText("foo")),
				triggeredOn,
			),
			`{
				"contact": {
					"created_on": "2018-10-18T14:20:30.000123456Z",
					"id": 0,
					"language": "eng",
					"name": "Bob",
					"urns": [],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"flow": {
					"name": "Registration",
					"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"
				},
				"params": [
					"foo"
				],
				"triggered_on": "2018-10-20T09:49:31.23456789Z",
				"type": "manual"
			}`,
		},
		{
			triggers.NewMsgTrigger(
				env,
				flow,
				contact,
				flows.NewMsgIn(flows.MsgUUID("c8005ee3-4628-4d76-be66-906352cb1935"), urns.URN("tel:+1234567890"), channel, "Hi there", nil),
				triggers.NewKeywordMatch(triggers.KeywordMatchTypeFirstWord, "hi"),
				triggeredOn,
			),
			`{
				"contact": {
					"created_on": "2018-10-18T14:20:30.000123456Z",
					"id": 0,
					"language": "eng",
					"name": "Bob",
					"urns": [],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
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
						"name": "Facebook",
						"uuid": "8cd472c4-bb85-459a-8c9a-c04708af799e"
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
		_, err = triggers.ReadTrigger(sessionAssets, triggerJSON)
		assert.NoError(t, err, "error reading trigger: %s", string(triggerJSON))
	}
}

func TestReadTrigger(t *testing.T) {
	// error if no type field
	_, err := triggers.ReadTrigger(nil, []byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = triggers.ReadTrigger(nil, []byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
