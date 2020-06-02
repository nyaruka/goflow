package engine_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluateTemplate(t *testing.T) {
	testFile, err := ioutil.ReadFile("testdata/templates.json")
	require.NoError(t, err)

	server := test.NewTestHTTPServer(49992)
	defer server.Close()
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)

	uuids.SetGenerator(uuids.NewSeededGenerator(123456))
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)))

	sessionWithURNs, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)
	sessionWithoutURNs, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyURNs)
	require.NoError(t, err)

	tests := []struct {
		Template   string `json:"template"`
		RedactURNs bool   `json:"redact_urns,omitempty"`

		Output     string          `json:"output,omitempty"`
		OutputJSON json.RawMessage `json:"output_json,omitempty"`
		Error      string          `json:"error,omitempty"`
	}{}

	err = jsonx.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	for i, tc := range tests {
		var run flows.FlowRun
		if tc.RedactURNs {
			run = sessionWithoutURNs.Runs()[0]
		} else {
			run = sessionWithURNs.Runs()[0]
		}

		eval, err := run.EvaluateTemplate(tc.Template)

		// clone test case and populate with actual values
		actual := tc
		if tc.OutputJSON != nil {
			actual.OutputJSON = []byte(eval)
		} else {
			actual.Output = eval
		}
		if err != nil {
			actual.Error = err.Error()
		}

		if !test.UpdateSnapshots {
			if tc.OutputJSON != nil {
				test.AssertEqualJSON(t, tc.OutputJSON, actual.OutputJSON, "output mismatch evaluating template: '%s'", tc.Template)
			} else {
				assert.Equal(t, tc.Output, actual.Output, "output mismatch evaluating template: '%s'", tc.Template)
			}
			assert.Equal(t, tc.Error, actual.Error, "error mismatch evaluating template: '%s'", tc.Template)
		} else {
			tests[i] = actual
		}
	}

	if test.UpdateSnapshots {
		actualJSON, err := jsonx.MarshalPretty(tests)
		require.NoError(t, err)

		err = ioutil.WriteFile("testdata/templates.json", actualJSON, 0666)
		require.NoError(t, err)
	}
}

func BenchmarkEvaluateTemplate(b *testing.B) {
	testFile, err := ioutil.ReadFile("testdata/templates.json")
	require.NoError(b, err)

	session, _, err := test.CreateTestSession("http://localhost", envs.RedactionPolicyNone)
	require.NoError(b, err)

	run := session.Runs()[0]

	tests := []struct {
		Template   string `json:"template"`
		RedactURNs bool   `json:"redact_urns,omitempty"`

		Output string `json:"output,omitempty"`
		Error  string `json:"error,omitempty"`
	}{}
	jsonx.Unmarshal(testFile, &tests)
	require.NoError(b, err)

	for n := 0; n < b.N; n++ {
		for _, tc := range tests {
			run.EvaluateTemplate(tc.Template)
		}
	}
}

