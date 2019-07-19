package actions_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var contactJSON = `{
	"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
	"name": "Ryan Lewis",
	"language": "eng",
	"timezone": "America/Guayaquil",
	"urns": [],
	"groups": [
		{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Testers"},
		{"uuid": "0ec97956-c451-48a0-a180-1ce766623e31", "name": "Males"}
	],
	"fields": {
		"gender": {
			"text": "Male"
		}
	},
	"created_on": "2018-06-20T11:40:30.123456789-00:00"
}`

func TestActionTypes(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/_assets.json")
	require.NoError(t, err)

	server := test.NewTestHTTPServer(49996)

	for typeName := range actions.RegisteredTypes() {
		testActionType(t, assetsJSON, typeName, server.URL)
	}
}

type inspectionResults struct {
	Templates    []string            `json:"templates"`
	Dependencies []string            `json:"dependencies"`
	Results      []*flows.ResultInfo `json:"results"`
}

func testActionType(t *testing.T, assetsJSON json.RawMessage, typeName string, testServerURL string) {
	testFile, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", typeName))
	require.NoError(t, err)

	tests := []struct {
		Description     string            `json:"description"`
		NoContact       bool              `json:"no_contact"`
		NoURNs          bool              `json:"no_urns"`
		NoInput         bool              `json:"no_input"`
		RedactURNs      bool              `json:"redact_urns"`
		Action          json.RawMessage   `json:"action"`
		Localization    json.RawMessage   `json:"localization"`
		InFlowType      flows.FlowType    `json:"in_flow_type"`
		ReadError       string            `json:"read_error"`
		ValidationError string            `json:"validation_error"`
		SkipValidation  bool              `json:"skip_validation"`
		Events          []json.RawMessage `json:"events"`
		ContactAfter    json.RawMessage   `json:"contact_after"`
		Inspection      json.RawMessage   `json:"inspection"`
	}{}

	err = json.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	defer dates.SetNowSource(dates.DefaultNowSource)
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	for _, tc := range tests {
		dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
		uuids.SetGenerator(uuids.NewSeededGenerator(12345))

		testName := fmt.Sprintf("test '%s' for action type '%s'", tc.Description, typeName)

		// pick a suitable "holder" flow in our assets JSON
		flowIndex := 0
		flowUUID := assets.FlowUUID("bead76f5-dac4-4c9d-996c-c62b326e8c0a")
		if tc.InFlowType == flows.FlowTypeVoice {
			flowIndex = 1
			flowUUID = assets.FlowUUID("7a84463d-d209-4d3e-a0ff-79f977cd7bd0")
		}

		// inject the action into a suitable node's actions in that flow
		actionsPath := []string{"flows", fmt.Sprintf("[%d]", flowIndex), "nodes", "[0]", "actions"}
		actionsJson := []byte(fmt.Sprintf("[%s]", string(tc.Action)))
		assetsJSON = test.JSONReplace(assetsJSON, actionsPath, actionsJson)

		// if we have a localization section, inject that too
		if tc.Localization != nil {
			localizationPath := []string{"flows", fmt.Sprintf("[%d]", flowIndex), "localization", "spa", "ad154980-7bf7-4ab8-8728-545fd6378912"}
			assetsJSON = test.JSONReplace(assetsJSON, localizationPath, tc.Localization)
		}

		// create session assets
		sa, err := test.CreateSessionAssets(assetsJSON, "")
		require.NoError(t, err, "unable to create session assets in %s", testName)

		// now try to read the flow, and if we expect a read error, check that
		flow, err := sa.Flows().Get(flowUUID)
		if tc.ReadError != "" {
			rootErr := errors.Cause(err)
			assert.EqualError(t, rootErr, tc.ReadError, "read error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected read error in %s", testName)
		}

		// if this action is expected to cause a validation error, check that
		err = flow.Validate(sa, nil)
		if tc.ValidationError != "" {
			rootErr := errors.Cause(err)
			assert.EqualError(t, rootErr, tc.ValidationError, "validation error mismatch in %s", testName)
			continue
		} else if !tc.SkipValidation {
			assert.NoError(t, err, "unexpected validation error in %s", testName)
		}

		// optionally load our contact
		var contact *flows.Contact
		if !tc.NoContact {
			contact, err = flows.ReadContact(sa, json.RawMessage(contactJSON), assets.PanicOnMissing)
			require.NoError(t, err)

			// optionally give our contact some URNs
			if !tc.NoURNs {
				channel := sa.Channels().Get("57f1078f-88aa-46f4-a59a-948a5739c03d")
				contact.AddURN(flows.NewContactURN(urns.URN("tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d&id=123"), channel))
				contact.AddURN(flows.NewContactURN(urns.URN("twitterid:54784326227#nyaruka"), nil))
			}

			// and switch their language
			if tc.Localization != nil {
				contact.SetLanguage(envs.Language("spa"))
			}
		}

		envBuilder := envs.NewEnvironmentBuilder().
			WithDefaultLanguage("eng").
			WithAllowedLanguages([]envs.Language{"eng", "spa"}).
			WithDefaultCountry("RW")

		if tc.RedactURNs {
			envBuilder.WithRedactionPolicy(envs.RedactionPolicyURNs)
		}

		env := envBuilder.Build()

		var trigger flows.Trigger
		ignoreEventCount := 0
		if tc.NoInput {
			var connection *flows.Connection
			if flow.Type() == flows.FlowTypeVoice {
				channel := sa.Channels().Get("57f1078f-88aa-46f4-a59a-948a5739c03d")
				connection = flows.NewConnection(channel.Reference(), urns.URN("tel:+12065551212"))
				trigger = triggers.NewManualVoiceTrigger(env, flow.Reference(), contact, connection, nil)
			} else {
				trigger = triggers.NewManualTrigger(env, flow.Reference(), contact, nil)
			}
		} else {
			msg := flows.NewMsgIn(
				flows.MsgUUID("aa90ce99-3b4d-44ba-b0ca-79e63d9ed842"),
				urns.URN("tel:+12065551212"),
				nil,
				"Hi everybody",
				[]utils.Attachment{
					"image/jpeg:http://http://s3.amazon.com/bucket/test.jpg",
					"audio/mp3:http://s3.amazon.com/bucket/test.mp3",
				},
			)
			trigger = triggers.NewMsgTrigger(env, flow.Reference(), contact, msg, nil)
			ignoreEventCount = 1 // need to ignore the msg_received event this trigger creates
		}

		// create session
		eng := engine.NewBuilder().WithDefaultUserAgent("goflow-testing").Build()
		session, _, err := eng.NewSession(sa, trigger)
		require.NoError(t, err)

		// check events are what we expected
		run := session.Runs()[0]
		runEvents := run.Events()
		actualEventsJSON, _ := json.Marshal(runEvents[ignoreEventCount:])
		expectedEventsJSON, _ := json.Marshal(tc.Events)
		test.AssertEqualJSON(t, expectedEventsJSON, actualEventsJSON, "events mismatch in %s", testName)

		// check contact is in the expected state
		if tc.ContactAfter != nil {
			contactJSON, _ := json.Marshal(session.Contact())

			test.AssertEqualJSON(t, tc.ContactAfter, contactJSON, "contact mismatch in %s", testName)
		}

		// try marshaling the action back to JSON
		actionJSON, err := json.Marshal(flow.Nodes()[0].Actions()[0])
		test.AssertEqualJSON(t, tc.Action, actionJSON, "marshal mismatch in %s", testName)

		// finally try inspecting this action
		if tc.Inspection != nil {
			dependencies := flow.ExtractDependencies()
			depStrings := make([]string, len(dependencies))
			for i := range dependencies {
				depStrings[i] = dependencies[i].String()
			}

			actual := &inspectionResults{
				Templates:    flow.ExtractTemplates(),
				Dependencies: depStrings,
				Results:      flow.ExtractResults(),
			}

			actualJSON, _ := json.Marshal(actual)
			test.AssertEqualJSON(t, tc.Inspection, actualJSON, "inspection mismatch in %s", testName)
		}
	}
}

