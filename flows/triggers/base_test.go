package triggers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTriggerTypes(t *testing.T) {
	assetsJSON, err := os.ReadFile("testdata/_assets.json")
	require.NoError(t, err)

	typeNames := make([]string, 0)
	for typeName := range triggers.RegisteredTypes() {
		typeNames = append(typeNames, typeName)
	}

	sort.Strings(typeNames)

	for _, typeName := range typeNames {
		testTriggerType(t, assetsJSON, typeName)
	}
}

func testTriggerType(t *testing.T, assetsJSON []byte, typeName string) {
	testPath := fmt.Sprintf("testdata/%s.json", typeName)
	testFile, err := os.ReadFile(testPath)
	require.NoError(t, err)

	tests := []struct {
		Description string          `json:"description"`
		Trigger     json.RawMessage `json:"trigger"`
		ReadError   string          `json:"read_error,omitempty"`
		Context     json.RawMessage `json:"context,omitempty"`
	}{}

	jsonx.MustUnmarshal(testFile, &tests)

	defer dates.SetNowFunc(time.Now)
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	for i, tc := range tests {
		dates.SetNowFunc(dates.NewFixedNow(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
		uuids.SetGenerator(uuids.NewSeededGenerator(12345, time.Now))

		testName := fmt.Sprintf("test '%s' for trigger type '%s'", tc.Description, typeName)

		// create session assets
		sa, err := test.CreateSessionAssets(assetsJSON, "")
		require.NoError(t, err, "unable to create session assets in %s", testName)

		contact, err := flows.ReadContact(sa, []byte(`{
                "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
                "name": "Bob",
                "status": "active",
                "created_on": "2018-01-01T12:00:00Z"
        }`), assets.PanicOnMissing)
		require.NoError(t, err)

		trigger, err := triggers.Read(sa, tc.Trigger, assets.PanicOnMissing)

		if tc.ReadError != "" {
			rootErr := test.RootError(err)
			assert.EqualError(t, rootErr, tc.ReadError, "read error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected read error in %s", testName)
		}

		// start a session with this trigger
		env := envs.NewBuilder().Build()
		eng := engine.NewBuilder().Build()
		session, _, err := eng.NewSession(context.Background(), sa, env, contact, trigger, nil)
		assert.NoError(t, err)

		assert.Equal(t, flows.FlowTypeMessaging, session.Type())
		assert.NotNil(t, session.Environment())

		// clone test case and populate with actual values
		actual := tc

		log := test.NewEventLog()
		actualContextJSON, _ := session.Runs()[0].EvaluateTemplate(`@(json(trigger))`, log.Log)
		assert.NoError(t, err)
		actual.Context = []byte(actualContextJSON)

		// re-marshal the trigger
		actual.Trigger, err = jsonx.Marshal(trigger)
		require.NoError(t, err)

		if !test.UpdateSnapshots {
			// check context representation
			test.AssertEqualJSON(t, tc.Context, actual.Context, "context mismatch in %s", testName)

			// check marshalled
			test.AssertEqualJSON(t, tc.Trigger, actual.Trigger, "marshal mismatch in %s", testName)
		} else {
			tests[i] = actual
		}
	}

	if test.UpdateSnapshots {
		actualJSON, err := jsonx.MarshalPretty(tests)
		require.NoError(t, err)

		err = os.WriteFile(testPath, actualJSON, 0666)
		require.NoError(t, err)
	}
}

var assetsJSON = `{
    "campaigns": [
        {
            "uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe",
            "name": "Reminders"
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
            "address": "+16055742523",
            "schemes": ["tel"],
            "roles": ["send", "receive"]
        }
    ],
    "flows": [
        {
            "uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc",
            "name": "Registration",
            "spec_version": "13.0.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        }
    ],
    "optins": [
        {
            "uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
            "name": "Joke Of The Day"
        }
    ],
    "users": [
        {
            "uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44",
            "name": "Bob McTickets",
            "email": "bob@nyaruka.com"
        }
    ]
}`

func TestTriggerMarshaling(t *testing.T) {
	test.MockUniverse()

	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc", "Registration")
	nexmo := sa.Channels().Get("3a05eaf5-cb1b-4246-bef1-f277419c83a7")
	channel := assets.NewChannelReference("3a05eaf5-cb1b-4246-bef1-f277419c83a7", "Nexmo")
	reminders := sa.Campaigns().Get("58e9b092-fe42-4173-876c-ff45a14a24fe")
	jotd := sa.OptIns().Get("248be71d-78e9-4d71-a6c4-9981d369e5cb")
	weather := sa.Topics().Get("472a7a73-96cb-4736-b567-056d987cc5b4")
	user := sa.Users().Get("0c78ef47-7d56-44d8-8f57-96e0f30e8f44")
	ticket := flows.NewTicket("276c2e43-d6f9-4c36-8e54-b5af5039acf6", weather, user)
	call := flows.NewCall("0198ce92-ff2f-7b07-b158-b21ab168ebba", nexmo, "tel:+12065551212")

	contact := flows.NewEmptyContact(sa, "Bob", i18n.Language("eng"), nil)
	contact.AddURN("tel:+12065551212")

	eng := engine.NewBuilder().Build()
	session, _, err := eng.NewSession(context.Background(), sa, env, contact, triggers.NewBuilder(flow).Manual().Build(), nil)
	require.NoError(t, err)

	history := flows.NewChildHistory(session)

	// can't create a trigger with invalid JSON
	assert.Panics(t, func() {
		triggers.NewBuilder(flow).FlowAction(history, json.RawMessage(`{"uuid"}`)).Build()
	})
	assert.Panics(t, func() {
		triggers.NewBuilder(flow).FlowAction(history, nil).Build()
	})

	triggerTests := []struct {
		trigger  flows.Trigger
		snapshot string
	}{
		{
			triggers.NewBuilder(flow).
				Call(events.NewCallReceived(call)).
				Build(),
			"call",
		},
		{
			triggers.NewBuilder(flow).
				Call(events.NewCallMissed()).
				Build(),
			"call_missed",
		},
		{
			triggers.NewBuilder(flow).
				Campaign(reminders, events.NewCampaignFired(reminders, "8d339613-f0be-48b7-92ee-155f4c7576f8")).
				Build(),
			"campaign",
		},
		{
			triggers.NewBuilder(flow).
				Chat(events.NewChatStarted(channel, nil)).
				Build(),
			"chat_new_conversation",
		},
		{
			triggers.NewBuilder(flow).
				Chat(events.NewChatStarted(channel, map[string]string{"referrer_id": "acme"})).
				Build(),
			"chat_referral",
		},
		{
			triggers.NewBuilder(flow).
				FlowAction(history, json.RawMessage(`{"uuid": "084e4bed-667c-425e-82f7-bdb625e6ec9e"}`)).
				Build(),
			"flow_action",
		},
		{
			triggers.NewBuilder(flow).
				FlowAction(history, json.RawMessage(`{"uuid": "084e4bed-667c-425e-82f7-bdb625e6ec9e"}`)).
				AsBatch().
				Build(),
			"flow_action_batch",
		},
		{
			triggers.NewBuilder(flow).
				Manual().
				WithParams(types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")})).
				WithUser(user).
				WithOrigin("api").
				AsBatch().
				Build(),
			"manual",
		},
		{
			triggers.NewBuilder(flow).
				Manual().
				Build(),
			"manual_minimal",
		},
		{
			triggers.NewBuilder(flow).
				Msg(events.NewMsgReceived(flows.NewMsgIn(urns.URN("tel:+1234567890"), channel, "Hi there", nil, "SMS1234"))).
				WithMatch(triggers.NewKeywordMatch(triggers.KeywordMatchTypeFirstWord, "hi")).
				Build(),
			"msg",
		},
		{
			triggers.NewBuilder(flow).
				OptIn(jotd, events.NewOptInStarted(jotd, channel)).
				Build(),
			"optin_started",
		},
		{
			triggers.NewBuilder(flow).
				OptIn(jotd, events.NewOptInStopped(jotd, channel)).
				Build(),
			"optin_stopped",
		},
		{
			triggers.NewBuilder(flow).
				Ticket(ticket, events.NewTicketClosed(ticket)).
				Build(),
			"ticket_closed",
		},
	}

	for _, tc := range triggerTests {
		triggerJSON, err := jsonx.MarshalPretty(tc.trigger)
		assert.NoError(t, err)

		test.AssertSnapshot(t, tc.snapshot, string(triggerJSON))

		// then try to read from the JSON
		_, err = triggers.Read(sa, triggerJSON, assets.PanicOnMissing)
		assert.NoError(t, err, "error reading trigger: %s", string(triggerJSON))
	}
}

func TestReadTrigger(t *testing.T) {
	env := envs.NewBuilder().Build()

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	sessionAssets, err := engine.NewSessionAssets(env, static.NewEmptySource(), nil)
	require.NoError(t, err)

	// error if no type field
	_, err = triggers.Read(sessionAssets, []byte(`{"foo": "bar"}`), missing)
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = triggers.Read(sessionAssets, []byte(`{"type": "do_the_foo", "foo": "bar"}`), missing)
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")

	trigger, err := triggers.Read(sessionAssets, []byte(`{
		"type": "channel",
		"flow": {
			"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc",
			"name": "Registration"
		},
		"triggered_on": "2018-10-20T09:49:31.23456789Z",
		"event": {
			"type": "incoming_call",
			"channel": {
				"uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7",
				"name": "Nexmo"
			}
		}
	}`), missing)
	assert.NoError(t, err)
	assert.NotNil(t, trigger)
	assert.Len(t, missingAssets, 0)
}

func TestTriggerSessionInitialization(t *testing.T) {
	env := envs.NewBuilder().WithDateFormat(envs.DateFormatMonthDayYear).Build()

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference(assets.FlowUUID("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"), "Registration")

	contact := flows.NewEmptyContact(sa, "Bob", i18n.Language("eng"), nil)
	contact.AddURN(urns.URN("tel:+12065551212"))

	params := types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")})

	trigger := triggers.NewBuilder(flow).Manual().WithParams(params).Build()

	assert.Equal(t, triggers.TypeManual, trigger.Type())
	assert.Equal(t, params, trigger.Params())

	eng := engine.NewBuilder().Build()
	session, _, err := eng.NewSession(context.Background(), sa, env, contact, trigger, nil)
	require.NoError(t, err)

	assert.Equal(t, flows.FlowTypeMessaging, session.Type())
	assert.Equal(t, contact, session.Contact())
	assert.Equal(t, env, session.Environment())
	assert.Equal(t, flow, session.Runs()[0].FlowReference())

	// params are optional
	trigger = triggers.NewBuilder(flow).Manual().Build()

	assert.Equal(t, triggers.TypeManual, trigger.Type())
	assert.Nil(t, trigger.Params())

	session, _, err = eng.NewSession(context.Background(), sa, env, contact, trigger, nil)
	require.NoError(t, err)

	assert.Equal(t, flows.FlowTypeMessaging, session.Type())
}

func TestTriggerContext(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference(assets.FlowUUID("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"), "Registration")
	user := sa.Users().Get("0c78ef47-7d56-44d8-8f57-96e0f30e8f44")

	params := types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")})
	trigger := triggers.NewBuilder(flow).
		Manual().
		WithParams(params).
		WithUser(user).
		WithOrigin("api").
		AsBatch().
		Build()

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"type":    types.NewXText("manual"),
		"params":  params,
		"keyword": types.XTextEmpty,
		"user": types.NewXObject(map[string]types.XValue{
			"__default__": types.NewXText("Bob McTickets"),
			"email":       types.NewXText("bob@nyaruka.com"),
			"name":        types.NewXText("Bob McTickets"),
			"first_name":  types.NewXText("Bob"),
		}),
		"optin":    nil,
		"origin":   types.NewXText("api"),
		"campaign": nil,
		"ticket":   nil,
	}), flows.Context(env, trigger))
}