func TestReadWithMissingAssets(t *testing.T) {
	// create standard test session and marshal to JSON
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	sessionJSON, err := jsonx.Marshal(session)
	require.NoError(t, err)

	// try to read it back but with no assets
	sessionAssets, err := engine.NewSessionAssets(session.Environment(), static.NewEmptySource(), nil)

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	eng := engine.NewBuilder().Build()
	_, err = eng.ReadSession(sessionAssets, sessionJSON, missing)
	require.NoError(t, err)
	assert.Equal(t, 16, len(missingAssets))
	assert.Equal(t, assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), ""), missingAssets[0])
	assert.Equal(t, assets.NewGroupReference(assets.GroupUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Testers"), missingAssets[1])
	assert.Equal(t, assets.NewGroupReference(assets.GroupUUID("4f1f98fc-27a7-4a69-bbdb-24744ba739a9"), "Males"), missingAssets[2])
	assert.Equal(t, assets.NewFlowReference(assets.FlowUUID("50c3706e-fedb-42c0-8eab-dda3335714b7"), "Registration"), missingAssets[13])
	assert.Equal(t, assets.NewFlowReference(assets.FlowUUID("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"), "Collect Age"), missingAssets[14])
}

func TestRunResuming(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/subflows.json")
	require.NoError(t, err)

	session, _, err := test.CreateSession(assetsJSON, assets.FlowUUID("72162f46-dce3-4798-9f19-384a2447efc5"))
	require.NoError(t, err)

	// each run should be marked as completed
	assert.Equal(t, 3, len(session.Runs()))
	assert.Equal(t, flows.RunStatusCompleted, session.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusCompleted, session.Runs()[1].Status())
	assert.Equal(t, flows.RunStatusCompleted, session.Runs()[2].Status())

	// change the UUID of the third flow so the nter_flow in the second flow will error
	assetsWithoutChildFlow := test.JSONReplace(assetsJSON, []string{"flows", "[2]", "uuid"}, []byte(`"653a3fa3-ff59-4a89-93c3-a8b9486ec479"`))

	session, _, err = test.CreateSession(assetsWithoutChildFlow, assets.FlowUUID("72162f46-dce3-4798-9f19-384a2447efc5"))
	require.NoError(t, err)

	// each run should be marked as failed
	assert.Equal(t, 2, len(session.Runs()))
	assert.Equal(t, flows.RunStatusFailed, session.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusFailed, session.Runs()[1].Status())
}

func TestResumeAfterWaitWithMissingFlowAssets(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("../../test/testdata/runner/subflow.json")
	require.NoError(t, err)

	session1, _, err := test.CreateSession(assetsJSON, assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"))
	require.NoError(t, err)

	assert.Equal(t, flows.SessionStatusWaiting, session1.Status())
	assert.Equal(t, flows.RunStatusActive, session1.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusWaiting, session1.Runs()[1].Status())

	// change the UUID of the child flow so it will effectively be missing
	assetsWithoutChildFlow := test.JSONReplace(assetsJSON, []string{"flows", "[1]", "uuid"}, []byte(`"653a3fa3-ff59-4a89-93c3-a8b9486ec479"`))

	session2, _, err := test.ResumeSession(session1, assetsWithoutChildFlow, "Hello")
	require.NoError(t, err)

	// should have a failed session
	assert.Equal(t, flows.SessionStatusFailed, session2.Status())
	assert.Equal(t, flows.RunStatusActive, session2.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusFailed, session2.Runs()[1].Status())

	// change the UUID of the parent flow so it will effectively be missing
	assetsWithoutParentFlow := test.JSONReplace(assetsJSON, []string{"flows", "[0]", "uuid"}, []byte(`"653a3fa3-ff59-4a89-93c3-a8b9486ec479"`))

	session3, _, err := test.ResumeSession(session1, assetsWithoutParentFlow, "Hello")
	require.NoError(t, err)

	// should have an failed session
	assert.Equal(t, flows.SessionStatusFailed, session3.Status())
	assert.Equal(t, flows.RunStatusActive, session3.Runs()[0].Status())
	assert.Equal(t, flows.RunStatusFailed, session3.Runs()[1].Status())
}

func TestWaitTimeout(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)

	t1 := time.Date(2018, 4, 11, 13, 24, 30, 123456000, time.UTC)
	dates.SetNowSource(dates.NewFixedNowSource(t1))

	assetsJSON, err := ioutil.ReadFile("testdata/timeout_test.json")
	require.NoError(t, err)

	session, sprint, err := test.CreateSession(assetsJSON, assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"))
	require.NoError(t, err)

	require.Equal(t, 1, len(session.Runs()[0].Path()))
	run := session.Runs()[0]

	require.Equal(t, 2, len(sprint.Events()))
	require.Equal(t, "msg_created", sprint.Events()[0].Type())
	require.Equal(t, "msg_wait", sprint.Events()[1].Type())

	// check our wait has a timeout
	waitEvent := run.Events()[1].(*events.MsgWaitEvent)
	require.Equal(t, 600, *waitEvent.TimeoutSeconds)

	_, err = session.Resume(resumes.NewWaitTimeout(nil, nil))
	require.NoError(t, err)

	require.Equal(t, flows.SessionStatusCompleted, session.Status())
	require.Equal(t, 2, len(run.Path()))
	require.Equal(t, 5, len(run.Events()))

	result := run.Results().Get("favorite_color")
	require.Equal(t, "Timeout", result.Category)
	require.Equal(t, "2018-04-11T13:24:30.123456Z", result.Value)
	require.Equal(t, "", result.Input)
}

func TestCurrentContext(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("../../test/testdata/runner/subflow_loop_with_wait.json")
	require.NoError(t, err)

	session, _, err := test.CreateSession(assetsJSON, assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"))
	require.NoError(t, err)

	assert.Equal(t, string(flows.SessionStatusWaiting), string(session.Status()))

	context := session.CurrentContext()
	assert.NotNil(t, context)

	runContext, _ := context.Get("run")
	flowContext, _ := runContext.(*types.XObject).Get("flow")
	flowName, _ := flowContext.(*types.XObject).Get("name")
	assert.Equal(t, types.NewXText("Child flow"), flowName)

	// check we can marshal it
	_, err = jsonx.Marshal(context)
	assert.NoError(t, err)

	// end it
	session.Resume(resumes.NewRunExpiration(nil, nil))
	assert.Equal(t, flows.SessionStatusCompleted, session.Status())

	// can still get context of completed session
	context = session.CurrentContext()
	assert.NotNil(t, context)

	runContext, _ = context.Get("run")
	flowContext, _ = runContext.(*types.XObject).Get("flow")
	flowName, _ = flowContext.(*types.XObject).Get("name")
	assert.Equal(t, types.NewXText("Parent Flow"), flowName)
}