func TestConstructors(t *testing.T) {
	actionUUID := flows.ActionUUID("ad154980-7bf7-4ab8-8728-545fd6378912")

	tests := []struct {
		action flows.Action
		json   string
	}{
		{
			actions.NewAddContactGroupsAction(
				actionUUID,
				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
					assets.NewVariableGroupReference("@(format_location(contact.fields.state)) Members"),
				},
			),
			`{
			"type": "add_contact_groups",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Members"
				}
			]
		}`,
		},
		{
			actions.NewAddContactURNAction(
				actionUUID,
				"tel",
				"+234532626677",
			),
			`{
			"type": "add_contact_urn",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"scheme": "tel",
			"path": "+234532626677"
		}`,
		},
		{
			actions.NewAddInputLabelsAction(
				actionUUID,
				[]*assets.LabelReference{
					assets.NewLabelReference(assets.LabelUUID("3f65d88a-95dc-4140-9451-943e94e06fea"), "Spam"),
					assets.NewVariableLabelReference("@(format_location(contact.fields.state)) Messages"),
				},
			),
			`{
			"type": "add_input_labels",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"labels": [
				{
					"uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
					"name": "Spam"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Messages"
				}
			]
		}`,
		},
		{
			actions.NewCallResthookAction(
				actionUUID,
				"new-registration",
				"My Result",
			),
			`{
			"type": "call_resthook",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"resthook": "new-registration",
			"result_name": "My Result"
		}`,
		},
		{
			actions.NewCallWebhookAction(
				actionUUID,
				"POST",
				"http://example.com/ping",
				map[string]string{
					"Authentication": "Token @fields.token",
				},
				`{"contact_id": 234}`, // body
				"Webhook Response",
			),
			`{
			"type": "call_webhook",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"method": "POST",
			"url": "http://example.com/ping",
			"headers": {
				"Authentication": "Token @fields.token"
			},
			"body": "{\"contact_id\": 234}",
			"result_name": "Webhook Response"
		}`,
		},
		{
			actions.NewPlayAudioAction(
				actionUUID,
				"http://uploads.temba.io/2353262.m4a",
			),
			`{
			"type": "play_audio",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"audio_url": "http://uploads.temba.io/2353262.m4a"
		}`,
		},
		{
			actions.NewSayMsgAction(
				actionUUID,
				"Hi @contact.name, are you ready to complete today's survey?",
				"http://uploads.temba.io/2353262.m4a",
			),
			`{
			"type": "say_msg",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"audio_url": "http://uploads.temba.io/2353262.m4a",
			"text": "Hi @contact.name, are you ready to complete today's survey?"
		}`,
		},
		{
			actions.NewRemoveContactGroupsAction(
				actionUUID,
				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
					assets.NewVariableGroupReference("@(format_location(contact.fields.state)) Members"),
				},
				false,
			),
			`{
			"type": "remove_contact_groups",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				},
				{
					"name_match": "@(format_location(contact.fields.state)) Members"
				}
			]
		}`,
		},
		{
			actions.NewSendBroadcastAction(
				actionUUID,
				"Hi there",
				[]string{"http://example.com/red.jpg"},
				[]string{"Red", "Blue"},
				[]urns.URN{"twitter:nyaruka"},
				[]*flows.ContactReference{
					flows.NewContactReference(flows.ContactUUID("cbe87f5c-cda2-4f90-b5dd-0ac93a884950"), "Bob Smith"),
				},
				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
				},
				nil,
			),
			`{
			"type": "send_broadcast",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"text": "Hi there",
			"attachments": ["http://example.com/red.jpg"],
			"quick_replies": ["Red", "Blue"],
			"urns": ["twitter:nyaruka"],
            "contacts": [
				{
					"uuid": "cbe87f5c-cda2-4f90-b5dd-0ac93a884950",
					"name": "Bob Smith"
				}
			],
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				}
			]
		}`,
		},
		{
			actions.NewSendEmailAction(
				actionUUID,
				[]string{"bob@example.com"},
				"Hi there",
				"So I was thinking...",
			),
			`{
			"type": "send_email",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"addresses": ["bob@example.com"],
			"subject": "Hi there",
			"body": "So I was thinking..."
		}`,
		},
		{
			actions.NewSendMsgAction(
				actionUUID,
				"Hi there",
				[]string{"http://example.com/red.jpg"},
				[]string{"Red", "Blue"},
				true,
			),
			`{
			"type": "send_msg",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"text": "Hi there",
			"attachments": ["http://example.com/red.jpg"],
			"quick_replies": ["Red", "Blue"],
			"all_urns": true
		}`,
		},
		{
			actions.NewSetContactChannelAction(
				actionUUID,
				assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
			),
			`{
			"type": "set_contact_channel",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"channel": {
				"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
				"name": "My Android Phone"
			}
		}`,
		},
		{
			actions.NewSetContactFieldAction(
				actionUUID,
				assets.NewFieldReference("gender", "Gender"),
				"Male",
			),
			`{
			"type": "set_contact_field",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"field": {
				"key": "gender",
				"name": "Gender"
			},
			"value": "Male"
		}`,
		},
		{
			actions.NewSetContactLanguageAction(
				actionUUID,
				"eng",
			),
			`{
			"type": "set_contact_language",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"language": "eng"
		}`,
		},
		{
			actions.NewSetContactNameAction(
				actionUUID,
				"Bob",
			),
			`{
			"type": "set_contact_name",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"name": "Bob"
		}`,
		},
		{
			actions.NewSetContactTimezoneAction(
				actionUUID,
				"Africa/Kigali",
			),
			`{
			"type": "set_contact_timezone",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"timezone": "Africa/Kigali"
		}`,
		},
		{
			actions.NewSetRunResultAction(
				actionUUID,
				"Response 1",
				"yes",
				"Yes",
			),
			`{
			"type": "set_run_result",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"name": "Response 1",
			"value": "yes",
			"category": "Yes"
		}`,
		},
		{
			actions.NewEnterFlowAction(
				actionUUID,
				assets.NewFlowReference(assets.FlowUUID("fece6eac-9127-4343-9269-56e88f391562"), "Parent"),
				true, // terminal
			),
			`{
			"type": "enter_flow",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"flow": {
				"uuid": "fece6eac-9127-4343-9269-56e88f391562",
				"name": "Parent"
			},
			"terminal": true
		}`,
		},
		{
			actions.NewStartSessionAction(
				actionUUID,
				assets.NewFlowReference(assets.FlowUUID("fece6eac-9127-4343-9269-56e88f391562"), "Parent"),
				[]urns.URN{"twitter:nyaruka"},
				[]*flows.ContactReference{
					flows.NewContactReference(flows.ContactUUID("cbe87f5c-cda2-4f90-b5dd-0ac93a884950"), "Bob Smith"),
				},
				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"),
				},
				nil,  // legacy vars
				true, // create new contact
			),
			`{
			"type": "start_session",
			"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
			"flow": {
				"uuid": "fece6eac-9127-4343-9269-56e88f391562",
				"name": "Parent"
			},
			"urns": ["twitter:nyaruka"],
            "contacts": [
				{
					"uuid": "cbe87f5c-cda2-4f90-b5dd-0ac93a884950",
					"name": "Bob Smith"
				}
			],
			"groups": [
				{
					"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
					"name": "Testers"
				}
			],
			"create_contact": true
		}`,
		},
	}

	for _, tc := range tests {
		// test marshaling the action
		actualJSON, err := json.Marshal(tc.action)
		assert.NoError(t, err)

		test.AssertEqualJSON(t, json.RawMessage(tc.json), actualJSON, "new action produced unexpected JSON")
	}
}

func TestReadAction(t *testing.T) {
	// error if no type field
	_, err := actions.ReadAction([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = actions.ReadAction([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
