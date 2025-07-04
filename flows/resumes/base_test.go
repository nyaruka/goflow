package resumes_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResumeTypes(t *testing.T) {
	assetsJSON, err := os.ReadFile("testdata/_assets.json")
	require.NoError(t, err)

	typeNames := make([]string, 0)
	for typeName := range resumes.RegisteredTypes() {
		typeNames = append(typeNames, typeName)
	}

	sort.Strings(typeNames)

	for _, typeName := range typeNames {
		testResumeType(t, assetsJSON, typeName)
	}
}

func testResumeType(t *testing.T, assetsJSON []byte, typeName string) {
	testPath := fmt.Sprintf("testdata/%s.json", typeName)
	testFile, err := os.ReadFile(testPath)
	require.NoError(t, err)

	tests := []struct {
		Description   string              `json:"description"`
		FlowUUID      assets.FlowUUID     `json:"flow_uuid"`
		Wait          json.RawMessage     `json:"wait,omitempty"`
		Resume        json.RawMessage     `json:"resume"`
		ReadError     string              `json:"read_error,omitempty"`
		ResumeError   string              `json:"resume_error,omitempty"`
		Events        json.RawMessage     `json:"events,omitempty"`
		RunStatus     flows.RunStatus     `json:"run_status,omitempty"`
		SessionStatus flows.SessionStatus `json:"session_status,omitempty"`
	}{}

	jsonx.MustUnmarshal(testFile, &tests)

	for i, tc := range tests {
		test.MockUniverse()

		testName := fmt.Sprintf("test '%s' for resume type '%s'", tc.Description, typeName)

		testAssetsJSON := assetsJSON
		if tc.Wait != nil {
			testAssetsJSON = test.JSONReplace(assetsJSON, []string{"flows", "[0]", "nodes", "[0]", "router", "wait"}, tc.Wait)
		}

		// create session assets
		sa, err := test.CreateSessionAssets(testAssetsJSON, "")
		require.NoError(t, err, "unable to create session assets in %s", testName)

		resume, err := resumes.Read(sa, tc.Resume, assets.PanicOnMissing)

		if tc.ReadError != "" {
			rootErr := test.RootError(err)
			assert.EqualError(t, rootErr, tc.ReadError, "read error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected read error in %s", testName)
		}

		flow, err := sa.Flows().Get(tc.FlowUUID)
		require.NoError(t, err)

		// start a waiting session
		env := envs.NewBuilder().Build()
		eng := engine.NewBuilder().Build()
		contact := flows.NewEmptyContact(sa, "Bob", i18n.Language("eng"), nil)
		tb := triggers.NewBuilder(flow.Reference(false)).Manual()
		var call *flows.Call
		if flow.Type() == flows.FlowTypeVoice {
			channel := sa.Channels().Get("a78930fe-6a40-4aa8-99c3-e61b02f45ca1")
			call = flows.NewCall("01978a2f-ad9a-7f2e-ad44-6e7547078cec", channel, urns.URN("tel:+12065551212"))
		}
		trigger := tb.Build()
		session, _, err := eng.NewSession(context.Background(), sa, env, contact, trigger, call)
		require.NoError(t, err)
		require.Equal(t, flows.SessionStatusWaiting, session.Status())

		// resume with our resume...
		sprint, err := session.Resume(context.Background(), resume)

		actual := tc
		actual.Resume = jsonx.MustMarshal(resume) // re-marshal the resume
		actual.RunStatus = session.Runs()[0].Status()
		actual.SessionStatus = session.Status()

		if err != nil {
			actual.ResumeError = err.Error()
		} else {
			actual.Events = jsonx.MustMarshal(sprint.Events())
		}

		if !test.UpdateSnapshots {
			// check resume marshalled correctly
			test.AssertEqualJSON(t, tc.Resume, actual.Resume, "marshal mismatch in %s", testName)

			// check statuses
			assert.Equal(t, tc.RunStatus, actual.RunStatus, "run status mismatch in %s", testName)
			assert.Equal(t, tc.SessionStatus, actual.SessionStatus, "session status mismatch in %s", testName)

			// check error or events generated by resume
			if actual.ResumeError != "" {
				assert.Equal(t, tc.ResumeError, actual.ResumeError, "resume error mismatch in %s", testName)
			} else {
				test.AssertEqualJSON(t, tc.Events, actual.Events, "events mismatch in %s", testName)
			}
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

func TestReadResume(t *testing.T) {
	env := envs.NewBuilder().Build()

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	sessionAssets, err := engine.NewSessionAssets(env, static.NewEmptySource(), nil)
	require.NoError(t, err)

	// error if no type field
	_, err = resumes.Read(sessionAssets, []byte(`{"foo": "bar"}`), missing)
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = resumes.Read(sessionAssets, []byte(`{"type": "do_the_foo", "foo": "bar"}`), missing)
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}

func TestResumeContext(t *testing.T) {
	env := envs.NewBuilder().Build()

	var resume flows.Resume = resumes.NewMsg(
		events.NewMsgReceived(flows.NewMsgIn(urns.URN("tel:1234567890"), nil, "Hello", nil, "SMS1234")),
	)

	assert.Equal(t, map[string]types.XValue{
		"type": types.NewXText("msg"),
		"dial": nil,
	}, resume.Context(env))

	resume = resumes.NewDial(events.NewDialEnded(flows.NewDial(flows.DialStatusNoAnswer, 5)))
	context := resume.Context(env)

	assert.Equal(t, types.NewXText("dial"), context["type"])
	assert.NotNil(t, context["dial"])
}
