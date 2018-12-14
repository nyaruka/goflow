package actions_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

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

func testActionType(t *testing.T, assetsJSON json.RawMessage, typeName string, testServerURL string) {
	testFile, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", typeName))
	require.NoError(t, err)

	tests := []struct {
		Description     string            `json:"description"`
		NoContact       bool              `json:"no_contact"`
		NoURNs          bool              `json:"no_urns"`
		NoInput         bool              `json:"no_input"`
		Action          json.RawMessage   `json:"action"`
		ValidationError string            `json:"validation_error"`
		Events          []json.RawMessage `json:"events"`
		ContactAfter    json.RawMessage   `json:"contact_after"`
	}{}

	err = json.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	defer utils.SetTimeSource(utils.DefaultTimeSource)
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	for _, tc := range tests {
		utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
		utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(12345))

		testName := fmt.Sprintf("test '%s' for action type '%s'", tc.Description, typeName)

		// create unstarted session from our assets
		session, err := test.CreateSession(assetsJSON, testServerURL)
		require.NoError(t, err)

		// read the action to be tested
		action, err := actions.ReadAction(tc.Action)
		require.NoError(t, err, "error loading action in %s", testName)
		assert.Equal(t, typeName, action.Type())

		// get a suitable "holder" flow
		var flowUUID assets.FlowUUID
		if len(action.AllowedFlowTypes()) == 1 && action.AllowedFlowTypes()[0] == flows.FlowTypeVoice {
			flowUUID = assets.FlowUUID("7a84463d-d209-4d3e-a0ff-79f977cd7bd0")
		} else {
			flowUUID = assets.FlowUUID("bead76f5-dac4-4c9d-996c-c62b326e8c0a")
		}

		flow, err := session.Assets().Flows().Get(flowUUID)
		require.NoError(t, err)

		// if this action is expected to fail validation, check that
		err = action.Validate(session.Assets(), flows.NewValidationContext())
		if tc.ValidationError != "" {
			assert.EqualError(t, err, tc.ValidationError, "validation error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected validation error in %s", testName)
		}

		// if not, add it to our flow
		flow.Nodes()[0].AddAction(action)

		// optionally load our contact
		var contact *flows.Contact
		if !tc.NoContact {
			contact, err = flows.ReadContact(session.Assets(), json.RawMessage(contactJSON), true)
			require.NoError(t, err)

			// optionally give our contact some URNs
			if !tc.NoURNs {
				channel, _ := session.Assets().Channels().Get("57f1078f-88aa-46f4-a59a-948a5739c03d")
				contact.AddURN(flows.NewContactURN(urns.URN("tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d&id=123"), channel))
				contact.AddURN(flows.NewContactURN(urns.URN("twitterid:54784326227#nyaruka"), nil))
			}
		}

		var trigger flows.Trigger
		ignoreEventCount := 0
		if tc.NoInput {
			trigger = triggers.NewManualTrigger(utils.NewDefaultEnvironment(), flow.Reference(), contact, nil, utils.Now())
		} else {
			msg := flows.NewMsgIn(flows.MsgUUID("aa90ce99-3b4d-44ba-b0ca-79e63d9ed842"), urns.URN("tel:+12065551212"), nil, "Hi everybody", nil)
			trigger = triggers.NewMsgTrigger(utils.NewDefaultEnvironment(), flow.Reference(), contact, msg, nil, utils.Now())
			ignoreEventCount = 1 // need to ignore the msg_received event this trigger creates
		}

		_, err = session.Start(trigger)
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
		actionJSON, err := json.Marshal(action)
		test.AssertEqualJSON(t, tc.Action, actionJSON, "marshal mismatch in %s", testName)
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
